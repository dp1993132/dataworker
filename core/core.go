package core

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/dp1993132/dataworker/worker"
)

// Exec 开始执行
func Exec() {
	path := os.Args[1]
	scripts := make([]worker.Script, 0)
	f, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	if f.IsDir() {
		dir, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		} else {
			for _, v := range dir {
				if strings.HasSuffix(v.Name(), ".lua") {
					if strings.HasSuffix(path, "/") == false {
						path = path + "/"
					}
					filename := path + v.Name()
					spt := worker.LoadLua(filename)
					scripts = append(scripts, spt)
				}
			}
		}
	} else {
		spt := worker.LoadLua(path)
		scripts = append(scripts, spt)
	}

	for _, spt := range scripts {
		go spt.Run()
	}
}
