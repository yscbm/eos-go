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

var PermissionQueryExceptionName = reflect.TypeOf(PermissionQueryException{}).Name()

type PermissionQueryException struct {
	_DatabaseException
	Elog log.Messages
}

func NewPermissionQueryException(parent _DatabaseException, message log.Message) *PermissionQueryException {
	return &PermissionQueryException{parent, log.Messages{message}}
}

func (e PermissionQueryException) Code() int64 {
	return 3060001
}

func (e PermissionQueryException) Name() string {
	return PermissionQueryExceptionName
}

func (e PermissionQueryException) What() string {
	return "Permission Query Exception"
}

func (e *PermissionQueryException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e PermissionQueryException) GetLog() log.Messages {
	return e.Elog
}

func (e PermissionQueryException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e PermissionQueryException) DetailMessage() string {
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

func (e PermissionQueryException) String() string {
	return e.DetailMessage()
}

func (e PermissionQueryException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3060001,
		Name: PermissionQueryExceptionName,
		What: "Permission Query Exception",
	}

	return json.Marshal(except)
}

func (e PermissionQueryException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*PermissionQueryException):
		callback(&e)
		return true
	case func(PermissionQueryException):
		callback(e)
		return true
	default:
		return false
	}
}
