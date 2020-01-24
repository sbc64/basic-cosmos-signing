package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

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
	// use this once your node is fully synced.
	clientRPC := rpc.NewRPCClient("localhost:26657", types.TestNetwork)

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
	// tbnb16hgptucs93skwsy6tdvl5p6kl3rq6x57s0jqdg
	//key, _ := keys.NewKeyStoreKeyManager("/home/sebas/Documents/binance.json", "")
	dat, err := ioutil.ReadFile("pk")
	if err != nil {
		panic(err)
	}
	key, _ := keys.NewPrivateKeyManager(string(dat))
	key.GetAddr()

	mess := []msg.Msg{
		msg.CreateSendMsg(key.GetAddr(), types.Coins{types.Coin{Denom: "BNB", Amount: 1}}, []msg.Transfer{{key.GetAddr(), types.Coins{types.Coin{Denom: "BNB", Amount: 1}}}}),
	}

	personalTxn := getRPCNodeInfo(key.GetAddr())

	fmt.Println(personalTxn)

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

	/*
		url := fmt.Sprintf("http://localhost:26657/broadcast_tx_sync?tx=0x%s", hex.EncodeToString(signed))
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
	*/
}
