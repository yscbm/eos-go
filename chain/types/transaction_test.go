package types

import (
	"encoding/json"
	"fmt"
	"github.com/eosspark/eos-go/common"
	"github.com/eosspark/eos-go/crypto"
	"github.com/eosspark/eos-go/crypto/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransaction(t *testing.T) {

	data := []byte{0x1, 0x0, 0x20, 0x60, 0x67, 0x3d, 0x9f, 0x33, 0xa8, 0x55, 0x8e, 0x1b, 0xd5, 0x42, 0x96, 0x79, 0xbc, 0xee, 0x2a, 0x51, 0x26, 0xa1, 0x99, 0x9a, 0x38, 0x73, 0x81, 0x6e, 0xa3, 0x6d, 0xe4, 0xdd, 0x44, 0xae, 0xbb, 0x39, 0x4f, 0x15, 0xfa, 0xd0, 0x6f, 0xdb, 0x6a, 0x6, 0xf8, 0xab, 0x69, 0x53, 0x9c, 0x6e, 0xcd, 0x8d, 0xd, 0xda, 0x32, 0x4f, 0x64, 0x91, 0x3a, 0xbb, 0x13, 0xc2, 0x7f, 0x84, 0x24, 0x94, 0xf4, 0x0, 0x0, 0x98, 0x1, 0x50, 0xeb, 0xc3, 0x5b, 0x4, 0x0, 0x9e, 0xd5, 0x72, 0xe4, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x5c, 0x5, 0xa3, 0xe1, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0}
	packedTrx := PackedTransaction{}
	err := rlp.DecodeBytes(data, &packedTrx)
	if err != nil {
		fmt.Println(err)
	}
	id := packedTrx.ID() //e97f9f1e4aaafe1b92feded9bdd140247465de773154bcccab86986e1806fa33
	fmt.Println(id)
	trx := packedTrx.GetTransaction()
	re, _ := json.Marshal(trx)
	fmt.Println("Trx:  ", string(re))

	fmt.Println(trx.ID())

	result, err := rlp.EncodeToBytes(packedTrx)
	assert.NoError(t, nil, err)
	assert.Equal(t, data, result)

}

