package debug

import "fmt"

type Category int

const (
	Sim   Category = 1 << iota
	Move
	Path
	Clock
	Job
	World
)

var (
	master bool
	flags  Category
)

func SetEnabled(v bool) {
	master = v
}

func IsEnabled(c Category) bool {
	return master && flags&c != 0
}

func Toggle(c Category) {
	flags ^= c
	fmt.Printf("[DEBUG] %s\n", ActiveString())
}

func ActiveString() string {
	s := "Debug: "
	if flags&Sim != 0 {
		s += "Sim "
	}
	if flags&Move != 0 {
		s += "Move "
	}
	if flags&Path != 0 {
		s += "Path "
	}
	if flags&Clock != 0 {
		s += "Clock "
	}
	if flags&Job != 0 {
		s += "Job "
	}
	if flags&World != 0 {
		s += "World "
	}
	if flags == 0 {
		s += "none"
	}
	return s
}
