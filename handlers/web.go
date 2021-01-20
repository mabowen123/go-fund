package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func ReptileTable(c *gin.Context) {
	c.HTML(http.StatusOK, "reptile.html", gin.H{})
}
