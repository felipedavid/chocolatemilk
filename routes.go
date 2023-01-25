package chocolatemilk

import (
	"net/http"
)

// TODO: Write a custom router to dispatch based on request method
func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./ui/static/"))))
	mux.HandleFunc("/", app.WelcomePage)

	return mux
}

func (app *App) WelcomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, "welcomepage.tmpl", nil)
}
