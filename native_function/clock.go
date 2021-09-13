package native_function

import (
	"time"
)

// ClockFunc is struct of clock function
type ClockFunc struct{}

// NewClockFunc is constructor of ClockFunc
func NewClockFunc() *ClockFunc {
	return &ClockFunc{}
}

// Arity returns 0
func (cf *ClockFunc) Arity() int {
	return 0
}

// Call return now unix time
func (cf *ClockFunc) Call(arguments []interface{}) (interface{}, error) {
	return float64(time.Now().UnixMilli()) / 1000.0, nil
}

func (cf *ClockFunc) String() string {
	return "<native fn>"
}
