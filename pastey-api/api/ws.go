package api

import (
	"errors"
	"net/http"

	"github.com/burakdrk/pastey/pastey-api/token"
	"github.com/burakdrk/pastey/pastey-api/ws"
	"github.com/gin-gonic/gin"
)

type wsRequest struct {
	DeviceID int64 `form:"device_id" binding:"required,min=1"`
}

func (server *Server) serveWs(ctx *gin.Context, hub *ws.Hub) {
	var req wsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	device, err := server.store.GetDeviceById(ctx, req.DeviceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if device.UserID != authPayload.UserID {
		err := errors.New("device doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if _, ok := hub.Clients[authPayload.UserID][req.DeviceID]; ok {
		err := errors.New("device is already connected")
		ctx.JSON(http.StatusConflict, errorResponse(err))
		return
	}

	conn, err := server.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	client := ws.NewClient(conn, authPayload.UserID, req.DeviceID)

	hub.Register <- client

	go client.Write()
	go client.Read(hub)
}
