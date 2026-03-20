package server

import (
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	mux *http.ServeMux
}

func NewApp() *App {
	return &App{
		mux: http.NewServeMux(),
	}
}

func (app *App) Handle(path string, handler http.Handler) {
	app.mux.Handle(path, handler)
}

func (app *App) RunServer(port string) error {

	app.printBanner(port)

	server := &http.Server{
		Addr:    port,
		Handler: LoggingMiddleware(CORSMiddleware(app.mux)),
	}

	return server.ListenAndServe()
}

func (app *App) printBanner(port string) {
	urlBase := fmt.Sprintf("http://localhost%s", port)

	fmt.Println("---------------------------------------------------")
	fmt.Printf("|%s|\n", textCenter("go_Inverntory v1.0.0", 51))
	fmt.Printf("|%s|\n", textCenter("Conexión a PostgreSQL establecida correctamente", 51))
	fmt.Printf("|%s|\n", textCenter(urlBase, 51))
	fmt.Printf("|%s|\n", textCenter(urlBase+"/playground", 51))
	fmt.Printf("|%s|\n", strings.Repeat(" ", 51))
	fmt.Println("---------------------------------------------------")
}

func textCenter(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}

	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}
