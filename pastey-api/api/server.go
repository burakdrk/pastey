package api

import (
	"net/http"

	db "github.com/burakdrk/pastey/pastey-api/db/sqlc"
	"github.com/burakdrk/pastey/pastey-api/token"
	"github.com/burakdrk/pastey/pastey-api/util"
	"github.com/burakdrk/pastey/pastey-api/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Server serves HTTP requests
type Server struct {
	store      *db.Store
	router     *gin.Engine
	upgrader   *websocket.Upgrader
	tokenMaker token.Maker
	config     util.Config
}

const (
	apiBasePath = "/v1"
)

// NewServer creates a new HTTP server and set up routing
func NewServer(store *db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSecret)
	if err != nil {
		return nil, err
	}

	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
		upgrader:   upgrader,
	}

	hub := ws.NewHub()
	go hub.Run()

	server.setupRouter(hub)

	return server, nil
}

func (server *Server) setupRouter(hub *ws.Hub) {
	router := gin.Default()

	router.POST(apiBasePath+"/users", server.createUser)
	router.POST(apiBasePath+"/users/login", server.loginUser)
	router.POST(apiBasePath+"/users/logout", server.logoutUser)
	router.POST(apiBasePath+"/token/refresh", server.renewAccessToken)

	authRouter := router.Group(apiBasePath).Use(authMiddleware(server.tokenMaker))

	authRouter.GET("/users/:id", server.getUser)
	authRouter.POST("/devices", server.createDevice)
	authRouter.GET("/devices", server.listUserDevices)
	authRouter.DELETE("/devices/:id", server.deleteDevice)
	authRouter.GET("/devices/:id/entries", server.listDeviceEntries)
	authRouter.DELETE("/entries/:id", server.deleteEntry)
	authRouter.GET("/ws", func(c *gin.Context) {
		server.serveWs(c, hub)
	})
	authRouter.POST("/copy", func(c *gin.Context) {
		server.copyEntry(c, hub)
	})

	server.router = router
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
