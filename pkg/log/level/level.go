package level

type XLevel uint8

const (
	InfoLevel XLevel = iota
	TraceLevel
	WarnLevel
	ErrorLevel
)

func (xlvl XLevel) String() string {
	switch xlvl {
	case InfoLevel:
		return "info"
	case TraceLevel:
		return "trace"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return ""
	}
}

func (xlvl *XLevel) Enable(lvl XLevel) bool {
	return lvl > *xlvl
}
