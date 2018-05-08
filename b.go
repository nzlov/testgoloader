package main

import (
	"errors"
	"os"
	"unsafe"

	"github.com/dearplain/goloader"
	"github.com/gin-gonic/gin"
)

const (
	NAME    = "a"
	VERSION = "0.0.1"
)

var symPtr = map[string]uintptr{}
var codeModule *goloader.CodeModule

func PluginLoad(objs ...interface{}) (string, string, error) {
	obj, ok := objs[0].(*gin.Engine)
	if !ok {
		return NAME, VERSION, errors.New("params 1 need *gin.Engine")
	}

	goloader.RegSymbol(symPtr)
	f, err := os.Open("b1.o")
	if err != nil {
		return NAME, VERSION, err
	}
	reloc, _ := goloader.ReadObj(f)
	f.Close()
	codeModule, err := goloader.Load(reloc, symPtr)
	if err != nil {
		return NAME, VERSION, err
	}
	addFuncPtr := codeModule.Syms["main.Add"]
	funcPtrContainer := (uintptr)(unsafe.Pointer(&addFuncPtr))
	runFunc := *(*func(*gin.Context))(unsafe.Pointer(&funcPtrContainer))

	obj.GET("/add/:a/:b", runFunc)
	return NAME, VERSION, nil
}

func PluginUnload(objs ...interface{}) error {
	codeModule.Unload()
	return nil
}
