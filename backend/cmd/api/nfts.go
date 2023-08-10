package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/flow"
	"github.com/hibiken/asynq"
	"github.com/ton-developer-program/internal/database"
	"github.com/ton-developer-program/internal/request"
	"github.com/ton-developer-program/internal/response"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

const (
	NFT_CONTRACT        = "b5ee9c720102130100033b000114ff00f4a413f4bcf2c80b0102016202030202ce04050201200f1004bd46c2220c700915be001d0d303fa4030f002f842b38e1c31f84301c705f2e195fa4001f864d401f866fa4030f86570f867f003e002d31f0271b0e30201d33f8210d0c3bfea5230bae302821004ded1485230bae3023082102fcb26a25220ba8060708090201200d0e00943031d31f82100524c7ae12ba8e39d33f308010f844708210c18e86d255036d804003c8cb1f12cb3f216eb39301cf179131e2c97105c8cb055004cf1658fa0213cb6accc901fb009130e200c26c12fa40d4d30030f847f841c8cbff5006cf16f844cf1612cc14cb3f5230cb0003c30096f8465003cc02de801078b17082100dd607e3403514804003c8cb1f12cb3f216eb39301cf179131e2c97105c8cb055004cf1658fa0213cb6accc901fb0000c632f8445003c705f2e191fa40d4d30030f847f841c8cbfff844cf1613cc12cb3f5210cb0001c30094f84601ccde801078b17082100524c7ae405503804003c8cb1f12cb3f216eb39301cf179131e2c97105c8cb055004cf1658fa0213cb6accc901fb0003fa8e4031f841c8cbfff843cf1680107082108b7717354015504403804003c8cb1f12cb3f216eb39301cf179131e2c97105c8cb055004cf1658fa0213cb6accc901fb00e082101f04537a5220bae30282106f89f5e35220ba8e165bf84501c705f2e191f847c000f2e193f823f867f003e08210d136d3b35220bae30230310a0b0c009231f84422c705f2e1918010708210d53276db102455026d830603c8cb1f12cb3f216eb39301cf179131e2c97105c8cb055004cf1658fa0213cb6accc901fb008b02f8648b02f865f003008e31f84422c705f2e191820afaf08070fb028010708210d53276db102455026d830603c8cb1f12cb3f216eb39301cf179131e2c97105c8cb055004cf1658fa0213cb6accc901fb00002082105fcc3d14ba93f2c19dde840ff2f000613b513434cfc07e187e90007e18dc3e188835d2708023859ffe18be90007e1935007e19be90007e1974cfcc3e19e44c38a000373e11fe11be107232cffe10f3c5be1133c5b33e1173c5b2cff27b55200201581112001dbc7e7f8017c217c20fc21fc227c234000db5631e005f08b0000db7b07e005f08f0"
	COLLECTION_CONTRACT = "b5ee9c724102140100021f000114ff00f4a413f4bcf2c80b0102016202030202cd04050201200e0f04e7d10638048adf000e8698180b8d848adf07d201800e98fe99ff6a2687d20699fea6a6a184108349e9ca829405d47141baf8280e8410854658056b84008646582a802e78b127d010a65b509e58fe59f80e78b64c0207d80701b28b9e382f970c892e000f18112e001718112e001f181181981e0024060708090201200a0b00603502d33f5313bbf2e1925313ba01fa00d43028103459f0068e1201a44343c85005cf1613cb3fccccccc9ed54925f05e200a6357003d4308e378040f4966fa5208e2906a4208100fabe93f2c18fde81019321a05325bbf2f402fa00d43022544b30f00623ba9302a402de04926c21e2b3e6303250444313c85005cf1613cb3fccccccc9ed54002c323401fa40304144c85005cf1613cb3fccccccc9ed54003c8e15d4d43010344130c85005cf1613cb3fccccccc9ed54e05f04840ff2f00201200c0d003d45af0047021f005778018c8cb0558cf165004fa0213cb6b12ccccc971fb008002d007232cffe0a33c5b25c083232c044fd003d0032c03260001b3e401d3232c084b281f2fff2742002012010110025bc82df6a2687d20699fea6a6a182de86a182c40043b8b5d31ed44d0fa40d33fd4d4d43010245f04d0d431d430d071c8cb0701cf16ccc980201201213002fb5dafda89a1f481a67fa9a9a860d883a1a61fa61ff480610002db4f47da89a1f481a67fa9a9a86028be09e008e003e00b01a500c6e"
)

