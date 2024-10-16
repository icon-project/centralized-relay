package bitcoin

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/icon-project/icon-bridge/common/codec"
)

func TestDecode(t *testing.T) {

	// hexString := types.new
	data, _ := hex.DecodeString("f884874465706f73697483303a31b83e74623170677a7838383079667237713864677a38647168773530736e6375346634686d7735636e3338303033353474757a6379396a783573687676377375b33078322e69636f6e2f687830316361383532383764363334323732326665373333633235363637363736623963663966386134823a9880")
	depositInfo := XCallMessage{}
	_, err := codec.RLP.UnmarshalFromBytes(data, &depositInfo)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(depositInfo)
	amount := new(big.Int).SetBytes(depositInfo.Amount)
	amountInDecimal := new(big.Float).SetInt(amount)
	fmt.Println(amountInDecimal)
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

	// // withdraw info data
	withdrawInfoWrapperV2 := CSMessageRequestV2{}
	_, err = codec.RLP.UnmarshalFromBytes(withdrawInfoWrapper.Payload, &withdrawInfoWrapperV2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(withdrawInfoWrapperV2)
	fmt.Println("\n =============")

	// // withdraw info
	withdrawInfo := XCallMessage{}
	_, err = codec.RLP.UnmarshalFromBytes(withdrawInfoWrapperV2.Data, &withdrawInfo)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(withdrawInfo)
}

func TestDecodeConnectionContract(t *testing.T) {
	data, _ := hex.DecodeString("f9013701b90133f90130b8463078322e6274632f74623170677a7838383079667237713864677a38647168773530736e6375346634686d7735636e3338303033353474757a6379396a783573687676377375aa63786663383665653736383765316266363831623535343862323636373834343438356330653731393287308a0a0000001501b886f884874465706f73697483303a31b83e74623170677a7838383079667237713864677a38647168773530736e6375346634686d7735636e3338303033353474757a6379396a783573687676377375b33078322e69636f6e2f687830316361383532383764363334323732326665373333633235363637363736623963663966386134823a9880ebaa637835373766356537353661626438396362636261333861353835303862363061313237353464326635")
	csMessage := CSMessage{}
	_, err := codec.RLP.UnmarshalFromBytes(data, &csMessage)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("csMessage: ", csMessage)
	csMessageV2 := CSMessageRequestV2{}
	_, err = codec.RLP.UnmarshalFromBytes(csMessage.Payload, &csMessageV2)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("csMessageV2 MessageType: ", csMessageV2.MessageType)
	fmt.Println("csMessageV2 Data: ", csMessageV2.Data)
	fmt.Println("csMessageV2 Protocols: ", csMessageV2.Protocols)
	fmt.Println("csMessageV2 From: ", csMessageV2.From)
	fmt.Println("csMessageV2 To: ", csMessageV2.To)
	fmt.Println("csMessageV2 Sn: ", csMessageV2.Sn)
}

func TestGetRuneBalanceAtIndex(t *testing.T) {
	//
	res, err := GetRuneTxIndex("https://open-api.unisat.io/v1/indexer/runes", "GET", os.Getenv("APIToken"), "60fa23d19c8116dbb09441bf3d1ee27067c3d2b3735caf2045db84ea8f76d436", 2)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("%+v", res)
}
