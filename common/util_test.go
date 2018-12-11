package common

import (
		"github.com/stretchr/testify/assert"
	"testing"
	)

func TestEmpty(t *testing.T) {
	type Action struct {
		ActionAccount uint64
		Data          []byte
	}
	type Transaction struct {
		Expiration         uint32
		NetUsageWords      uint
		MaxCPUUsageMs      uint8
		DelaySec           uint
		ContextFreeActions []*Action
	}

	action1 := &Action{
		ActionAccount: 9876543,
		Data:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}}
	action2 := &Action{
		ActionAccount: 987654321,
		Data:          []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}}
	test := &Transaction{Expiration: 100,
		NetUsageWords:      9,
		MaxCPUUsageMs:      199,
		DelaySec:           99999,
		ContextFreeActions: []*Action{action1, action2},
	}
	assert.Equal(t, false, Empty(test))
	test = nil
	assert.Equal(t, true, Empty(test))

	test1 := Transaction{}
	assert.Equal(t, true, Empty(test1))
}

type Data struct {
	a string
	b []string
	c []interface{}
	d map[string]interface{}
	e chan interface{}
	child ChildData
}

func (d *Data) IsEmpty() bool {
	return d.a == "" && len(d.b) == 0 && len(d.c) == 0 &&
		len(d.d) == 0 && len(d.e) == 0 && d.child.IsEmpty()
}

type ChildData struct {
	a string
	b []string
	c []interface{}
	d map[string]interface{}
	e chan interface{}
}

func (d *ChildData) IsEmpty() bool {
	return d.a == "" && len(d.b) == 0 && len(d.c) == 0 &&
		len(d.d) == 0 && len(d.e) == 0
}

var d = Data{}

func BenchmarkEmpty(b *testing.B) {
	for i:=0; i<b.N; i++ {
		if Empty(d) {}
	}
}

func BenchmarkEmpty2(b *testing.B) {
	for i:=0; i<b.N; i++ {
		if empty(d) {}
	}
}

func Test_empty(t *testing.T) {
	assert.Equal(t, true, empty(uint8(0)))
	assert.Equal(t, true, empty(uint32(0)))
	assert.Equal(t, true, empty(uint64(0)))
	assert.Equal(t, true, empty(int64(0)))
	assert.Equal(t, true, empty(int32(0)))
	assert.Equal(t, true, empty(int(0)))
	assert.Equal(t, true, empty(false))
}

