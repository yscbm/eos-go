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

var EosioAssertCodeExceptionName = reflect.TypeOf(EosioAssertCodeException{}).Name()

type EosioAssertCodeException struct {
	_ActionValidateException
	Elog log.Messages
}

func NewEosioAssertCodeException(parent _ActionValidateException, message log.Message) *EosioAssertCodeException {
	return &EosioAssertCodeException{parent, log.Messages{message}}
}

func (e EosioAssertCodeException) Code() int64 {
	return 3050004
}

func (e EosioAssertCodeException) Name() string {
	return EosioAssertCodeExceptionName
}

func (e EosioAssertCodeException) What() string {
	return "eosio_assert_code assertion failure"
}

func (e *EosioAssertCodeException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e EosioAssertCodeException) GetLog() log.Messages {
	return e.Elog
}

func (e EosioAssertCodeException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e EosioAssertCodeException) DetailMessage() string {
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

func (e EosioAssertCodeException) String() string {
	return e.DetailMessage()
}

func (e EosioAssertCodeException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3050004,
		Name: EosioAssertCodeExceptionName,
		What: "eosio_assert_code assertion failure",
	}

	return json.Marshal(except)
}

func (e EosioAssertCodeException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*EosioAssertCodeException):
		callback(&e)
		return true
	case func(EosioAssertCodeException):
		callback(e)
		return true
	default:
		return false
	}
}