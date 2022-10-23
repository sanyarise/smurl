package routeropenapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/sanyarise/smurl/internal/infrastructure/api/handler"
	"github.com/sanyarise/smurl/internal/infrastructure/api/helpers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type RouterOpenAPI struct {
	*chi.Mux
	hs     *handler.Handlers
	logger *zap.Logger
	url    string
}

func NewRouterOpenAPI(hs *handler.Handlers, l *zap.Logger, url string) *RouterOpenAPI {
	r := chi.NewRouter()

	ret := &RouterOpenAPI{
		hs:     hs,
		logger: l,
		url:    url,
	}

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Mount("/", Handler(ret))

	swg, err := GetSwagger()
	if err != nil {
		msg := fmt.Sprintf("error:%s", err)
		l.Error("swagger fail",
			zap.String(" ", msg))
	}

	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		enc := json.NewEncoder(w)
		_ = enc.Encode(swg)
	})

	ret.Mux = r
	return ret
}

type Smurl handler.Smurl

func (Smurl) Bind(r *http.Request) error {
	return nil
}

func (Smurl) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Get method displaying the start page
func (roa *RouterOpenAPI) Get(w http.ResponseWriter, r *http.Request) {
	l := roa.logger
	l.Debug("Enter in router func Get()")
	w.WriteHeader(http.StatusOK)

	tmpl, err := template.ParseFiles("./static/home.tmpl")
	if err != nil {
		msg := fmt.Sprintf("error on template parse home page files:%s", err)
		l.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		pErr := ErrorPage(w, "./static/500.tmpl", roa.url, l)
		if pErr != nil {
			msg = fmt.Sprintf("error on error page func: %s", err)
			l.Error(msg)

			render.Render(w, r, ErrRender(err))
		}
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		msg := fmt.Sprintf("error on execute home page files: %s", err)
		l.Error(msg)
		pErr := ErrorPage(w, "./static/500.tmpl", roa.url, l)
		if pErr != nil {
			msg = fmt.Sprintf("error on execute home page files: %s", err)
			l.Error(msg)
			render.Render(w, r, ErrRender(err))
		}
		return
	}
}

// PostCreate creating a minified url
func (roa *RouterOpenAPI) PostCreate(w http.ResponseWriter, r *http.Request) {
	l := roa.logger
	l.Debug("Enter in router func PostCreate()")
	// Reading the long address from the request body
	longURL := r.FormValue("long_url")
	msg := fmt.Sprintf("LongURL: %s", longURL)
	l.Debug(msg)
	// Checking the validity of a long address
	ok := helpers.CheckURL(longURL, l)
	if !ok {
		l.Error("incorrect long url")
		pErr := ErrorPage(w, "./static/400.tmpl", roa.url, l)
		if pErr != nil {
			msg := fmt.Sprint(pErr)
			l.Error(msg)
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("incorrect long url")))
		}
		return
	}
	l.Debug("URL check sucess")
	// Create a request with JSON
	url := roa.url + "create"
	l.Debug(url)
	str := fmt.Sprintf(`{"long_url":"%s"}`, longURL)
	l.Debug(str)
	jsonStr := []byte(str)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		msg := fmt.Sprintf("error on create new request: %s", err)
		l.Error(msg)
		pErr := ErrorPage(w, "./static/500.tmpl", roa.url, l)
		if pErr != nil {
			render.Render(w, r, ErrRender(err))
		}

		return
	}

	hsu := Smurl{LongURL: longURL}
	// Calling the handler to create a reduced url
	hss, err := roa.hs.CreateSmurlHandle(r.Context(), handler.Smurl(hsu))
	if err != nil {
		msg := fmt.Sprintf("create smurl error %s: ", err)
		l.Error(msg)
		pErr := ErrorPage(w, "./static/400.tmpl", roa.url, l)
		if pErr != nil {
			msg := fmt.Sprint(pErr)
			l.Error(msg)
			render.Render(w, r, ErrInvalidRequest(err))
		}
		return
	}
	l.Debug("Create smurl successfull")
	w.WriteHeader(http.StatusCreated)
	// Write the full address to the resulting structure
	hss.SmallURL = roa.url + hss.SmallURL
	// Call the function to render the page with the result
	err = ResultPage(w, "./static/result.tmpl", hss, roa.url, l)
	if err != nil {
		l.Error("",
			zap.Error(err))
		render.Render(w, r, Smurl(hss))
	}
}

