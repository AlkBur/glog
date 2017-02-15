package glog

type Printer interface {
	init(tag string)
	D(message string, args ...interface{})
	E(message string, args ...interface{})
	W(message string, args ...interface{})
	I(message string, args ...interface{})
	V(message string, args ...interface{})
	WTF(message string, args ...interface{})
	log(priority int, msg string, args ...interface{})
}
