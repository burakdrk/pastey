package api

import (
	"errors"
	"net/http"
	"slices"

	db "github.com/burakdrk/pastey/pastey-api/db/sqlc"
	"github.com/burakdrk/pastey/pastey-api/token"
	"github.com/burakdrk/pastey/pastey-api/ws"
	"github.com/gin-gonic/gin"
)

type copyEntryRequest struct {
	FromDeviceID int64 `json:"from_device_id" binding:"required,min=1"`
	Copies       []struct {
		ToDeviceID    int64  `json:"to_device_id" binding:"required,min=1"`
		EncryptedData string `json:"encrypted_data" binding:"required"`
	} `json:"copies" binding:"gt=0,dive,required"`
}

func (server *Server) copyEntry(ctx *gin.Context, hub *ws.Hub) {
	var req copyEntryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	devices, err := server.store.ListUserDevices(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if !slices.ContainsFunc(devices, func(d db.Device) bool {
		return d.ID == req.FromDeviceID
	}) {
		err := errors.New("from_device_id doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	isDuplicate := make(map[int64]bool)
	for _, copy := range req.Copies {
		if _, ok := isDuplicate[copy.ToDeviceID]; ok {
			err := errors.New("duplicate to_device_id")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		isDuplicate[copy.ToDeviceID] = true

		if !slices.ContainsFunc(devices, func(d db.Device) bool {
			return d.ID == copy.ToDeviceID
		}) {
			err := errors.New("to_device_id doesn't belong to the user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
	}

	arg := db.SaveCopyParams{
		UserID:       authPayload.UserID,
		FromDeviceID: req.FromDeviceID,
		Copies: []struct {
			ToDeviceID    int64  `json:"to_device_id"`
			EncryptedData string `json:"encrypted_data"`
		}(req.Copies),
	}

	entries, err := server.store.SaveCopy(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for _, entry := range entries {
		hub.Broadcast <- &ws.Message{
			UserID:        entry.UserID,
			FromDeviceID:  entry.FromDeviceID,
			ToDeviceID:    entry.ToDeviceID,
			EncryptedData: entry.EncryptedData,
		}
	}

	ctx.JSON(http.StatusOK, entries)
}
