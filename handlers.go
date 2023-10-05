package main

import (
	"log"
	"net/http"
	"text/template"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, app.sessionManager.GetString(r.Context(), "inside"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "admin" && password == "password" {
		app.sessionManager.Put(r.Context(), "inside", "yes")
	}

	http.Redirect(w, r, "/special", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Put(r.Context(), "inside", "no")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) special(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/special.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, app.sessionManager.GetString(r.Context(), "inside"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}
