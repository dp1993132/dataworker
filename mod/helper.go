package mod

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

// CompileString 编译lua字符串
func CompileString(source string) (*lua.FunctionProto, error) {
	reader := strings.NewReader(source)
	chunk, err := parse.Parse(reader, source)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, source)
	if err != nil {
		return nil, err
	}
	return proto, nil
}
