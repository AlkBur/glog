package glog

//
//import (
//	"fmt"
//	"github.com/AlkBur/ColorAnsi"
//	"reflect"
//	"syscall"
//	//"unsafe"
//)
//
//
////
//var (
//	kernel32 = syscall.NewLazyDLL("kernel32.dll")
//	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
//
//)
////
////func init() {
////	var st uint32
////	r, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(v.Fd()), uintptr(unsafe.Pointer(&st)), 0)
////	if r != 0 && e == 0 {
////		fmt.Println("Termn")
////	}
////}
//
//type FormatterStr struct {
//	ColorAnsi.Color
//	rf reflect.Value
//	val interface{}
//}
//
//type arrFormatterStr []*FormatterStr
//
//func NewFormatterStr(v ...interface{}) *FormatterStr {
//	return &FormatterStr{
//		rf: reflect.ValueOf(v),
//		val: v,
//	}
//}
//
//func (lvl Levels) String() string {
//	switch lvl {
//	case LevelDebug:
//		return "debug"
//	case LevelInfo:
//		return "info"
//	case LevelWarn:
//		return "warning"
//	case LevelError:
//		return "error"
//	case LevelFatal:
//		return "fatal"
//	}
//	return "unknown"
//}
//
//func (arr arrFormatterStr) String() []string {
//	rez := make([]string, 0, len(arr))
//
//
//	return rez
//}
//
//func stringFormat(v []interface{}) []string {
//	for i, arg := range v {
//		switch arg.(type) {
//		case error:
//			v[i] = fmt.Sprint(arg)
//		}
//	}
//
//	formated := formatArgs(v)
//
//	return formated.String()
//}
//
//func formatArgs(args []interface{}) arrFormatterStr {
//	formatted := make(arrFormatterStr, 0, len(args))
//	for _, a := range args {
//		fs := NewFormatterStr(a)
//		fs.SetForegroundLightCyan()
//		formatted = append(formatted, fs)
//	}
//	return formatted
//}
//
//func (fs *FormatterStr)String() string {
//	w := fmt.Sprint(reflect.ValueOf(fs.rf).Interface())
//	return fs.Color.GetString(fmt.Sprint(w))
//}
//
//func ColorPrint(s string, i int) {
//	kernel32 := syscall.NewLazyDLL("kernel32.dll")
//	proc := kernel32.NewProc("SetConsoleTextAttribute")
//	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
//	fmt.Print(s)
//	handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7))
//	CloseHandle := kernel32.NewProc("CloseHandle")
//	CloseHandle.Call(handle)
//}
