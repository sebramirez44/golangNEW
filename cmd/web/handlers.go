package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

//hacemos que todas las funciones sean miembros del struct application
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w) //usamos un helper
		return
	}
	//inicializamos slice con ambos arhcivos
	//base template tiene que ser el primero en la slice.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	//leemos los archivos y guardamos los templates en un template set
	//usamos un variadic parameter
	ts, err := template.ParseFiles(files...)
	if err != nil {
		//escribimos el log message a nuestro logger custom definido en application
		app.serverError(w, err)
		return
	}
	//decimos que queremos enviar la base template, esta a su vez invoca las
	//templates de title y main
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
