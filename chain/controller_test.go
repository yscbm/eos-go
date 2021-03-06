package chain

import (
	"fmt"
	"github.com/eosspark/eos-go/chain/types"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/ecc"
	"github.com/eosspark/eos-go/entity"
	"github.com/eosspark/eos-go/exception"
	"github.com/eosspark/eos-go/exception/try"
	"github.com/eosspark/eos-go/log"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestController_ProduceProcess(t *testing.T) {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			produceProcess()
		}
	}

}

func produceProcess() {
	con := GetControllerInstance()
	/*basetester := newBaseTester(con)
	basetester.ProduceBlock(common.Milliseconds(common.DefaultConfig.BlockIntervalMs), 0)*/

	signatureProviders := make(map[ecc.PublicKey]signatureProviderType)
	//con := GetControllerInstance()
	con.AbortBlock()
	now := common.Now()
	var base common.TimePoint
	if now > con.HeadBlockTime() {
		base = now
	} else {
		base = con.HeadBlockTime()
	}
	minTimeToNextBlock := common.DefaultConfig.BlockIntervalUs - (int64(base.TimeSinceEpoch()) % common.DefaultConfig.BlockIntervalUs)
	blockTime := base.AddUs(common.Microseconds(minTimeToNextBlock))

	if blockTime.Sub(now) < common.Microseconds(common.DefaultConfig.BlockIntervalUs/10) { // we must sleep for at least 50ms
		blockTime = blockTime.AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
	}
	con.StartBlock(types.NewBlockTimeStamp(blockTime), 0)
	unappliedTrxs := con.GetUnappliedTransactions()
	if len(unappliedTrxs) > 0 {
		for _, trx := range unappliedTrxs {
			trace := con.PushTransaction(trx, common.MaxTimePoint(), 0)
			if trace.Except != nil {
				log.Error("produce is failed isExhausted=true")
			} else {
				con.DropUnappliedTransaction(trx)
			}
		}
	}
	con.Pending.PendingBlockState = con.PendingBlockState()
	con.FinalizeBlock()
	pubKey, err := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	if err != nil {
		log.Error("produceLoop NewPublicKey is error :%s", err.Error())
	}
	priKey, err2 := ecc.NewPrivateKey("5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3")
	if err2 != nil {
		log.Error("produceLoop NewPrivateKey is error :%s", err.Error())
	}
	pbs := con.PendingBlockState()

	signatureProviders[pubKey] = makeKeySignatureProvider(priKey)
	a := signatureProviders[pbs.BlockSigningKey]
	con.SignBlock(func(d crypto.Sha256) ecc.Signature {
		return a(d)
	})

	con.CommitBlock(true)
}

type signatureProviderType = func(sha256 crypto.Sha256) ecc.Signature

func makeKeySignatureProvider(key *ecc.PrivateKey) signatureProviderType {
	signFunc := func(digest crypto.Sha256) ecc.Signature {
		sign, err := key.Sign(digest.Bytes())
		if err != nil {
			panic(err)
		}
		return sign
	}
	return signFunc
}

func TestPopBlock(t *testing.T) {
	con := GetControllerInstance()
	con.PopBlock()
}

func TestAbortBlock(t *testing.T) {
	con := GetControllerInstance()
	con.AbortBlock()
}

func CallBackApplayHandler(p *ApplyContext) {
	fmt.Println("SetApplyHandler CallBack")
}

func CallBackApplayHandler2(p *ApplyContext) {
	fmt.Println("SetApplyHandler CallBack2")
}
func TestSetApplyHandler(t *testing.T) {
	con := GetControllerInstance()

	//applyCon := ApplyContext{}
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.ScopeName(common.N("eosio")), common.ActionName(common.N("newaccount")), CallBackApplayHandler)
	con.SetApplayHandler(common.AccountName(common.N("eosio")), common.ScopeName(common.N("eosio")), common.ActionName(common.N("setcode")), CallBackApplayHandler2)

	handler1 := con.FindApplyHandler(common.AccountName(common.N("eosio")), common.ScopeName(common.N("eosio")), common.ActionName(common.N("newaccount")))
	handler1(nil)

	handler2 := con.FindApplyHandler(common.AccountName(common.N("eosio")), common.ScopeName(common.N("eosio")), common.ActionName(common.N("setcode")))
	handler2(nil)

	fmt.Println(len(con.ApplyHandlers))

}

