package delivery

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sanyarise/smurl/internal/usecase"
	"go.uber.org/zap"
)

type Router struct {
	*chi.Mux
	usecase usecase.Usecase
	logger  *zap.Logger
	url     string
}

func NewRouter(usecase usecase.Usecase, logger *zap.Logger, url string) *Router {
	r := chi.NewRouter()

	router := &Router{
		usecase: usecase,
		logger:  logger,
		url:     url,
	}

	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Group(func(r chi.Router) {
		r.Post("/create", router.Create)
		r.Get("/r/{smallUrl}", router.Redirect)
		r.Get("/s/{adminUrl}", router.GetStat)
		r.Get("/", router.HomePage)
	})
	router.Mux = r
	return router
}

func (Smurl) Bind(r *http.Request) error {
	return nil
}

func (Smurl) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
