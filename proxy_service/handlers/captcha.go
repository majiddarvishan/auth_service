package handlers

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

func Newcaptcha(c *gin.Context) {
	captchaId := captcha.NewLen(6)
	c.JSON(http.StatusOK, gin.H{"captchaId": captchaId})
}

func NewcaptchaImage(c *gin.Context) {
	captchaId := c.Param("captchaId")
	c.Header("Content-Type", "image/png")
	if err := captcha.WriteImage(c.Writer, captchaId, 240, 80); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
