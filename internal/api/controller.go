package api

import "github.com/go-chi/chi/v5"

type Controller interface {
	RegisterRoutes(router chi.Router)
}

func RegisterRoutes(r *chi.Mux, controllers ...Controller) *chi.Mux {
	for _, ctrl := range controllers {
		ctrl.RegisterRoutes(r)
	}

	return r
}
