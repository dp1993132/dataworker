package httpmod

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

var exports = map[string]lua.LGFunction{
	"get": func(l *lua.LState) int {
		uri := l.Get(1)
		res, err := http.Get(uri.String())
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LString(fmt.Sprintf("%s", buf)))
		defer res.Body.Close()
		return 1
	},
	"post": func(l *lua.LState) int {
		uri := l.Get(1)
		params := l.Get(2).String()
		res, err := http.Post(uri.String(), "application/json", bytes.NewBuffer([]byte(params)))
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
		}
		l.Push(lua.LString(fmt.Sprintf("%s", buf)))
		defer res.Body.Close()
		return 1
	},
}

func Load(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), exports)
	l.Push(mod)
	return 1
}
