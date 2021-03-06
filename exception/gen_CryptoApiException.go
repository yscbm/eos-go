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

var CryptoApiExceptionName = reflect.TypeOf(CryptoApiException{}).Name()

type CryptoApiException struct {
	_ContractApiException
	Elog log.Messages
}

func NewCryptoApiException(parent _ContractApiException, message log.Message) *CryptoApiException {
	return &CryptoApiException{parent, log.Messages{message}}
}

func (e CryptoApiException) Code() int64 {
	return 3230001
}

func (e CryptoApiException) Name() string {
	return CryptoApiExceptionName
}

func (e CryptoApiException) What() string {
	return "Crypto API exception"
}

func (e *CryptoApiException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e CryptoApiException) GetLog() log.Messages {
	return e.Elog
}

func (e CryptoApiException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e CryptoApiException) DetailMessage() string {
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

func (e CryptoApiException) String() string {
	return e.DetailMessage()
}

func (e CryptoApiException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3230001,
		Name: CryptoApiExceptionName,
		What: "Crypto API exception",
	}

	return json.Marshal(except)
}

func (e CryptoApiException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*CryptoApiException):
		callback(&e)
		return true
	case func(CryptoApiException):
		callback(e)
		return true
	default:
		return false
	}
}
