package bitcoin

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/icon-project/icon-bridge/common/codec"
)

func TestDecode(t *testing.T) {
	data, _ := hex.DecodeString("f90106874465706f736974aa307830303030303030303030303030303030303030303030303030303030303030303030303030303030aa307843393736333336343766363634664645443041613933363937374437353638363664433733433037b33078312e69636f6e2f637832316539346330386330336461656538306332356438656533656132326132303738366563323331890db7148f8c6dd40000b8687b226d6574686f64223a225f73776170222c22706172616d73223a7b2270617468223a5b5d2c227265636569766572223a223078312e69636f6e2f687830343232393361393061303433656136383933653731663565343662663464646366323334633233227d7d")
	depositInfo := XCallMessage{}
	_, err := codec.RLP.UnmarshalFromBytes(data, &depositInfo)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(depositInfo)
}

func TestGetRuneBalanceAtIndex(t *testing.T) {
	//
	res, err := GetRuneTxIndex("https://open-api.unisat.io/v1/indexer/runes", "GET", os.Getenv("APIToken"), "60fa23d19c8116dbb09441bf3d1ee27067c3d2b3735caf2045db84ea8f76d436", 2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%+v", res)
}