const (
	DEPLOYED_NFT_POSTFIX        = "/v1/deployed-nft/n/"
	DEPLOYED_COLLECTION_POSTFIX = "/v1/deployed-nft/c/"
)

func (app *application) insertExistingCollectionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CollectionAddress string `json:"friendly_address"`
		DefaultWeight     string `json:"default_weight"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.logger.Error(err, nil)
		app.badRequest(w, r, err)
		return
	}

	weightInt64, err := strconv.ParseInt(input.DefaultWeight, 10, 64)
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	collection := &database.SBTCollection{
		FriendlyAddress: input.CollectionAddress,
		DefaultWeight:   weightInt64,
	}

	payload, err := json.Marshal(collection)

	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	runGetCollection := asynq.NewTask(database.TYPE_ADD_COLLECTION, payload)

	info, err := app.asynqClient.Enqueue(runGetCollection, asynq.ProcessIn(20*time.Second), asynq.MaxRetry(5), asynq.ProcessIn(5*time.Second), asynq.Retention(10*time.Minute), asynq.Queue(database.PRIORITY_URGENT))
	if err != nil {
		app.logger.Error(fmt.Errorf("error: %s", err), nil)
		return
	}

	app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
	})
}

func (app *application) insertCollectionHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	var input struct {
		DisplayName        string                   `json:"display_name"`
		Description        string                   `json:"description"`
		BannerImageUrl     *[]database.FileListItem `json:"banner_image_url"`
		CollectionImageUrl []database.FileListItem  `json:"collection_image_url"`
		DefaultWeight      string                   `json:"default_weight"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.logger.Error(err, nil)
		app.badRequest(w, r, err)
		return
	}

	codeCellBytesCollection, err := hex.DecodeString(COLLECTION_CONTRACT)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	codeCellCollection, err := cell.FromBOC(codeCellBytesCollection)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}
	nftItemCellBocBytes, err := hex.DecodeString(NFT_CONTRACT)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	nftItemCell, err := cell.FromBOC(nftItemCellBocBytes)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	royalty := cell.BeginCell().
		MustStoreUInt(0, 16).
		MustStoreUInt(0, 16).
		MustStoreAddr(address.MustParseAddr(app.config.Ton.AdminWallet)).
		EndCell()
	// generate hash based on displayName and current timestamp
	base64Data := base64.StdEncoding.EncodeToString([]byte(input.DisplayName + input.Description + strconv.FormatInt(time.Now().Unix(), 10)))

	// url safe
	base64Data = strings.Replace(base64Data, "+", "-", -1)
	base64Data = strings.Replace(base64Data, "/", "_", -1)

	// collection data
	collectionContent := nft.ContentOffchain{URI: app.config.App.BaseUrl + DEPLOYED_COLLECTION_POSTFIX + base64Data + "/meta.json"}
	collectionContentCell, err := collectionContent.ContentCell()
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	var coverImgUrl *string

	if input.BannerImageUrl != nil {
		for _, v := range *input.BannerImageUrl {
			coverImgUrl = &v.Response.URL
		}
	}

	metaJson := &database.CollectionMetadata{
		Base64:      base64Data,
		Name:        input.DisplayName,
		Description: input.Description,
		CoverImage:  coverImgUrl,
		Image:       input.CollectionImageUrl[0].Response.URL,
		ExternalURL: app.config.App.BaseUrl,
		Marketplace: app.config.App.DomainName,
	}

	// insert meta json to db
	err = app.sqlModels.Nfts.InsertCollectionMetadata(metaJson)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	// prefix for NFTs data
	uri := app.config.App.BaseUrl + DEPLOYED_NFT_POSTFIX
	commonContentCell := cell.BeginCell().MustStoreStringSnake(uri).EndCell()

	contentRef := cell.BeginCell().
		MustStoreRef(collectionContentCell).
		MustStoreRef(commonContentCell).
		EndCell()

	dataCell := cell.BeginCell().MustStoreAddr(address.MustParseAddr(user.FriendlyAddress)).
		MustStoreUInt(0, 64).
		MustStoreRef(contentRef).
		MustStoreRef(nftItemCell).
		MustStoreRef(royalty).
		EndCell()

	state := &tlb.StateInit{
		Data: dataCell,
		Code: codeCellCollection,
	}

	stateCell, err := tlb.ToCell(state)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	contractAddr := address.NewAddress(0, 0, stateCell.Hash())

	weightInt64, err := strconv.ParseInt(input.DefaultWeight, 10, 64)
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	collection := &database.SBTCollection{
		FriendlyAddress: contractAddr.String(),
		DefaultWeight:   weightInt64,
	}

	payload, err := json.Marshal(collection)

	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	runGetCollection := asynq.NewTask(database.TYPE_ADD_COLLECTION, payload)

	info, err := app.asynqClient.Enqueue(runGetCollection, asynq.ProcessIn(20*time.Second), asynq.MaxRetry(5), asynq.ProcessIn(5*time.Second), asynq.Retention(10*time.Minute), asynq.Queue(database.PRIORITY_URGENT))
	if err != nil {
		app.logger.Error(fmt.Errorf("error: %s", err), nil)
		return
	}

	app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))

	// url safe base64

	base64SignedDeployMsg := base64.URLEncoding.EncodeToString(stateCell.ToBOC())

	response.JSON(w, http.StatusAccepted, map[string]interface{}{
		"msg_body":         base64SignedDeployMsg,
		"contract_address": contractAddr.String(),
	})

}

