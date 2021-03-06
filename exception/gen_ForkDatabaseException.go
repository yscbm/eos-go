// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var ForkDatabaseExceptionName = reflect.TypeOf(ForkDatabaseException{}).Name()

type ForkDatabaseException struct {
	_ForkDatabaseException
	Elog log.Messages
}

func NewForkDatabaseException(parent _ForkDatabaseException, message log.Message) *ForkDatabaseException {
	return &ForkDatabaseException{parent, log.Messages{message}}
}

func (e ForkDatabaseException) Code() int64 {
	return 3020000
}

func (e ForkDatabaseException) Name() string {
	return ForkDatabaseExceptionName
}

func (e ForkDatabaseException) What() string {
	return "Fork database exception"
}

func (e *ForkDatabaseException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ForkDatabaseException) GetLog() log.Messages {
	return e.Elog
}

func (e ForkDatabaseException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ForkDatabaseException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, l := range e.Elog {
		buffer.WriteString("[")
		buffer.WriteString(l.GetMessage())
		buffer.WriteString("]")
		buffer.WriteString("\n")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e ForkDatabaseException) String() string {
	return e.DetailMessage()
}

func (e ForkDatabaseException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3020000,
		Name: ForkDatabaseExceptionName,
		What: "Fork database exception",
	}

	return json.Marshal(except)
}

func (e ForkDatabaseException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ForkDatabaseException):
		callback(&e)
		return true
	case func(ForkDatabaseException):
		callback(e)
		return true
	default:
		return false
	}
}
