package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	NAME    = "a"
	VERSION = "0.0.1"
)

func PluginInit(objs ...interface{}) (string, string, error) {
	obj := objs[0].(*gin.Engine)
	obj.GET("/add/:a/:b", Add)
	return NAME, VERSION, nil
}

func Add(c *gin.Context) {
	as := c.Param("a")
	a, err := strconv.ParseInt(as, 10, 64)
	if err != nil {
		c.JSON(401, map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
		return
	}
	bs := c.Param("b")
	b, err := strconv.ParseInt(bs, 10, 64)
	if err != nil {
		c.JSON(401, map[string]interface{}{
			"success": false,
			"data":    err.Error(),
		})
		return
	}
	c.JSON(200, map[string]interface{}{
		"success": false,
		"data":    (a + b) * 1,
	})
}
