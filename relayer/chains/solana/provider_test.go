package solana

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/alt"
	"github.com/icon-project/centralized-relay/relayer/chains/solana/types"
	"github.com/near/borsh-go"
	"github.com/stretchr/testify/assert"
)

type IDLAccount struct {
	Authority solana.PublicKey
	DataLen   uint32
}

func TestIDL(t *testing.T) {
	dataraw := []byte{24, 70, 98, 191, 58, 144, 123, 158, 93, 50, 105, 173, 228, 54, 251, 71, 232, 235, 85, 94, 2, 72, 159, 136, 154, 169, 126, 147, 158, 138, 34, 80, 203, 137, 250, 29, 223, 82, 226, 157, 70, 2, 0, 0, 120, 156, 197, 84, 127, 111, 218, 48, 16, 253, 42, 85, 254, 182, 170, 216, 249, 65, 224, 63, 198, 218, 161, 106, 155, 208, 58, 105, 221, 170, 42, 50, 137, 27, 44, 136, 195, 108, 3, 67, 136, 239, 190, 151, 100, 162, 166, 131, 180, 147, 42, 13, 9, 37, 241, 221, 189, 187, 123, 247, 206, 59, 143, 231, 185, 22, 198, 120, 3, 239, 90, 222, 77, 175, 62, 100, 239, 103, 215, 147, 241, 183, 36, 28, 221, 148, 227, 234, 203, 84, 23, 223, 231, 239, 134, 87, 219, 201, 205, 199, 94, 49, 156, 76, 190, 6, 227, 254, 143, 91, 143, 120, 165, 176, 60, 231, 150, 123, 131, 157, 167, 120, 41, 128, 240, 43, 227, 139, 5, 76, 107, 161, 141, 172, 20, 78, 252, 75, 122, 233, 227, 196, 44, 69, 230, 124, 230, 194, 100, 90, 46, 109, 235, 52, 210, 130, 91, 145, 95, 108, 164, 157, 93, 12, 85, 54, 171, 180, 183, 39, 158, 84, 198, 234, 85, 86, 59, 161, 188, 251, 67, 150, 66, 216, 180, 201, 148, 26, 139, 184, 26, 78, 214, 112, 165, 84, 220, 34, 116, 112, 31, 244, 9, 13, 19, 130, 71, 18, 17, 198, 8, 237, 37, 132, 250, 140, 68, 236, 129, 120, 60, 203, 170, 149, 178, 71, 144, 43, 35, 52, 112, 54, 90, 90, 62, 93, 224, 4, 137, 5, 170, 150, 133, 130, 161, 249, 218, 147, 227, 54, 255, 36, 223, 215, 136, 186, 168, 209, 240, 166, 133, 93, 233, 186, 218, 29, 90, 124, 148, 74, 228, 14, 59, 119, 117, 216, 109, 27, 181, 119, 224, 164, 146, 246, 68, 19, 140, 249, 36, 234, 19, 230, 247, 8, 11, 98, 212, 159, 16, 22, 249, 36, 236, 225, 213, 127, 147, 62, 204, 214, 88, 81, 166, 75, 93, 21, 154, 151, 136, 123, 210, 2, 125, 225, 231, 157, 161, 227, 133, 220, 14, 89, 14, 128, 22, 217, 58, 45, 145, 152, 23, 167, 166, 25, 130, 4, 10, 46, 192, 65, 80, 207, 146, 208, 8, 127, 76, 148, 5, 111, 62, 206, 231, 65, 78, 197, 7, 255, 212, 232, 12, 142, 118, 187, 108, 72, 180, 90, 170, 194, 229, 35, 53, 234, 201, 188, 138, 67, 215, 214, 44, 204, 193, 56, 221, 90, 97, 32, 33, 119, 38, 66, 229, 29, 92, 160, 243, 208, 39, 65, 216, 74, 58, 129, 206, 125, 18, 67, 227, 241, 127, 97, 194, 86, 93, 60, 148, 166, 56, 209, 234, 153, 66, 157, 229, 248, 187, 107, 26, 82, 66, 105, 132, 193, 67, 10, 1, 37, 12, 28, 196, 49, 84, 17, 146, 126, 216, 64, 138, 181, 120, 6, 56, 2, 222, 167, 179, 60, 82, 22, 213, 107, 132, 203, 1, 79, 74, 73, 143, 97, 209, 144, 36, 137, 27, 89, 57, 93, 124, 22, 155, 243, 48, 61, 130, 141, 164, 97, 76, 130, 0, 241, 152, 3, 150, 148, 66, 168, 81, 91, 85, 221, 123, 71, 81, 45, 53, 59, 111, 46, 85, 222, 50, 136, 219, 14, 231, 143, 82, 44, 242, 163, 56, 45, 126, 166, 50, 255, 87, 85, 157, 235, 226, 245, 121, 45, 38, 142, 219, 86, 9, 187, 169, 244, 188, 107, 216, 29, 154, 63, 173, 3, 199, 225, 104, 246, 175, 47, 174, 35, 229, 73, 190, 144, 243, 97, 255, 27, 192, 185, 55, 57, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	data := dataraw[8:]

	var idlAc IDLAccount
	err := borsh.Deserialize(&idlAc, data)
	assert.NoError(t, err)

	compressedBytes := data[36 : 36+idlAc.DataLen]

	fmt.Println("Authority:", idlAc.Authority.Bytes())
	fmt.Println("Data:", data)
	fmt.Println("Compressed Bytes:", compressedBytes)

	decompressedBytes, err := decompress(compressedBytes)
	assert.NoError(t, err)

	var idlData IDL

	err = json.Unmarshal(decompressedBytes, &idlData)
	assert.NoError(t, err)

	fmt.Printf("%+v", idlData)
}

func TestIDLAddress(t *testing.T) {
	progPubkey, err := solana.PublicKeyFromBase58("FiXbEGcDhFPHW84CJmHoRbrgYkBAEyPJL7gAPPT3H9ZS")
	assert.NoError(t, err)

	idlPubkey, err := solana.PublicKeyFromBase58("3fFhJNrxpdnKcxsY9sem81bg3VPQL5FySwzg99354spR")
	assert.NoError(t, err)

	signer, _, err := solana.FindProgramAddress([][]byte{}, progPubkey)
	assert.NoError(t, err)

	calculatedIdlAddr, err := solana.CreateWithSeed(signer, "anchor:idl", progPubkey)
	assert.NoError(t, err)

	assert.Equal(t, idlPubkey, calculatedIdlAddr)
}

func TestEventLogParse(t *testing.T) {
	eventData := "HDQnaQjSWwkIAAAAMHgzLmljb24BAAAAAAAAAA=="

	// Decode the Base64 encoded string
	eventBytes, err := base64.StdEncoding.DecodeString(eventData)
	assert.NoError(t, err)

	fmt.Printf("Decoded bytes: %v\n", eventBytes)

	// Extract the first 8 bytes as discriminator
	if len(eventBytes) < 8 {
		t.Fatalf("Decoded bytes too short to contain discriminator: %v", eventBytes)
	}
	discriminator := eventBytes[:8]
	fmt.Printf("Discriminator: %v\n", discriminator)

	// Remaining bytes are the serialized event
	remainingBytes := eventBytes[8:]

	ev := struct {
		To string
		Sn uint64
	}{}

	// Deserialize the remaining bytes into the TestEvent struct
	err = borsh.Deserialize(&ev, remainingBytes)
	assert.NoError(t, err)

	fmt.Printf("TestEvent: %+v\n", ev)
}

func TestCsMessageDecode(t *testing.T) {
	msgLog := "AAEzAAAAMHgzLmljb24vaHhiNmI1NzkxYmUwYjVlZjY3MDYzYjNjMTBiODQwZmI4MTUxNGRiMmZkLAAAADhRNEZ2c0hDV0s2OEV6WXRzc3RkRll3VUwxU0hDaXVMUFJESmsxZ2FLaVE4AwAAAAAAAAAAAAAAAAAAAAABAAAAkAEAAAAsAAAAR3Z6ZnlQQnNuaXdWbm40S2pHNG9nS2VxOXl5cXhIamR5aU11WlJtaHBZOWUA"

	msgBytes, err := base64.StdEncoding.DecodeString(msgLog)
	assert.NoError(t, err)

	fmt.Println("Msg Bytes:", msgBytes)

	msg := types.CsMessage{}

	err = borsh.Deserialize(&msg, msgBytes[:])

	assert.NoError(t, err)
	fmt.Printf("\nDecoded Msg Request: %+v\n", msg.Request)
	fmt.Printf("\nDecoded Msg Result: %+v\n", msg.Result)
}

func TestBigInt(t *testing.T) {
	xcallProgID, err := solana.PublicKeyFromBase58("7Ya5FjxuYScJhkjRs4WK9cxFXnvfw6Eby89W3xNkUQyF")
	assert.NoError(t, err)
	configAddr, _, err := solana.FindProgramAddress([][]byte{[]byte("config")}, xcallProgID)
	assert.NoError(t, err)

	fmt.Println("Xcall Config Addr: ", configAddr)

	pda := types.PDA{SeedPrefix: "config", ProgramID: xcallProgID}

	addr, err := pda.GetAddress()
	assert.NoError(t, err)

	fmt.Println("Addr: ", addr)
}

func TestExtendLookupTable(t *testing.T) {
	tableAddr, _ := solana.PublicKeyFromBase58("HJ6JRbBAPFfeUtiiD2VKAoTH9w7ZCyCGZSaevFFCZtsJ")
	authority, _ := solana.PublicKeyFromBase58("FUarP2p5EnxD66vVDL4PWRoWMzA56ZVHG24hpEDFShEz")

	addr1, _ := solana.PublicKeyFromBase58("9aE476sH92Vz7DMPyq5WLPkrKWivxeuTKEFKd2sZZcde")
	addr2, _ := solana.PublicKeyFromBase58("2xNweLHLqrbx4zo1waDvgWJHgsUpPj8Y8icbAFeR4a8i")
	addresses := solana.PublicKeySlice{addr1, addr2}

	expectedData := []byte{
		2, 0, 0, 0,
		2, 0, 0, 0, 0, 0, 0, 0,
		127, 96, 107, 250, 152, 133, 208, 224, 73, 251, 113, 151, 128, 139, 86, 80, 101, 70, 138, 50, 141, 153, 218, 110, 56, 39, 122, 181, 120, 55, 86, 185,
		29, 11, 113, 4, 101, 239, 39, 167, 201, 112, 156, 239, 236, 36, 251, 140, 76, 199, 150, 228, 218, 214, 20, 123, 180, 181, 103, 160, 71, 251, 237, 123,
	}

	ins := alt.ExtendLookupTable(tableAddr, authority, nil, addresses)

	actualData, err := ins.Data()
	assert.NoError(t, err)

	assert.Equal(t, expectedData, actualData)
}

func TestGetLookupAccount(t *testing.T) {
	acData := []byte{
		1, 0, 0, 0, 255, 255, 255, 255, 255, 255, 255, 255, 84, 176, 0, 0, 0, 0, 0, 0, 0, 1, 93, 50,
		105, 173, 228, 54, 251, 71, 232, 235, 85, 94, 2, 72, 159, 136, 154, 169, 126, 147, 158, 138,
		34, 80, 203, 137, 250, 29, 223, 82, 226, 157, 0, 0, 201, 58, 234, 217, 141, 202, 49, 220, 114,
		198, 42, 55, 222, 248, 182, 167, 185, 171, 252, 249, 154, 158, 84, 243, 9, 165, 106, 55, 136,
		250, 196, 132, 255, 86, 146, 250, 112, 188, 234, 19, 220, 98, 78, 56, 46, 136, 210, 141, 229,
		7, 195, 223, 133, 136, 183, 200, 219, 183, 84, 77, 213, 114, 25, 221,
	}

	ac, err := alt.DeserializeLookupTable(acData)
	assert.NoError(t, err)

	fmt.Printf("Account %+v", ac)
}

func TestTest(t *testing.T) {
	res, err := reateLookupTableAccount(context.Background())
	assert.NoError(t, err)

	fmt.Println("Received Pubkey: ", res.String())
	time.Sleep(5 * time.Second)
}

func reateLookupTableAccount(ctx context.Context) (*solana.PublicKey, error) {
	createFunc := func(rs uint64) (*solana.PublicKey, error) {
		time.Sleep(time.Duration(rs * uint64(time.Second)))
		if rs%2 == 0 {
			return &solana.SystemProgramID, nil
		} else {
			return nil, fmt.Errorf("sending error")
		}
	}

	recentSlot := uint64(1)

	type itemType struct {
		ac  *solana.PublicKey
		err error
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	n := uint64(6)

	responseChan := make(chan itemType, n)
	wg := &sync.WaitGroup{}

	for slot := recentSlot; slot < recentSlot+n; slot++ {
		wg.Add(1)
		go func(ctx context.Context, sl uint64) {
			defer wg.Done()
			resKey, err := createFunc(sl)
			select {
			case <-ctx.Done():
				fmt.Printf("\nI am cancelled %d\n", sl)
			default:
				responseChan <- itemType{ac: resKey, err: err}
				fmt.Printf("\nI wrote %d\n", sl)
			}
		}(ctx, slot)
	}

	go func() {
		wg.Wait()
		close(responseChan)
		fmt.Println("all go rountine completed and channel closed")
	}()

	var lastErr error
	for item := range responseChan {
		if item.err == nil {
			cancel()
			return item.ac, nil
		} else {
			lastErr = item.err
		}
	}

	return nil, lastErr
}
