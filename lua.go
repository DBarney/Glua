package glua

import (
	"fmt"
	"io"
	"path"
	"reflect"

	"github.com/yuin/gopher-lua"
)

// Lua allows for storage of multiple lua states. Since lua is single threaded we need
// multiple states to be able to not bottle neck on a single process.
// all functions on this struct are safe to use across different processes.
type Lua struct {
	unused chan *lua.LState
	reuse  bool
}

// New creates and configures the holder for lua states.
func New() *Lua {
	return &Lua{
		reuse:  false,
		unused: make(chan *lua.LState, 4),
	}
}

func newLua() *lua.LState {
	L := lua.NewState()
	L.OpenLibs()

	// we want our default functions available
	L.DoString(`require("runtime")`)
	return L
}

// toLua converts most golang values into formats that can be used inside of lua
// this allows creating complex multi leveled objects from golang values.
// Currently we support:
// numbers
// strings
// maps
// slices
// TODO: add support for converting objects to maps. I tried before but it got
// too complex for what I was looking for
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

// mapToTable converts a golang map to a lua table
func mapToTable(L *lua.LState, data reflect.Value) lua.LValue {
	table := L.NewTable()
	for _, k := range data.MapKeys() {
		v := data.MapIndex(k).Interface()
		table.RawSetString(k.Interface().(string), toLua(L, v))
	}
	return table
}

// sliveToTable converts a golang slice to a lua table
func sliceToTable(L *lua.LState, data reflect.Value) lua.LValue {
	table := L.NewTable()
	for i := 0; i < data.Len(); i++ {
		table.RawSetInt(i, toLua(L, data.Index(i).Interface()))
	}
	return table
}

// getLua gets a stashed lua state, or if one is not available, it creates and initializeds a new one
func (lt *Lua) getLua() *lua.LState {
	select {
	case L := <-lt.unused:
		return L
	default:
		return newLua()
	}
}

// putLua caches a lua state so that it can be reused. If we already have enough cached, the state is closed
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

// Render is used to render the html page. It accepts `data interface{}` which allows setting arbitrary data in the template
// and it also accepts `name` as a parameter which is the name of the lua template to run.
// the template is written to the `w io.Writer` interface.
func (lt *Lua) Render(w io.Writer, data interface{}, name string) error {
	L := lt.getLua()
	defer lt.putLua(L)

	// this is the write function that is called from inside lua to write the html data.
	write := func(L *lua.LState) int {
		// get number of arguments that were passed to the function
		n := L.GetTop()
		// write each argument to the io.Writer
		for i := 1; i <= n; i++ {
			s := L.ToString(i)
			w.Write([]byte(s))
		}
		// there are no return values, so we return 0
		return 0
	}

	// TODO: these should be configurable, right now we expect the template files to be under the `lua/endpoints` folder.
	// we locad them from disk so that they can be edited and worked on without needing to recompile the golang binary
	name = path.Clean(name)
	dir, entry := path.Split(name)
	if entry == "" {
		name = path.Join(dir, "index")
	}

	// load the correct code using the lua built in require
	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("require"),
		NRet:    1,
		Protect: true,
	}, lua.LString(name))

	if err != nil {
		return err
	}

	// the require function returns a single value, which is a function representing the file that was required
	template := L.Get(-1)
	L.Pop(1)

	// call the template function and pass the data to be rendered into it.
	err = L.CallByParam(lua.P{
		Fn:      template,
		NRet:    1,
		Protect: true,
	}, toLua(L, data))

	if err != nil {
		return err
	}

	// there is only 1 return value, which is an object that represents the html data to be rendered:
	// {__name="html",{__name="body",{__name="h1","this is my site}}}
	structure := L.Get(-1)
	L.Pop(1)

	// render the html by calling the lua render function and passing in the structure and the write function
	// as arguments
	err = L.CallByParam(lua.P{
		Fn:      L.GetGlobal("render"),
		NRet:    0,
		Protect: true,
	}, structure, L.NewFunction(write))

	if err != nil {
		return err
	}

	// add a newline to make it appear better in terminal
	w.Write([]byte{'\n'})
	return nil
}
