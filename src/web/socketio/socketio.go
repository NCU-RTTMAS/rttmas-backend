/**
* SocketIO implementation
* It supports function-based message emission
* 2023/10/17
*
* @author Henry C. (gsauce.work@gmail.com)
 */

/*
* HOW TO USE
*
* 1. Add this to the import section: omc_socketio "golang/middleware/socket"
*
* 2. Call this to emit a message: omc_socketio.EmitMessage(ROOM_NAME, EVENT_NAME, PAYLOAD)
*      (all three payloads are of the type string)
 */

package socketio

import (
	"fmt"
	"math/rand"
	"net/http"
	logger "rttmas-backend/utils/logger"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// This is a global SocketIO server instance
var (
	ioServer *SocketIO
)

// A wrapper struct
type SocketIO struct {
	*socketio.Server // This is the real SocketIO server
}

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

// Create a new SocketIO server instance
func CreateSocketIOServer() *SocketIO {
	ioServer := new(SocketIO)
	ioServer.Server = socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{CheckOrigin: allowOriginFunc},
			&websocket.Transport{CheckOrigin: allowOriginFunc},
		},
	})

	ioServer.OnConnect("", func(s socketio.Conn) error {
		// s.SetContext("")
		// s.SetContext("foo")
		// s.Join("user-report")
		s.Join("")
		s.Join("rttmas")
		return nil
	})

	ioServer.OnEvent("/", "echo", func(s socketio.Conn, msg string) {
		logger.Info("Emitting...")

		time := time.Now().Unix()
		alarmType := rand.Intn(7-1) + 1
		alarmLevel := rand.Intn(3-1) + 1
		payload := fmt.Sprintf(`{"id":"12345678","alarm_datetime":"%d","resolved_datetime":0,"level":%d,"maintainer_id":"0","psn":"PSN","status":0,"type":%d,"is_hidden":false,"created_at":0,"last_modified_at":0,"maintainer_name":""}`, time, alarmLevel, alarmType)

		s.Emit("newAlarm", payload)
		logger.Info("Emit success.")
	})

	ioServer.OnEvent("/", "foo", func(s socketio.Conn, msg string) {
		logger.Info("foo trigger has been detected")

	})

	return ioServer
}

// Get a reference to the global server variable
// Used in the main function
func GetServerInstance() *socketio.Server {

	if ioServer == nil {
		ioServer = CreateSocketIOServer()
		go func() {
			err := ioServer.Serve()
			if err != nil {
				logger.Fatal(err)
			}
		}()

		logger.Info("ioserver is served ")
	}
	return ioServer.Server
}

// Call this function to emit a message to a room
func EmitMessage(room string, event string, payload string) {
	// logger.Info("Emitting message...")

	if ioServer != nil {
		ioServer.Server.BroadcastToRoom("", room, event, payload)
	} else {
	}
}
