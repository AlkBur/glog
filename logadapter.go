package glog

type LogAdapter interface {
	d(tag, message string)
	e(tag, message string)
	w(tag, message string)
	i(tag, message string)
	v(tag, message string)
	wtf(tag, message string)
}
