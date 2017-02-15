package glog

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"sync"
)

const (
	CHUNK_SIZE             = 4000
	MIN_STACK_OFFSET       = 3
	TOP_LEFT_CORNER        = "╔"
	BOTTOM_LEFT_CORNER     = "╚"
	MIDDLE_CORNER          = "╟"
	HORIZONTAL_DOUBLE_LINE = "║"
	DOUBLE_DIVIDER         = "════════════════════════════════════════════"
	SINGLE_DIVIDER         = "────────────────────────────────────────────"
)

var (
	TOP_BORDER    string
	MIDDLE_BORDER string
	BOTTOM_BORDER string
)

var bufferPool *sync.Pool

func init() {
	TOP_BORDER = strings.Join([]string{TOP_LEFT_CORNER, DOUBLE_DIVIDER, DOUBLE_DIVIDER}, "")
	MIDDLE_BORDER = strings.Join([]string{MIDDLE_CORNER, SINGLE_DIVIDER, SINGLE_DIVIDER}, "")
	BOTTOM_BORDER = strings.Join([]string{BOTTOM_LEFT_CORNER, DOUBLE_DIVIDER, DOUBLE_DIVIDER}, "")

	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
}

type LoggerPrinter struct {
	tag              string
	localMethodCount ThreadLocal
	localTag         ThreadLocal
	settings         *Settings
	currentThread    string
}

func NewLoggerPrinter() *LoggerPrinter {
	l := &LoggerPrinter{
		currentThread: "main",
		settings:      NewSettings(),
	}
	l.init(DEFAULT_TAG)
	return l
}

func (log *LoggerPrinter) setThread(tag string) {
	if tag == "" {
		panic("Thread may not be null")
	}
	if len(strings.TrimSpace(tag)) == 0 {
		panic("Thread may not be empty")
	}
	log.currentThread = tag
}

func (log *LoggerPrinter) init(tag string) {
	if tag == "" {
		panic("tag may not be null")
	}
	if len(strings.TrimSpace(tag)) == 0 {
		panic("tag may not be empty")
	}
	log.tag = tag
}

func (l *LoggerPrinter) d(message string, args ...interface{}) {
	l.log(LevelDEBUG, message, args...)
}

func (l *LoggerPrinter) e(message string, args ...interface{}) {
	l.log(LevelERROR, message, args...)
}

func (l *LoggerPrinter) w(message string, args ...interface{}) {
	l.log(LevelWARN, message, args...)
}

func (l *LoggerPrinter) v(message string, args ...interface{}) {
	l.log(LevelVERBOSE, message, args...)
}

func (l *LoggerPrinter) i(message string, args ...interface{}) {
	l.log(LevelINFO, message, args...)
}

func (l *LoggerPrinter) wtf(message string, args ...interface{}) {
	l.log(LevelASSERT, message, args...)
}

func (l *LoggerPrinter) getTag() string {
	return l.tag
}

func (l *LoggerPrinter) log(priority int, msg string, args ...interface{}) {
	if l.settings.getLogLevel() == LogLevelNONE {
		return
	}
	arrStr := make([]string, 0, 100)

	tag := l.getTag()
	message := createMessage(msg, args...)

	//if (throwable != null && message != null) {
	//message += " : \n\r\t\t\t" + strings.Join(CallerInfo(), "\n\r\t\t\t")
	////message += " : " + Helper.getStackTraceString(throwable)
	//}
	if message == "" {
		message = "No message/exception is set"
	}
	methodCount := l.getMethodCount()

	arrStr = append(arrStr, l.logTopBorder(priority, tag))
	arrStr = append(arrStr, l.logHeaderContent(priority, tag, methodCount)...)

	//get bytes of message with system's default charset (which is UTF-8 for Android)
	//bytes := bytes.NewBufferString(message)
	length := len(message)
	if length <= CHUNK_SIZE {
		if methodCount > 0 {
			arrStr = append(arrStr, l.logDivider(priority, tag))
		}
		arrStr = append(arrStr, l.logContent(priority, tag, message)...)
		arrStr = append(arrStr, l.logBottomBorder(priority, tag))
		l.logChunk(priority, tag, strings.Join(arrStr, "\n"))
		return
	}
	if methodCount > 0 {
		arrStr = append(arrStr, l.logDivider(priority, tag))
	}
	for i := 0; i < length; i += CHUNK_SIZE {
		count := int(math.Min(float64(length-i), CHUNK_SIZE))
		//create a new String with system's default charset (which is UTF-8)
		arrStr = append(arrStr, l.logContent(priority, tag, message[i:i+count])...)
	}
	arrStr = append(arrStr, l.logBottomBorder(priority, tag))
	l.logChunk(priority, tag, strings.Join(arrStr, "\n"))
}

