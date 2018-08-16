package new_vm

import (
	"testing"

	"fmt"

	"os"

	"github.com/iost-official/Go-IOS-Protocol/account"
	"github.com/iost-official/Go-IOS-Protocol/common"
	"github.com/iost-official/Go-IOS-Protocol/core/contract"
	"github.com/iost-official/Go-IOS-Protocol/core/new_block"
	"github.com/iost-official/Go-IOS-Protocol/core/new_tx"
	"github.com/iost-official/Go-IOS-Protocol/db"
	"github.com/iost-official/Go-IOS-Protocol/new_vm/database"
)

var testID = []string{
	"IOST4wQ6HPkSrtDRYi2TGkyMJZAB3em26fx79qR3UJC7fcxpL87wTn", "EhNiaU4DzUmjCrvynV3gaUeuj2VjB1v2DCmbGD5U2nSE",
	"IOST558jUpQvBD7F3WTKpnDAWg6HwKrfFiZ7AqhPFf4QSrmjdmBGeY", "8dJ9YKovJ5E7hkebAQaScaG1BA8snRUHPUbVcArcTVq6",
	"IOST7ZNDWeh8pHytAZdpgvp7vMpjZSSe5mUUKxDm6AXPsbdgDMAYhs", "7CnwT7BXkEFAVx6QZqC7gkDhQwbvC3d2CkMZvXHZdDMN",
	"IOST54ETA3q5eC8jAoEpfRAToiuc6Fjs5oqEahzghWkmEYs9S9CMKd", "Htarc5Sp4trjqY4WrTLtZ85CF6qx87v7CRwtV4RRGnbF",
	"IOST7GmPn8xC1RESMRS6a62RmBcCdwKbKvk2ZpxZpcXdUPoJdapnnh", "Bk8bAyG4VLBcrsoRErPuQGhwCy4C1VxfKE4jjX9oLhv",
	"IOST7ZGQL4k85v4wAxWngmow7JcX4QFQ4mtLNjgvRrEnEuCkGSBEHN", "546aCDG9igGgZqVZeybajaorP5ZeF9ghLu2oLncXk3d6",
	"IOST59uMX3Y4ab5dcq8p1wMXodANccJcj2efbcDThtkw6egvcni5L9", "DXNYRwG7dRFkbWzMNEbKfBhuS8Yn51x9J6XuTdNwB11M",
	"IOST8mFxe4kq9XciDtURFZJ8E76B8UssBgRVFA5gZN9HF5kLUVZ1BB", "AG8uECmAwFis8uxTdWqcgGD9tGDwoP6CxqhkhpuCdSeC",
	"IOST7uqa5UQPVT9ongTv6KmqDYKdVYSx4DV2reui4nuC5mm5vBt3D9", "GJt5WSSv5WZi1axd3qkb1vLEfxCEgKGupcXf45b5tERU",
	"IOST6wYBsLZmzJv22FmHAYBBsTzmV1p1mtHQwkTK9AjCH9Tg5Le4i4", "7U3uwEeGc2TF3Xde2oT66eTx1Uw15qRqYuTnMd3NNjai",
}

var systemContract = &contract.Contract{
	ID:   "iost.system",
	Code: "codes",
	Info: &contract.Info{
		Lang:        "native",
		VersionCode: "1.0.0",
		Abis: []*contract.ABI{
			{
				Name:     "Transfer",
				Payment:  0,
				GasPrice: int64(1000),
				Limit:    contract.NewCost(100, 100, 100),
				Args:     []string{"string", "string", "number"},
			},
		},
	},
}

func TestIntergration_Transfer(t *testing.T) {

	mvccdb, err := db.NewMVCCDB("mvcc")
	if err != nil {
		t.Fatal(err)
	}

	//mvccdb := database.NewDatabase()

	defer os.RemoveAll("mvcc")

	vi := database.NewVisitor(0, mvccdb)

	vi.SetBalance(testID[0], 1000000)

	vi.SetContract(systemContract)

	bh := &block.BlockHead{
		ParentHash: []byte("abc"),
		Number:     10,
		Witness:    "witness",
		Time:       123456,
	}

	e := NewEngine(bh, mvccdb)

	e.SetUp("js_path", jsPath)
	e.SetUp("log_level", "debug")
	e.SetUp("log_enable", "")

	act := tx.NewAction("iost.system", "Transfer", fmt.Sprintf(`["%v","%v",%v]`, testID[0], testID[2], "100"))

	trx := tx.NewTx([]tx.Action{act}, nil, int64(10000), int64(1), int64(10000000))

	ac, err := account.NewAccount(common.Base58Decode(testID[1]))
	if err != nil {
		t.Fatal(err)
	}
	trx, err = tx.SignTx(trx, ac)
	if err != nil {
		t.Fatal(err)
	}
	//
	//	cpl := contract.Compiler{}
	//
	//	code := `
	//class Contract {
	// constructor() {
	//
	// }
	// hello(someone) {
	//  return "hello "+ someone + "!";
	// }
	//}
	//
	//module.exports = Contract;
	//`
	//
	//	abi := `
	//{
	//  "lang": "javascript",
	//  "version": "1.0.0",
	//  "abi": [{
	//    "name": "hello",
	//    "args": ["string"],
	//    "payment": 0,
	//    "cost_limit": [1,1,1],
	//    "price_limit": 1
	//  }
	//  ]
	//}
	//`
	//
	//	c, err := cpl.Parse("contract", code, abi)
	//	if err != nil {
	//		t.Fatal(err)
	//	}

	//acSet := tx.NewAction("iost.system", "SetCode", fmt.Sprintf(`["%v","%v",%v]`, testID[0], testID[2], "100"))
	//
	//trxSet := tx.NewTx([]tx.Action{act}, nil, int64(10000), int64(1), int64(10000000))

	t.Log(e.Exec(&trx))
	t.Log("balance of sender :", vi.Balance(testID[0]))
	t.Log("balance of receiver :", vi.Balance(testID[2]))
}
