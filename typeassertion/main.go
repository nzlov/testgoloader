package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dearplain/goloader"
	"github.com/nzlov/testgoloader/engine"
)

var symPtr = map[string]uintptr{}
var plugin *engine.Plugin

func main() {

	genSymPtr()

	plugin = engine.NewPlugin(os.Args[1])

	fmt.Println(reload())
}

func genSymPtr() {
	goloader.RegSymbol(symPtr)
	goloader.RegTypes(symPtr, &engine.ErrInitParams{})
	goloader.RegTypes(symPtr, strconv.ParseInt)
	goloader.RegTypes(symPtr, map[string]interface{}{})
}

func reload() error {

	err := plugin.Load(symPtr, map[string]interface{}{
		"Check": 1,
	})
	if err != nil {
		return err
	}

	return nil
}
