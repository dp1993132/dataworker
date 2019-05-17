package worker

import (
	"bytes"
	"dataworker/config"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Worker .
type Worker struct {
	sync.Mutex
	Client         *http.Client
	requestList    []*http.Request
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

// CreateRequest 创建请求队列
func (wk *Worker) CreateRequest(method string, urls []string, body io.Reader) error {
	reqs := make([]*http.Request, len(urls))
	for k, v := range urls {
		req, err := http.NewRequest(method, v, body)
		if err != nil {
			return err
		}
		reqs[k] = req
	}
	wk.requestList = reqs
	return nil
}

// SetInterval 设置时间间隔
func (wk *Worker) SetInterval(t time.Duration) {
	wk.interval = t
}

// AddRequest 添加请求
func (wk *Worker) AddRequest(req *config.Request) error {
	if wk.requestList == nil {
		wk.requestList = make([]*http.Request, 0)
	}
	r, err := http.NewRequest(req.Method, req.URL, bytes.NewBuffer([]byte(req.Body)))
	if err != nil {
		return err
	} else {
		wk.requestList = append(wk.requestList, r)
	}
	return nil
}

// SetRequestList 添加任务
func (wk *Worker) SetRequestList(requestList ...*http.Request) {
	wk.requestList = requestList
}

// Do 运行
func (wk *Worker) Do() {
	wk.undo = len(wk.requestList)
	wk.datas = nil
	for _, request := range wk.requestList {
		go func(req *http.Request) {
			resp, err := wk.Client.Do(req)
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
	}
}

// OnError 错误事件
func (wk *Worker) OnError(handler func(err error)) {
	for {
		err := <-wk.errors
		wk.done()
		handler(err)
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
		} else {
			wk.errors <- fmt.Errorf("request error: %d", resp.StatusCode)
		}
		handler(string(data))
		wk.done()
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
