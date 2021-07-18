package p2p

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yoonhero/ohpotatocoin/utils"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	for {
		fmt.Println("Waiting for message ...")
		_, p, err := conn.ReadMessage()
		fmt.Println("Message arrived")
		utils.HandleErr(err)
		fmt.Printf("%s\n\n", p)
	}
}
