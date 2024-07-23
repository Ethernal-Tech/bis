package app

import (
	"bisgo/app/DB"
	p2pserver "bisgo/app/P2P/server"
	provingserver "bisgo/app/proving/server"
	webserver "bisgo/app/web/server"
	"bisgo/config"
	"bisgo/errlog"
	"fmt"
	"net/http"
	"strings"
)

type app struct {
	p2pServer     *p2pserver.P2PServer
	webServer     *webserver.WebServer
	provingServer *provingserver.ProvingServer
	*http.Server
}

func Run() {
	app := app{
		p2pServer:     p2pserver.GetP2PServer(),
		webServer:     webserver.GetWebServer(),
		provingServer: provingserver.GetProvingServer(),
	}

	app.Server = &http.Server{
		Addr:    config.ResolveServerPort(),
		Handler: app.Mux(),
	}

	defer DB.Close()

	fmt.Println("Starting server at port", strings.Split(app.Addr, ":")[1])

	err := app.ListenAndServe()
	if err != nil {
		errlog.Println(err)
	}
}

func (a *app) Mux() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/p2p"):
			a.p2pServer.Mux().ServeHTTP(w, r)
		case strings.Contains(r.URL.Path, "/proof"):
			a.provingServer.Mux().ServeHTTP(w, r)
		default:
			a.webServer.Routes().ServeHTTP(w, r)
		}
	})
}
