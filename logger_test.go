package glog

import "testing"

func TestDebugLog(t *testing.T) {
	d("message")
}

func TestVerboseLog(t *testing.T) {
	v("message")
}

func TestWarningLog(t *testing.T) {
	w("message")
}

func TestErrorLog(t *testing.T) {
	e("message")
}
