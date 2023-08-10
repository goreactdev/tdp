package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ton-developer-program/internal/database"
	"github.com/ton-developer-program/internal/funcs"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/ton/nft"
)

func (app *application) getCollection(collection database.SBTCollection) (*database.SBTCollection, error) {

	collectionAddr, err := address.ParseAddr(collection.FriendlyAddress)
	if err != nil {
		return nil, fmt.Errorf("error parsing address: %v", err)
	}

	app.logger.Info(fmt.Sprintf("getting collections for address %v", collection.FriendlyAddress))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collectionClient := nft.NewCollectionClient(app.tonLiteClient, collectionAddr)

	collectionData, err := collectionClient.GetCollectionData(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting collection data: %v", err)
	}

	rawCollectionAddress, err := funcs.NewAddress(collection.FriendlyAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating address: %v", err)
	}

	var collectionUrl string

	if onchainCollection, ok := collectionData.Content.(*nft.ContentOnchain); ok {
		collectionUrl = onchainCollection.Image
		if err != nil {
			return nil, fmt.Errorf("error getting nft content cell: %v", err)
		}
	}

	// check if nft's content is offchain
	if offchainCollection, ok := collectionData.Content.(*nft.ContentOffchain); ok {
		collectionUrl = offchainCollection.URI
		if err != nil {
			return nil, fmt.Errorf("error getting nft content cell: %v", err)
		}

	}

	// check if nft's content is semichain
	if semichainCollection, ok := collectionData.Content.(*nft.ContentSemichain); ok {
		collectionUrl = semichainCollection.URI
		if err != nil {
			return nil, fmt.Errorf("error getting nft content cell: %v", err)
		}

	}

	collectionMetadata, collectionBody, err := app.getMetaData(collectionUrl)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error get metadata: %v", err))
		return nil, err
	}

	// raw orwner address
	rawOwnerAddress, err := funcs.NewAddress(collectionData.OwnerAddress.String())
	if err != nil {
		return nil, fmt.Errorf("error creating address: %v", err)
	}

	insertCollection := &database.SBTCollection{
		FriendlyAddress:      collection.FriendlyAddress,
		RawAddress:           rawCollectionAddress.ToString(),
		FriendlyOwnerAddress: collectionData.OwnerAddress.String(),
		RawOwnerAddress:      rawOwnerAddress.ToString(),
		NextItemIndex:        collectionData.NextItemIndex.Int64(),
		ContentUri:           collectionUrl,
		DefaultWeight:        collection.DefaultWeight,
		Name:                 collectionMetadata.Name,
		Description:          collectionMetadata.Description,
		Image:                collectionMetadata.Image,
		ContentJson:          collectionBody,
	}

	// insert collection into database
	insertedCollection, err := app.sqlModels.Nfts.InsertCollection(insertCollection)
	if err != nil {
		return nil, fmt.Errorf("error inserting collection: %v", err)
	}

	app.logger.Info(fmt.Sprintf("inserted collection %v", insertedCollection.RawAddress))

	return insertedCollection, nil

}


func (app *application) getNFTByCollection(collectionAddress string, index int64) (*database.SBTToken, error) {
	parseCollectionAddr, err := address.ParseAddr(collectionAddress)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error parsing collection address: %v", err))
		return nil, err
	}

	app.logger.Info(fmt.Sprintf("migrating nfts for collection %v", collectionAddress))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nftAddr, err := app.getNftAddressByIndex(ctx, parseCollectionAddr, index)
	if err != nil {
		app.logger.Warning(fmt.Sprintf("error getting nft address by appex: %v", err))
		return nil, err
	}

	app.logger.Info(fmt.Sprintf("nft address %v", nftAddr.String()))

	item := nft.NewItemClient(app.tonLiteClient, nftAddr)

	nftData, err := item.GetNFTData(context.Background())
	if err != nil {
		return nil, err
	}

	collection := nft.NewCollectionClient(app.tonLiteClient, nftData.CollectionAddress)

	if err != nil {
		return nil, err
	}

	// get collection by address
	collectionDB, err := app.sqlModels.Nfts.GetCollectionByAddress(collectionAddress)
	if err != nil {
		return nil, err
	}

	if nftData.Initialized {
		// get full nft's content url using collection method that will merge base url with nft's data

		nftContent, err := collection.GetNFTContent(context.Background(), nftData.Index, nftData.Content)
		if err != nil {
			return nil, err
		}

		var urlNft string

		// check if nft's content is onchain
		if onchain, ok := nftContent.(*nft.ContentOnchain); ok {
			urlNft = onchain.Image

		}
		// check if nft's content is offchain
		if offchain, ok := nftContent.(*nft.ContentOffchain); ok {
			urlNft = offchain.URI
		}

		// check if nft's content is semichain
		if semichain, ok := nftContent.(*nft.ContentSemichain); ok {
			urlNft = semichain.URI
		}

		// convert nft
		nftIndex := nftData.Index.Int64()

		// get nfts	metadata through client

		metadata, body, err := app.getMetaData(urlNft)
		if err != nil {
			app.logger.Warning(fmt.Sprintf("error get metadata: %v", err))
			return nil, err
		}

		rawAddress, err := funcs.NewAddress(nftAddr.String())
		if err != nil {
			return nil, fmt.Errorf("error creating address: %v", err)
		}

		rawOwnerAddress, err := funcs.NewAddress(nftData.OwnerAddress.String())
		if err != nil {
			return nil, fmt.Errorf("error creating address: %v", err)
		}
		rawOwnerAddressString := rawOwnerAddress.ToString()

		var weight int64 = collectionDB.DefaultWeight

		if strings.Contains(urlNft, app.config.App.DomainName) {
			base64FromUrl := strings.Split(urlNft, "/")[6]

			metadataWeight, err := app.sqlModels.Nfts.GetWeightByBase64(base64FromUrl)
			if err != nil {
				return nil, err
			}
	
			if metadataWeight != 0 {
				weight = metadataWeight
			}	
		}


		var insertNFT = &database.SBTToken{
			FriendlyAddress:     nftAddr.String(),
			RawAddress:          rawAddress.ToString(),
			SBTCollectionID:     collectionDB.ID,
			RawOwnerAddress:     rawOwnerAddressString,
			FriendlyOwnerAddress: nftData.OwnerAddress.String(),
			ContentUri:         urlNft,
			Name:                metadata.Name,
			Description:         metadata.Description,
			ContentJson:         body,
			Weight: 			 weight,
			Index:               nftIndex,
			Image:               metadata.Image,
		}

		tx, err := app.sqlModels.Nfts.DB.BeginTxx(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating transaction: %v", err)
		}

		err = app.sqlModels.Nfts.InsertToken(tx, insertNFT)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error inserting nfts: %v", err)
		}

		err = tx.Commit()
		if err != nil {
			return nil, fmt.Errorf("error commiting transaction: %v", err)
		}

		app.logger.Info(fmt.Sprintf("nft %v upserted", nftAddr.String()))

		return insertNFT, nil
	} else {
		app.logger.Info(fmt.Sprintf("nft %v is not initialized", nftAddr.String()))
	}

	app.logger.Info(fmt.Sprintf("migrated %v nft for collection %v", nftAddr, collectionAddress))

	return nil, nil
}


