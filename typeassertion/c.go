package main

import (
	"fmt"
	"reflect"

	"github.com/nzlov/testgoloader/engine"
)

func PluginLoad(args ...interface{}) (string, string, error) {
	if len(args) != 1 {
		return "", "", engine.NewErrInitParams("agrs need 1 map[string]interface{}")
	}
	switch t := args[0].(type) {
	case map[string]interface{}:
		if _, ok := t["Check"]; !ok {
			return "", "", engine.NewErrInitParams("agrs need 1 map[string]interface{} has key Check")
		}
		fmt.Println(t["Check"])
	default:
		return "", "", engine.NewErrInitParams(fmt.Sprintf("agrs need 1 is map[string]interface{} but it's %+v\n", reflect.TypeOf(t)))
	}
	return "ok", "ok", nil
}
