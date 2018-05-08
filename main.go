package main

import (
	"os"
	"strconv"

	"github.com/nzlov/testgoloader/engine"

	"github.com/dearplain/goloader"
	"github.com/gin-gonic/gin"
)

var app *gin.Engine

var plugin *engine.Plugin

var symPtr = map[string]uintptr{}

func main() {
	app = gin.Default()

	genSymPtr()

	plugin = engine.NewPlugin(os.Args[1])

	err := reload()
	if err != nil {
		panic(err)
	}

	app.GET("/reload", func(c *gin.Context) {

		err := reload()
		if err != nil {
			c.JSON(502, map[string]interface{}{
				"success": false,
				"data":    err.Error(),
			})
			return
		}
		c.JSON(200, map[string]interface{}{
			"success": true,
			"data":    "",
		})
	})

	app.Run(":5555")
}

func genSymPtr() {
	goloader.RegSymbol(symPtr)
	goloader.RegTypes(symPtr, app)
	goloader.RegTypes(symPtr, gin.HandlerFunc(nil))
	goloader.RegTypes(symPtr, strconv.ParseInt)
}

func reload() error {

	err := plugin.Load(symPtr, app)
	if err != nil {
		return err
	}

	return nil
}
