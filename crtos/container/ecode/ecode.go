package ecode

import (
	"fmt"
	"strconv"
	"sync/atomic"
)

var (
	_messages atomic.Value         // NOTE: stored map[int]string
	_codes    = map[int]struct{}{} // register codes.
)

//Code ..
type Code int

// Codes ecode error interface which has a code & message.
type Codes interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	//Detail get error detail,it may be nil.
	Details() []interface{}
}

// Int parse code int to error.
func Int(i int) Code { return Code(i) }

func add(e int) Code {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = struct{}{}
	return Int(e)
}

func (e Code) Error() string {
	return strconv.FormatInt(int64(e), 10)
}

// Code return error code
func (e Code) Code() int { return int(e) }

//Message ..
func (e Code) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[e.Code()]; ok {
			return msg
		}
	}
	return e.Error()
}

//Details return details.
func (e Code) Details() []interface{} { return nil }
