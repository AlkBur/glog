package glog

const DEFAULT_TAG = "PRETTYLOGGER"
const classNameLogger = "glog"

const (
	LevelVERBOSE = 2
	LevelDEBUG   = 3
	LevelINFO    = 4
	LevelWARN    = 5
	LevelERROR   = 6
	LevelASSERT  = 7
)

var printer Printer

func init() {
	printer = NewLoggerPrinter()
}

func d(v ...interface{}) {
	message, args := getMessage(v)
	printer.d(message, args...)
}

func e(v ...interface{}) {
	message, args := getMessage(v)
	printer.e(message, args...)
}

func v(v ...interface{}) {
	message, args := getMessage(v)
	printer.v(message, args...)
}

func w(v ...interface{}) {
	message, args := getMessage(v)
	printer.w(message, args...)
}

func getMessage(v []interface{}) (message string, arg []interface{}) {
	if len(v) > 0 {
		switch t := v[0].(type) {
		case string:
			message = t
			arg = v[1:]
		default:
			arg = v[:]
		}
	}
	return
}
