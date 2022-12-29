package delivery

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sanyarise/smurl/internal/models"
	"go.uber.org/zap"
)

const (
	status200 = http.StatusOK
	status201 = http.StatusCreated
	status400 = http.StatusBadRequest
	status500 = http.StatusInternalServerError
	page200   = "./static/result.tmpl"
	pageStat  = "./static/statistics.tmpl"
	page400   = "./static/400.tmpl"
	page500   = "./static/500.tmpl"
)

type Smurl struct {
	CreatedAt  string
	ModifiedAt string
	SmallURL   string
	LongURL    string
	AdminURL   string
	IPInfo     []string
	Count      string
	URL        string
}

// Get method displaying the start page
func (router *Router) HomePage(w http.ResponseWriter, r *http.Request) {
	router.logger.Debug("Enter in delivery HomePage()")
	w.WriteHeader(http.StatusOK)

	tmpl, err := template.ParseFiles("./static/home.tmpl")
	if err != nil {
		router.logger.Error(fmt.Sprintf("error on template parse home page files:%s", err))
		w.WriteHeader(http.StatusInternalServerError)
		err = router.ErrorPage(w, page500, status500)
		if err != nil {
			router.logger.Error(fmt.Sprintf("error on error page func: %s", err))

			render.Render(w, r, ErrRender(err))
		}
		return
	}

	err = tmpl.Execute(w, nil)

	if err != nil {
		router.logger.Error(fmt.Sprintf("error on execute home page files: %s", err))
		err = router.ErrorPage(w, page500, status500)
		if err != nil {
			router.logger.Error(fmt.Sprintf("error on execute home page files: %s", err))
			render.Render(w, r, ErrRender(err))
		}
		return
	}
}

// PostCreate creating a minified url
func (router *Router) Create(w http.ResponseWriter, r *http.Request) {
	router.logger.Debug("Enter in delivery Create()")
	// Reading the long address from the request body
	longURL := r.FormValue("long_url")
	router.logger.Debug(fmt.Sprintf("LongURL: %s", longURL))
	// Checking the validity of a long address
	ok := router.helpers.CheckURL(longURL)
	if !ok {
		router.logger.Error("incorrect long url")
		err := router.ErrorPage(w, page400, status400)
		if err != nil {
			router.logger.Error(err.Error())
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("incorrect long url")))
		}
		return
	}
	router.logger.Debug("URL check sucess")

	// Calling usecase method to create a reduced url
	newSmurl, err := router.usecase.Create(context.Background(), longURL)
	if err != nil {
		router.logger.Error(fmt.Sprintf("create smurl error %s: ", err))
		err := router.ErrorPage(w, page500, status500)
		if err != nil {
			router.logger.Error(err.Error())
			render.Render(w, r, ErrRender(err))
		}
		return
	}
	router.logger.Debug("Create smurl success")

	// Call the function to render the page with the result
	err = router.ResultPage(w, page200, newSmurl, status201)
	if err != nil {
		router.logger.Error("",
			zap.Error(err))
	}
}

// GetSmallUrl following the reduced url received from the query string
func (router *Router) Redirect(w http.ResponseWriter, r *http.Request) {
	router.logger.Debug("Enter in delivery Redirect()")
	smallUrl := chi.URLParam(r, "smallUrl")
	ctx := context.Background()

	// Search for a small url in the database
	smurl, err := router.usecase.FindURL(ctx, smallUrl)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			router.logger.Debug(fmt.Sprintf("smallUrl %s is not exist", smallUrl))
			err := router.ErrorPage(w, page400, status400)
			if err != nil {
				router.logger.Error(err.Error())
				render.Render(w, r, ErrInvalidRequest(fmt.Errorf("incorrect small url")))
				return
			}
		} else {
			err := router.ErrorPage(w, page500, status500)
			if err != nil {
				router.logger.Error(err.Error())
				render.Render(w, r, ErrRender(err))
				return
			}
		}
		return
	}
	// Getting information about IP
	ip := router.helpers.GetIP(r)
	// Call the handler to search for a small url,
	// search for the corresponding long url, update
	// statistics
	smurl.IPInfo = append(smurl.IPInfo, ip)
	err = router.usecase.UpdateStat(ctx, *smurl)
	if err != nil {
		router.logger.Error(err.Error())
		err = router.ErrorPage(w, page500, status500)
		if err != nil {
			router.logger.Error(err.Error())
			render.Render(w, r, ErrRender(err))
			return
		}
	}
	// Redirect to the found long address
	http.Redirect(w, r, smurl.LongURL, http.StatusTemporaryRedirect)
	router.logger.Info("Redirect on long url success")
}

