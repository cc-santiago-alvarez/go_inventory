package explorer

import (
	"html/template"
	"net/http"
)

var page = template.Must(template.New("graphiql-explorer").Parse(`<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>{{.title}}</title>
    <style>
      body {
        height: 100%;
        margin: 0;
        width: 100%;
        overflow: hidden;
      }
      #graphiql {
        height: 100vh;
      }
    </style>
    <script
      src="https://cdn.jsdelivr.net/npm/react@18.2.0/umd/react.production.min.js"
      integrity="sha256-S0lp+k7zWUMk2ixteM6HZvu8L9Eh//OVrt+ZfbCpmgY="
      crossorigin="anonymous"
    ></script>
    <script
      src="https://cdn.jsdelivr.net/npm/react-dom@18.2.0/umd/react-dom.production.min.js"
      integrity="sha256-IXWO0ITNDjfnNXIu5POVfqlgYoop36bDzhodR6LW5Pc="
      crossorigin="anonymous"
    ></script>
    <link
      rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/graphiql@3.0.6/graphiql.min.css"
      integrity="sha256-wTzfn13a+pLMB5rMeysPPR1hO7x0SwSeQI+cnw7VdbE="
      crossorigin="anonymous"
    />
    <link
      rel="stylesheet"
      href="https://cdn.jsdelivr.net/npm/@graphiql/plugin-explorer@0.3.5/dist/style.css"
      crossorigin="anonymous"
    />
  </head>
  <body>
    <div id="graphiql">Cargando...</div>

    <script
      src="https://cdn.jsdelivr.net/npm/graphiql@3.0.6/graphiql.min.js"
      integrity="sha256-eNxH+Ah7Z9up9aJYTQycgyNuy953zYZwE9Rqf5rH+r4="
      crossorigin="anonymous"
    ></script>
    <script
      src="https://cdn.jsdelivr.net/npm/@graphiql/plugin-explorer@0.3.5/dist/index.umd.js"
      crossorigin="anonymous"
    ></script>

    <script>
      const url = location.protocol + '//' + location.host + {{.endpoint}};
      const fetcher = GraphiQL.createFetcher({ url: url });
      const explorerPlugin = GraphiQLPluginExplorer.explorerPlugin();

      ReactDOM.render(
        React.createElement(GraphiQL, {
          fetcher: fetcher,
          plugins: [explorerPlugin],
          isHeadersEditorEnabled: true,
          shouldPersistHeaders: true,
        }),
        document.getElementById('graphiql'),
      );
    </script>
  </body>
</html>
`))

// Handler returns an http.HandlerFunc that serves GraphiQL with the Explorer plugin.
func Handler(title, endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		err := page.Execute(w, map[string]any{
			"title":    title,
			"endpoint": endpoint,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
