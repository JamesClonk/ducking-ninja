package main

import (
	"net/http"
	"strings"

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
		r.HTML(w, http.StatusOK, "execute", &Page{Active: "Commands"})
	}
}
