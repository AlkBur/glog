package glog

type Settings struct {
	methodCount    int
	showThreadInfo bool
	methodOffset   int
	logAdapter     LogAdapter
	logLevel       LogLevel
}

func NewSettings() *Settings {
	return &Settings{
		methodCount:    2,
		showThreadInfo: true,
		methodOffset:   0,
		logLevel:       LogLevelFULL,
	}
}

func (s *Settings) setLogLevel(lv LogLevel) {
	s.logLevel = lv
}

func (s *Settings) getLogLevel() LogLevel {
	return s.logLevel
}

func (s *Settings) getMethodCount() int {
	return s.methodCount
}

func (s *Settings) getLogAdapter() LogAdapter {
	if s.logAdapter == nil {
		s.logAdapter = NewConsoleLogAdapter()
	}
	return s.logAdapter
}

func (s *Settings) setlogAdapter(logAdapter LogAdapter) {
	s.logAdapter = logAdapter
}

func (s *Settings) isShowThreadInfo() bool {
	return s.showThreadInfo
}

func (s *Settings) hideThreadInfo() {
	s.showThreadInfo = false
}

func (s *Settings) reset() {
	s.methodCount = 2
	s.methodOffset = 0
	s.showThreadInfo = true
	s.logLevel = LogLevelFULL
}

func (s *Settings) getMethodOffset() int {
	return s.methodOffset
}
