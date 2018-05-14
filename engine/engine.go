package engine

import (
	"errors"
	"os"
	"unsafe"

	"github.com/dearplain/goloader"
)

const (
	PluginLoad   = "main.PluginLoad"
	PluginUnload = "main.PluginUnLoad"
)

// 插件初始化方法 返回 Name,Version,error
type PluginLoadFunc = *func(...interface{}) (string, string, error)
type PluginUnloadFunc = *func(...interface{}) error

var (
	ErrNoPluginInitFunc = errors.New("Plugin file don't defined PluginLoad Func.")
)

var DefaultSymPtr = make(map[string]uintptr)

func init() {
	goloader.RegSymbol(DefaultSymPtr)
}

type Plugin struct {
	Name    string
	Version string

	File       string
	codeModule *goloader.CodeModule
}

func NewPlugin(file string) *Plugin {
	return &Plugin{
		File: file,
	}
}

func (p *Plugin) Load(symPtr map[string]uintptr, objs ...interface{}) error {
	p.Unload()
	f, err := os.Open(p.File)
	if err != nil {
		return err
	}
	defer f.Close()
	reloc, err := goloader.ReadObj(f)
	if err != nil {
		return err
	}
	p.codeModule, err = goloader.Load(reloc, symPtr)
	if err != nil {
		return err
	}

	initFuncPtr, ok := p.codeModule.Syms[PluginLoad]
	if !ok {
		return ErrNoPluginInitFunc
	}
	funcPtrContainer := (uintptr)(unsafe.Pointer(&initFuncPtr))
	initFunc := *(PluginLoadFunc)(unsafe.Pointer(&funcPtrContainer))
	name, version, err := initFunc(objs...)
	if err != nil {
		return err
	}
	p.Name = name
	p.Version = version
	return nil
}

func (p *Plugin) Func(name string) (unsafe.Pointer, bool) {
	if v, ok := p.codeModule.Syms[name]; ok {
		pc := (uintptr)(unsafe.Pointer(&v))
		return unsafe.Pointer(&pc), true
	}
	return unsafe.Pointer(nil), false
}

func (p *Plugin) Unload(args ...interface{}) error {
	if p.codeModule != nil {
		unloadFuncPtr, ok := p.codeModule.Syms[PluginUnload]
		if ok {
			funcPtrContainer := (uintptr)(unsafe.Pointer(&unloadFuncPtr))
			unloadFunc := *(PluginUnloadFunc)(unsafe.Pointer(&funcPtrContainer))
			err := unloadFunc(args...)
			if err != nil {
				return err
			}
		}

		p.codeModule.Unload()
		p.codeModule = nil
		p.Name = "Unknown"
		p.Version = "Unknown"
	}
	return nil
}
