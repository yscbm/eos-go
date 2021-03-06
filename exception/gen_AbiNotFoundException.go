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

var AbiNotFoundExceptionName = reflect.TypeOf(AbiNotFoundException{}).Name()

type AbiNotFoundException struct {
	_AbiException
	Elog log.Messages
}

func NewAbiNotFoundException(parent _AbiException, message log.Message) *AbiNotFoundException {
	return &AbiNotFoundException{parent, log.Messages{message}}
}

func (e AbiNotFoundException) Code() int64 {
	return 3150001
}

func (e AbiNotFoundException) Name() string {
	return AbiNotFoundExceptionName
}

func (e AbiNotFoundException) What() string {
	return "No ABI Found"
}

func (e *AbiNotFoundException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e AbiNotFoundException) GetLog() log.Messages {
	return e.Elog
}

func (e AbiNotFoundException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e AbiNotFoundException) DetailMessage() string {
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

func (e AbiNotFoundException) String() string {
	return e.DetailMessage()
}

func (e AbiNotFoundException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3150001,
		Name: AbiNotFoundExceptionName,
		What: "No ABI Found",
	}

	return json.Marshal(except)
}

func (e AbiNotFoundException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*AbiNotFoundException):
		callback(&e)
		return true
	case func(AbiNotFoundException):
		callback(e)
		return true
	default:
		return false
	}
}
