package app

import (
	"bisgo/app/DB"
	p2pserver "bisgo/app/P2P/server"
	webserver "bisgo/app/web/server"
	"bisgo/config"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type app struct {
	p2pServer *p2pserver.P2PServer
	webServer *webserver.WebServer
	*http.Server
}

func Run() {
	app := app{
		p2pServer: p2pserver.GetP2PServer(),
		webServer: webserver.GetWebServer(),
	}

	app.Server = &http.Server{
		Addr:    config.ResolveServerPort(),
		Handler: app.Mux(),
	}

	defer DB.Close()

	fmt.Println("The server starts at", app.Addr)
	err := app.ListenAndServe()

	if err != nil {
		log.Fatalf("\033[31m" + "Failed to start the server!" + "\033[31m")
	}
}

func (a *app) Mux() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/p2p"):
			a.p2pServer.Mux().ServeHTTP(w, r)
		case strings.Contains(r.URL.Path, "/proof"):
			// proof service
		default:
			a.webServer.Routes().ServeHTTP(w, r)
		}
	})
}