func TestSignedBlock(t *testing.T) {
	data := []byte{0x66, 0x4f, 0xad, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0xec, 0x3b, 0x3f, 0x3a, 0xe, 0x67, 0x96, 0xc6, 0xc5, 0x0, 0x43, 0xd1, 0x47, 0xac, 0xe2, 0x31, 0x93, 0xa6, 0x6e, 0x4b, 0x88, 0x55, 0x7b, 0x81, 0x50, 0x93, 0xa5, 0xf6, 0x3a, 0x1c, 0x77, 0x50, 0xf6, 0x33, 0x7f, 0x9e, 0x46, 0x91, 0xa4, 0xca, 0x2e, 0x32, 0x55, 0x48, 0x7, 0x2c, 0xc2, 0x82, 0x7a, 0xae, 0x7f, 0xba, 0x5f, 0xaa, 0x17, 0xb0, 0x38, 0xd5, 0xf9, 0xb, 0x44, 0x48, 0xb9, 0x47, 0x64, 0x11, 0x34, 0x75, 0x7b, 0xc0, 0x15, 0x27, 0xeb, 0x52, 0x14, 0x3b, 0x4d, 0x61, 0xf9, 0xd6, 0x49, 0x24, 0xf2, 0x4b, 0x7f, 0x19, 0x20, 0x7c, 0x2f, 0x46, 0xc2, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x42, 0x22, 0x7c, 0xd4, 0x29, 0x7, 0x34, 0x3d, 0x90, 0x3d, 0xd2, 0x6a, 0x10, 0xbb, 0x41, 0x93, 0x57, 0x4c, 0xf, 0xad, 0xec, 0x90, 0x27, 0x66, 0xa9, 0xe5, 0x4f, 0x4c, 0xde, 0xfd, 0x3, 0xc3, 0x25, 0x3f, 0x7d, 0x66, 0x77, 0x1e, 0x14, 0xf3, 0x5f, 0x9c, 0xd9, 0xc, 0xcf, 0xe9, 0x6a, 0x5b, 0x3d, 0xfa, 0x80, 0x8f, 0xf, 0x6c, 0xea, 0xf7, 0x9b, 0xdf, 0x2f, 0x74, 0xab, 0x6f, 0x47, 0x9e, 0x1, 0x0, 0x2d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x60, 0x67, 0x3d, 0x9f, 0x33, 0xa8, 0x55, 0x8e, 0x1b, 0xd5, 0x42, 0x96, 0x79, 0xbc, 0xee, 0x2a, 0x51, 0x26, 0xa1, 0x99, 0x9a, 0x38, 0x73, 0x81, 0x6e, 0xa3, 0x6d, 0xe4, 0xdd, 0x44, 0xae, 0xbb, 0x39, 0x4f, 0x15, 0xfa, 0xd0, 0x6f, 0xdb, 0x6a, 0x6, 0xf8, 0xab, 0x69, 0x53, 0x9c, 0x6e, 0xcd, 0x8d, 0xd, 0xda, 0x32, 0x4f, 0x64, 0x91, 0x3a, 0xbb, 0x13, 0xc2, 0x7f, 0x84, 0x24, 0x94, 0xf4, 0x0, 0x0, 0x98, 0x1, 0x50, 0xeb, 0xc3, 0x5b, 0x4, 0x0, 0x9e, 0xd5, 0x72, 0xe4, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x5c, 0x5, 0xa3, 0xe1, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
	signedBlock := SignedBlock{}
	_ = rlp.DecodeBytes(data, &signedBlock)
	data, err := json.Marshal(signedBlock)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Receive P2PMessag ", string(data))

	//result ,err := rlp.EncodeToBytes(signedBlock)
	//assert.NoError(t,nil,err)
	//assert.Equal(t,data,result)

}

func TestReceiveSignedBlock(t *testing.T) {
	data1 := []byte{0x9f, 0x1, 0x0, 0x0, 0x7, 0x66, 0x4f, 0xad, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0xec, 0x3b, 0x3f, 0x3a, 0xe, 0x67, 0x96, 0xc6, 0xc5, 0x0, 0x43, 0xd1, 0x47, 0xac, 0xe2, 0x31, 0x93, 0xa6, 0x6e, 0x4b, 0x88, 0x55, 0x7b, 0x81, 0x50, 0x93, 0xa5, 0xf6, 0x3a, 0x1c, 0x77, 0x50, 0xf6, 0x33, 0x7f, 0x9e, 0x46, 0x91, 0xa4, 0xca, 0x2e, 0x32, 0x55, 0x48, 0x7, 0x2c, 0xc2, 0x82, 0x7a, 0xae, 0x7f, 0xba, 0x5f, 0xaa, 0x17, 0xb0, 0x38, 0xd5, 0xf9, 0xb, 0x44, 0x48, 0xb9, 0x47, 0x64, 0x11, 0x34, 0x75, 0x7b, 0xc0, 0x15, 0x27, 0xeb, 0x52, 0x14, 0x3b, 0x4d, 0x61, 0xf9, 0xd6, 0x49, 0x24, 0xf2, 0x4b, 0x7f, 0x19, 0x20, 0x7c, 0x2f, 0x46, 0xc2, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x42, 0x22, 0x7c, 0xd4, 0x29, 0x7, 0x34, 0x3d, 0x90, 0x3d, 0xd2, 0x6a, 0x10, 0xbb, 0x41, 0x93, 0x57, 0x4c, 0xf, 0xad, 0xec, 0x90, 0x27, 0x66, 0xa9, 0xe5, 0x4f, 0x4c, 0xde, 0xfd, 0x3, 0xc3, 0x25, 0x3f, 0x7d, 0x66, 0x77, 0x1e, 0x14, 0xf3, 0x5f, 0x9c, 0xd9, 0xc, 0xcf, 0xe9, 0x6a, 0x5b, 0x3d, 0xfa, 0x80, 0x8f, 0xf, 0x6c, 0xea, 0xf7, 0x9b, 0xdf, 0x2f, 0x74, 0xab, 0x6f, 0x47, 0x9e, 0x1, 0x0, 0x2d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x60, 0x67, 0x3d, 0x9f, 0x33, 0xa8, 0x55, 0x8e, 0x1b, 0xd5, 0x42, 0x96, 0x79, 0xbc, 0xee, 0x2a, 0x51, 0x26, 0xa1, 0x99, 0x9a, 0x38, 0x73, 0x81, 0x6e, 0xa3, 0x6d, 0xe4, 0xdd, 0x44, 0xae, 0xbb, 0x39, 0x4f, 0x15, 0xfa, 0xd0, 0x6f, 0xdb, 0x6a, 0x6, 0xf8, 0xab, 0x69, 0x53, 0x9c, 0x6e, 0xcd, 0x8d, 0xd, 0xda, 0x32, 0x4f, 0x64, 0x91, 0x3a, 0xbb, 0x13, 0xc2, 0x7f, 0x84, 0x24, 0x94, 0xf4, 0x0, 0x0, 0x98, 0x1, 0x50, 0xeb, 0xc3, 0x5b, 0x4, 0x0, 0x9e, 0xd5, 0x72, 0xe4, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x5c, 0x5, 0xa3, 0xe1, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
	data2 := []byte{0x9f, 0x1, 0x0, 0x0, 0x7, 0x76, 0x4f, 0xad, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x15, 0xfc, 0xc0, 0xb8, 0x41, 0x5, 0x50, 0x37, 0xe7, 0xa7, 0xf2, 0xe3, 0x99, 0xed, 0xad, 0x14, 0x9a, 0x58, 0xa1, 0xb, 0x38, 0xfc, 0x75, 0xae, 0x7c, 0x19, 0xc7, 0xcc, 0xd7, 0x9, 0x5a, 0x6c, 0x4e, 0x90, 0x6d, 0xe1, 0x2f, 0x21, 0x87, 0x3, 0x7, 0x45, 0x69, 0xfc, 0xb6, 0xb3, 0xe4, 0xc1, 0x89, 0x64, 0x3d, 0x5e, 0xce, 0x6c, 0x13, 0x41, 0xb5, 0x5b, 0xe8, 0x81, 0xdd, 0x9f, 0xf3, 0x7e, 0xdc, 0xf6, 0x77, 0x79, 0x16, 0x1b, 0xe4, 0xc, 0x78, 0x93, 0x72, 0xb1, 0x63, 0xac, 0xe6, 0x4, 0xb7, 0x9d, 0x30, 0xa4, 0x27, 0xff, 0x91, 0x1c, 0x46, 0x63, 0xa, 0xe9, 0xd0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0xd, 0x90, 0x27, 0xae, 0xb, 0x4f, 0x8, 0x2e, 0x72, 0x26, 0x49, 0xe2, 0xb0, 0xcd, 0xc4, 0x64, 0x1e, 0x3d, 0xb0, 0xd9, 0x63, 0x7, 0xce, 0x36, 0x49, 0x71, 0xf6, 0xcb, 0x35, 0x3b, 0xe0, 0xe6, 0x31, 0x48, 0x9f, 0x46, 0x7c, 0xaa, 0x12, 0x7b, 0xf2, 0x22, 0x72, 0xf9, 0x4f, 0x9a, 0x2b, 0x8c, 0x14, 0xeb, 0x99, 0xb0, 0x35, 0x26, 0xb0, 0x23, 0x4e, 0x99, 0xf5, 0xa8, 0xc3, 0x67, 0x3e, 0xa9, 0x1, 0x0, 0xd2, 0x0, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x37, 0x3, 0x2f, 0x8f, 0x54, 0x1c, 0xd3, 0x36, 0x82, 0x9a, 0x88, 0xa7, 0x63, 0x20, 0xe7, 0xb3, 0xea, 0xee, 0xa7, 0x50, 0x89, 0x26, 0x19, 0x9f, 0x91, 0x30, 0x3f, 0x1d, 0x5b, 0xac, 0xc1, 0xf6, 0x11, 0x7e, 0x7d, 0x9e, 0x48, 0x95, 0xba, 0xdf, 0xe3, 0x8, 0x99, 0x1e, 0xcc, 0xf1, 0xb3, 0xfa, 0xb9, 0x32, 0xd7, 0x78, 0x52, 0x73, 0x11, 0x6d, 0xea, 0x79, 0x7c, 0x8, 0x68, 0x2e, 0x79, 0xbe, 0x0, 0x0, 0x98, 0x1, 0x58, 0xeb, 0xc3, 0x5b, 0x14, 0x0, 0x4c, 0xc1, 0x4b, 0x27, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1a, 0xa3, 0x6a, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0x72, 0x47, 0xd0, 0x91, 0xa5, 0xb0, 0x20, 0xe8, 0x74, 0x50, 0xcf, 0xf9, 0x1, 0x4e, 0x38, 0xc0, 0xf4, 0x16, 0x8b, 0xca, 0xd5, 0x99, 0xb4, 0x5d, 0x1d, 0xfa, 0xa7, 0x19, 0x37, 0xe6, 0x35, 0x16, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
	fmt.Println(data1)
	fmt.Println(data2)
}

func TestTransactionID(t *testing.T) {
	data := []byte{0xa6, 0xe4, 0xb7, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x38, 0x1e, 0xef, 0x63, 0x5b, 0x5c, 0x2, 0xe6, 0xdf, 0xb1, 0xab, 0xa5, 0x78, 0xc3, 0x23, 0x60, 0xb3, 0x51, 0x7c, 0xae, 0xd7, 0xd0, 0x47, 0xbb, 0x86, 0x4c, 0x11, 0xe, 0xcd, 0xc9, 0x8f, 0xe6, 0x7f, 0x35, 0x88, 0x97, 0x60, 0x58, 0x62, 0xd4, 0xe9, 0xa4, 0x13, 0x63, 0x28, 0x2f, 0x5c, 0xe4, 0x36, 0xde, 0x7, 0x9a, 0x8b, 0xcd, 0x90, 0x1, 0xa1, 0x5a, 0xc3, 0x86, 0x4d, 0x85, 0x47, 0xa5, 0x32, 0x7f, 0x4, 0xca, 0xfc, 0x37, 0x43, 0x7c, 0x2c, 0x1f, 0xde, 0x6c, 0x3b, 0x7a, 0x6e, 0x1f, 0x46, 0x14, 0xd2, 0x21, 0x65, 0xa9, 0x49, 0xe4, 0x16, 0xcd, 0x65, 0xb8, 0xd9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x34, 0x22, 0x5d, 0x1f, 0xbf, 0x89, 0xf, 0xa9, 0xb5, 0xf8, 0xd1, 0x9a, 0xbc, 0xc7, 0x31, 0x54, 0xf9, 0xac, 0x30, 0xae, 0xc8, 0xf4, 0x1e, 0x48, 0xf7, 0xf2, 0xc2, 0xc9, 0xd5, 0xe2, 0x93, 0x92, 0x74, 0x58, 0xf1, 0xb, 0xb0, 0xf4, 0x68, 0x70, 0x3b, 0x70, 0x8c, 0x3d, 0x5f, 0x60, 0x44, 0x27, 0x6e, 0xde, 0xf0, 0xb0, 0x19, 0xe6, 0x6d, 0xea, 0xf4, 0xfe, 0x74, 0x95, 0xda, 0xf6, 0x10, 0xd0, 0x1, 0x0, 0x1d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x44, 0x9c, 0x14, 0x28, 0x79, 0x8c, 0x5d, 0x4d, 0xfa, 0xcc, 0xfc, 0xf1, 0xdf, 0x4, 0x7f, 0x6b, 0xef, 0xd0, 0x99, 0x92, 0xbe, 0x38, 0x8d, 0xed, 0x3b, 0x74, 0xfe, 0xae, 0xe0, 0xf, 0x4e, 0x1b, 0x1d, 0x12, 0x8e, 0xc1, 0xa5, 0x18, 0x14, 0xe4, 0x16, 0xc9, 0xf6, 0x16, 0xc6, 0x13, 0xac, 0x11, 0x90, 0xe4, 0xfb, 0x40, 0xda, 0xef, 0x27, 0xd0, 0x9c, 0xcf, 0x4e, 0x66, 0xdd, 0x83, 0x54, 0xe7, 0x0, 0x0, 0x98, 0x1, 0xf0, 0x35, 0xc9, 0x5b, 0x37, 0x0, 0x38, 0xa, 0xd3, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0xa6, 0x82, 0x34, 0x3, 0xea, 0x30, 0x55, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
	signedBlock := SignedBlock{}
	_ = rlp.DecodeBytes(data, &signedBlock)
	data, err := json.Marshal(signedBlock)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Receive P2PMessag ", string(data))

	transactions := signedBlock.Transactions
	//fmt.Println(transactions[0].TransactionReceiptHeader)

	for _, TrxReceipt := range transactions {

		if TrxReceipt.Trx.TransactionID == common.TransactionIdType(*crypto.NewSha256Nil()) {

			packedTrx := TrxReceipt.Trx.PackedTransaction

			enc, _ := rlp.EncodeToBytes(packedTrx)
			fmt.Printf("%#v\n", enc)
			data, err := json.Marshal(packedTrx.PackedTrx)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("encode result : packedTrx:   ", string(data))

			//
			//fmt.Printf("1           %#v\n",packedTrx.PackedTrx)

			trx := packedTrx.GetTransaction()
			data, err = json.Marshal(trx)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("trx receive  ", string(data))

			//trx.ID()

			fmt.Println(trx.ID())

			//signedTrx := packedTrx.GetSignedTransaction()
			//fmt.Println(signedTrx.Transaction)
			//data, err = json.Marshal(signedTrx.Transaction)
			//if err != nil {
			//	fmt.Println(err)
			//}
			//fmt.Println("trx signed  ", string(data))
			//
			//newPackedTrx := types.NewPackedTransactionBySignedTrx(signedTrx,common.CompressionNone)
			//
			//fmt.Printf("2            %#v\n",newPackedTrx.PackedTrx)
			//
			//fmt.Println( bytes.Compare(newPackedTrx.PackedTrx,packedTrx.PackedTrx))
			//
			//fmt.Println(signedTrx.Transaction)

		}
	}

}

//receive data:  []byte{0x9f, 0x1, 0x0, 0x0, 0x7, 0xa6, 0xe4, 0xb7, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x38, 0x1e, 0xef, 0x63, 0x5b, 0x5c, 0x2, 0xe6, 0xdf, 0xb1, 0xab, 0xa5, 0x78, 0xc3, 0x23, 0x60, 0xb3, 0x51, 0x7c, 0xae, 0xd7, 0xd0, 0x47, 0xbb, 0x86, 0x4c, 0x11, 0xe, 0xcd, 0xc9, 0x8f, 0xe6, 0x7f, 0x35, 0x88, 0x97, 0x60, 0x58, 0x62, 0xd4, 0xe9, 0xa4, 0x13, 0x63, 0x28, 0x2f, 0x5c, 0xe4, 0x36, 0xde, 0x7, 0x9a, 0x8b, 0xcd, 0x90, 0x1, 0xa1, 0x5a, 0xc3, 0x86, 0x4d, 0x85, 0x47, 0xa5, 0x32, 0x7f, 0x4, 0xca, 0xfc, 0x37, 0x43, 0x7c, 0x2c, 0x1f, 0xde, 0x6c, 0x3b, 0x7a, 0x6e, 0x1f, 0x46, 0x14, 0xd2, 0x21, 0x65, 0xa9, 0x49, 0xe4, 0x16, 0xcd, 0x65, 0xb8, 0xd9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x34, 0x22, 0x5d, 0x1f, 0xbf, 0x89, 0xf, 0xa9, 0xb5, 0xf8, 0xd1, 0x9a, 0xbc, 0xc7, 0x31, 0x54, 0xf9, 0xac, 0x30, 0xae, 0xc8, 0xf4, 0x1e, 0x48, 0xf7, 0xf2, 0xc2, 0xc9, 0xd5, 0xe2, 0x93, 0x92, 0x74, 0x58, 0xf1, 0xb, 0xb0, 0xf4, 0x68, 0x70, 0x3b, 0x70, 0x8c, 0x3d, 0x5f, 0x60, 0x44, 0x27, 0x6e, 0xde, 0xf0, 0xb0, 0x19, 0xe6, 0x6d, 0xea, 0xf4, 0xfe, 0x74, 0x95, 0xda, 0xf6, 0x10, 0xd0, 0x1, 0x0, 0x1d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x44, 0x9c, 0x14, 0x28, 0x79, 0x8c, 0x5d, 0x4d, 0xfa, 0xcc, 0xfc, 0xf1, 0xdf, 0x4, 0x7f, 0x6b, 0xef, 0xd0, 0x99, 0x92, 0xbe, 0x38, 0x8d, 0xed, 0x3b, 0x74, 0xfe, 0xae, 0xe0, 0xf, 0x4e, 0x1b, 0x1d, 0x12, 0x8e, 0xc1, 0xa5, 0x18, 0x14, 0xe4, 0x16, 0xc9, 0xf6, 0x16, 0xc6, 0x13, 0xac, 0x11, 0x90, 0xe4, 0xfb, 0x40, 0xda, 0xef, 0x27, 0xd0, 0x9c, 0xcf, 0x4e, 0x66, 0xdd, 0x83, 0x54, 0xe7, 0x0, 0x0, 0x98, 0x1, 0xf0, 0x35, 0xc9, 0x5b, 0x37, 0x0, 0x38, 0xa, 0xd3, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0xa6, 0x82, 0x34, 0x3, 0xea, 0x30, 0x55, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
//signed Block Num: 57
//---------------*************----------------
//Receive signedBlock:    {"timestamp":"2018-10-19T01:39:31.000","producer":"eosio","confirmed":0,"previous":"000000381eef635b5c02e6dfb1aba578c32360b3517caed7d047bb864c110ecd","transaction_mroot":"c98fe67f358897605862d4e9a41363282f5ce436de079a8bcd9001a15ac3864d","action_mroot":"8547a5327f04cafc37437c2c1fde6c3b7a6e1f4614d22165a949e416cd65b8d9","schedule_version":0,"new_producers":null,"header_extensions":[],"producer_signature":"SIG_K1_K25NdAjgDn4niYaJZFTD4dESHqLMVCEdqqazQhAyCkmre1LtW6VfZ2tT5vKFutjbfXC8hYba2iTZcoz1M5k4agrSDZSUYh","transactions":[{"status":"executed","cpu_usage_us":285,"net_usage_words":25,"trx":[{"signatures":["SIG_K1_KdivyH4YyTLkPcqW7GDCovLkAixmMuUgWHnDjySkuA3iRPFXS3hr3BMUtcCVamRgrgYVSDF9RXjGa6Dfd6V9SeRa7tM4Pi"],"compression":"none","packed_context_free_data":"","packed_trx":"f035c95b3700380ad3d100000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232660000000000ea305500a6823403ea305501000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000000","UnpackedTrx":null},"0000000000000000000000000000000000000000000000000000000000000000"]}],"block_extensions":[]}
//encode result: []byte{0xa6, 0xe4, 0xb7, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x38, 0x1e, 0xef, 0x63, 0x5b, 0x5c, 0x2, 0xe6, 0xdf, 0xb1, 0xab, 0xa5, 0x78, 0xc3, 0x23, 0x60, 0xb3, 0x51, 0x7c, 0xae, 0xd7, 0xd0, 0x47, 0xbb, 0x86, 0x4c, 0x11, 0xe, 0xcd, 0xc9, 0x8f, 0xe6, 0x7f, 0x35, 0x88, 0x97, 0x60, 0x58, 0x62, 0xd4, 0xe9, 0xa4, 0x13, 0x63, 0x28, 0x2f, 0x5c, 0xe4, 0x36, 0xde, 0x7, 0x9a, 0x8b, 0xcd, 0x90, 0x1, 0xa1, 0x5a, 0xc3, 0x86, 0x4d, 0x85, 0x47, 0xa5, 0x32, 0x7f, 0x4, 0xca, 0xfc, 0x37, 0x43, 0x7c, 0x2c, 0x1f, 0xde, 0x6c, 0x3b, 0x7a, 0x6e, 0x1f, 0x46, 0x14, 0xd2, 0x21, 0x65, 0xa9, 0x49, 0xe4, 0x16, 0xcd, 0x65, 0xb8, 0xd9, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1f, 0x34, 0x22, 0x5d, 0x1f, 0xbf, 0x89, 0xf, 0xa9, 0xb5, 0xf8, 0xd1, 0x9a, 0xbc, 0xc7, 0x31, 0x54, 0xf9, 0xac, 0x30, 0xae, 0xc8, 0xf4, 0x1e, 0x48, 0xf7, 0xf2, 0xc2, 0xc9, 0xd5, 0xe2, 0x93, 0x92, 0x74, 0x58, 0xf1, 0xb, 0xb0, 0xf4, 0x68, 0x70, 0x3b, 0x70, 0x8c, 0x3d, 0x5f, 0x60, 0x44, 0x27, 0x6e, 0xde, 0xf0, 0xb0, 0x19, 0xe6, 0x6d, 0xea, 0xf4, 0xfe, 0x74, 0x95, 0xda, 0xf6, 0x10, 0xd0, 0x1, 0x0, 0x1d, 0x1, 0x0, 0x0, 0x19, 0x1, 0x1, 0x0, 0x20, 0x44, 0x9c, 0x14, 0x28, 0x79, 0x8c, 0x5d, 0x4d, 0xfa, 0xcc, 0xfc, 0xf1, 0xdf, 0x4, 0x7f, 0x6b, 0xef, 0xd0, 0x99, 0x92, 0xbe, 0x38, 0x8d, 0xed, 0x3b, 0x74, 0xfe, 0xae, 0xe0, 0xf, 0x4e, 0x1b, 0x1d, 0x12, 0x8e, 0xc1, 0xa5, 0x18, 0x14, 0xe4, 0x16, 0xc9, 0xf6, 0x16, 0xc6, 0x13, 0xac, 0x11, 0x90, 0xe4, 0xfb, 0x40, 0xda, 0xef, 0x27, 0xd0, 0x9c, 0xcf, 0x4e, 0x66, 0xdd, 0x83, 0x54, 0xe7, 0x0, 0x0, 0x98, 0x1, 0xf0, 0x35, 0xc9, 0x5b, 0x37, 0x0, 0x38, 0xa, 0xd3, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0xa6, 0x82, 0x34, 0x3, 0xea, 0x30, 0x55, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0}
//[]byte{0x1, 0x0, 0x20, 0x44, 0x9c, 0x14, 0x28, 0x79, 0x8c, 0x5d, 0x4d, 0xfa, 0xcc, 0xfc, 0xf1, 0xdf, 0x4, 0x7f, 0x6b, 0xef, 0xd0, 0x99, 0x92, 0xbe, 0x38, 0x8d, 0xed, 0x3b, 0x74, 0xfe, 0xae, 0xe0, 0xf, 0x4e, 0x1b, 0x1d, 0x12, 0x8e, 0xc1, 0xa5, 0x18, 0x14, 0xe4, 0x16, 0xc9, 0xf6, 0x16, 0xc6, 0x13, 0xac, 0x11, 0x90, 0xe4, 0xfb, 0x40, 0xda, 0xef, 0x27, 0xd0, 0x9c, 0xcf, 0x4e, 0x66, 0xdd, 0x83, 0x54, 0xe7, 0x0, 0x0, 0x98, 0x1, 0xf0, 0x35, 0xc9, 0x5b, 0x37, 0x0, 0x38, 0xa, 0xd3, 0xd1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x40, 0x9e, 0x9a, 0x22, 0x64, 0xb8, 0x9a, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0xa8, 0xed, 0x32, 0x32, 0x66, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0xa6, 0x82, 0x34, 0x3, 0xea, 0x30, 0x55, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x2, 0xc0, 0xde, 0xd2, 0xbc, 0x1f, 0x13, 0x5, 0xfb, 0xf, 0xaa, 0xc5, 0xe6, 0xc0, 0x3e, 0xe3, 0xa1, 0x92, 0x42, 0x34, 0x98, 0x54, 0x27, 0xb6, 0x16, 0x7c, 0xa5, 0x69, 0xd1, 0x3d, 0xf4, 0x35, 0xcf, 0x1, 0x0, 0x0, 0x0, 0x0}
//encode result : packedTrx:    "f035c95b3700380ad3d100000000010000000000ea305500409e9a2264b89a010000000000ea305500000000a8ed3232660000000000ea305500a6823403ea305501000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000000"
//trx receive   {"expiration":1539913200,"ref_block_num":55,"ref_block_prefix":3520268856,"max_net_usage_words":0,"max_cpu_usage_ms":0,"delay_sec":0,"context_free_actions":[],"actions":[{"account":"eosio","name":"newaccount","authorization":[{"actor":"eosio","permission":"active"}],"data":"0000000000ea305500a6823403ea305501000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf0100000001000000010002c0ded2bc1f1305fb0faac5e6c03ee3a1924234985427b6167ca569d13df435cf01000000"}],"transaction_extensions":[]}
//cbb52b1177b73e47f0282d8fb706158acadf78f872f0bb93e94acad159dec50a
//---------------*************----------------
//receive data:  []byte{0xb9, 0x0, 0x0, 0x0, 0x7, 0xa7, 0xe4, 0xb7, 0x46, 0x0, 0x0, 0x0, 0x0, 0x0, 0xea, 0x30, 0x55, 0x0, 0x0, 0x0, 0x0, 0x0, 0x39, 0x98, 0x9c, 0x66, 0x39, 0x71, 0x71, 0x73, 0xfb, 0x68, 0xf7, 0xec, 0xca, 0x7a, 0x24, 0xb3, 0xd, 0xc4, 0xca, 0x81, 0x17, 0x2a, 0x37, 0x39, 0x51, 0x2c, 0x26, 0xa2, 0xa1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xad, 0x5e, 0x40, 0x20, 0x86, 0x89, 0xc2, 0xae, 0x4e, 0x96, 0x9d, 0x33, 0x7d, 0xdd, 0x49, 0xac, 0xdb, 0xe6, 0xe7, 0xaa, 0xa0, 0xf0, 0x8, 0x2d, 0xb4, 0xfd, 0x80, 0x1c, 0x4a, 0xaf, 0x45, 0xa2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x7f, 0x48, 0x9d, 0xf0, 0x2a, 0x33, 0xf5, 0x2b, 0xdc, 0xe1, 0x35, 0x66, 0xee, 0x4, 0xb7, 0xf8, 0xb1, 0x3d, 0x4e, 0x39, 0x53, 0xd2, 0xa2, 0xed, 0xd5, 0x29, 0x2e, 0x5, 0x3a, 0xb3, 0x13, 0xfb, 0xf, 0x83, 0x8, 0xa4, 0x28, 0xce, 0xe3, 0x29, 0xcd, 0x7c, 0x91, 0xf7, 0x4a, 0x35, 0x50, 0x6e, 0x89, 0x6, 0x5c, 0x5a, 0x84, 0xe4, 0x9b, 0x90, 0xf7, 0x99, 0xd9, 0xc, 0xb8, 0xf2, 0x5f, 0xd6, 0x0, 0x0}
//signed Block Num: 58
