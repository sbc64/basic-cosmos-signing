package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/binance-chain/go-sdk/client"
	rpc "github.com/binance-chain/go-sdk/client/rpc"
	"github.com/binance-chain/go-sdk/common/types"
	keys "github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
)

type PersonalTranscation struct {
	sequence int64
	account  int64
}

func getSDKNodeInfo(key keys.KeyManager) PersonalTranscation {
	clientSDK, err := sdk.NewDexClient("testnet-dex.binance.org", types.TestNetwork, key)
	status, err := clientSDK.GetNodeInfo()
	fmt.Printf("Testnet: %+v\n", status.NodeInfo.Network)
	if err != nil {
		panic(err)
	}
	account, err := clientSDK.GetAccount(key.GetAddr().String())
	if err != nil {
		panic(err)
	}

	return PersonalTranscation{
		sequence: account.Sequence,
		account:  account.Number,
	}

}

func getRPCNodeInfo(address types.AccAddress) PersonalTranscation {
	clientRPC := rpc.NewRPCClient("http://localhost:27147", types.TestNetwork)

	account, err := clientRPC.GetAccount(address)
	if err != nil {
		panic(err)
	}
	return PersonalTranscation{
		sequence: account.GetSequence(),
		account:  account.GetAccountNumber(),
	}
}

func main() {
	dat, err := ioutil.ReadFile("pk")
	if err != nil {
		panic(err)
	}
	key, _ := keys.NewPrivateKeyManager(string(dat))
	key.GetAddr()

	mess := []msg.Msg{
		msg.CreateSendMsg(key.GetAddr(), types.Coins{types.Coin{Denom: "BNB", Amount: 1}}, []msg.Transfer{{key.GetAddr(), types.Coins{types.Coin{Denom: "BNB", Amount: 1}}}}),
	}
	fmt.Println("Sup")

	//personalTxn := getRPCNodeInfo(key.GetAddr())
	personalTxn := getSDKNodeInfo(key)

	fmt.Println(personalTxn)
	fmt.Println("Sup")

	m := tx.StdSignMsg{
		Msgs:          mess,
		Source:        0,
		Sequence:      personalTxn.sequence,
		AccountNumber: personalTxn.account,
		ChainID:       "Binance-Chain-Nile",
	}
	signed, err := key.Sign(m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Signed: %v\n", hex.EncodeToString(signed))

	if err != nil {
		panic(err)
	}

	url := fmt.Sprintf("http://localhost:27147/broadcast_tx_sync?tx=0x%s", hex.EncodeToString(signed))
	fmt.Println(url)
	var buffer []byte
	httpResponse, err := http.Post(url, "application/json", bytes.NewBuffer(buffer))
	if err != nil {
		panic(err)
	}
	defer httpResponse.Body.Close()
	bodyBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %+v\n", string(bodyBytes))
	fmt.Printf("Buffer: %+v\n", hex.EncodeToString(buffer))
}
