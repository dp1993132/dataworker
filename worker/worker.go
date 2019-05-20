package worker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Request 请求
type Request struct {
	Method string
	URL    string
	Body   string
}

// Worker .
type Worker struct {
	sync.Mutex
	Client         *http.Client
	requestList    []*Request
	responseC      chan *http.Response
	errors         chan error
	datas          []string
	undo           int
	interval       time.Duration
	allDoneHandler func(data string)
}

// NewWorker 创建新的worker
func NewWorker() *Worker {
	wk := new(Worker)
	wk.Client = &http.Client{}
	wk.interval = 10 * time.Second
	wk.Client.Timeout = 10 * time.Second
	wk.responseC = make(chan *http.Response, 10)
	wk.errors = make(chan error, 10)
	return wk
}

// SetInterval 设置时间间隔
func (wk *Worker) SetInterval(t time.Duration) {
	wk.interval = t
}

// AddRequest 添加请求
func (wk *Worker) AddRequest(req *Request) error {
	wk.requestList = append(wk.requestList, req)
	return nil
}

// SetRequestList 添加任务
func (wk *Worker) SetRequestList(requestList ...*Request) {
	wk.requestList = requestList
}

// Do 运行
func (wk *Worker) Do() {
	wk.undo = len(wk.requestList)
	for _, request := range wk.requestList {
		go func(req *Request) {
			httpreq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(req.Body))
			if err != nil {
				wk.errors <- err
			}
			resp, err := wk.Client.Do(httpreq)
			if err != nil {
				wk.errors <- err
			} else {
				wk.responseC <- resp
			}
		}(request)
	}
}

func (wk *Worker) done() {
	wk.Lock()
	defer wk.Unlock()
	wk.undo--
	if wk.undo <= 0 && wk.allDoneHandler != nil {
		wk.allDoneHandler("[" + strings.Join(wk.datas, ",") + "]")
		wk.datas = make([]string, 0)
	}
}

// OnError 错误事件
func (wk *Worker) OnError(handler func(err error)) {
	for {
		err := <-wk.errors
		handler(err)
		wk.done()
	}
}

// OnSuccess 成功事件
func (wk *Worker) OnSuccess(handler func(res string)) {
	for {
		resp := <-wk.responseC
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			wk.errors <- err
		} else if resp.StatusCode == 200 {
			wk.datas = append(wk.datas, string(data))
			handler(string(data))
			wk.done()
		} else {
			wk.errors <- fmt.Errorf("request error: %d", resp.StatusCode)
		}
		resp.Body.Close()
	}
}

// OnAllDone 所有的都执行完成事件
func (wk *Worker) OnAllDone(handler func(datas string)) {
	wk.allDoneHandler = handler
}

// Run 循环运行
func (wk *Worker) Run() {
	for {
		if wk.undo == 0 {
			wk.Do()
		}
		if wk.interval <= 0 {
			break
		}
		<-time.After(wk.interval)
	}
}