func (app *application) getCollectionsHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	pagination, err := getPagination(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}


	var collections []*database.SBTCollection

	// get role by user 

	role, err := app.sqlModels.Permissions.GetUserRoles(user.ID)
	

	if role.Name == "admin" {
		collections, err = app.sqlModels.Nfts.GetAllCollections(pagination)
		if err != nil {
			app.logger.Error(err, nil)
			return
		}
	} else {
		collections, err = app.sqlModels.Nfts.GetCollectionsByOwnerAddress(pagination, user.FriendlyAddress)
		if err != nil {
			app.logger.Error(err, nil)
			return
		}
	}

	// add header x-total-count

	headers := http.Header{
		"x-total-count": []string{strconv.Itoa(len(collections))},
	}

	response.JSONWithHeaders(w, http.StatusOK, collections, headers)
}

// get all collections and tokens

func (app *application) getTokensHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	pagination, err := getPagination(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}


	
	role, err := app.sqlModels.Permissions.GetUserRoles(user.ID)

	var tokens []*database.SBTToken

	if role.Name == "admin" {
		tokens, err = app.sqlModels.Nfts.GetTokens(pagination, "")
		if err != nil {
			app.logger.Error(err, nil)
			return
		}
	} else {

		tokens, err = app.sqlModels.Nfts.GetTokens(pagination, user.FriendlyAddress)
		if err != nil {
			app.logger.Error(err, nil)
			return
		}
	}



	totalCount, err := app.sqlModels.Nfts.GetTotal(database.SBT_TOKENS_TABLE, "*", "")
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	// could be added to response
	headers := http.Header{
		"x-total-count":                 []string{strconv.Itoa(totalCount)},
		"Access-Control-Expose-Headers": []string{"X-Total-Count"},
	}

	response.JSONWithHeaders(w, http.StatusOK, tokens, headers)
}

