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

var UnknownHostExceptionName = reflect.TypeOf(UnknownHostException{}).Name()

type UnknownHostException struct {
	Exception
	Elog log.Messages
}

func NewUnknownHostException(parent Exception, message log.Message) *UnknownHostException {
	return &UnknownHostException{parent, log.Messages{message}}
}

func (e UnknownHostException) Code() int64 {
	return UnknownHostExceptionCode
}

func (e UnknownHostException) Name() string {
	return UnknownHostExceptionName
}

func (e UnknownHostException) What() string {
	return "Unknown Host"
}

func (e *UnknownHostException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e UnknownHostException) GetLog() log.Messages {
	return e.Elog
}

func (e UnknownHostException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e UnknownHostException) DetailMessage() string {
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

func (e UnknownHostException) String() string {
	return e.DetailMessage()
}

func (e UnknownHostException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: UnknownHostExceptionCode,
		Name: UnknownHostExceptionName,
		What: "Unknown Host",
	}

	return json.Marshal(except)
}

func (e UnknownHostException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*UnknownHostException):
		callback(&e)
		return true
	case func(UnknownHostException):
		callback(e)
		return true
	default:
		return false
	}
}
