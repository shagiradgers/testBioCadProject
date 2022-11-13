package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) getModels(ctx *gin.Context) {
	unitGuid := ctx.Query("unit_guid")
	if unitGuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "err",
			"msg":    "missed unit_guid parameter",
		})
		return
	}

	limit, err := -1, error(nil)
	if limitStr := ctx.Query("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status": "err",
				"msg":    "bad limit parameter",
			})
			return
		}
	}

	models, err := h.Db.GetModels(unitGuid, limit)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "err",
			"msg":    "Server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"res":    models,
	})
}
