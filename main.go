package main

import (	
	"log"
	"net/http"

	"vintz.fr/nonogram/server"
)

func main() {
	real_main()
}

func real_main() {

//	imagesHandler := http.StripPrefix("/images/",http.FileServer(http.Dir("images/")))
//	http.HandleFunc("/hello", helloHandler)

	http.Handle ("/images/", server.StaticFileHandler( "/images/", "images/") )
	http.Handle ("/js/", server.StaticFileHandler( "/js/", "js/") )
	http.HandleFunc("/partie", server.MakeNonoHandler())
	

	addr := ":8080"
	log.Printf( "Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
