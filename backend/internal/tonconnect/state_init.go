package tonconnect

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ton-developer-program/util"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/wallet"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
)

var knownHashes = make(map[string]wallet.Version)

func init() {
	for i := wallet.Version(0); i <= wallet.V4R2; i++ {
		ver := wallet.GetCodeHashByVer(i)
		knownHashes[hex.EncodeToString(ver[:])] = i
	}
}

func ParseStateInit(stateInit string) ([]byte, error) {
	cells, err := boc.DeserializeBocBase64(stateInit)
	if err != nil || len(cells) != 1 {
		return nil, err
	}
	var state tlb.StateInit
	err = tlb.Unmarshal(cells[0], &state)
	if err != nil {
		return nil, err
	}
	if !state.Data.Exists || !state.Code.Exists {
		return nil, fmt.Errorf("empty init state")
	}
	codeHash, err := state.Code.Value.Value.HashString()
	if err != nil {
		return nil, err
	}
	version, prs := knownHashes[codeHash]
	if !prs {
		return nil, err
	}
	var pubKey tlb.Bits256
	switch version {
	case wallet.V1R1, wallet.V1R2, wallet.V1R3, wallet.V2R1, wallet.V2R2:
		var data wallet.DataV1V2
		err = tlb.Unmarshal(&state.Data.Value.Value, &data)
		if err != nil {
			return nil, err
		}
		pubKey = data.PublicKey
	case wallet.V3R1, wallet.V3R2, wallet.V4R1, wallet.V4R2:
		var data wallet.DataV3
		err = tlb.Unmarshal(&state.Data.Value.Value, &data)
		if err != nil {
			return nil, err
		}
		pubKey = data.PublicKey
	}

	return pubKey[:], nil
}



func ConvertToFriendlyAddr(s string) *address.Address {
	// divide string to 2 parts by :
	var arr = strings.Split(s, ":")
	if len(arr) != 2 {
		panic("Invalid address " + s)
	}

	// decode hex to bytes in second part
	decoded, err := hex.DecodeString(arr[1])
	if err != nil {
		panic("Invalid address hex " + s)
	}

	return address.NewAddress(0, 0, decoded)
}



func NewTonConnection(pool *liteclient.ConnectionPool, config util.Config) (*ton.APIClient, error) {
	if config.Ton.PublicConfig != "" {
		err := pool.AddConnectionsFromConfigUrl(context.Background(), config.Ton.PublicConfig)
		if err != nil {
			return nil, err
		}

		return ton.NewAPIClient(pool), nil
	}
	
	err := pool.AddConnection(context.Background(), config.Ton.NodeAddress, config.Ton.ApiKey)
	if err != nil {
		return nil, err
	}

	return ton.NewAPIClient(pool), nil
}
