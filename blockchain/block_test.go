package blockchain

import (
	"reflect"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1, "abcd")
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Errorf("createBlock() should return an instance of a block")
	}
}