func (app *application) getPrototypeTokensHandler(w http.ResponseWriter, r *http.Request) {

	pagination, err := getPagination(r)

	nameLike := r.URL.Query().Get("name_like")

	// create filters
	filters := []string{}
	if nameLike != "" {
		filters = append(filters, fmt.Sprintf("name ILIKE '%%%s%%'", nameLike))
	}

	filter := ""
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s", strings.Join(filters, " OR "))
	}

	tokens, err := app.sqlModels.Nfts.GetPrototypes(pagination, filter, "id DESC")
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	resPrototype := make([]struct {
		ID          int64                       `json:"id"`
		DisplayName string                      `json:"display_name"`
		Description string                      `json:"description"`
		Image       string                      `json:"image"`
		Attributes  database.AttributeJsonArray `json:"attributes"`
		Weight      int64                       `json:"weight"`
		Base64      string                      `json:"base64"`
	}, len(tokens))

	// get metadata
	for i, token := range tokens {
		metadata, err := app.sqlModels.Nfts.GetNFTMetadataByID(token.ID)
		if err != nil {
			app.logger.Error(err, nil)
			return
		}

		resPrototype[i].ID = token.ID
		resPrototype[i].DisplayName = metadata.Name
		resPrototype[i].Description = metadata.Description
		resPrototype[i].Image = metadata.Image
		resPrototype[i].Attributes = metadata.Attributes
		resPrototype[i].Weight = token.Weight
		resPrototype[i].Base64 = metadata.Base64
	}

	totalCount, err := app.sqlModels.Nfts.GetTotal(database.SBT_PROTOTYPE_TABLE, "*", "")
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	// could be added to response
	headers := http.Header{
		"x-total-count":                 []string{strconv.Itoa(totalCount)},
		"Access-Control-Expose-Headers": []string{"X-Total-Count"},
	}

	response.JSONWithHeaders(w, http.StatusOK, resPrototype, headers)
}