// PostStat displaying statistics when receiving a request from the admin url
func (router *Router) GetStat(w http.ResponseWriter, r *http.Request) {
	router.logger.Debug("Enter in delivery GetStat()")

	// Read the admin url from the request
	adminURL := chi.URLParam(r, "adminUrl")

	// Сall the handler to get statistics
	smurl, err := router.usecase.ReadStat(context.Background(), adminURL)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			router.logger.Debug(fmt.Sprintf("adminUrl %s is not exist", adminURL))
			err := router.ErrorPage(w, page400, status400)
			if err != nil {
				router.logger.Error(err.Error())
				render.Render(w, r, ErrInvalidRequest(fmt.Errorf("incorrect small url")))
				return
			}
		} else {
			err := router.ErrorPage(w, page500, status500)
			if err != nil {
				router.logger.Error(err.Error())
				render.Render(w, r, ErrRender(err))
				return
			}
		}
		return
	}
	// Call the function to display the result
	err = router.ResultPage(w, pageStat, smurl, status200)
	if err != nil {
		router.logger.Error(fmt.Sprintf("create smurl error %s: ", err))
		err := router.ErrorPage(w, page500, status500)
		if err != nil {
			router.logger.Error(err.Error())
			render.Render(w, r, ErrRender(err))
		}
		return
	}
}

// ResultPage display result page
func (router *Router) ResultPage(w http.ResponseWriter, page string, smurl *models.Smurl, status int) error {
	router.logger.Debug("Enter in delivery ResultPage()")
	// Write the full address to the resulting structure
	smurl.SmallURL = router.url + "r/" + smurl.SmallURL
	smurl.AdminURL = router.url + "s/" + smurl.AdminURL
	var outSmurl Smurl
	if status == status201 {
		outSmurl.AdminURL = smurl.AdminURL
		outSmurl.SmallURL = smurl.SmallURL
		outSmurl.URL = router.url
	} else if status == status200 {
		outSmurl.AdminURL = smurl.AdminURL
		outSmurl.SmallURL = smurl.SmallURL
		outSmurl.CreatedAt = smurl.CreatedAt.String()
		outSmurl.ModifiedAt = smurl.ModifiedAt.String()
		outSmurl.LongURL = smurl.LongURL
		outSmurl.Count = fmt.Sprint(smurl.Count)
		outSmurl.IPInfo = smurl.IPInfo
		outSmurl.URL = router.url
	}
	router.logger.Debug(fmt.Sprintf("smurlWithServerUrl: %v \n", outSmurl))
	w.WriteHeader(status)
	ts, err := template.ParseFiles(page)
	if err != nil {
		router.logger.Error(err.Error())
		err = router.ErrorPage(w, page500, status500)
		if err != nil {
			router.logger.Error(err.Error())
			return err
		}
		return nil
	}
	router.logger.Debug("Parse template success")
	err = ts.Execute(w, outSmurl)
	if err != nil {
		router.logger.Error(err.Error())
		err = router.ErrorPage(w, page500, status500)
		if err != nil {
			router.logger.Error(err.Error())
			return err
		}
		return nil
	}
	router.logger.Debug("Execute template success")
	return nil
}

// ErrorPage display the error page
func (router *Router) ErrorPage(w http.ResponseWriter, page string, status int) error {
	router.logger.Debug("Enter in delivery ErrorPage()")
	w.WriteHeader(status)
	ts, err := template.ParseFiles(page)
	if err != nil {
		router.logger.Error(fmt.Sprintf("error on parse template file: %v", err))
		return err
	}
	err = ts.Execute(w, router.url)
	if err != nil {
		router.logger.Error(fmt.Sprintf("error on execute template file: %v", err))
		return err
	}
	router.logger.Debug("ErrorPage template execute success")
	return nil
}
