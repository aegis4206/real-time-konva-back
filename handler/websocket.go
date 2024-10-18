package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	Conn  *websocket.Conn
	Uuid  string
	Mutex sync.Mutex
}

type Message struct {
	ID          string    `json:"id"`
	Shape       string    `json:"shape"`
	FillColor   string    `json:"fillColor"`
	StrokeColor string    `json:"strokeColor"`
	X           float32   `json:"x,omitempty"`
	Y           float32   `json:"y,omitempty"`
	Height      float32   `json:"height,omitempty"`
	Width       float32   `json:"width,omitempty"`
	Radius      float32   `json:"radius,omitempty"`
	Points      []float32 `json:"points,omitempty"`
	ScaleX      float32   `json:"scaleX,omitempty"`
	ScaleY      float32   `json:"scaleY,omitempty"`
	Text        string    `json:"text,omitempty"`
	Image       string    `json:"image,omitempty"`
}

func (ws *WebSocketConnection) WriteJSON(content interface{}) error {
	ws.Mutex.Lock()
	defer ws.Mutex.Unlock()
	return ws.Conn.WriteJSON(content)
}

var connections = make([]*WebSocketConnection, 0)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://192.168.6.87:5173"
	},
}

func InitWs(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header.Get("Origin"))
	conn, _ := upgrader.Upgrade(w, r, nil)
	uuid, _ := uuid.NewV4()
	id := uuid.String()
	currentConn := &WebSocketConnection{Conn: conn, Uuid: id}
	connections = append(connections, currentConn)
	handleMessage(currentConn)
}

func handleMessage(currentConn *WebSocketConnection) {
	for {
		var message []*Message
		err := currentConn.Conn.ReadJSON(&message)
		log.Println(message)
		if err != nil {
			var tempConnections = make([]*WebSocketConnection, 0)

			for _, connection := range connections {
				if connection.Uuid != currentConn.Uuid {
					tempConnections = append(tempConnections, connection)
				}
			}
			connections = tempConnections
			log.Println("Current Connection: ", len(connections))
			log.Println("temp Connection: ", len(tempConnections))

			currentConn.Conn.Close()
			return
		}

		for _, connection := range connections {
			if connection.Conn != currentConn.Conn {
				connection.WriteJSON(message)
			}
		}
	}
}
