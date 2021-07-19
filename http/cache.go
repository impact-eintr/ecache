package http

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/impact-eintr/ecache/cache"
	"github.com/impact-eintr/ecache/errors"
)

type cacheHandler struct {
	cache.Cache
}

func (ch *cacheHandler) GetHandler(c *gin.Context) {
	key := c.Param("key")
	if len(key) == 0 {
		c.JSON(http.StatusBadRequest, errors.NewerrMsg(
			errors.CodeInvalidPath, nil))
		return
	}

	b, err := ch.Get(key)
	if len(b) == 0 || err != nil {
		if err != nil {
			log.Println(err)
		}
		c.JSON(http.StatusNotFound, errors.NewerrMsg(
			errors.CodeKeyNotFound, err))
		return
	}

	c.JSON(http.StatusOK, errors.NewerrMsg(
		errors.CodeSuccess, string(b)))

}

func (ch *cacheHandler) PutHandler(c *gin.Context) {
	key := c.Param("key")
	if len(key) == 0 {
		c.JSON(http.StatusBadRequest, errors.NewerrMsg(
			errors.CodeInvalidPath, nil))
		return
	}

	b, _ := ioutil.ReadAll(c.Request.Body)
	if len(b) != 0 {
		err := ch.Set(key, b)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, errors.NewerrMsg(
				errors.CodeKeySetFaild, err))
		}
	}

}

func (ch *cacheHandler) DelHandler(c *gin.Context) {
	key := c.Param("key")
	if len(key) == 0 {
		c.JSON(http.StatusBadRequest, errors.NewerrMsg(
			errors.CodeInvalidPath, nil))
		return
	}

	err := ch.Del(key)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, errors.NewerrMsg(
			errors.CodeKeyDelFaild, nil))
	}

}

func NewCacheHandler(cache cache.Cache) *cacheHandler {
	return &cacheHandler{cache}
}
