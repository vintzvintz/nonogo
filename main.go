package main

import (
	"io"
	"log"
	"net/http"

	"vintz.fr/nonogram/server"
)

func main() {
	real_main()
}

func real_main() {


	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}

//	imagesHandler := http.StripPrefix("/images/",http.FileServer(http.Dir("images/")))

	http.HandleFunc("/hello", helloHandler)
	http.Handle ("/images/", server.ImagesHandler( "/images/", "images/") )
	http.HandleFunc("/", server.MakeNonoHandler())
	

	addr := ":8080"
	log.Printf( "Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
