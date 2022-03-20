package smurlapi

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/sanyarise/smurl/internal/infrastructure/api/handler"
	"github.com/sanyarise/smurl/internal/infrastructure/api/helpers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

//Роутер chi с написанными вручную методами, практически идентичен сгенерированному
type RouterChi struct {
	*chi.Mux
	hs     *handler.Handlers
	logger *zap.Logger
	url    string
}

func NewRouterChi(hs *handler.Handlers, l *zap.Logger) *RouterChi {
	r := chi.NewRouter()
	ret := &RouterChi{
		hs:     hs,
		logger: l,
	}
	r.Get("/", ret.Home)
	r.Post("/create", ret.CreateSmurl)
	r.Get("/{smallURL}", ret.Redirect)
	r.Post("/stat", ret.GetStat)

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

func (rc *RouterChi) Home(w http.ResponseWriter, r *http.Request) {
	l := rc.logger
	l.Debug("enter in router home")

	tmpl, err := template.ParseFiles("./static/home.tmpl")
	if err != nil {
		msg := fmt.Sprintf("error on template parse home page files:%s", err)
		l.Error(msg)

		pErr := ErrorPage(w, r, "./static/500.tmpl")
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
		pErr := ErrorPage(w, r, "./static/500.tmpl")
		if pErr != nil {
			msg = fmt.Sprintf("error on execute home page files: %s", err)
			l.Error(msg)
			render.Render(w, r, ErrRender(err))
		}
		return
	}
}

func (rc *RouterChi) CreateSmurl(w http.ResponseWriter, r *http.Request) {
	l := rc.logger
	l.Debug("router post create")

	longURL := r.FormValue("long_url")
	msg := fmt.Sprintf("longURL: %s", longURL)
	l.Debug(msg)

	url := rc.url + "create"

	str := fmt.Sprintf(`{"long_url":"%s"}`, longURL)
	jsonStr := []byte(str)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		msg := fmt.Sprintf("error on create new request: %s", err)
		l.Error(msg)
		pErr := ErrorPage(w, r, "./static/500.tmpl")
		if pErr != nil {
			render.Render(w, r, ErrRender(err))
		}

		return
	}

	hsu := Smurl{LongURL: longURL}
	if err := render.Bind(r, &hsu); err != nil {
		zap.Error(err)
		err = ErrorPage(w, r, "./static/500.tmpl")
		if err != nil {
			render.Render(w, r, ErrRender(err))
		}
		return
	}
	hss, err := rc.hs.CreateSmurlHandle(r.Context(), handler.Smurl(hsu))
	if err != nil {
		msg := fmt.Sprintf("create smurl error %s: ", err)
		l.Error(msg)
		pErr := ErrorPage(w, r, "./static/400.tmpl")
		if pErr != nil {
			zap.Error(pErr)
			render.Render(w, r, ErrInvalidRequest(err))
		}
		return
	}
	l.Debug("create smurl successfull")
	err = ResultPage(w, r, "./static/result.tmpl", hss)
	if err != nil {
		l.Error("",
			zap.Error(err))
		render.Render(w, r, Smurl(hss))
	}
}

func (rc *RouterChi) Redirect(w http.ResponseWriter, r *http.Request) {
	l := rc.logger
	l.Debug("router redirect")
	smallURL := string(chi.URLParam(r, "smallURL"))
	if smallURL == "" {
		return
	}

	msg := fmt.Sprintf("smallURL:%s", smallURL)
	l.Debug(msg)

	ip := helpers.GetIP(r) + " "

	l.Debug("ip",
		zap.String(":", ip))

	red, err := rc.hs.RedirectHandle(r.Context(), smallURL, ip)
	if err != nil {
		l.Error("",
			zap.Error(err))
		pErr := ErrorPage(w, r, "./static/400.tmpl")
		if pErr != nil {
			l.Error("",
				zap.Error(pErr))
			render.Render(w, r, ErrRender(err))
		}
		return
	}
	l.Debug("redirect successfull")
	http.Redirect(w, r, red.LongURL, http.StatusTemporaryRedirect)
}

func (rc *RouterChi) GetStat(w http.ResponseWriter, r *http.Request) {
	l := rc.logger
	l.Debug("router post stat")

	adminURL := r.FormValue("admin_url")

	url := rc.url + "stat"

	str := fmt.Sprintf(`{"admin_url":"%s"}`, adminURL)
	jsonStr := []byte(str)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		l.Error("",
			zap.Error(err))
		pErr := ErrorPage(w, r, "./static/500.tmpl")
		if pErr != nil {
			render.Render(w, r, ErrRender(err))
		}
	}

	smurl := Smurl{AdminURL: adminURL}
	if err := render.Bind(r, &smurl); err != nil {
		l.Error("",
			zap.Error(err))

		pErr := ErrorPage(w, r, "./static/500.tmpl")
		if pErr != nil {
			render.Render(w, r, ErrRender(err))
		}
		return
	}
	gd, err := rc.hs.GetStatHandle(r.Context(), handler.Smurl(smurl))
	if err != nil {
		l.Error("",
			zap.Error(err))
		pErr := ErrorPage(w, r, "./static/400.tmpl")
		if pErr != nil {
			render.Render(w, r, ErrRender(err))
		}

		return
	}

	err = ResultPage(w, r, "./static/statistics.tmpl", gd)
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

func ErrorPage(w http.ResponseWriter, r *http.Request, page string) error {
	ts, err := template.ParseFiles(page)
	if err != nil {
		return err
	}
	err = ts.Execute(w, nil)
	if err != nil {
		return err
	}
	return nil
}

func ResultPage(w http.ResponseWriter, r *http.Request, page string, smurl handler.Smurl) error {
	ts, err := template.ParseFiles(page)
	if err != nil {
		return err
	}
	err = ts.Execute(w, smurl)
	if err != nil {
		return err
	}
	return nil
}