var IrreversibleBlock chan types.BlockState = make(chan types.BlockState)

func TestController_CreateNativeAccount(t *testing.T) {
	//CreateNativeAccount(name common.AccountName,owner types.Authority,active types.Authority,isPrivileged bool)
	control := GetControllerInstance()
	name := common.AccountName(common.N("eos"))

	owner := types.Authority{}
	owner.Threshold = 2
	active := types.Authority{}
	active.Threshold = 1
	control.CreateNativeAccount(name, owner, active, false)
	fmt.Println(name)
	result := entity.AccountObject{}
	result.Name = name
	//control.DB.Find("name", result)

	fmt.Println("check account name:", strings.Compare(name.String(), "eos"))
	assert.Equal(t, "eos", name.String())
	control.Close()
}

func TestController_GetWasmInterface(t *testing.T) {
	control := GetControllerInstance()
	log.Info("%#v", control.WasmIf)
	//assert.Equal(t, nil, control.WasmIf)
}

func TestController_GetGlobalProperties(t *testing.T) {
	c := GetControllerInstance()
	result := c.GetGlobalProperties()
	gp := entity.GlobalPropertyObject{}
	gp.ID = common.IdType(1)
	err := c.DB.Find("ID", gp, &gp)
	if err != nil {
		assert.Error(t, err, gp)
	}
	assert.Equal(t, false, common.Empty(result)) //GlobalProperties not initialized
	assert.Equal(t, false, result == &gp)
	c.Close()
}

func TestController_GetDynamicGlobalProperties(t *testing.T) {
	c := GetControllerInstance()
	result := c.GetDynamicGlobalProperties()
	dgpo := entity.DynamicGlobalPropertyObject{}
	dgpo.ID = 1
	assert.Equal(t, &dgpo, result)
	//fmt.Println("*******", result)
}

func TestController_GetBlockIdForNum_NotFound(t *testing.T) {
	c := GetControllerInstance()
	try.Try(func() {
		c.GetBlockIdForNum(10)
	}).Catch(func(ex exception.Exception) { //TODO catch exception code
		assert.Equal(t, 3100002, int(ex.Code()))
	}).End()

}

func TestController_StartBlock(t *testing.T) {
	c := GetControllerInstance()
	w := types.NewBlockTimeStamp(common.Now())
	s := types.Irreversible
	c.StartBlock(w, uint16(s))
	c.Close()
}

func TestController_Close(t *testing.T) {
	c := GetControllerInstance()
	c.Close()
}

func TestController_UpdateProducersAuthority(t *testing.T) {
	c := GetControllerInstance()
	c.AbortBlock()
	now := common.Now()
	var base common.TimePoint
	if now > c.HeadBlockTime() {
		base = now
	} else {
		base = c.HeadBlockTime()
	}
	minTimeToNextBlock := common.DefaultConfig.BlockIntervalUs - (int64(base.TimeSinceEpoch()) % common.DefaultConfig.BlockIntervalUs)
	blockTime := base.AddUs(common.Microseconds(minTimeToNextBlock))

	if blockTime.Sub(now) < common.Microseconds(common.DefaultConfig.BlockIntervalUs/10) { // we must sleep for at least 50ms
		blockTime = blockTime.AddUs(common.Microseconds(common.DefaultConfig.BlockIntervalUs))
	}
	c.StartBlock(types.NewBlockTimeStamp(blockTime), 0)
	c.updateProducersAuthority()
	c.Close()
}

func BenchmarkGetControllerInstance(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	GetControllerInstance()
	b.StopTimer()
	/*i := c.LastIrreversibleBlockId()
	fmt.Println(i)*/
}
