package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/yoonhero/ohpotatocoin/blockchain"
)

// set constant
// port is http port string
// templateDir is directory of templates

const (
	port        string = ":4000"
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
	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}

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
		// parses the raw query from URL
		// updates r.Form
		r.ParseForm()

		// get form input named blockData
		data := r.Form.Get("blockData")

		// addblock data is data from r.Form.Get("blockchain")
		blockchain.GetBlockchain().AddBlock(data)

		// redirect http
		// writer is http.ResponseWrite
		// url is "/"
		// when redirect
		// it redirects not temporary but perminantely
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start() {
	// Must is a helper that wraps a call to a function returning (*Template, error)
	// ParseGlob creates a new Template and parses the template definitions from the files identified by the pattern.
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	// if url is "/"
	http.HandleFunc("/", home)

	// if url is "/add"
	http.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)

	log.Fatal(http.ListenAndServe(port, nil))
}
