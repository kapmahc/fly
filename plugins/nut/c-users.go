package nut

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUsersSignIn(c *gin.Context) {
	c.HTML(http.StatusOK, "nut.users.sign-in", c.Keys)
}
