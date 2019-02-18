package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func checkErr(f func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := f(c)
		if err != nil {
			logrus.Error(err)
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
		}
	}
}
