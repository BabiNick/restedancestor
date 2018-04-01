//Package server is used to separate the handlers from other functions of the API.
package server

import (
	"encoding/json"
	"net/http"
	"github.com/bruno-chavez/restedancestor/database"
	"github.com/bruno-chavez/restedancestor/lib"
)

//RquoteHandler takes care of the rquote route.
func RquoteHandler(w http.ResponseWriter, r *http.Request) {

	slice := database.Parser()
	marshaledData, _ := json.MarshalIndent(lib.Random(slice), "", "")

	w.Header().Set("Content-Type", "application/json")
	w.Write(marshaledData)
}