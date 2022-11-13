package server

import (
	"github.com/gin-gonic/gin"
	"testBaioCadProject/pkg/database"
)

type Handler struct {
	Db *database.Database
	*gin.Engine
}

// NewHandler create new handler
func NewHandler(db *database.Database) *Handler {
	h := new(Handler)
	router := gin.New()

	get := router.Group("get")
	{
		get.GET("models", h.getModels)
	}

	h.Db = db
	h.Engine = router
	return h
}
