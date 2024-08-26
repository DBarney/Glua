package glua

import (
	"fmt"
	"io"
	"path"
	"reflect"

	"github.com/yuin/gopher-lua"
)

type Lua struct {
	unused chan *lua.LState
	reuse  bool
}

func New() *Lua {
	return &Lua{
		reuse:  false,
		unused: make(chan *lua.LState, 4),
	}
}

func newLua() *lua.LState {
	L := lua.NewState()
	L.OpenLibs()

	L.DoString(`require("lua/runtime")`)
	return L
}

func toLua(L *lua.LState, v interface{}) lua.LValue {
	switch t := v.(type) {
	case int:
		return lua.LNumber(float64(t))
	case int64:
		return lua.LNumber(float64(t))
	case uint:
		return lua.LNumber(float64(t))
	case uint64:
		return lua.LNumber(float64(t))
	case float32:
		return lua.LNumber(float64(t))
	case float64:
		return lua.LNumber(t)
	case string:
		return lua.LString(t)
	case []byte:
		return lua.LString(string(t))
	default:
		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Map:
			return mapToTable(L, val)
		case reflect.Slice:
			return sliceToTable(L, val)
		default:
			fmt.Println("unknown type")
			return lua.LNumber(0)
		}
	}
}

func mapToTable(L *lua.LState, data reflect.Value) lua.LValue {
	table := L.NewTable()
	for _, k := range data.MapKeys() {
		v := data.MapIndex(k).Interface()
		table.RawSetString(k.Interface().(string), toLua(L, v))
	}
	return table
}

func sliceToTable(L *lua.LState, data reflect.Value) lua.LValue {
	table := L.NewTable()
	for i := 0; i < data.Len(); i++ {
		table.RawSetInt(i, toLua(L, data.Index(i).Interface()))
	}
	return table
}

func (lt *Lua) getLua() *lua.LState {
	select {
	case L := <-lt.unused:
		return L
	default:
		return newLua()
	}
}

func (lt *Lua) putLua(L *lua.LState) {
	if !lt.reuse {
		L.Close()
		return
	}
	select {
	case lt.unused <- L:
	default:
		L.Close()
	}
}

func (lt *Lua) Render(w io.Writer, data interface{}, name string) error {
	L := lt.getLua()
	defer lt.putLua(L)

	write := func(L *lua.LState) int {
		n := L.GetTop()
		for i := 1; i <= n; i++ {
			s := L.ToString(i)
			w.Write([]byte(s))
		}
		return 0
	}

	name = path.Clean(name)
	dir, entry := path.Split(name)
	if entry == "" {
		name = path.Join(dir, "index")
	}
	name = path.Join("lua", "endpoints", name)

	// load the correct code
	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("require"),
		NRet:    1,
		Protect: true,
	}, lua.LString(name))

	if err != nil {
		return err
	}

	// build the html
	v := L.Get(-1)
	L.Pop(1)
	err = L.CallByParam(lua.P{
		Fn:      v,
		NRet:    1,
		Protect: true,
	}, toLua(L, data))

	if err != nil {
		return err
	}

	// render the html
	v = L.Get(-1)
	L.Pop(1)
	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal("render"),
		NRet:    0,
		Protect: true,
	}, v, L.NewFunction(write))

	if err != nil {
		return err
	}

	// add a newline to make it appear better in terminal
	w.Write([]byte{'\n'})
	return nil
}
