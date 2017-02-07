package glog

type Levels int

const (
	LevelFatal Levels = iota
	LevelError
	LevelInfo
	LevelWarn
	LevelDebug
)
