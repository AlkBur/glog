package glog

import (
	"fmt"
	"sync"
)

type ConsoleLogAdapter struct {
	//out io.Writer
	//err io.Writer
	mu sync.RWMutex
}

func NewConsoleLogAdapter() *ConsoleLogAdapter {
	return &ConsoleLogAdapter{
	//out: os.Stderr,
	//err: os.Stderr,
	}
}

func (ad *ConsoleLogAdapter) d(tag, message string) {
	ad.log(LevelDEBUG, tag, message)
}

func (ad *ConsoleLogAdapter) e(tag, message string) {
	ad.log(LevelERROR, tag, message)
}

func (ad *ConsoleLogAdapter) w(tag, message string) {
	ad.log(LevelWARN, tag, message)
}

func (ad *ConsoleLogAdapter) i(tag, message string) {
	ad.log(LevelINFO, tag, message)
}

func (ad *ConsoleLogAdapter) v(tag, message string) {
	ad.log(LevelVERBOSE, tag, message)
}

func (ad *ConsoleLogAdapter) wtf(tag, message string) {
	ad.log(LevelASSERT, tag, message)
}

func (ad *ConsoleLogAdapter) log(lv int, tag, message string) {
	ad.mu.Lock()
	defer ad.mu.Unlock()

	fmt.Println(message)
}
