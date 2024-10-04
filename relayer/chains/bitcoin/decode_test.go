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
	fmt.Println("\n =============")

	// test decode 2
	data, _ = hex.DecodeString("f9011d01b90119f90116b33078322e69636f6e2f637866633836656537363837653162663638316235353438623236363738343434383563306537313932b83e74623170677a7838383079667237713864677a38647168773530736e6375346634686d7735636e3338303033353474757a6379396a783573687676377375821e8501b852f8508a5769746864726177546f83303a30b83e74623170677a7838383079667237713864677a38647168773530736e6375346634686d7735636e3338303033353474757a6379396a78357368767637737564f848b8463078322e6274632f74623170677a7838383079667237713864677a38647168773530736e6375346634686d7735636e3338303033353474757a6379396a783573687676377375")
	withdrawInfoWrapper := CSMessage{}
	_, err = codec.RLP.UnmarshalFromBytes(data, &withdrawInfoWrapper)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(withdrawInfoWrapper)
	fmt.Println("\n =============")

	// withdraw info data
	withdrawInfoWrapperV2 := CSMessageRequestV2{}
	_, err = codec.RLP.UnmarshalFromBytes(withdrawInfoWrapper.Payload, &withdrawInfoWrapperV2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(withdrawInfoWrapperV2)
	fmt.Println("\n =============")

	// withdraw info
	withdrawInfo := XCallMessage{}
	_, err = codec.RLP.UnmarshalFromBytes(withdrawInfoWrapperV2.Data, &withdrawInfo)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(withdrawInfo)
}

func TestGetRuneBalanceAtIndex(t *testing.T) {
	//
	res, err := GetRuneTxIndex("https://open-api.unisat.io/v1/indexer/runes", "GET", os.Getenv("APIToken"), "60fa23d19c8116dbb09441bf3d1ee27067c3d2b3735caf2045db84ea8f76d436", 2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%+v", res)
}
