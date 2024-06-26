package api

import (
	"errors"
	"net/http"
	"time"

	db "github.com/burakdrk/pastey/pastey-api/db/sqlc"
	"github.com/burakdrk/pastey/pastey-api/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createDeviceRequest struct {
	DeviceName string `json:"device_name" binding:"required,min=1"`
	PublicKey  string `json:"public_key" binding:"required,min=1"`
}

type createDeviceResponse struct {
	db.Device             `json:"device"`
	SessionId             uuid.UUID `json:"session_id"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (server *Server) createDevice(ctx *gin.Context) {
	var req createDeviceRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if !server.canHaveMoreDevices(ctx, authPayload.UserID) {
		return
	}

	arg := db.CreateDeviceParams{
		UserID:     authPayload.UserID,
		DeviceName: req.DeviceName,
		PublicKey:  req.PublicKey,
	}

	device, err := server.store.CreateDevice(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, rpayload, err := server.tokenMaker.CreateToken(
		device.UserID,
		server.config.RefreshTokenDuration,
		true,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           rpayload.ID,
		RefreshToken: refreshToken,
		UserID:       device.UserID,
		UserAgent:    ctx.Request.UserAgent(),
		IpAddress:    ctx.ClientIP(),
		ExpiresAt:    rpayload.ExpiresAt,
		DeviceID:     device.ID,
	},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createDeviceResponse{
		Device:                device,
		SessionId:             session.ID,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: rpayload.ExpiresAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) canHaveMoreDevices(ctx *gin.Context, userID int64) bool {
	devices, err := server.store.ListUserDevices(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	user, err := server.store.GetUserById(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	var maxDevices int

	if user.Ispremium {
		maxDevices = 10
	} else {
		maxDevices = 2
	}

	if len(devices) >= maxDevices {
		err := errors.New("maximum number of devices reached")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return false
	}

	return true
}

func (server *Server) listUserDevices(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	devices, err := server.store.ListUserDevices(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, devices)
}

type deleteDeviceRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteDevice(ctx *gin.Context) {
	var req deleteDeviceRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	device, err := server.store.GetDeviceById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if device.UserID != authPayload.UserID {
		err := errors.New("device doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteDevice(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "device deleted")
}

type listDeviceEntriesRequest struct {
	DeviceID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) listDeviceEntries(ctx *gin.Context) {
	var req listDeviceEntriesRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
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

	entries, err := server.store.GetEntriesForDevice(ctx, req.DeviceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, entries)
}