func (app *application) getMetaJsonNft(w http.ResponseWriter, r *http.Request) {
	base64String := flow.Param(r.Context(), "base64")

	// get meta json by hash
	metaJson, err := app.sqlModels.Nfts.GetNFTMetadataByBase64(base64String)
	if err != nil {

		if err == sql.ErrNoRows {
			app.notFound(w, r)
			return
		}
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	//

	if metaJson == nil {
		app.notFound(w, r)
		return
	}

	response.JSON(w, http.StatusOK, metaJson)

}

func (app *application) getMetaJsonCollection(w http.ResponseWriter, r *http.Request) {
	base64String := flow.Param(r.Context(), "base64")

	// get meta json by hash
	metaJson, err := app.sqlModels.Nfts.GetCollectionMetadataByBase64(base64String)
	if err != nil {

		if err == sql.ErrNoRows {
			app.notFound(w, r)
			return
		}

		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	//

	if metaJson == nil {
		app.notFound(w, r)
		return
	}

	response.JSON(w, http.StatusOK, metaJson)

}

func (app *application) mintHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserFriendlyAddress       *string                 `json:"user_address"`
		CollectionFriendlyAddress string                  `json:"collection_address"`
		MetaJSONID                string                  `json:"meta_json_id"`
		CSVFile                   []database.FileListItem `json:"csv_file"`
	}

	err := request.DecodeJSON(w, r, &input)

	if err != nil {
		app.logger.Error(err, nil)
		app.badRequest(w, r, err)
		return
	}

	// metajson_id looks like this: "33: TON Lisbon HUB"
	// ID: NAME
	// extract ID by first colon
	metajsonString := strings.Split(input.MetaJSONID, ":")[0]
	// convert to int64
	metaJsonID, err := strconv.ParseInt(metajsonString, 10, 64)
	if err != nil {
		app.logger.Error(err, nil)
		app.badRequest(w, r, err)
		return
	}
	

	collectionAddr := address.MustParseAddr(input.CollectionFriendlyAddress)
	collection := nft.NewCollectionClient(app.tonLiteClient, collectionAddr)

	collectionData, err := collection.GetCollectionData(context.Background())
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	// convert to int64
	// get meta json by base64
	metaJson, err := app.sqlModels.Nfts.GetNFTMetadataByID(metaJsonID)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	// // insert meta json to db
	err = app.sqlModels.Nfts.UpdateAttributesMetadata(metaJson)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	if input.CollectionFriendlyAddress == "" {
		app.badRequest(w, r, errors.New("wallet is required"))
		return
	}

	if input.CSVFile == nil && input.UserFriendlyAddress == nil {
		app.badRequest(w, r, errors.New("csv_file_url or user_address is required"))
		return
	}

	if input.CSVFile != nil {

		// make request to csv file
		resp, err := http.Get(input.CSVFile[0].Response.URL)
		if err != nil {
			app.logger.Error(err, nil)
			app.serverError(w, r, err) 
			 app.logger.Error(err,nil)
			return
		}

		defer resp.Body.Close()

		// read csv file. it has format: address
		reader := csv.NewReader(resp.Body)
		reader.Comma = ','
		reader.FieldsPerRecord = 1

		// read all lines
		lines, err := reader.ReadAll()
		if err != nil {
			app.logger.Error(err, nil)
			app.serverError(w, r, err) 
			 app.logger.Error(err,nil)
			return
		}

		// check if csv file is empty
		if len(lines) == 0 {
			app.badRequest(w, r, errors.New("csv file is empty"))
			return
		}

		offchainCon := &nft.ContentOffchain{
			URI: metaJson.Base64 + "/meta.json",
		}

		con := cell.BeginCell().MustStoreStringSnake(offchainCon.URI).EndCell()

		// Adjust the loop to only process a maximum of 200 lines
		if len(lines) > 100 {
			app.errorMessage(w, r, http.StatusBadRequest, "csv file has more than 100 lines", nil)
			return
		}

		dict := cell.NewDict(64)

		// print all lines except first
		for i, addr := range lines {

			dict.Set(cell.BeginCell().MustStoreUInt(collectionData.NextItemIndex.Uint64()+uint64(i), 64).EndCell(), cell.BeginCell().
				MustStoreCoins(tlb.MustFromTON("0.04").NanoTON().Uint64()).
				MustStoreRef(
					cell.BeginCell().
						MustStoreAddr(address.MustParseAddr(addr[0])). // owner
						MustStoreRef(con).
						MustStoreAddr(address.MustParseAddr(app.config.Ton.AdminWallet)). // editor address: admin wallet
						EndCell()).
				EndCell())

		}

		dataCell := cell.BeginCell().
			MustStoreUInt(2, 32).             // op code for mint batch
			MustStoreUInt(rand.Uint64(), 64). // query id
			MustStoreRef(dict.MustToCell()).
			EndCell()

		base64SignedDeployMsg := base64.StdEncoding.EncodeToString(dataCell.ToBOC())

		response.JSON(w, http.StatusOK, map[string]interface{}{
			"collection_address": collectionAddr,
			"msg_body":           base64SignedDeployMsg,
			"fee_for_tx":         tlb.MustFromTON(fmt.Sprint(0.06 * float64(len(lines)))).String(),
		})
		return
	}


	dict := cell.NewDict(64)

	offchainCon := &nft.ContentOffchain{
		URI: metaJson.Base64 + "/meta.json",
	}

	// get user byt username

	userforMint, err := app.sqlModels.Users.GetByUsername(*input.UserFriendlyAddress)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	con := cell.BeginCell().MustStoreStringSnake(offchainCon.URI).EndCell()

	dict.Set(cell.BeginCell().MustStoreUInt(collectionData.NextItemIndex.Uint64(), 64).EndCell(), cell.BeginCell().
		MustStoreCoins(tlb.MustFromTON("0.025").NanoTON().Uint64()).
		MustStoreRef(
			cell.BeginCell().
				MustStoreAddr(address.MustParseAddr(userforMint.FriendlyAddress)). // owner
				MustStoreRef(con).
				MustStoreAddr(address.MustParseAddr(app.config.Ton.AdminWallet)). // editor address
				EndCell()).
		EndCell())

	dataCell := cell.BeginCell().
		MustStoreUInt(2, 32).             // op code for mint batch
		MustStoreUInt(rand.Uint64(), 64). // query id
		MustStoreRef(dict.MustToCell()).
		EndCell()

	base64SignedDeployMsg := base64.URLEncoding.EncodeToString(dataCell.ToBOC())

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"collection_address": input.CollectionFriendlyAddress,
		"user_address":       input.UserFriendlyAddress,
		"msg_body":           base64SignedDeployMsg,
	})

}
func (app *application) insertPrototypeNft(w http.ResponseWriter, r *http.Request) {

	var input struct {
		DisplayName string                       `json:"display_name"`
		Description string                       `json:"description"`
		Image       []database.FileListItem      `json:"image"`
		Weight      string                       `json:"weight"`
		Attributes  *database.AttributeJsonArray `json:"attributes"`
	}

	err := request.DecodeJSON(w, r, &input)

	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	if input.DisplayName == "" {
		app.badRequest(w, r, errors.New("display name is required"))
		return
	}

	if input.Description == "" {
		app.badRequest(w, r, errors.New("description is required"))
		return
	}

	var attributes database.AttributeJsonArray

	if input.Attributes != nil {
		attributes = *input.Attributes
	}

	var weight int64 = 0

	if input.Weight != "" {
		weight, err = strconv.ParseInt(input.Weight, 10, 64)
		if err != nil {
			app.logger.Error(err, nil)
			app.serverError(w, r, err) 
			 app.logger.Error(err,nil)
			return
		}
	}

	base64Data := base64.StdEncoding.EncodeToString([]byte(input.DisplayName + input.Description + strconv.FormatInt(time.Now().Unix(), 10)))

	// url safe base64
	base64Data = strings.ReplaceAll(base64Data, "+", "-")
	base64Data = strings.ReplaceAll(base64Data, "/", "_")

	metaJson := &database.NFTMetadata{
		Base64:      base64Data,
		Name:        input.DisplayName,
		Description: input.Description,
		Image:       input.Image[0].Response.URL,
		Attributes:  attributes,
		ExternalURL: app.config.App.BaseUrl + DEPLOYED_NFT_POSTFIX + base64Data + "/meta.json",
		Marketplace: app.config.App.DomainName,
	}

	// // insert meta json to db
	err = app.sqlModels.Nfts.InsertNFTMetadata(metaJson)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	sbtPrototype := &database.SBTPrototype{
		MetadataID: metaJson.ID,
		Weight:     weight,
	}

	err = app.sqlModels.Nfts.InsertPrototype(sbtPrototype)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	response.JSON(w, http.StatusOK, sbtPrototype)
}