func (l *LoggerPrinter) getMethodCount() int {
	count := converInterfaceToInt(l.localMethodCount.get())
	result := l.settings.getMethodCount()
	if count != 0 {
		l.localMethodCount.remove()
		result = count
	}
	if result < 0 {
		panic("methodCount cannot be negative")
	}
	return result
}

func (l *LoggerPrinter) logTopBorder(logType int, tag string) string {
	return TOP_BORDER
}

func (l *LoggerPrinter) logChunk(logType int, tag, chunk string) {
	finalTag := l.formatTag(tag)
	switch logType {
	case LevelERROR:
		l.settings.getLogAdapter().e(finalTag, chunk)
	case LevelINFO:
		l.settings.getLogAdapter().i(finalTag, chunk)
	case LevelVERBOSE:
		l.settings.getLogAdapter().v(finalTag, chunk)
	case LevelWARN:
		l.settings.getLogAdapter().w(finalTag, chunk)
	case LevelASSERT:
		l.settings.getLogAdapter().wtf(finalTag, chunk)
	case LevelDEBUG:
	// Fall through, log debug by default
	default:
		l.settings.getLogAdapter().d(finalTag, chunk)
	}
}

func (l *LoggerPrinter) formatTag(tag string) string {
	if tag == "" && !strings.EqualFold(l.tag, tag) {
		return l.tag + "-" + tag
	}
	return l.tag
}

func (l *LoggerPrinter) logHeaderContent(logType int, tag string, methodCount int) []string {
	arrStr := make([]string, 0)

	trace := getStackTrace()
	if l.settings.isShowThreadInfo() {
		arrStr = append(arrStr, HORIZONTAL_DOUBLE_LINE+" Thread: "+l.currentThread)
		arrStr = append(arrStr, l.logDivider(logType, tag))
	}
	level := ""
	stackOffset := getStackOffset(trace) + l.settings.getMethodOffset()

	//corresponding method count with the current stack may exceeds the stack trace. Trims the count
	if methodCount+stackOffset > len(trace) {
		methodCount = len(trace) - stackOffset - 1
	}

	for i := methodCount; i > 0; i-- {
		stackIndex := i + stackOffset
		if stackIndex >= len(trace) {
			continue
		}
		builder := bufferPool.Get().(*bytes.Buffer)
		builder.Reset()
		defer bufferPool.Put(builder)

		//builder := new(bytes.Buffer)
		builder.WriteString("║ ")
		builder.WriteString(level)
		builder.WriteString(getSimpleClassName(trace[stackIndex].getClassName()))
		builder.WriteByte('.')
		builder.WriteString(trace[stackIndex].getMethodName())
		builder.WriteByte(' ')
		builder.WriteString(" (")
		builder.WriteString(trace[stackIndex].getFileName())
		builder.WriteByte(':')
		builder.WriteString(trace[stackIndex].getLineNumber())
		builder.WriteByte(')')
		level += "\t"

		arrStr = append(arrStr, builder.String())
	}
	return arrStr
}

func (l *LoggerPrinter) logDivider(logType int, tag string) string {
	return MIDDLE_BORDER
}

func (l *LoggerPrinter) logContent(logType int, tag, chunk string) []string {
	lines := strings.Split(chunk, "\n")
	for i := range lines {
		lines[i] = HORIZONTAL_DOUBLE_LINE + lines[i]
	}
	return lines
}

func (l *LoggerPrinter) logBottomBorder(logType int, tag string) string {
	return BOTTOM_BORDER
}

func getStackOffset(trace []*StackTraceElement) int {
	for i := MIN_STACK_OFFSET; i < len(trace); i++ {
		e := trace[i]
		name := e.getClassName()
		if !strings.EqualFold(name, classNameLogger) {
			return i - 1
		}
	}
	return -1
}

func createMessage(message string, args ...interface{}) string {
	if len(args) > 0 {
		if strings.Contains(message, "%") && !strings.Contains(message, "%%") {
			return fmt.Sprintf(message, args...)
		} else {
			return fmt.Sprintf(message+strings.Repeat(" %v", len(args)), args...)
		}
	}
	return message
}

func getSimpleClassName(name string) string {
	lastIndex := strings.LastIndex(name, ".")
	return name[lastIndex+1:]
}
