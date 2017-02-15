package glog

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
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

const DefaultTimestampFormat = time.RFC3339

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
}

func NewLoggerPrinter() *LoggerPrinter {
	l := &LoggerPrinter{
		settings: NewSettings(),
	}
	l.init(DEFAULT_TAG)
	return l
}

func (log *LoggerPrinter) SetTag(tag string) {
	if tag == "" {
		panic("Thread may not be null")
	}
	if len(strings.TrimSpace(tag)) == 0 {
		panic("Thread may not be empty")
	}
	log.tag = tag
}

func (log *LoggerPrinter) GetTag() string {
	return log.tag
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

func (l *LoggerPrinter) D(message string, args ...interface{}) {
	l.log(LevelDEBUG, message, args...)
}

func (l *LoggerPrinter) E(message string, args ...interface{}) {
	l.log(LevelERROR, message, args...)
}

func (l *LoggerPrinter) W(message string, args ...interface{}) {
	l.log(LevelWARN, message, args...)
}

func (l *LoggerPrinter) V(message string, args ...interface{}) {
	l.log(LevelVERBOSE, message, args...)
}

func (l *LoggerPrinter) I(message string, args ...interface{}) {
	l.log(LevelINFO, message, args...)
}

func (l *LoggerPrinter) WTF(message string, args ...interface{}) {
	l.log(LevelASSERT, message, args...)
}

func (l *LoggerPrinter) log(priority int, msg string, args ...interface{}) {
	if l.settings.GetLogLevel() == LogLevelNONE {
		return
	}
	arrStr := make([]string, 0, 100)

	message := createMessage(msg, args...)

	//if (throwable != null && message != null) {
	//message += " : \n\r\t\t\t" + strings.Join(CallerInfo(), "\n\r\t\t\t")
	////message += " : " + Helper.getStackTraceString(throwable)
	//}
	if message == "" {
		message = "No message/exception is set"
	}
	methodCount := l.getMethodCount()

	arrStr = append(arrStr, l.logTopBorder(priority))
	arrStr = append(arrStr, l.logHeaderContent(priority, methodCount)...)

	//get bytes of message with system's default charset (which is UTF-8 for Android)
	//bytes := bytes.NewBufferString(message)
	length := len(message)
	if length <= CHUNK_SIZE {
		if methodCount > 0 {
			arrStr = append(arrStr, l.logDivider(priority))
		}
		arrStr = append(arrStr, l.logContent(priority, message)...)
		arrStr = append(arrStr, l.logBottomBorder(priority))
		l.logChunk(priority, strings.Join(arrStr, "\n"))
		return
	}
	if methodCount > 0 {
		arrStr = append(arrStr, l.logDivider(priority))
	}
	for i := 0; i < length; i += CHUNK_SIZE {
		count := int(math.Min(float64(length-i), CHUNK_SIZE))
		//create a new String with system's default charset (which is UTF-8)
		arrStr = append(arrStr, l.logContent(priority, message[i:i+count])...)
	}
	arrStr = append(arrStr, l.logBottomBorder(priority))
	l.logChunk(priority, strings.Join(arrStr, "\n"))
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

func (l *LoggerPrinter) logTopBorder(logType int) string {
	return TOP_BORDER
}

func (l *LoggerPrinter) logChunk(logType int, chunk string) {
	finalTag := l.GetTag()
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

func (l *LoggerPrinter) logHeaderContent(logType int, methodCount int) []string {
	arrStr := make([]string, 0)

	trace := getStackTrace()
	if l.settings.IsShowThreadInfo() {
		arrStr = append(arrStr, HORIZONTAL_DOUBLE_LINE+" Tag: "+l.tag+"\t time="+time.Now().Format(DefaultTimestampFormat))
		arrStr = append(arrStr, l.logDivider(logType))
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

func (l *LoggerPrinter) logDivider(logType int) string {
	return MIDDLE_BORDER
}

func (l *LoggerPrinter) logContent(logType int, chunk string) []string {
	lines := strings.Split(chunk, "\n")
	for i := range lines {
		lines[i] = HORIZONTAL_DOUBLE_LINE + lines[i]
	}
	return lines
}

func (l *LoggerPrinter) logBottomBorder(logType int) string {
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
