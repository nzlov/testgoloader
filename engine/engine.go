package engine

import (
	"errors"
	"os"
	"unsafe"

	"github.com/dearplain/goloader"
)

const PluginInit = "main.PluginInit"

// 插件初始化方法 返回 Name,Version,error
type PluginInitFunc = *func(...interface{}) (string, string, error)

var (
	ErrNoPluginInitFunc = errors.New("Plugin file don't defined PluginInit Func.")
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

	initFuncPtr, ok := p.codeModule.Syms[PluginInit]
	if !ok {
		return ErrNoPluginInitFunc
	}
	funcPtrContainer := (uintptr)(unsafe.Pointer(&initFuncPtr))
	initFunc := *(PluginInitFunc)(unsafe.Pointer(&funcPtrContainer))
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

func (p *Plugin) Unload() {
	if p.codeModule != nil {
		p.codeModule.Unload()
		p.codeModule = nil
		p.Name = "Unknown"
		p.Version = "Unknown"
	}
}
