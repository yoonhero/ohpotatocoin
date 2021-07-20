package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yoonhero/ohpotatocoin/blockchain"
	"github.com/yoonhero/ohpotatocoin/p2p"
	"github.com/yoonhero/ohpotatocoin/utils"
	"github.com/yoonhero/ohpotatocoin/wallet"
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

// // Addblockbody struct
// // which used when post a data
// // data looks like
// // {"message": "data"}
// type addBlockBody struct {
// 	Message string `json:"message"`
// }

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To     string
	Amount int
}

type addPeerPayload struct {
	Address, Port string
}

type loadWalletPayload struct {
	Address string
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
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the Status of the blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
		},
		{
			URL:         url("/latestblocks"),
			Method:      "GET",
			Description: "See A Latest Blocks",
		},
		{
			URL:         url("/latesttransactions"),
			Method:      "GET",
			Description: "See A Latest Transactions",
		},
		{
			URL:         url("/blocks/{hash]"),
			Method:      "Get",
			Description: "See A Block",
			Payload:     "data:hash",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
		{
			URL:         url("/transaction"),
			Method:      "POST",
			Description: "Make a T ransaction",
		},
		{
			URL:         url("/mempool"),
			Method:      "GET",
			Description: "See A Unconfirmed Transactions",
		},
		{
			URL:         url("/ws"),
			Method:      "GET",
			Description: "Upgrade to WebSockets",
		},
	}
	// add content json type
	// rw.Header().Add("Content-Type", "application/json")

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
		// rw.Header().Add("Content-Type", "application/json")

		// send all blocks
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blockchain())))

		// when POST
	case "POST":
		// {"message":"myblockdata"}

		// // new variable struct AddBlockBody
		// var addBlockBody addBlockBody

		// // send pointers and set variable a posted data
		// utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))

		// add block whose data is addBlockBody.Message
		newBlock := blockchain.Blockchain().AddBlock()

		p2p.BroadcastNewBlock(newBlock)

		// send a 201 sign
		rw.WriteHeader(http.StatusCreated)
	}

}

func block(rw http.ResponseWriter, r *http.Request) {
	// get mux var from http.Request
	// shape looks like
	// map[id:1]
	vars := mux.Vars(r)

	// get only id
	// id := vars["height"]

	// strconv.Atoi convert string to int
	hash := vars["hash"]

	// handle err
	// utils.HandleErr(err)

	// FindBlock by id
	block, err := blockchain.FindBlock(hash)

	encoder := json.NewEncoder(rw)

	// if err founded
	if err == blockchain.ErrNotFound {
		// format err to string
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		// send the block
		encoder.Encode(block)
	}

}

// func add json content type
func jsonContentTypeMiddleWare(next http.Handler) http.Handler {
	// make a type of http.Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// add content json type
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func loggerMiddleWare(next http.Handler) http.Handler {
	// make a type of http.Handler
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// add content json type
		fmt.Println(r.URL)
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	blockchain.Status(blockchain.Blockchain(), rw)
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(address, blockchain.Blockchain())
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(address, blockchain.Blockchain())))
	}
}

func latestblocks(rw http.ResponseWriter, r *http.Request) {
	blockchain.LatestBlock(blockchain.Blockchain(), rw)
	// utils.HandleErr(json.NewEncoder(rw).Encode())
}

func latesttransactions(rw http.ResponseWriter, r *http.Request) {
	blockchain.GetLatestTransactions(blockchain.Blockchain(), rw)
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool().Txs))
}

func transaction(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	tx, err := blockchain.Mempool().AddTx(payload.To, payload.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	p2p.BroadcastNewTx(tx)
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["address"]
	// json.NewEncoder(rw).Encode(myWalletResponse{Address: address})

	utils.HandleErr(json.NewEncoder(rw).Encode(wallet.RestApiWallet(key)))
}

func createKey(rw http.ResponseWriter, r *http.Request) {
	utils.HandleErr(json.NewEncoder(rw).Encode(wallet.CreatePrivKey()))
}

func peers(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var payload addPeerPayload
		json.NewDecoder(r.Body).Decode(&payload)
		p2p.AddPeer(payload.Address, payload.Port, port[1:], true)
		rw.WriteHeader(http.StatusOK)
	case "GET":
		json.NewEncoder(rw).Encode(p2p.AllPeers(&p2p.Peers))
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	// use NewServeMux() to fix the err
	// which occurs when we try to run various http server
	router := mux.NewRouter()

	// add json content type
	router.Use(jsonContentTypeMiddleWare, loggerMiddleWare)

	// when  get or post "/" url
	router.HandleFunc("/", documentation).Methods("GET")

	router.HandleFunc("/status", status).Methods("GET")

	// when get or post "/blocks" url
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")

	// get parameter using mux
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/latestblocks", latestblocks).Methods("GET")
	router.HandleFunc("/latesttransactions", latesttransactions).Methods("GET")

	router.HandleFunc("/balance/{address}", balance).Methods("GET")

	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet/{key}", myWallet).Methods("GET")
	router.HandleFunc("/createkey", createKey).Methods("GET")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")

	router.HandleFunc("/transactions", transaction).Methods("POST")

	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)

	// print if err exist
	log.Fatal(http.ListenAndServe(port, router))
}
