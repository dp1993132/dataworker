package core

import (
	"dataworker/worker"
	"io/ioutil"
	"os"
	"strings"
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
					filename := path + v.Name()
					spt := worker.LoadLua(filename)
					spt.Stdout(os.Stdout)
					spt.Stderr(os.Stderr)
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