// GetSmallUrl following the reduced url received from the query string
func (roa *RouterOpenAPI) GetSmallUrl(w http.ResponseWriter, r *http.Request, u string) {
	l := roa.logger
	l.Debug("Enter in router func GetSmallUrl()")
	// Cut off incorrect addresses
	if len(u) != 8 {
		l.Debug(fmt.Sprintf("Incorrect small url %s", u))
		err := ErrorPage(w, "./static/400.tmpl", roa.url, l)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
		}
		return
	}
	// Getting information about IP
	ip := helpers.GetIP(r, l) + "\n"

	// Call the handler to search for a small url,
	// search for the corresponding long url, update
	// statistics
	red, err := roa.hs.RedirectHandle(r.Context(), u, ip)
	if err != nil {
		l.Error("",
			zap.Error(err))
		err = ErrorPage(w, "./static/400.tmpl", roa.url, l)
		if err != nil {
			l.Error("",
				zap.Error(err))
			render.Render(w, r, ErrInvalidRequest(err))
		}
		return
	}
	msg := fmt.Sprintf("Redirect on url %s successfull", red.LongURL)
	l.Debug(msg)
	// Redirect to the found long address
	http.Redirect(w, r, red.LongURL, http.StatusTemporaryRedirect)
}

// PostStat displaying statistics when receiving a request from the admin url
func (roa *RouterOpenAPI) PostStat(w http.ResponseWriter, r *http.Request) {
	l := roa.logger
	l.Debug("Enter in router func PostStat()")
	// Read the admin url from the request
	adminURL := r.FormValue("admin_url")
	// Form a request with JSON
	url := roa.url + "stat"

	str := fmt.Sprintf(`{"admin_url":"%s"}`, adminURL)
	jsonStr := []byte(str)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		l.Error("",
			zap.Error(err))
		pErr := ErrorPage(w, "./static/500.tmpl", roa.url, l)
		if pErr != nil {
			render.Render(w, r, ErrRender(err))
		}
	}

	smurl := Smurl{AdminURL: adminURL}
	// Сall the handler to get statistics
	gd, err := roa.hs.GetStatHandle(r.Context(), handler.Smurl(smurl))
	if err != nil {
		l.Error("",
			zap.Error(err))
		pErr := ErrorPage(w, "./static/400.tmpl", roa.url, l)
		if pErr != nil {
			render.Render(w, r, ErrInvalidRequest(err))
		}

		return
	}
	gd.SmallURL = roa.url + gd.SmallURL
	w.WriteHeader(http.StatusOK)
	// Call the function to display the result
	err = ResultPage(w, "./static/statistics.tmpl", gd, roa.url, l)
	if err != nil {
		l.Error("",
			zap.Error(err))
		err = render.Render(w, r, Smurl(gd))
		if err != nil {
			l.Error("",
				zap.Error(err))
		}
	}
}

// ResultPage display result page
func ResultPage(w http.ResponseWriter, page string, smurl handler.Smurl, url string, l *zap.Logger) error {
	l.Debug("Enter in router func ResultPage()")
	type smurlWithServerUrl struct {
		SmallURL string
		LongURL  string
		AdminURL string
		IPInfo   string
		Count    string
		URL      string
	}
	swsu := smurlWithServerUrl{
		SmallURL: smurl.SmallURL,
		LongURL:  smurl.LongURL,
		AdminURL: smurl.AdminURL,
		IPInfo:   smurl.IPInfo,
		Count:    smurl.Count,
		URL:      url,
	}
	l.Debug(fmt.Sprintf("smurlWithServerUrl: %v \n", swsu))

	ts, err := template.ParseFiles(page)
	if err != nil {
		return err
	}
	l.Debug("Parse template success")
	err = ts.Execute(w, swsu)
	if err != nil {
		return err
	}
	l.Debug("Execute template success")
	return nil
}

// ErrorPage display the error page
func ErrorPage(w http.ResponseWriter, page string, url string, l *zap.Logger) error {
	l.Debug("Enter in router func ErrorPage")
	ts, err := template.ParseFiles(page)
	if err != nil {
		return err
	}
	err = ts.Execute(w, url)
	if err != nil {
		return err
	}
	l.Debug("ErrorPage template execute success")
	return nil
}
