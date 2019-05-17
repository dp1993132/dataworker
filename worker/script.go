package worker

import (
	"bytes"
	"dataworker/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
)

// Script 脚本接口
type Script interface {
	Run()
	Stdout(io.Writer)
	Stderr(io.Writer)
}

// script 脚本
type script struct {
	stdout io.Writer
	stderr io.Writer
	*Worker
}

func (spt *script) Stdout(w io.Writer) {
	spt.stdout = w
}
func (spt *script) Stderr(w io.Writer) {
	spt.stderr = w
}

// LoadLua 加载脚本
func LoadLua(filename string) Script {
	var spt script
	spt.Stdout(bytes.NewBuffer([]byte{}))
	spt.Stderr(bytes.NewBuffer([]byte{}))
	wk := NewWorker()
	L := lua.NewState()
	L.SetGlobal("load", L.NewFunction(func(l *lua.LState) int {
		s := LoadLua(l.Get(1).String())
		s.Stdout(os.Stdout)
		s.Stderr(os.Stderr)
		go s.Run()
		return 0
	}))
	L.SetGlobal("print", L.NewFunction(func(l *lua.LState) int {
		spt.stdout.Write([]byte(l.Get(1).String() + "\n"))
		return 0
	}))
	L.SetGlobal("addRequest", L.NewFunction(func(l *lua.LState) int {
		request := new(config.Request)
		v := l.Get(1)
		gluamapper.Map(v.(*lua.LTable), request)
		wk.AddRequest(request)
		return 0
	}))
	L.SetGlobal("setInterval", L.NewFunction(func(l *lua.LState) int {
		v := l.Get(1)
		iv, err := strconv.ParseInt(v.String(), 10, 32)
		if err != nil {
			panic("设置间隔时间失败")
		}
		wk.SetInterval(time.Duration(iv) * time.Second)
		return 0
	}))
	L.SetGlobal("send", L.NewFunction(func(l *lua.LState) int {
		url := l.Get(1).String()
		body := l.Get(2).String()
		resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
		if err != nil {
			L.Push(lua.LString("发送请求失败:" + err.Error()))
			return 1
		} else if resp.StatusCode != 200 {
			L.Push(lua.LString(fmt.Sprintf("服务器返回错误:%d\n", resp.StatusCode)))
			return 1
		}
		return 0
	}))
	err := L.DoFile(filename)
	if err != nil {
		spt.stderr.Write([]byte(err.Error()))
	}

	// 设置成功回调
	onSuccess := L.GetGlobal("onSuccess")
	if lua.LVAsBool(onSuccess) {
		go wk.OnSuccess(func(res string) {
			L.CallByParam(lua.P{
				Fn: onSuccess,
			}, lua.LString(res))
		})
	}
	// 设置错误回调
	onError := L.GetGlobal("onError")
	if lua.LVAsBool(onError) {
		go wk.OnError(func(err error) {
			L.CallByParam(lua.P{
				Fn: onError,
			}, lua.LString(err.Error()))
		})
	}
	// 设置全部完成回调
	onAllDone := L.GetGlobal("onAllDone")
	if lua.LVAsBool(onAllDone) {
		go wk.OnAllDone(func(res string) {
			L.CallByParam(lua.P{
				Fn: onAllDone,
			}, lua.LString(res))
		})
	}
	spt.Worker = wk
	return &spt
}
