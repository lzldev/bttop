package proctype

type ProcType int

const (
	_ ProcType = iota
	Proccess
	Thread
)

func (p ProcType) String() string {
	switch p {
	case Proccess:
		return "Proccess"
	case Thread:
		return "Thread"
	default:
		return "Unknown"
	}
}
