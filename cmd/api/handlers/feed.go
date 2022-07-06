package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Feed(c *gin.Context) {

	c.JSON(http.StatusOK, feedResponse)
}
