package p2p

import (
	"fmt"
	"net/http"
	"time"

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
		_, p, err := conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s\n\n", p)
		time.Sleep(1 * time.Second)
		message := fmt.Sprintf("We also think that: %s", p)
		utils.HandleErr(conn.WriteMessage(websocket.TextMessage, []byte(message)))
	}
}
