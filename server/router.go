package server

import "net/http"

type HandlerFunc func(c *Context)

// Metodo para todas las rutas que utilicen el metodo GET
func (app *App) Get(path string, handler func(*Context)) {
	app.mux.HandleFunc("GET "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
}

// Metodo para todas las rutas que utilicen el metodo POST
func (app *App) Post(path string, handler func(*Context)) {
	app.mux.HandleFunc("POST "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
}

// Metodo para todas las rutas que utilicen el metodo PUT
func (app *App) Put(path string, handler func(*Context)) {
	app.mux.HandleFunc("PUT "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
}

// Metodo para todas las rutas que utilicen el metodo DELETE
func (app *App) Delete(path string, handler func(*Context)) {
	app.mux.HandleFunc("DELETE "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
}
