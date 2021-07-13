package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/impact-eintr/ecache/errors"
)

func (ch *cacheHandler) StatusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, errors.NewerrMsg(
		errors.CodeSuccess, ch.GetStat()))

}
