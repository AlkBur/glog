package glog

import "testing"

func TestDebugLog(t *testing.T) {
	D("message")
}

func TestVerboseLog(t *testing.T) {
	V("message")
}

func TestWarningLog(t *testing.T) {
	W("message")
}

func TestErrorLog(t *testing.T) {
	E("message")
}

func TestInfoLog(t *testing.T) {
	I("message")
}
