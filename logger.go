package glog

const DEFAULT_TAG = "main"
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

func D(v ...interface{}) {
	message, args := getMessage(v)
	printer.D(message, args...)
}

func E(v ...interface{}) {
	message, args := getMessage(v)
	printer.E(message, args...)
}

func V(v ...interface{}) {
	message, args := getMessage(v)
	printer.V(message, args...)
}

func W(v ...interface{}) {
	message, args := getMessage(v)
	printer.W(message, args...)
}

func I(v ...interface{}) {
	message, args := getMessage(v)
	printer.I(message, args...)
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
