package api

import (
	"errors"
	"net/http"

	"github.com/burakdrk/pastey/pastey-api/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type deleteEntryRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteEntry(ctx *gin.Context) {
	var req deleteEntryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	entry, err := server.store.GetEntryByEntryId(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if entry[0].UserID != authPayload.UserID {
		err := errors.New("entry doesn't belong to the user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.store.DeleteEntry(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "entry deleted")
}