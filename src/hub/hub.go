package hub

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Token              string
	Ip_address         string
	Port               string
	Api_version        string
	Api_base_url       string
	Websocket_base_url string
	Wsapp              Websocketapp // well see how this is used and build
}

type Eventlistener struct {
	On_open         any
	On_message      any
	On_error        any
	On_close        any
	On_ping         any
	On_pong         any
	On_data         any
	On_cont_message any
	Ping_intervall  int `default: 60`
}

type Websocketapp struct {
	Websocket websocket.Conn
}

func New_hub(token string, ip_address string) Hub {
	return Hub{
		Token:              token,
		Ip_address:         ip_address,
		Port:               "8443",
		Api_version:        "v1",
		Api_base_url:       fmt.Sprintf("https://%s:%s/%s", ip_address, "8443", "v1"),
		Websocket_base_url: fmt.Sprintf("wss://%s:%s/%s", ip_address, "8443", "v1"),
	}
}

func Headers(h *Hub) map[string]string {
	header := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", h.Token)}
	return header
}

// the websocketlogic probably needs to be custom written..
func New_websocketapp(h *Hub) Websocketapp {
	url := h.Websocket_base_url
	con, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error occured creating websocket app: %s \n", err)
	}
	defer con.Close()
	return Websocketapp{Websocket: *con}
}
