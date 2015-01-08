package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type Page struct {
	Active string
	Data   interface{}
}

func WebIndex(r *render.Render) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "index", &Page{Active: "Home"})
	}
}

func WebShowLogs(r *render.Render) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		logs := getLogs()
		data := make([][]string, len(logs))
		for _, line := range logs {
			data = append(data, strings.Split(line, ";"))
		}
		r.HTML(w, http.StatusOK, "logs", &Page{Active: "Logs", Data: data})
	}
}

func WebListCommands(r *render.Render, commands Commands) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r.HTML(w, http.StatusOK, "commands", &Page{Active: "Commands", Data: commands})
	}
}

func WebExecuteCommand(r *render.Render, commands Commands) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["command"]
		command, output, err := executeCommand(id, commands)
		if err != nil {
			r.HTML(w, http.StatusOK, "execute",
				&Page{
					Active: "Commands",
					Data: &CommandResponse{
						Id:      id,
						Command: command,
						Output:  strings.Split(output, "\n"),
						Error:   fmt.Sprint(err),
					},
				})
			return
		}
		r.HTML(w, http.StatusOK, "execute",
			&Page{
				Active: "Commands",
				Data: &CommandResponse{
					Id:      id,
					Command: command,
					Output:  strings.Split(output, "\n"),
				},
			})
	}
}
