package glog

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type ThreadLocal struct {
	value interface{}
	sync.RWMutex
}

func (tl *ThreadLocal) get() interface{} {
	tl.RLock()
	defer tl.RUnlock()
	return tl.value
}

func (tl *ThreadLocal) set(val interface{}) {
	tl.Lock()
	defer tl.Unlock()
	tl.value = val
}

func (tl *ThreadLocal) remove() {
	tl.Lock()
	defer tl.Unlock()
	tl.value = nil
}

func StackTrace(all bool) []string {
	// Reserve 10K buffer at first
	buf := make([]byte, 10240)

	for {
		size := runtime.Stack(buf, all)
		// The size of the buffer may be not enough to hold the stacktrace,
		// so double the buffer size
		if size == len(buf) {
			buf = make([]byte, len(buf)<<1)
			continue
		}
		break
	}

	return formatStack(buf)
}

func formatStack(buf []byte) []string {
	buf = bytes.Trim(buf, "\x00")
	str := string(buf)
	str = strings.Replace(str, "\t", "", -1)

	stack := strings.Split(str, "\n")

	return stack[1 : len(stack)-1]
}

func converInterfaceToInt(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	default:
		return 0
	}
}

func CallerInfo2() []*StackTraceElement {
	strArr := StackTrace(false)
	arr := make([]*StackTraceElement, 0, len(strArr))
	for _, str := range strArr {
		el := &StackTraceElement{
			file: str,
		}
		arr = append(arr, el)
	}
	return arr
}

func Display(x interface{}) string {
	var count int
	return display(reflect.ValueOf(x), &count)
}

func display(v reflect.Value, count *int) string {
	*count++
	if *count > 10000 {
		return ""
	}

	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Slice, reflect.Array:
		arr := make([]string, v.Len()+2)
		arr[0] = "["
		for i := 0; i < v.Len(); i++ {
			arr[i+1] = display(v.Index(i), count)
		}
		arr[len(arr)-1] = "]"
		return strings.Join(arr, "\n\r")
	case reflect.Struct:
		arr := make([]string, v.NumField()+2)
		arr[0] = "{"
		for i := 0; i < v.NumField(); i++ {
			arr[i+1] = fmt.Sprintf("%s = %s", v.Type().Field(i).Name, display(v.Field(i), count))
		}
		arr[len(arr)-1] = "}"
		return strings.Join(arr, "\n\r")
	case reflect.Map:
		arr := make([]string, len(v.MapKeys())+2)
		arr[0] = "{"
		for i, key := range v.MapKeys() {
			arr[i+1] = fmt.Sprintf("%s = %s", formatAtom(key), display(v.MapIndex(key), count))
		}
		arr[len(arr)-1] = "}"
		return strings.Join(arr, "\n\r")
	case reflect.Ptr:
		if v.IsNil() {
			return "nil"
		} else {
			return display(v.Elem(), count)
		}
	case reflect.Interface:
		if v.IsNil() {
			return "nil"
		} else {
			return fmt.Sprintf("type = %s; val = %s", v.Elem().Type(), display(v.Elem(), count))
		}
	default:
		return formatAtom(v)
	}
}

func formatAtom(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Invalid:
		return "invalid"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', 6, 64)
	//case reflect.Complex64, reflect.Complex128:
	//	return string(v.Complex())
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.String:
		return strconv.Quote(v.String())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return v.Type().String() + " ©x" + strconv.FormatUint(uint64(v.Pointer()), 16)
	default:
		return v.Type().String() + " value"
	}
}
