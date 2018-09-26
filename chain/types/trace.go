package types

import "github.com/eosspark/eos-go/common"

type BaseActionTrace struct {
	Receipt       ActionReceipt
	Act           Action
	Elapsed       common.Microseconds
	CpuUsage      uint64
	Console       string
	TotalCpuUsage uint64                   /// total of inline_traces[x].cpu_usage + cpu_usage
	TrxId         common.TransactionIDType ///< the transaction that generated this action
}

type ActionTrace struct {
	BaseActionTrace
	InlineTraces []ActionTrace
}

type TransactionTrace struct {
	Id           common.TransactionIDType
	Receipt      *TransactionReceiptHeader
	Elapsed      common.Microseconds
	NetUsage     uint64
	Scheduled    bool
	ActionTraces []ActionTrace ///< disposable

	FailedDtrxTrace *TransactionTrace
	Except          error
}