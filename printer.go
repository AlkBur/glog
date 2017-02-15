package glog

type Printer interface {
	init(tag string)
	d(message string, args ...interface{})
	e(message string, args ...interface{})
	w(message string, args ...interface{})
	i(message string, args ...interface{})
	v(message string, args ...interface{})
	wtf(message string, args ...interface{})
	log(priority int, msg string, args ...interface{})
}
