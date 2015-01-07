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

func main() {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, http.StatusOK, "It's a ducking ninja!")
	})
	router.HandleFunc("/logs", ShowLogs(r)).Methods("GET")

	loggedRouter := mux.NewRouter()
	commands := readCommands(commandsFilename)
	loggedRouter.HandleFunc("/commands", ListCommands(r, commands)).Methods("GET")
	loggedRouter.HandleFunc("/do/{command}", ExecuteCommand(r, commands)).Methods("GET")

	// run logging middleware only for loggedRouter
	router.PathPrefix("/").Handler(negroni.New(
		logging(),
		negroni.Wrap(loggedRouter),
	))

	n := negroni.Classic()
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
