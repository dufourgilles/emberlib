package errors

import (
	"fmt"
	"runtime"
)

type ErrorStruct struct {
	Message error
	Stack []string
}

type Error *ErrorStruct

func New(format string, params ...interface{}) Error {
	_,file,line,_ := runtime.Caller(1)
	stack := []string{fmt.Sprintf("%s:%d", file,line)}
	e := ErrorStruct{Message: fmt.Errorf(format, params...), Stack: stack}
	return &e
}

func NewError(err error) Error {
	if err == nil { return nil }
	_,file,line,_ := runtime.Caller(1)
	stack := []string{fmt.Sprintf("%s:%d", file,line)}
	e := ErrorStruct{Message: err, Stack: stack}
	return &e
}

func Update(e Error) Error {
	if e == nil { return nil}
	_,file,line,_ := runtime.Caller(1)
	stack := fmt.Sprintf("%s:%d", file,line)
	e.Stack = append(e.Stack, stack)
	return e
}