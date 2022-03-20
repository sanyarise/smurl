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

//Метод для отображения стартовой страницы
func (roa *RouterOpenAPI) Get(w http.ResponseWriter, r *http.Request) {
	l := roa.logger
	l.Debug("enter in router get")
	w.WriteHeader(http.StatusOK)

	tmpl, err := template.ParseFiles("./static/home.tmpl")
	if err != nil {
		msg := fmt.Sprintf("error on template parse home page files:%s", err)
		l.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		pErr := ErrorPage(w, "./static/500.tmpl")
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
		pErr := ErrorPage(w, "./static/500.tmpl")
		if pErr != nil {
			msg = fmt.Sprintf("error on execute home page files: %s", err)
			l.Error(msg)
			render.Render(w, r, ErrRender(err))
		}
		return
	}
}

//Метод для создания уменьшенного урл.
func (roa *RouterOpenAPI) PostCreate(w http.ResponseWriter, r *http.Request) {
	l := roa.logger
	l.Debug("router post create")
	//Считываем длинный адрес из тела запроса
	longURL := r.FormValue("long_url")
	msg := fmt.Sprintf("longURL: %s", longURL)
	l.Debug(msg)
	//Проверяем валидность длинного адреса
	ok := helpers.CheckURL(longURL)
	if !ok {
		l.Error("incorrect long url")
		pErr := ErrorPage(w, "./static/400.tmpl")
		if pErr != nil {
			msg := fmt.Sprint(pErr)
			l.Error(msg)
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("incorrect long url")))
		}
		return
	}
	//Создаем реквест с JSON
	url := roa.url + "create"

	str := fmt.Sprintf(`{"long_url":"%s"}`, longURL)
	jsonStr := []byte(str)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		msg := fmt.Sprintf("error on create new request: %s", err)
		l.Error(msg)
		pErr := ErrorPage(w, "./static/500.tmpl")
		if pErr != nil {
			render.Render(w, r, ErrRender(err))
		}

		return
	}

	hsu := Smurl{LongURL: longURL}
	//Вызываем обработчик для создания уменьшенного урл
	hss, err := roa.hs.CreateSmurlHandle(r.Context(), handler.Smurl(hsu))
	if err != nil {
		msg := fmt.Sprintf("create smurl error %s: ", err)
		l.Error(msg)
		pErr := ErrorPage(w, "./static/400.tmpl")
		if pErr != nil {
			msg := fmt.Sprint(pErr)
			l.Error(msg)
			render.Render(w, r, ErrInvalidRequest(err))
		}
		return
	}
	l.Debug("create smurl successfull")
	w.WriteHeader(http.StatusCreated)
	//Записываем полный адрес в полученную структуру
	hss.SmallURL = roa.url + hss.SmallURL
	//Вызываем функцию для рендеринга страницы с результатом
	err = ResultPage(w, "./static/result.tmpl", hss)
	if err != nil {
		l.Error("",
			zap.Error(err))
		render.Render(w, r, Smurl(hss))
	}
}

//Метод для перехода по уменьшенному урл, полученному из строки запроса
func (roa *RouterOpenAPI) GetSmallUrl(w http.ResponseWriter, r *http.Request, u string) {
	l := roa.logger
	l.Debug("router get small url")
	//Отсекаем некорректные адреса
	if len(u) != 8 {
		l.Debug(fmt.Sprintf("incorrect small url %s", u))
		err := ErrorPage(w, "./static/400.tmpl")
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
		}
		return
	}
	//Получаем информацию об IP
	ip := helpers.GetIP(r) + " "
	//Вызываем обработчик для поиска маленького урл,
	//поиска соответствующего длинного урл, обновления
	//статистики
	red, err := roa.hs.RedirectHandle(r.Context(), u, ip)
	if err != nil {
		l.Error("",
			zap.Error(err))
		err = ErrorPage(w, "./static/400.tmpl")
		if err != nil {
			l.Error("",
				zap.Error(err))
			render.Render(w, r, ErrInvalidRequest(err))
		}
		return
	}
	msg := fmt.Sprintf("redirect on url %s successfull", red.LongURL)
	l.Debug(msg)
	//Осуществляем перенаправление по найденному длинному адресу
	http.Redirect(w, r, red.LongURL, http.StatusTemporaryRedirect)
}

//Метод для отображения статистики при получении запроса с админским
//урл
func (roa *RouterOpenAPI) PostStat(w http.ResponseWriter, r *http.Request) {
	l := roa.logger
	l.Debug("router post stat")
	//Считываем админский урл из реквеста
	adminURL := r.FormValue("admin_url")
	//Формируем реквест с JSON
	url := roa.url + "stat"

	str := fmt.Sprintf(`{"admin_url":"%s"}`, adminURL)
	jsonStr := []byte(str)
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	if err != nil {
		l.Error("",
			zap.Error(err))
		pErr := ErrorPage(w, "./static/500.tmpl")
		if pErr != nil {
			render.Render(w, r, ErrRender(err))
		}
	}

	smurl := Smurl{AdminURL: adminURL}
	//Вызываем обработчик для получения статистики
	gd, err := roa.hs.GetStatHandle(r.Context(), handler.Smurl(smurl))
	if err != nil {
		l.Error("",
			zap.Error(err))
		pErr := ErrorPage(w, "./static/400.tmpl")
		if pErr != nil {
			render.Render(w, r, ErrInvalidRequest(err))
		}

		return
	}
	gd.SmallURL = roa.url + gd.SmallURL
	w.WriteHeader(http.StatusOK)
	//Вызываем функцию для отображения результата
	err = ResultPage(w, "./static/statistics.tmpl", gd)
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

//Вспомогательная функция для отображения страницы с ошибкой
func ErrorPage(w http.ResponseWriter, page string) error {
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

//Вспомогательная функция для отображения страницы с результатом
func ResultPage(w http.ResponseWriter, page string, smurl handler.Smurl) error {
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
