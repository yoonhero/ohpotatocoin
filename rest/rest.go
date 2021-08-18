package rest

import (
	"encoding/hex"
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

type addBlockBody struct {
	From string `json:"from"`
}

type addTxPayload struct {
	Privkey string
	To      string
	Amount  int
}

type addPeerPayload struct {
	Address, Port string
}

type walletPayload struct {
	Key string `json:"key"`
}

type createKeyAddressPayload struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

type InfoMining struct {
	Block      *blockchain.Block `json:"block"`
	Hash       string            `json:"hash"`
	Difficulty int               `json:"difficulty"`
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
			Payload:     "{'from':''(miner)}",
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
			URL:         url("/transactions"),
			Method:      "POST",
			Description: "Make a T ransaction",
			Payload:     "{'privkey':''(from), 'to':'', 'amount':''}",
		},

		{
			URL:         url("/mempool"),
			Method:      "GET",
			Description: "See A Unconfirmed Transactions",
		},
		{
			URL:         url("/createkey"),
			Method:      "GET",
			Description: "Make a Random Private Key and Public Key",
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
		var addBlockBody addBlockBody
		json.NewDecoder(r.Body).Decode(&addBlockBody)
		// {"message":"myblockdata"}

		// // new variable struct AddBlockBody
		// var addBlockBody addBlockBody

		// // send pointers and set variable a posted data
		// utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))

		// add block whose data is addBlockBody.Message
		newBlock := blockchain.Blockchain().AddBlock(addBlockBody.From)

		p2p.BroadcastNewBlock(newBlock)

		// send a 201 sign
		rw.WriteHeader(http.StatusCreated)
	}

}

func mining(rw http.ResponseWriter, r *http.Request) {
	var addBlockBody addBlockBody
	json.NewDecoder(r.Body).Decode(&addBlockBody)

	block, hash := blockchain.Blockchain().SendInfoOfMining(addBlockBody.From)
	fmt.Println(utils.Hash("JnsgMDAwMDAwZGQwZTNkOGYxNDM1N2QyZTIyMTQ4ZGJlM2U3ZjAxMDhkMDFkYmVlNTM0ZTgyOGMxNTU2Njg4OTM4ZiAzMiA2IDAgMCBbMHhjMDAwMjhjYTAwXX0"))
	json.NewEncoder(rw).Encode(InfoMining{
		Block:      block,
		Hash:       hash,
		Difficulty: block.Difficulty,
	})

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
		utils.AllowConnection(rw)
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
	switch r.Method {
	case "POST":
		var payload addTxPayload
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))

		tx, err := blockchain.Mempool().AddTx(payload.Privkey, payload.To, payload.Amount)
		if err != nil {
			json.NewEncoder(rw).Encode(errorResponse{err.Error()})
			return
		}
		p2p.BroadcastNewTx(tx)
		rw.WriteHeader(http.StatusCreated)
	case "GET":
		blockchain.Transactions(blockchain.Blockchain(), rw)
	}

}

func findTx(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	blockchain.FindTransactions(blockchain.Blockchain(), rw, hash)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	var payload walletPayload
	json.NewDecoder(r.Body).Decode(&payload)
	// json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
	bytes, err := hex.DecodeString(payload.Key)
	utils.HandleErr(err)
	json.NewEncoder(rw).Encode(wallet.RestApiWallet(bytes))
	rw.WriteHeader(http.StatusOK)
}

func createKey(rw http.ResponseWriter, r *http.Request) {
	address, key := wallet.RestApiCreatePrivKey()
	utils.HandleErr(json.NewEncoder(rw).Encode(createKeyAddressPayload{address, fmt.Sprintf("%x", key)}))
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
	router.HandleFunc("/mining", mining).Methods("POST")
	router.HandleFunc("/latestblocks", latestblocks).Methods("GET")
	router.HandleFunc("/latesttransactions", latesttransactions).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/wallet", myWallet).Methods("POST")
	router.HandleFunc("/createkey", createKey).Methods("GET")
	router.HandleFunc("/ws", p2p.Upgrade).Methods("GET")
	router.HandleFunc("/transactions", transaction).Methods("POST", "GET")
	router.HandleFunc("/transaction/{hash:[a-f0-9]+}", findTx).Methods("GET")
	router.HandleFunc("/peers", peers).Methods("GET", "POST")
	fmt.Printf("Listening on http://localhost%s\n", port)

	// print if err exist
	log.Fatal(http.ListenAndServe(port, router))
}
