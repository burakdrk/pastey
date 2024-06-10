package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/burakdrk/pastey/pastey-api/db/sqlc"
	"github.com/burakdrk/pastey/pastey-api/token"
	"github.com/burakdrk/pastey/pastey-api/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userResponse struct {
	ID              int64     `json:"id"`
	Email           string    `json:"email"`
	Ispremium       bool      `json:"ispremium"`
	Isemailverified bool      `json:"isemailverified"`
	CreatedAt       time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		ID:              user.ID,
		Email:           user.Email,
		Ispremium:       user.Ispremium,
		Isemailverified: user.Isemailverified,
		CreatedAt:       user.CreatedAt,
	}
}

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if req.ID != authPayload.UserID {
		err := errors.New("user ID doesn't match with the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	user, err := server.store.GetUserById(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

func (server *Server) getUserMe(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	user, err := server.store.GetUserById(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	DeviceId int64  `json:"device_id" binding:"-"`
}

type loginUserResponseWithoutDevice struct {
	AccessToken          string       `json:"access_token"`
	AccessTokenExpiresAt time.Time    `json:"access_token_expires_at"`
	User                 userResponse `json:"user"`
}

type loginUserResponseWithDevice struct {
	SessionId             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.ComparePassword(req.Password, user.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, apayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
		false,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponseWithoutDevice{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: apayload.ExpiresAt,
		User:                 newUserResponse(user),
	}

	if req.DeviceId > 0 {
		device, err := server.store.GetDeviceById(ctx, req.DeviceId)
		if err != nil {
			if err == sql.ErrNoRows {
				err := errors.New("device not found")
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}

			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if device.UserID != user.ID {
			err := errors.New("device ID doesn't match with the authenticated user")
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		err = server.store.DeleteSessionsByDevice(ctx, req.DeviceId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		refreshToken, rpayload, err := server.tokenMaker.CreateToken(
			user.ID,
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
			UserID:       user.ID,
			UserAgent:    ctx.Request.UserAgent(),
			IpAddress:    ctx.ClientIP(),
			ExpiresAt:    rpayload.ExpiresAt,
			DeviceID:     req.DeviceId,
		},
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		rsp := loginUserResponseWithDevice{
			AccessToken:           accessToken,
			AccessTokenExpiresAt:  apayload.ExpiresAt,
			User:                  newUserResponse(user),
			SessionId:             session.ID,
			RefreshToken:          refreshToken,
			RefreshTokenExpiresAt: rpayload.ExpiresAt,
		}

		ctx.JSON(http.StatusOK, rsp)
	} else {
		ctx.JSON(http.StatusOK, rsp)
	}
}

type logoutUserRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) logoutUser(ctx *gin.Context) {
	var req logoutUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if !payload.IsRefresh {
		err := errors.New("token is not a refresh token")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteSession(ctx, payload.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "logged out"})
}