func (app *application) insertNFTbyAddr(tx *sqlx.Tx, nftAddr string) (*database.SBTToken, error) {
	app.logger.Info(fmt.Sprintf("nft address %v", nftAddr))

	item := nft.NewItemClient(app.tonLiteClient, address.MustParseAddr(nftAddr))

	nftData, err := item.GetNFTData(context.Background())
	if err != nil {
		return nil, err
	}

	collection := nft.NewCollectionClient(app.tonLiteClient, nftData.CollectionAddress)

	if err != nil {
		return nil, err
	}

	// get collection by address
	collectionDB, err := app.sqlModels.Nfts.GetCollectionByAddress(nftData.CollectionAddress.String())
	if err != nil {
		return nil, err
	}

	if nftData.Initialized {
		// get full nft's content url using collection method that will merge base url with nft's data

		nftContent, err := collection.GetNFTContent(context.Background(), nftData.Index, nftData.Content)
		if err != nil {
			return nil, err
		}

		var urlNft string

		// check if nft's content is onchain
		if onchain, ok := nftContent.(*nft.ContentOnchain); ok {
			urlNft = onchain.Image

		}
		// check if nft's content is offchain
		if offchain, ok := nftContent.(*nft.ContentOffchain); ok {
			urlNft = offchain.URI
		}

		// check if nft's content is semichain
		if semichain, ok := nftContent.(*nft.ContentSemichain); ok {
			urlNft = semichain.URI
		}

		// convert nft
		nftIndex := nftData.Index.Int64()

		// get nfts	metadata through client

		metadata, body, err := app.getMetaData(urlNft)
		if err != nil {
			app.logger.Warning(fmt.Sprintf("error get metadata: %v", err))
		}

		rawAddress, err := funcs.NewAddress(nftAddr)
		if err != nil {
			return nil, fmt.Errorf("error creating address: %v", err)
		}

		rawOwnerAddress, err := funcs.NewAddress(nftData.OwnerAddress.String())
		if err != nil {
			return nil, fmt.Errorf("error creating address: %v", err)
		}
		rawOwnerAddressString := rawOwnerAddress.ToString()


		var weight int64 = collectionDB.DefaultWeight

		if strings.Contains(urlNft, app.config.App.DomainName) {
			base64FromUrl := strings.Split(urlNft, "/")[6]

			metadataWeight, err := app.sqlModels.Nfts.GetWeightByBase64(base64FromUrl)
			if err != nil {
				return nil, err
			}
	
			if metadataWeight != 0 {
				weight = metadataWeight
			}	
		}


		var insertNFT = &database.SBTToken{
			FriendlyAddress:     nftAddr,
			RawAddress:          rawAddress.ToString(),
			SBTCollectionID:     collectionDB.ID,
			RawOwnerAddress:     rawOwnerAddressString,
			FriendlyOwnerAddress: nftData.OwnerAddress.String(),
			ContentUri:         urlNft,
			Name:                metadata.Name,
			Description:         metadata.Description,
			ContentJson:         body,
			Weight: 			 weight,
			Index:               nftIndex,
			Image:               metadata.Image,
		}

		err = app.sqlModels.Nfts.InsertToken(tx, insertNFT)
		if err != nil {
			return nil, fmt.Errorf("error inserting nfts: %v", err)
		}

		app.logger.Info(fmt.Sprintf("nft %v upserted", nftAddr))

		return insertNFT, nil
	} else {
		app.logger.Info(fmt.Sprintf("nft %v is not initialized", nftAddr))
	}

	app.logger.Info(fmt.Sprintf("migrated %v nft for collection %v", nftAddr, collectionDB.FriendlyAddress))

	return nil, nil
}
