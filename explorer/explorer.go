package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yoonhero/ohpotatocoin/blockchain"
)

// set constant
// port is http port string
// templateDir is directory of templates

const (
	templateDir string = "explorer/templates/"
)

// variable templates type Template struct
var templates *template.Template

// set homeData struct
// pagetitle is string that shows on the html title, h1, ext
// Blocks are slice of blockchain.block
type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

// when get or post from link "/"
func home(rw http.ResponseWriter, r *http.Request) {
	// set data
	// pagetitle is "Home"
	// Blocks is blockchain's allblocks
	var blocks []*blockchain.Block
	for _, v := range blockchain.Blocks(blockchain.Blockchain()) {
		h := fmt.Sprintf("%s", v.Hash[0:7]) + "..."
		v.Hash = h
		if len(v.PrevHash) > 7 {
			ph := fmt.Sprintf("%s", v.PrevHash[0:7]) + "..."
			v.PrevHash = ph
		}

		blocks = append(blocks, v)
	}
	data := homeData{"OhPotato", blocks[0:6]}
	// execute
	// writer is http.ResponseWrite
	// templates is "home.gohtml"
	// data is data
	templates.ExecuteTemplate(rw, "home", data)
}

// when get or post from link "/add"
func add(rw http.ResponseWriter, r *http.Request) {
	// switch r.Method that is "GET" or "POST" or "PUT" ext
	switch r.Method {
	// if r.Method is "GET"
	case "GET":
		// execute
		// writer is http.ResponseWrite
		// templates is "add.gohtml"
		// data is nil
		templates.ExecuteTemplate(rw, "add", nil)

		// if r.Method is "POST"
	case "POST":
		// // parses the raw query from URL
		// // updates r.Form
		// r.ParseForm()

		// // get form input named blockData
		// data := r.Form.Get("blockData")

		// addblock data is data from r.Form.Get("blockchain")
		blockchain.Blockchain().AddBlock("")

		// redirect http
		// writer is http.ResponseWrite
		// url is "/"
		// when redirect
		// it redirects not temporary but perminantely
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(port int) {
	// use NewServeMux() to fix the err
	// which occurs when we try to run various http server
	router := mux.NewRouter()

	// Must is a helper that wraps a call to a function returning (*Template, error)
	// ParseGlob creates a new Template and parses the template definitions from the files identified by the pattern.
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	cssHandler := http.FileServer(http.Dir("./css/"))

	// if url is "/"
	router.HandleFunc("/", home)
	http.Handle("/css/", http.StripPrefix("/css/", cssHandler))
	// if url is "/add"
	// router.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost:%d\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}
