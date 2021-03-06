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

var BinaryenExceptionName = reflect.TypeOf(BinaryenException{}).Name()

type BinaryenException struct {
	_WasmException
	Elog log.Messages
}

func NewBinaryenException(parent _WasmException, message log.Message) *BinaryenException {
	return &BinaryenException{parent, log.Messages{message}}
}

func (e BinaryenException) Code() int64 {
	return 3070005
}

func (e BinaryenException) Name() string {
	return BinaryenExceptionName
}

func (e BinaryenException) What() string {
	return "binaryen exception"
}

func (e *BinaryenException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BinaryenException) GetLog() log.Messages {
	return e.Elog
}

func (e BinaryenException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e BinaryenException) DetailMessage() string {
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

func (e BinaryenException) String() string {
	return e.DetailMessage()
}

func (e BinaryenException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3070005,
		Name: BinaryenExceptionName,
		What: "binaryen exception",
	}

	return json.Marshal(except)
}

func (e BinaryenException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BinaryenException):
		callback(&e)
		return true
	case func(BinaryenException):
		callback(e)
		return true
	default:
		return false
	}
}