// deleteCollectionHandler deletes a collection from the database.
func (app *application) deleteCollectionHandler(w http.ResponseWriter, r *http.Request) {

	id := flow.Param(r.Context(), "id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	err = app.sqlModels.Nfts.DeleteCollection(idInt64)
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

// deleteTokenHandler deletes a token from the database.
func (app *application) deleteTokenHandler(w http.ResponseWriter, r *http.Request) {

	id := flow.Param(r.Context(), "id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	err = app.sqlModels.Nfts.DeleteToken(idInt64)
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

func (app *application) deletePrototypeHandler(w http.ResponseWriter, r *http.Request) {

	id := flow.Param(r.Context(), "id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	err = app.sqlModels.Nfts.DeletePrototype(idInt64)
	if err != nil {
		app.logger.Error(err, nil)
		return
	}

	response.JSON(w, http.StatusOK, nil)
}

// edit
func (app *application) updateCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// get by id
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	collection, err := app.sqlModels.Nfts.GetCollectionById(idInt64)
	if err != nil {
		app.errorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	var input struct {
		RawAddress      *string        `db:"raw_address" json:"raw_address"`
		FriendlyAddress *string        `db:"friendly_address" json:"friendly_address"`
		Name            *string        `db:"name" json:"name"`
		Description     *string        `db:"description" json:"description"`
		Image           *string        `db:"image" json:"image"`
		ContentJson     database.JSONB `db:"content_json" json:"content_json"`
		DefaultWeight   *int64         `db:"default_weight" json:"default_weight"`
		CreatedAt       *int64         `db:"created_at" json:"created_at"`
		UpdatedAt       *int64         `db:"updated_at" json:"updated_at"`
	}

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	metadata := &database.CollectionMetadata{}

	if input.RawAddress != nil {
		collection.RawAddress = *input.RawAddress
	}

	if input.FriendlyAddress != nil {
		collection.FriendlyAddress = *input.FriendlyAddress
	}

	if input.Name != nil {
		collection.Name = input.Name
		metadata.Name = *input.Name
	}

	if input.Description != nil {
		collection.Description = input.Description
		metadata.Description = *input.Description
	}

	if input.Image != nil {
		collection.Image = input.Image
		metadata.Image = *input.Image
	}

	if input.ContentJson != nil {
		collection.ContentJson = input.ContentJson
	}

	if input.DefaultWeight != nil {
		collection.DefaultWeight = *input.DefaultWeight
	}

	if input.CreatedAt != nil {
		collection.CreatedAt = *input.CreatedAt
	}

	if input.UpdatedAt != nil {
		collection.UpdatedAt = *input.UpdatedAt
	}

	collection, err = app.sqlModels.Nfts.UpdateCollection(collection)
	if err != nil {
		app.errorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// update collection metadata
	err = app.sqlModels.Nfts.UpdateCollectionMetadata(metadata, strings.Split(collection.ContentUri, "/")[6])
	if err != nil {
		app.errorMessage(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	err = response.JSON(w, http.StatusOK, collection)

	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
	}

}

func (app *application) updateTokenHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// get by id
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	nft, err := app.sqlModels.Nfts.GetTokenByID(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	var input struct {
		RawAddress           *string        `db:"raw_address" json:"raw_address"`
		FriendlyAddress      *string        `db:"friendly_address" json:"friendly_address"`
		ContentUri           *string        `db:"content_uri" json:"content_uri"`
		RawOwnerAddress      *string        `db:"raw_owner_address" json:"raw_owner_address"`
		FriendlyOwnerAddress *string        `db:"friendly_owner_address" json:"friendly_owner_address"`
		Name                 *string        `db:"name" json:"name"`
		Description          *string        `db:"description" json:"description"`
		Image                *string        `db:"image" json:"image"`
		ContentJson          database.JSONB `db:"content_json" json:"content_json"`
		Weight               *string         `db:"weight" json:"weight"`
	}
	

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if input.RawAddress != nil {
		nft.RawAddress = *input.RawAddress
	}

	if input.FriendlyAddress != nil {
		nft.FriendlyAddress = *input.FriendlyAddress
	}

	if input.ContentUri != nil {
		nft.ContentUri = *input.ContentUri
	}

	if input.RawOwnerAddress != nil {
		nft.RawOwnerAddress = *input.RawOwnerAddress
	}

	if input.FriendlyOwnerAddress != nil {
		nft.FriendlyOwnerAddress = *input.FriendlyOwnerAddress
	}

	if input.ContentJson != nil {
		nft.ContentJson = input.ContentJson
	}

	if input.Weight != nil {
		// convert weight to int64
		weightInt64, err := strconv.ParseInt(*input.Weight, 10, 64)
		if err != nil {
			app.badRequest(w, r, errors.New("weight must be an integer"))
			return
		}		
		nft.Weight = weightInt64
	}

	nft, err = app.sqlModels.Nfts.UpdateToken(nft)

	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	// update metadata
	err = response.JSON(w, http.StatusOK, nft)

	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
	}

}

func (app *application) updatePrototypeHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// get by id
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	nft, err := app.sqlModels.Nfts.GetNFTMetadataByPrototypeID(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Image       *string `json:"image"`
		Weight      *int64  `json:"weight"`
		Attributes  *database.AttributeJsonArray `json:"attributes"`
	}

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}
	

	if input.Name != nil {
		nft.Name = *input.Name
	}

	if input.Description != nil {
		nft.Description = *input.Description
	}

	if input.Image != nil {
		nft.Image = *input.Image
	}

	if input.Attributes != nil {
		nft.Attributes = *input.Attributes
	}

	if input.Weight != nil {
		err = app.sqlModels.Nfts.UpdatePrototypeWeight(idInt64, *input.Weight)
	}

	err = app.sqlModels.Nfts.UpdateNFTMetadata(nft)

	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, nft)

	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
	}

}

// get collection
func (app *application) getCollectionHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// get by id
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	collection, err := app.sqlModels.Nfts.GetCollectionById(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, collection)

	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
	}

}

// get token

func (app *application) getTokenHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// get by id
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	token, err := app.sqlModels.Nfts.GetTokenByID(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, token)

	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
	}

}

func (app *application) getMetadataNftsHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// get by id
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	token, err := app.sqlModels.Nfts.GetNFTMetadataByID(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, token)

	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
	}

}

func (app *application) getMetaJsonNFT(w http.ResponseWriter, r *http.Request) {
	base64String := flow.Param(r.Context(), "base64")

	// get meta json by base64
	metaJson, err := app.sqlModels.Nfts.GetNFTMetadataByBase64(base64String)
	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, metaJson)

	if err != nil {
		app.serverError(w, r, err) 
		 app.logger.Error(err,nil)
	}

}
