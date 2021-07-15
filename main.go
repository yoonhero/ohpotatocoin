package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// constant post string
const port string = ":5000"

//`json:"name"` => return name not Name
//'json:"omitempty"` => don't send if field is empty
// url, method, description, payload in type URLDescription struct
type URLDescription struct {
	URL         string `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

// when link is "/"
func documentation(rw http.ResponseWriter, r *http.Request) {

	// []URLDescription struct slice
	data := []URLDescription{
		{
			URL:         "/",
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         "http://localhost:4000/blocks",
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
	}
	// add content json type

	rw.Header().Add("Content-Type", "application/json")

	// json.NewEncoder(rw).Encode(data)
	// is same
	// b, err := json.Marshal(data)
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data)
}

func main() {
	// when "/" link get or post
	http.HandleFunc("/", documentation)

	fmt.Printf("Listening on http://localhost%s\n", port)

	// print if err exist
	log.Fatal(http.ListenAndServe(port, nil))
}
