package chocolatemilk

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func DefaultPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Chocolate Milk! Your sweetest Go web framework :3")
}

func (c *ChocolateMilk) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	if c.Debug {
		mux.Use(middleware.Logger)
	}
	mux.Use(middleware.Recoverer)

	mux.Get("/", DefaultPage)

	return mux
}
