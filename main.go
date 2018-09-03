// Package main is the start of the API, creates all the routes and sets a handler for each one.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bruno-chavez/restedancestor/handler"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.NewRoute().
		Path("/random").
		HandlerFunc(handler.RandomHandler).
		Methods("GET", "OPTIONS")
	router.NewRoute().
		Path("/all").
		HandlerFunc(handler.AllHandler).
		Methods("GET", "OPTIONS")
	router.NewRoute().
		Path("/search/{word}").
		HandlerFunc(handler.SearchHandler).
		Methods("GET", "OPTIONS")
	router.NewRoute().
		Path("/one/{uuid}").
		HandlerFunc(handler.OneHandler).
		Methods("GET")
	router.NewRoute().
		Path("/top").
		HandlerFunc(handler.TopHandler).
		Methods("GET")
	// Ça n'a rien à faire dans un ordre GET (+1 pour le suivant)
	router.NewRoute().
		Path("/one/{uuid}/dislike").
		HandlerFunc(handler.DislikeHandler).
		Methods("GET")
	router.NewRoute().
		Path("/one/{uuid}/like").
		HandlerFunc(handler.LikeHandler).
		Methods("GET")
	/*router.NewRoute().
	Path("/update").
	HandlerFunc(handler.UpdateHandler).
	Methods("GET", "OPTIONS")
	*/

	fmt.Println("Welcome to restedancestor, the API is running in a maddening fashion!")

	fmt.Println("The Ancestor is waiting and listening on localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
