package api

import (
	db "github.com/burakdrk/pastey/pastey-api/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Server serves HTTP requests
type Server struct {
	store    *db.Store
	router   *gin.Engine
	upgrader *websocket.Upgrader
}

// NewServer creates a new HTTP server and set up routing
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)

	router.GET("/ws", server.handleWebSocket)

	server.router = router
	server.upgrader = upgrader
	return server
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
