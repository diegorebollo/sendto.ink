package main

import (
	"log"
	"net/http"
)

type app struct {
	host string
	port string
}

func (app *app) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Sendto.ink"))
}

func (app *app) magicLink(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("magiclink"))
}

func main() {
	app := &app{
		host: "0.0.0.0",
		port: "3821",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.index)
	mux.HandleFunc("/magiclink", app.magicLink)

	log.Print("server running on " + app.port)
	err := http.ListenAndServe(app.host+":"+app.port, mux)
	log.Fatal(err)
}
