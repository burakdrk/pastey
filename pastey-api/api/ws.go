package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) handleWebSocket(ctx *gin.Context) {
	conn, err := server.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		fmt.Println("Received message:", string(p), "from ", conn.RemoteAddr(), "type: ", messageType)

		err = conn.WriteMessage(messageType, p)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
}
