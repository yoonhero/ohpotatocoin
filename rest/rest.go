package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yoonhero/ohpotatocoin/blockchain"
	"github.com/yoonhero/ohpotatocoin/utils"
)

// variable post string
var port string

// new type URL
type url string

// type URL's interface
func (u url) MarshalText() ([]byte, error) {
	// var url is http://localhost + port + URL
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

//`json:"name"` => return name not Name
//'json:"omitempty"` => don't send if field is empty
// url, method, description, payload in type URLDescription struct
type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

// URLDescription all string to return value
// func (u URLDescription) String() string {
// 	return "Hello I'm the URL Description"
// }

// Addblockbody struct
// which used when post a data
// data looks like
// {"message": "data"}
type addBlockBody struct {
	Message string
}

// when url is "/"
func documentation(rw http.ResponseWriter, r *http.Request) {

	// []URLDescription struct slice
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{id]"),
			Method:      "Get",
			Description: "See A Block",
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

// when get or post url /blocks
func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// when GET
	case "GET":
		// recognize that this content is json
		rw.Header().Add("Content-Type", "application/json")

		// send all blocks
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())

		// when POST
	case "POST":
		// {"message":"myblockdata"}

		// new variable struct AddBlockBody
		var addBlockBody addBlockBody

		// send pointers and set variable a posted data
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))

		// add block whose data is addBlockBody.Message
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)

		// send a 201 sign
		rw.WriteHeader(http.StatusCreated)
	}
}

func Start(aPort int) {
	// use NewServeMux() to fix the err
	// which occurs when we try to run various http server
	handler := http.NewServeMux()

	port = fmt.Sprintf(":%d", aPort)
	// when  get or post "/" url
	handler.HandleFunc("/", documentation)

	// when get or post "/blocks" url
	handler.HandleFunc("/blocks", blocks)

	fmt.Printf("Listening on http://localhost%s\n", port)

	// print if err exist
	log.Fatal(http.ListenAndServe(port, handler))
}
