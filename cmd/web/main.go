package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	//primero el nombre del flag, -addr, valod default, string explicando lo que hace
	addr := flag.String("addr", ":4000", "HTTP network address")
	//leemos el valor del flag, si no hacemos esto siempre es defualt
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//inicializamos un application
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}
	
	//el valor que nos regresa flag.String es un pointer entonces le hacemos dereference.
	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
