package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	commandsFilename = "commands.json"
	authFilename     = "auth.json"
)

type Page struct {
	Active string
	Output string
}

func main() {
	r := render.New(render.Options{
		IndentJSON: true,
		Layout:     "layout",
		Extensions: []string{".html"},
	})

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, http.StatusOK, "It's a ducking ninja!")
	})

	router.HandleFunc("/web/", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "index", &Page{Active: "Home"})
	})
	router.HandleFunc("/web/logs", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "logs", &Page{Active: "Logs"})
	})
	router.HandleFunc("/web/commands", func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "commands", &Page{Active: "Commands"})
	})

	router.HandleFunc("/api/logs", ShowLogs(r)).Methods("GET")

	loggedRouter := mux.NewRouter()
	commands := readCommands(commandsFilename)
	loggedRouter.HandleFunc("/api/commands", ListCommands(r, commands)).Methods("GET")
	loggedRouter.HandleFunc("/api/do/{command}", ExecuteCommand(r, commands)).Methods("GET")

	// run logging middleware only for loggedRouter
	router.PathPrefix("/api").Handler(negroni.New(
		logging(),
		negroni.Wrap(loggedRouter),
	))

	//n := negroni.Classic()
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

func ShowLogs(r *render.Render) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, http.StatusOK, getLogs())
	}
}

func ListCommands(r *render.Render, commands Commands) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, http.StatusOK, commands)
	}
}

type CommandResponse struct {
	Id      string
	Command string
	Output  []string
	Error   string
}

func ExecuteCommand(r *render.Render, commands Commands) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["command"]

		// check if valid command
		if command, found := commands[id]; found {
			cmd := strings.Split(command, " ")
			output, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
			if err != nil {
				r.JSON(w, http.StatusInternalServerError,
					CommandResponse{
						Id:      id,
						Command: command,
						Output:  strings.Split(string(output), "\n"),
						Error:   fmt.Sprint(err),
					})
				return
			}
			r.JSON(w, http.StatusOK,
				CommandResponse{
					Id:      id,
					Command: command,
					Output:  strings.Split(string(output), "\n"),
				})
			return
		}
		r.JSON(w, http.StatusBadRequest,
			CommandResponse{
				Id:    id,
				Error: "Unknown command!",
			})
	}
}
