package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Request 请求配置
type Request struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
	Body   string `yaml:"body"`
}

// CallBack 回调函数
type CallBack struct {
	Shell string `ymal:"shell"`
	URL   string `yaml:"url"`
}

// Worker 工作进程配置
type Worker struct {
	RequestList []*Request `yaml:"request"`
	Success     *CallBack  `yaml:"success"`
	Error       *CallBack  `yaml:"error"`
	AllDone     *CallBack  `yaml:"allDone"`
}

// Config 配置
type Config struct {
	WorkerList []Worker `yaml:"workerList"`
}

// Load 读配置
func Load(filename string) (*Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, err
}
