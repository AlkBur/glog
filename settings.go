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

func (s *Settings) SetLogLevel(lv LogLevel) {
	s.logLevel = lv
}

func (s *Settings) GetLogLevel() LogLevel {
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

func (s *Settings) IsShowThreadInfo() bool {
	return s.showThreadInfo
}

func (s *Settings) HideThreadInfo() {
	s.showThreadInfo = false
}

func (s *Settings) Reset() {
	s.methodCount = 2
	s.methodOffset = 0
	s.showThreadInfo = true
	s.logLevel = LogLevelFULL
}

func (s *Settings) getMethodOffset() int {
	return s.methodOffset
}

func (s *Settings) SetMethodOffset(skip int) {
	s.methodOffset = skip
}
