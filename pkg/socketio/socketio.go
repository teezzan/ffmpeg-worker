package socketio

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

var Server *socketio.Server

func init() {
	Server = socketio.NewServer(nil)
	Server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		id := s.ID()
		log.Println("connected:", id)
		return nil
	})

	Server.OnEvent("/", "status", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
	})

	Server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	Server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := Server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()

}

func EmitCompleted(uuid string, UserID string) {
	Server.ForEach("/", UserID, func(conn socketio.Conn) {
		log.Println("id ", UserID)
		conn.Emit("completed", uuid)
	})
}

// func IsOnline(UserID string) bool{

// 	Server.ForEach("/", UserID, func(conn socketio.Conn) {
// 		return true
// 	})
// }
