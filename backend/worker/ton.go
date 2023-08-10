package main

import (
	"context"
	"fmt"

	"github.com/xssnick/tonutils-go/address"
)

func (app *application) getNftAddressByIndex(ctx context.Context, collectionAddress *address.Address, nextItemIndex int64) (*address.Address, error) {
	b, err := app.tonLiteClient.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}

	res, err := app.tonLiteClient.RunGetMethod(ctx, b, collectionAddress, "get_nft_address_by_index", nextItemIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to run get_nft_address_by_index method: %w", err)
	}

	nftAddressRes, err := res.Slice(0)
	if err != nil {
		return nil, fmt.Errorf("err get nftAddressRes slice value: %w", err)
	}

	nftAddress, err := nftAddressRes.LoadAddr()
	if err != nil {
		return nil, fmt.Errorf("failed to load nftAddress from result slice: %w", err)
	}

	return nftAddress, nil
}
