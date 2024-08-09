package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"log"
	"os"
	"testing"

	"github.com/icon-project/icon-bridge/common/codec"
	"github.com/stretchr/testify/assert"
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

func TestEncodeToXcallMessage(t *testing.T) {
	//
	test1 := multisig.RadFiProvideLiquidityMsg{
		Fee:       30,
		UpperTick: -173940,
		LowerTick: -178320,
		Min0:      0,
		Min1:      0,
	}

	protocols := []string{
		"0x932e088453515720B8eD50c1999C4Bc7bc11991F",
		"0x526C5Bcd376FAD738780e099E9723A62044D0319",
	}

	res, err := ToXCallMessage(
		test1,
		"0x3.BTC/bc1qvqkshkdj67uwvlwschyq8wja6df4juhewkg5fg",
		"0x2a68F967bFA230780a385175d0c86AE4048d3096",
		2,
		protocols,
		common.HexToAddress("0x000013938B55EDFBF3380656CC321770cCF470E1"),
		common.HexToAddress("0x66f2A9220C8479d73eE84df0932f38C496e8E9e3"),
		common.HexToAddress("0x7a4a1aF7B59c5FF522D5F11336d5b20d5116c7cb"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	resReadable := hex.EncodeToString(res)
	fmt.Printf("%+v", resReadable)

	assert.Equal(t, resReadable, "f902a201b9029ef9029bb23078332e4254432f6263317176716b73686b646a36377577766c77736368797138776a61366466346a756865776b67356667aa3078326136384639363762464132333037383061333835313735643063383641453430343864333039360200b901e0000000000000000000000000000013938b55edfbf3380656cc321770ccf470e1000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001648831645600000000000000000000000066f2a9220c8479d73ee84df0932f38c496e8e9e30000000000000000000000007a4a1af7b59c5ff522d5f11336d5b20d5116c7cb0000000000000000000000000000000000000000000000000000000000000bb8fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd4770fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd588c0000000000000000000000000000000000000000000fb768105935a2f1a5b649000000000000000000000000000000000000000000000000077cf984c325302e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002a68f967bfa230780a385175d0c86ae4048d3096000000000000000000000000000000000000000000000000000000003b9aca0000000000000000000000000000000000000000000000000000000000f856aa307839333265303838343533353135373230423865443530633139393943344263376263313139393146aa307835323643354263643337364641443733383738306530393945393732334136323034344430333139")
}
