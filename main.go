package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "It's a ducking ninja!")
	})

	n := negroni.Classic()
	n.UseHandler(router)

	n.Run(":3000")
}
