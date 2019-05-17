package mod

import (
	"github.com/yuin/gopher-lua"
)

// Load 加载模块
func Load(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), exports)
	// l.SetField(mod, )
	l.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"newWorker": func(l *lua.LState) int {
		return 0
	},
}
