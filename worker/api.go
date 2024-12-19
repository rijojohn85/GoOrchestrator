package worker

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Api struct {
	Address string
	Port    int
	Worker  *Worker
	Router  *chi.Mux
}
type ErrResponse struct {
	HTTPStatusCode int
	Message        string
}

func (a *Api) initRouter() {
	a.Router = chi.NewRouter()

	a.Router.Route("/stats", func(r chi.Router) {
		r.Get("/", a.GetStatsHandler)
	})

	a.Router.Route("/tasks", func(r chi.Router) {
		r.Post("/", a.StartTaskHandler)
		r.Get("/", a.GetTaskHandler)
		r.Route("/{taskID}", func(r chi.Router) {
			r.Delete("/", a.StopTaskHandler)
		})
	})
}

func (a *Api) Start() {
	a.initRouter()
	log.Printf("Listening on %s:%d", a.Address, a.Port)
	log.Println(a.Address)
	log.Println(a.Port)
	http.ListenAndServe(":5555", a.Router)
}
