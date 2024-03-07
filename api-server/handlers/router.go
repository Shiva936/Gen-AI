package handlers

import (
	"api-server/dtos"
	"api-server/services/userquery"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		if c.Request.Method != http.MethodPost {
			c.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}
		c.Next()
	})

	router.POST("/ask-anything", GetUserQuery)

	return router
}

func GetUserQuery(c *gin.Context) {
	req := &dtos.Request{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()

	res, err := userquery.NewUserQuery().AskAnything(&ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
