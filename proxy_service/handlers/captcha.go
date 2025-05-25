package handlers

import (
	"net/http"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

// CaptchaResponse defines the response model for a new CAPTCHA
type CaptchaResponse struct {
	CaptchaID string `json:"captchaId" example:"a1b2c3"`

}

// Newcaptcha godoc
// @Summary      Generate a new CAPTCHA ID
// @Description  Returns a new CAPTCHA ID to be used for verification
// @Tags         captcha
// @Produce      json
// @Success      200  {object}  CaptchaResponse
// @Router       /captcha/new [get]
func Newcaptcha(c *gin.Context) {
	captchaId := captcha.NewLen(6)
	c.JSON(http.StatusOK, gin.H{"captchaId": captchaId})
}

// NewcaptchaImage godoc
// @Summary      Serve CAPTCHA image
// @Description  Returns a CAPTCHA image in PNG format for the given CAPTCHA ID.
// @Description  Example: GET /captcha/image/a1b2c3
// @Tags         captcha
// @Param        captchaId  path      string  true  "CAPTCHA ID"
// @Produce      image/png
// @Success      200  {file}  binary
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /captcha/image/{captchaId} [get]
func NewcaptchaImage(c *gin.Context) {
	captchaId := c.Param("captchaId")
	c.Header("Content-Type", "image/png")
	if err := captcha.WriteImage(c.Writer, captchaId, 240, 80); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
