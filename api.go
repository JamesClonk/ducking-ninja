package main

import (
	"errors"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func ApiShowLogs(r *render.Render) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		r.JSON(w, http.StatusOK, getLogs())
	}
}

func ApiListCommands(r *render.Render, commands Commands) func(http.ResponseWriter, *http.Request) {
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

func ApiExecuteCommand(r *render.Render, commands Commands) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		id := mux.Vars(req)["command"]
		command, output, err := executeCommand(id, commands)
		if err != nil {
			r.JSON(w, http.StatusOK,
				CommandResponse{
					Id:      id,
					Command: command,
					Output:  strings.Split(output, "\n"),
					Error:   fmt.Sprint(err),
				})
			return
		}
		r.JSON(w, http.StatusOK,
			CommandResponse{
				Id:      id,
				Command: command,
				Output:  strings.Split(output, "\n"),
			})
	}
}

func executeCommand(id string, commands Commands) (command, output string, err error) {
	// check if valid command
	if command, found := commands[id]; found {
		cmd := strings.Split(command, " ")
		output, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		if err != nil {
			return command, string(output), err
		}
		return command, string(output), nil
	}
	return command, output, errors.New("Unknown command!")
}
