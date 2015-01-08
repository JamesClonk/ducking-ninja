package main

import (
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	commandsFilename = "commands.json"
	authFilename     = "auth.json"
)

func main() {
	commands := readCommands(commandsFilename)

	r := render.New(render.Options{
		IndentJSON: true,
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/web/", http.StatusMovedPermanently)
	})
	router.HandleFunc("/web/", WebIndex(r))
	router.HandleFunc("/web/logs", WebShowLogs(r))
	router.HandleFunc("/api/logs", ApiShowLogs(r)).Methods("GET")

	loggedRouter := mux.NewRouter()
	loggedRouter.HandleFunc("/web/commands", WebListCommands(r, commands))
	loggedRouter.HandleFunc("/web/do/{command}", WebExecuteCommand(r, commands))
	loggedRouter.HandleFunc("/api/commands", ApiListCommands(r, commands)).Methods("GET")
	loggedRouter.HandleFunc("/api/do/{command}", ApiExecuteCommand(r, commands)).Methods("GET")

	// run logging middleware only for loggedRouter
	router.PathPrefix("/").Handler(negroni.New(
		logging(),
		negroni.Wrap(loggedRouter),
	))

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		&negroni.Static{
			Dir:       http.Dir("assets"),
			Prefix:    "/web",
			IndexFile: "index.html",
		},
	)
	n.Use(authenticate(r, readAuth(authFilename)))
	n.UseHandler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}
	n.Run(":" + port)
}
