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

var ActorWhitelistExceptionName = reflect.TypeOf(ActorWhitelistException{}).Name()

type ActorWhitelistException struct {
	_WhitelistBlacklistException
	Elog log.Messages
}

func NewActorWhitelistException(parent _WhitelistBlacklistException, message log.Message) *ActorWhitelistException {
	return &ActorWhitelistException{parent, log.Messages{message}}
}

func (e ActorWhitelistException) Code() int64 {
	return 3130001
}

func (e ActorWhitelistException) Name() string {
	return ActorWhitelistExceptionName
}

func (e ActorWhitelistException) What() string {
	return "Authorizing actor of transaction is not on the whitelist"
}

func (e *ActorWhitelistException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ActorWhitelistException) GetLog() log.Messages {
	return e.Elog
}

func (e ActorWhitelistException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ActorWhitelistException) DetailMessage() string {
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

func (e ActorWhitelistException) String() string {
	return e.DetailMessage()
}

func (e ActorWhitelistException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3130001,
		Name: ActorWhitelistExceptionName,
		What: "Authorizing actor of transaction is not on the whitelist",
	}

	return json.Marshal(except)
}

func (e ActorWhitelistException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ActorWhitelistException):
		callback(&e)
		return true
	case func(ActorWhitelistException):
		callback(e)
		return true
	default:
		return false
	}
}