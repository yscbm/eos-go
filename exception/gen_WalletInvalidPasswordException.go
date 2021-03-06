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

var WalletInvalidPasswordExceptionName = reflect.TypeOf(WalletInvalidPasswordException{}).Name()

type WalletInvalidPasswordException struct {
	_WalletException
	Elog log.Messages
}

func NewWalletInvalidPasswordException(parent _WalletException, message log.Message) *WalletInvalidPasswordException {
	return &WalletInvalidPasswordException{parent, log.Messages{message}}
}

func (e WalletInvalidPasswordException) Code() int64 {
	return 3120005
}

func (e WalletInvalidPasswordException) Name() string {
	return WalletInvalidPasswordExceptionName
}

func (e WalletInvalidPasswordException) What() string {
	return "Invalid wallet password"
}

func (e *WalletInvalidPasswordException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e WalletInvalidPasswordException) GetLog() log.Messages {
	return e.Elog
}

func (e WalletInvalidPasswordException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e WalletInvalidPasswordException) DetailMessage() string {
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

func (e WalletInvalidPasswordException) String() string {
	return e.DetailMessage()
}

func (e WalletInvalidPasswordException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3120005,
		Name: WalletInvalidPasswordExceptionName,
		What: "Invalid wallet password",
	}

	return json.Marshal(except)
}

func (e WalletInvalidPasswordException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*WalletInvalidPasswordException):
		callback(&e)
		return true
	case func(WalletInvalidPasswordException):
		callback(e)
		return true
	default:
		return false
	}
}
