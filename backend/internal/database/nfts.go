package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type NftsModel struct {
	DB *sqlx.DB
}

type JSONB map[string]interface{}

func (a JSONB) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

const (
	SBT_COLLECTIONS_TABLE = "sbt_collections"
	SBT_TOKENS_TABLE      = "sbt_tokens"
	SBT_PROTOTYPE_TABLE      = "sbt_prototype"

)

type SBTCollection struct {
	ID                   int64   `db:"id" json:"id"`
	RawAddress           string  `db:"raw_address" json:"raw_address"`
	FriendlyAddress      string  `db:"friendly_address" json:"friendly_address"`
	NextItemIndex        int64   `db:"next_item_index" json:"next_item_index"`
	ContentUri           string  `db:"content_uri" json:"content_uri"`
	RawOwnerAddress      string  `db:"raw_owner_address" json:"raw_owner_address"`
	FriendlyOwnerAddress string  `db:"friendly_owner_address" json:"friendly_owner_address"`
	Name                 *string `db:"name" json:"name"`
	Description          *string `db:"description" json:"description"`
	Image                *string `db:"image" json:"image"`
	ContentJson          JSONB   `db:"content_json" json:"content_json"`
	DefaultWeight        int64   `db:"default_weight" json:"default_weight"`
	CreatedAt            int64   `db:"created_at" json:"created_at"`
	UpdatedAt            int64   `db:"updated_at" json:"updated_at"`
	Version              int64   `db:"version" json:"version"`
}

type Metadata struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}

type SBTToken struct {
	ID                   int64   `db:"id" json:"id"`
	RawAddress           string  `db:"raw_address" json:"raw_address"`
	FriendlyAddress      string  `db:"friendly_address" json:"friendly_address"`
	SBTCollectionID      int64   `db:"sbt_collections_id " json:"sbt_collections_id"`
	ContentUri           string  `db:"content_uri" json:"content_uri"`
	RawOwnerAddress      string  `db:"raw_owner_address" json:"raw_owner_address"`
	FriendlyOwnerAddress string  `db:"friendly_owner_address" json:"friendly_owner_address"`
	IsPinned 			bool 	`db:"is_pinned" json:"is_pinned"`
	Name                 *string `db:"name" json:"name"`
	Description          *string `db:"description" json:"description"`
	Image                *string `db:"image" json:"image"`
	ContentJson          JSONB   `db:"content_json" json:"content_json"`
	Weight               int64   `db:"weight" json:"weight"`
	Index                int64   `db:"index" json:"index"`
	CreatedAt            int64   `db:"created_at" json:"created_at"`
	UpdatedAt            int64   `db:"updated_at" json:"updated_at"`
	Version              int64   `db:"version" json:"version"`
}

// get sbt token by content uri
func (m *NftsModel) GetSbtTokenByContentUri(contentUri string) (*SBTToken, error) {
	var sbtToken SBTToken

	err := m.DB.Get(&sbtToken, "SELECT * FROM sbt_tokens WHERE content_uri = $1", contentUri)
	if err != nil {
		return nil, err
	}

	return &sbtToken, nil
}

type File struct {
	UID string `json:"uid"`
}

type OriginFileObj struct {
	UID string `json:"uid"`
}


// CREATE TABLE IF NOT EXISTS sbt_prototype (
//     id BIGSERIAL NOT NULL PRIMARY KEY,
//     name TEXT NOT NULL,
//     description TEXT,
//     image TEXT,
//     weight INTEGER NOT NULL,
//     created_at BIGINT NOT NULL,
//     updated_at BIGINT NOT NULL,
//     version BIGINT NOT NULL DEFAULT 1
// );


type CollectionMetadata struct {
	ID           int64   `db:"id" json:"id"`
	Base64 	 string  `db:"base64" json:"base64"`
	Name         string  `db:"name" json:"name"`
	Description  string  `db:"description" json:"description"`
	Image        string  `db:"image" json:"image"`
	CoverImage   *string  `db:"cover_image" json:"cover_image"`
	ExternalURL  string  `db:"external_url" json:"external_url"`
	Marketplace  string  `db:"marketplace" json:"marketplace"`
	CreatedAt    int64   `db:"created_at" json:"created_at"`
	UpdatedAt    int64   `db:"updated_at" json:"updated_at"`
	Version      int64   `db:"version" json:"version"`
}


type AttributeJson struct {
	TraitType string `json:"trait_type"`
	ValueN     string `json:"value"`
}



type AttributeJsonArray []AttributeJson


func (sla AttributeJsonArray) Value() (driver.Value, error) {
    // implement the Value method for inserting into the database
	return json.Marshal(sla)
}

func (sla *AttributeJsonArray) Scan(value interface{}) error {
    var data = []byte(value.([]uint8))
    return json.Unmarshal(data, &sla)
}


type FileListItem struct {
	LastModified     int64         `json:"lastModified"`
	LastModifiedDate string        `json:"lastModifiedDate"`
	Name             string        `json:"name"`
	OriginFileObj    OriginFileObj `json:"originFileObj"`
	Percent          int64         `json:"percent"`
	Response         struct {
		URL string `json:"url"`
	} `json:"response"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	Status   string `json:"status"`
	ThumbURL string `json:"thumbUrl"`
	Type     string `json:"type"`
	UID      string `json:"uid"`
	Xhr      struct {
	} `json:"xhr"`
}

type ImageData struct {
	File      File           `json:"file"`
	FileList  []FileListItem `json:"fileList"`
}


// insert sbt token into database
func (m *NftsModel) InsertToken(tx *sqlx.Tx, token *SBTToken) error {
	query := `
		INSERT INTO sbt_tokens (raw_address, friendly_address, sbt_collections_id, content_uri, raw_owner_address, friendly_owner_address, name, description, image, content_json, weight, index, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, 1)
		RETURNING id, raw_address, friendly_address, sbt_collections_id, content_uri, raw_owner_address, friendly_owner_address, is_pinned, name, description, image, content_json, weight, index, created_at, updated_at, version	
		`

	row := tx.QueryRow(query, token.RawAddress, token.FriendlyAddress, token.SBTCollectionID, token.ContentUri, token.RawOwnerAddress, token.FriendlyOwnerAddress, token.Name, token.Description, token.Image, token.ContentJson, token.Weight, token.Index, time.Now().Unix(), time.Now().Unix())

	err := row.Scan(
		&token.ID,
		&token.RawAddress,
		&token.FriendlyAddress,
		&token.SBTCollectionID,
		&token.ContentUri,
		&token.RawOwnerAddress,
		&token.FriendlyOwnerAddress,
		&token.IsPinned,
		&token.Name,
		&token.Description,
		&token.Image,
		&token.ContentJson,
		&token.Weight,
		&token.Index,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.Version,
	)

	if err != nil {
		return err
	}

	return nil
}


type SBTPrototype struct {
	ID          int64   `db:"id" json:"id"`
	MetadataID  int64   `db:"metadata_id" json:"metadata_id"`
	Metadata  NFTMetadata `db:"metadata" json:"metadata"`
	Weight      int64   `db:"weight" json:"weight"`
	CreatedAt   int64   `db:"created_at" json:"created_at"`
	UpdatedAt   int64   `db:"updated_at" json:"updated_at"`
	Version     int64   `db:"version" json:"version"`
}

// update weight of prototype
func (m *NftsModel) UpdatePrototypeWeight(id int64, weight int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE sbt_prototype
		SET weight = $1
		WHERE id = $2
		`

	_, err := m.DB.ExecContext(ctx, query, weight, id)
	if err != nil {
		return err
	}

	return nil
}

// delete prototype
func (m *NftsModel) DeletePrototype(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		DELETE FROM activities
		WHERE  sbt_prototype_id = $1	
		`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

   query = `
		DELETE FROM sbt_prototype
		WHERE id = $1
		`

	_, err = m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}


// insert sbt prototype into database
func (m *NftsModel) InsertPrototype(prototype *SBTPrototype) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO sbt_prototype (metadata_id, weight, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, 1)
		RETURNING *`

	row := m.DB.QueryRowContext(ctx, query, prototype.MetadataID, prototype.Weight, time.Now().Unix(), time.Now().Unix())

	err := row.Scan(
		&prototype.ID,
		&prototype.MetadataID,
		&prototype.Weight,
		&prototype.CreatedAt,
		&prototype.UpdatedAt,
		&prototype.Version,
	)

	if err != nil {
		return err
	}

	return nil
}




	
// // get tokens from database
// func (m *NftsModel) GetTokens(start, end int, filters, sort string) ([]*SBTToken, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
// 	defer cancel()

// 	query := fmt.Sprintf(`
// 		SELECT *
// 		FROM sbt_tokens
// 		%s
// 		ORDER BY %s
// 		LIMIT $1 OFFSET $2
// 		`, filters, sort)

// 	rows, err := m.DB.QueryContext(ctx, query, end-start, start)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()

// 	tokens := []*SBTToken{}

// 	for rows.Next() {
// 		var token SBTToken

// 		err := rows.Scan(
// 			&token.ID,
// 			&token.RawAddress,
// 			&token.FriendlyAddress,
// 			&token.SBTCollectionID,
// 			&token.ContentUri,
// 			&token.RawOwnerAddress,
// 			&token.FriendlyOwnerAddress,
// 			&token.Name,
// 			&token.Description,
// 			&token.Image,
// 			&token.ContentJson,
// 			&token.Weight,
// 			&token.Index,
// 			&token.CreatedAt,
// 			&token.UpdatedAt,
// 			&token.Version,
// 		)

// 		if err != nil {
// 			return nil, err
// 		}

// 		tokens = append(tokens, &token)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return tokens, nil

// }

// get prototypes from database
func (m *NftsModel) GetPrototypes(pagination *Pagination, filters, sort string) ([]*SBTPrototype, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT sbt_prototype.*
		FROM sbt_prototype
		LEFT JOIN nft_metadata ON sbt_prototype.metadata_id = nft_metadata.id
		%s
		ORDER BY %s
		LIMIT $1 OFFSET $2
		`, filters, sort)

	rows, err := m.DB.QueryContext(ctx, query, pagination.End-pagination.Start, pagination.Start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	prototypes := []*SBTPrototype{}

	for rows.Next() {
		var prototype SBTPrototype

		err := rows.Scan(
			&prototype.ID,
			&prototype.MetadataID,
			&prototype.Weight,
			&prototype.CreatedAt,
			&prototype.UpdatedAt,
			&prototype.Version,
		)

		if err != nil {
			return nil, err
		}

		prototypes = append(prototypes, &prototype)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prototypes, nil
}








type NFTMetadata struct {
	ID           int64   `db:"id" json:"id"`
	Base64 	 string  `db:"base64" json:"base64"`
	Name         string  `db:"name" json:"name"`
	Description  string  `db:"description" json:"description"`
	Attributes   AttributeJsonArray   `db:"attributes" json:"attributes"`
	ExternalURL  string  `db:"external_url" json:"external_url"`
	Image        string  `db:"image" json:"image"`
	Marketplace  string  `db:"marketplace" json:"marketplace"`
	CreatedAt    int64   `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt    int64   `db:"updated_at" json:"updated_at,omitempty"`
	Version      int64   `db:"version" json:"version,omitempty"`
}

func (m *NftsModel) InsertNFTMetadata(metadata *NFTMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO nft_metadata (base64, name, description, attributes, external_url, image, marketplace, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1)
		RETURNING *`

	row := m.DB.QueryRowContext(ctx, query, metadata.Base64, metadata.Name, metadata.Description, metadata.Attributes, metadata.ExternalURL, metadata.Image, metadata.Marketplace, time.Now().Unix(), time.Now().Unix())

	err := row.Scan(
		&metadata.ID,
		&metadata.Base64,
		&metadata.Name,
		&metadata.Description,
		&metadata.Attributes,
		&metadata.ExternalURL,
		&metadata.Image,
		&metadata.Marketplace,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
		&metadata.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

// update nft metadata
func (m *NftsModel) UpdateNFTMetadata(metadata *NFTMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE nft_metadata
		SET base64 = $1, name = $2, description = $3, attributes = $4, external_url = $5, image = $6, marketplace = $7, updated_at = $8
		WHERE id = $9
		RETURNING *`

	row := m.DB.QueryRowContext(ctx, query, metadata.Base64, metadata.Name, metadata.Description, metadata.Attributes, metadata.ExternalURL, metadata.Image, metadata.Marketplace, time.Now().Unix(), metadata.ID)

	err := row.Scan(
		&metadata.ID,
		&metadata.Base64,
		&metadata.Name,
		&metadata.Description,
		&metadata.Attributes,
		&metadata.ExternalURL,
		&metadata.Image,
		&metadata.Marketplace,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
		&metadata.Version,
	)

	if err != nil {
		return err
	}

	return nil
}




// type CollectionMetadata struct {
// 	ID           int64   `db:"id" json:"id"`
// 	Base64 	 string  `db:"base64" json:"base64"`
// 	Name         string  `db:"name" json:"name"`
// 	Description  string  `db:"description" json:"description"`
// 	Image        string  `db:"image" json:"image"`
// 	CoverImage   *string  `db:"cover_image" json:"cover_image"`
// 	ExternalURL  string  `db:"external_url" json:"external_url"`
// 	Marketplace  string  `db:"marketplace" json:"marketplace"`
// 	CreatedAt    int64   `db:"created_at" json:"created_at"`
// 	UpdatedAt    int64   `db:"updated_at" json:"updated_at"`
// 	Version      int64   `db:"version" json:"version"`
// }

// update collection metadata: name, description and image
func (m *NftsModel) UpdateCollectionMetadata(metadata *CollectionMetadata, base64 string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE collection_metadata
		SET name = $1, description = $2, image = $3, updated_at = $4
		WHERE base64 = $5
		`

	_, err := m.DB.ExecContext(ctx, query, metadata.Name, metadata.Description, metadata.Image, time.Now().Unix(), base64)

	if err != nil {
		return err
	}

	return nil
}

// insert collection metadata
func (m *NftsModel) InsertCollectionMetadata(metadata *CollectionMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO collection_metadata (base64, name, description, image, cover_image, external_url, marketplace, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1)
		RETURNING *`
		
	row := m.DB.QueryRowContext(ctx, query, metadata.Base64, metadata.Name, metadata.Description, metadata.Image, metadata.CoverImage, metadata.ExternalURL, metadata.Marketplace, time.Now().Unix(), time.Now().Unix())

	err := row.Scan(
		&metadata.ID,
		&metadata.Base64,
		&metadata.Name,
		&metadata.Description,
		&metadata.Image,
		&metadata.CoverImage,
		&metadata.ExternalURL,
		&metadata.Marketplace,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
		&metadata.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

// get collection metadata by base64
func (m *NftsModel) GetCollectionMetadataByBase64(base64 string) (*CollectionMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT *
		FROM collection_metadata
		WHERE base64 = $1`


	row := m.DB.QueryRowContext(ctx, query, base64)

	var metadata CollectionMetadata

	err := row.Scan(
		&metadata.ID,
		&metadata.Base64,
		&metadata.Name,
		&metadata.Description,
		&metadata.Image,
		&metadata.CoverImage,
		&metadata.ExternalURL,
		&metadata.Marketplace,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
		&metadata.Version,
	)

	if err != nil {
		return nil, err
	}

	return &metadata, nil
}




// update metadata
func (m *NftsModel) UpdateAttributesMetadata(metadata *NFTMetadata) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE nft_metadata
		SET attributes = $1, updated_at = $2, version = version + 1
		WHERE id = $3
		RETURNING *`
	
	row := m.DB.QueryRowContext(ctx, query, metadata.Attributes, time.Now().Unix(), metadata.ID)

	err := row.Scan(
		&metadata.ID,
		&metadata.Base64,
		&metadata.Name,
		&metadata.Description,
		&metadata.Attributes,
		&metadata.ExternalURL,
		&metadata.Image,
		&metadata.Marketplace,
		&metadata.CreatedAt,
		&metadata.UpdatedAt,
		&metadata.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

		

// get metadata by base64
func (m *NftsModel) GetNFTMetadataByBase64(base64String string) (*NFTMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT id, base64, name, description, attributes, external_url, image, marketplace
		FROM nft_metadata
		WHERE base64 = $1
		LIMIT 1
	`

	row := m.DB.QueryRowContext(ctx, query, base64String)

	var metadata NFTMetadata

	err := row.Scan(
		&metadata.ID,
		&metadata.Base64,
		&metadata.Name,
		&metadata.Description,
		&metadata.Attributes,
		&metadata.ExternalURL,
		&metadata.Image,
		&metadata.Marketplace,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &metadata, nil
}

// get prototype weight by base64 of nft
func (m *NftsModel) GetWeightByBase64(base64String string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT weight
		FROM sbt_prototype
		WHERE metadata_id = (SELECT id FROM nft_metadata WHERE base64 = $1)
		LIMIT 1
		`
	
	row := m.DB.QueryRowContext(ctx, query, base64String)

	var weight int64

	err := row.Scan(
		&weight,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return weight, nil

}

func (m *NftsModel) GetNFTMetadataByPrototypeID(prototypeID int64) (*NFTMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT m.id, m.base64, m.name, m.description, m.attributes, m.external_url, m.image, m.marketplace
		FROM nft_metadata m
		INNER JOIN sbt_prototype p ON p.metadata_id = m.id
		WHERE p.id = $1
		LIMIT 1
	`

	row := m.DB.QueryRowContext(ctx, query, prototypeID)

	var metadata NFTMetadata

	err := row.Scan(
		&metadata.ID,
		&metadata.Base64,
		&metadata.Name,
		&metadata.Description,
		&metadata.Attributes,
		&metadata.ExternalURL,
		&metadata.Image,
		&metadata.Marketplace,
	)

	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

// get prototype by user rating
func (m *NftsModel) GetPrototypesByRating(userId int64) ([]*NFTMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
	SELECT nm.*
	FROM nft_metadata nm
	JOIN sbt_prototype sp ON nm.ID = sp.metadata_id
	JOIN activities a ON sp.id = a.id
	JOIN users u ON u.id = $1
	WHERE u.rating >= a.token_threshold AND NOT EXISTS (
		SELECT 1 
		FROM sbt_tokens
		WHERE friendly_owner_address = u.friendly_address
		  AND (content_json->>'id')::NUMERIC = nm.id
	) AND NOT EXISTS (
		SELECT 1
		FROM stored_rewards
		WHERE user_address = u.friendly_address
		  AND base64_metadata = nm.base64
	)
	ORDER BY a.token_threshold DESC
	`

	rows, err := m.DB.QueryContext(ctx, query, userId)
	
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var prototypes []*NFTMetadata

	for rows.Next() {
		var metadata NFTMetadata

		err := rows.Scan(
			&metadata.ID,
			&metadata.Base64,
			&metadata.Name,
			&metadata.Description,
			&metadata.Attributes,
			&metadata.ExternalURL,
			&metadata.Image,
			&metadata.Marketplace,
			&metadata.CreatedAt,
			&metadata.UpdatedAt,
			&metadata.Version,
		)

		if err != nil {
			return nil, err
		}

		prototypes = append(prototypes, &metadata)
	}

	if err = rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return prototypes, nil




}

// check if user has certain nft by checking content_json (jsonb) field

func (m *NftsModel) CheckIfUserHasNFTByContentJSON(addr string, metadataID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT COUNT(*) FROM sbt_tokens
		WHERE friendly_owner_address = $1 AND content_json->>'id' = $2
	`

	row := m.DB.QueryRowContext(ctx, query, addr, metadataID)

	var count int64

	err := row.Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
		


// get nft metadata by id
func (m *NftsModel) GetNFTMetadataByID(id int64) (*NFTMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT id, base64, name, description, attributes, external_url, image, marketplace
		FROM nft_metadata
		WHERE id = $1
		LIMIT 1
	`

	row := m.DB.QueryRowContext(ctx, query, id)

	var metadata NFTMetadata

	err := row.Scan(
		&metadata.ID,
		&metadata.Base64,
		&metadata.Name,
		&metadata.Description,
		&metadata.Attributes,
		&metadata.ExternalURL,
		&metadata.Image,
		&metadata.Marketplace,
	)

	if err != nil {
		return nil, err
	}

	return &metadata, nil
}



// insert collection into database
func (m *NftsModel) InsertCollection(collection *SBTCollection) (*SBTCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO sbt_collections (raw_address, friendly_address, next_item_index, content_uri, raw_owner_address, friendly_owner_address, name, description, image, content_json, default_weight, created_at, updated_at, version)	
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING *
		`

	row := m.DB.QueryRowContext(ctx, query, collection.RawAddress, collection.FriendlyAddress, collection.NextItemIndex, collection.ContentUri, collection.RawOwnerAddress, collection.FriendlyOwnerAddress, collection.Name, collection.Description, collection.Image, collection.ContentJson, collection.DefaultWeight, time.Now().Unix(), time.Now().Unix(), 1)

	err := row.Scan(
		&collection.ID,
		&collection.RawAddress,
		&collection.FriendlyAddress,
		&collection.NextItemIndex,
		&collection.ContentUri,
		&collection.RawOwnerAddress,
		&collection.FriendlyOwnerAddress,
		&collection.Name,
		&collection.Description,
		&collection.Image,
		&collection.ContentJson,
		&collection.DefaultWeight,
		&collection.CreatedAt,
		&collection.UpdatedAt,
		&collection.Version,
	)

	if err != nil {
		return nil, err
	}

	return collection, nil
}
	

func (m *NftsModel) GetCollectionsByOwnerAddress(pagination *Pagination, ownerAddr string) ([]*SBTCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT *
		FROM sbt_collections
		WHERE friendly_owner_address = $3
		ORDER BY id
		LIMIT $1 OFFSET $2
		`)

	rows, err := m.DB.QueryContext(ctx, query, pagination.End-pagination.Start, pagination.Start, ownerAddr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collections := []*SBTCollection{}

	for rows.Next() {
		var collection SBTCollection

		err := rows.Scan(
			&collection.ID,
			&collection.RawAddress,
			&collection.FriendlyAddress,
			&collection.NextItemIndex,
			&collection.ContentUri,
			&collection.RawOwnerAddress,
			&collection.FriendlyOwnerAddress,
			&collection.Name,
			&collection.Description,
			&collection.Image,
			&collection.ContentJson,
			&collection.DefaultWeight,
			&collection.CreatedAt,
			&collection.UpdatedAt,
			&collection.Version,
		)

		if err != nil {
			return nil, err
		}

		collections = append(collections, &collection)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}
func (m *NftsModel) GetAllCollections(pagination *Pagination) ([]*SBTCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT *
		FROM sbt_collections
		ORDER BY id
		LIMIT $1 OFFSET $2
		`)

	rows, err := m.DB.QueryContext(ctx, query, pagination.End-pagination.Start, pagination.Start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	collections := []*SBTCollection{}

	for rows.Next() {
		var collection SBTCollection

		err := rows.Scan(
			&collection.ID,
			&collection.RawAddress,
			&collection.FriendlyAddress,
			&collection.NextItemIndex,
			&collection.ContentUri,
			&collection.RawOwnerAddress,
			&collection.FriendlyOwnerAddress,
			&collection.Name,
			&collection.Description,
			&collection.Image,
			&collection.ContentJson,
			&collection.DefaultWeight,
			&collection.CreatedAt,
			&collection.UpdatedAt,
			&collection.Version,
		)

		if err != nil {
			return nil, err
		}

		collections = append(collections, &collection)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}


// get collection by address

func (m *NftsModel) GetCollectionByAddress(address string) (*SBTCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT *
		FROM sbt_collections
		WHERE friendly_address = $1
		`

	row := m.DB.QueryRowContext(ctx, query, address)

	var collection SBTCollection

	err := row.Scan(
		&collection.ID,
		&collection.RawAddress,
		&collection.FriendlyAddress,
		&collection.NextItemIndex,
		&collection.ContentUri,
		&collection.RawOwnerAddress,
		&collection.FriendlyOwnerAddress,
		&collection.Name,
		&collection.Description,
		&collection.Image,
		&collection.ContentJson,
		&collection.DefaultWeight,
		&collection.CreatedAt,
		&collection.UpdatedAt,
		&collection.Version,
	)

	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// get collection by nft address
func (m *NftsModel) GetCollectionByNftAddress(address string) (*SBTCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT *
		FROM sbt_collections
		WHERE id = (
			SELECT sbt_collections_id FROM sbt_tokens WHERE friendly_address = $1
		)
		`

	row := m.DB.QueryRowContext(ctx, query, address)

	var collection SBTCollection

	err := row.Scan(
		&collection.ID,
		&collection.RawAddress,
		&collection.FriendlyAddress,
		&collection.NextItemIndex,
		&collection.ContentUri,
		&collection.RawOwnerAddress,
		&collection.FriendlyOwnerAddress,
		&collection.Name,
		&collection.Description,
		&collection.Image,
		&collection.ContentJson,
		&collection.DefaultWeight,
		&collection.CreatedAt,
		&collection.UpdatedAt,
		&collection.Version,
	)

	if err != nil {
		return nil, err
	}

	return &collection, nil
}


// get tokens from database
func (m *NftsModel) GetTokens(pagination *Pagination, addr string) ([]*SBTToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var rows *sql.Rows
	var err error

	if addr == "" {
		query := `
		SELECT *
		FROM sbt_tokens
		ORDER BY id
		LIMIT $1 OFFSET $2
		`
		rows, err = m.DB.QueryContext(ctx, query, pagination.End-pagination.Start, pagination.Start)
		if err != nil {
			return nil, err
		}
	} else {
		query := `
		SELECT *
		FROM sbt_tokens
		WHERE friendly_owner_address = $3
		ORDER BY id
		LIMIT $1 OFFSET $2
		`
		rows, err = m.DB.QueryContext(ctx, query, pagination.End-pagination.Start, pagination.Start, addr)
		if err != nil {
			return nil, err
		}
	}


	defer rows.Close()

	tokens := []*SBTToken{}

	for rows.Next() {
		var token SBTToken		

		err := rows.Scan(
			&token.ID,
			&token.RawAddress,
			&token.FriendlyAddress,
			&token.SBTCollectionID,
			&token.ContentUri,
			&token.RawOwnerAddress,
			&token.FriendlyOwnerAddress,
			&token.IsPinned,
			&token.Name,
			&token.Description,
			&token.Image,
			&token.ContentJson,
			&token.Weight,
			&token.Index,
			&token.CreatedAt,
			&token.UpdatedAt,
			&token.Version,
		)

		if err != nil {
			return nil, err
		}

		tokens = append(tokens, &token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil

}


// get token by id
func (m *NftsModel) GetTokenByID(id int64) (*SBTToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT *
		FROM sbt_tokens
		WHERE id = $1
		`

	row := m.DB.QueryRowContext(ctx, query, id)

	var token SBTToken

	

	err := row.Scan(
		&token.ID,
		&token.RawAddress,
		&token.FriendlyAddress,
		&token.SBTCollectionID,
		&token.ContentUri,
		&token.RawOwnerAddress,
		&token.FriendlyOwnerAddress,
		&token.IsPinned,
		&token.Name,
		&token.Description,
		&token.Image,
		&token.ContentJson,
		&token.Weight,
		&token.Index,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.Version,		
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// get token by address
func (m *NftsModel) GetTokenByAddress(address string) (*SBTToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT *
		FROM sbt_tokens
		WHERE friendly_address = $1
		`

	row := m.DB.QueryRowContext(ctx, query, address)

	var token SBTToken


	err := row.Scan(
		&token.ID,
		&token.RawAddress,
		&token.FriendlyAddress,
		&token.SBTCollectionID,
		&token.ContentUri,
		&token.RawOwnerAddress,
		&token.FriendlyOwnerAddress,
		&token.IsPinned,
		&token.Name,
		&token.Description,
		&token.Image,
		&token.ContentJson,
		&token.Weight,
		&token.Index,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.Version,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

// check if user has sbt with name Auth NFT

func (m *NftsModel) HasAuthNFT(address, base64 string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	contentUri := fmt.Sprintf("https://tdp.tonbuilders.com/v1/deployed-nft/n/%s/meta.json", base64)
	

	query := `
		SELECT COUNT(*)
		FROM sbt_tokens
		WHERE friendly_owner_address = $1 AND content_uri = $2
		`

	row := m.DB.QueryRowContext(ctx, query, address, contentUri)

	var count int

	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// pin nft
func (m *NftsModel) PinNFT(id int64, pin bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE sbt_tokens
		SET is_pinned = $1
		WHERE id = $2
		`

	_, err := m.DB.ExecContext(ctx, query, pin, id)
	if err != nil {
		return err
	}

	return nil
}


// get sbt tokens by owner address 
func (m *NftsModel) GetTokensByOwnerAddress(address string, pagination *Pagination) ([]*SBTToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT *
		FROM sbt_tokens
		WHERE friendly_owner_address = $1
		ORDER BY is_pinned DESC, id DESC
		LIMIT $2 OFFSET $3
		`

	rows, err := m.DB.QueryContext(ctx, query, address, pagination.End-pagination.Start, pagination.Start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tokens := []*SBTToken{}



	for rows.Next() {
		var token SBTToken

		err := rows.Scan(
			&token.ID,
			&token.RawAddress,
			&token.FriendlyAddress,
			&token.SBTCollectionID,
			&token.ContentUri,
			&token.RawOwnerAddress,
			&token.FriendlyOwnerAddress,
			&token.IsPinned,
			&token.Name,
			&token.Description,
			&token.Image,
			&token.ContentJson,
			&token.Weight,
			&token.Index,
			&token.CreatedAt,
			&token.UpdatedAt,
			&token.Version,
		)

		if err != nil {
			return nil, err
		}

		tokens = append(tokens, &token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

// get date of last token created
func (m *NftsModel) GetLastTokenCreated(address string) (string, string, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT friendly_address, name, weight
		FROM sbt_tokens
		WHERE friendly_owner_address = $1
		ORDER BY created_at DESC
		LIMIT 1
		`

	row := m.DB.QueryRowContext(ctx, query, address)

	var name string
	var weight int64
	var friendlyAddress string

	err := row.Scan(&friendlyAddress, &name, &weight)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", 0, nil
		}
		return "", "", 0, err
	}

	return name, friendlyAddress, weight, nil
	
}

// get total number of tokens
func (m *NftsModel) GetTotalTokensByOwnerAddress(address string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT COUNT(*)	
		FROM sbt_tokens
		WHERE friendly_owner_address = $1
		`	

	row := m.DB.QueryRowContext(ctx, query, address)

	var total int

	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}



// get total number of tokens
func (m *NftsModel) GetTotal(table, column, filters string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := fmt.Sprintf(`
		SELECT COUNT(%s)
		FROM %s
		%s
		`, column, table, filters)

	row := m.DB.QueryRowContext(ctx, query)

	var total int

	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// delete collection from database
func (m *NftsModel) DeleteCollection(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// and delete all tokens from this collection
	query := `
		DELETE FROM sbt_tokens
		WHERE sbt_collections_id = $1
		`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	query = `
		DELETE FROM sbt_collections
		WHERE id = $1
		`

	_, err = m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// delete token from database
func (m *NftsModel) DeleteToken(id int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		DELETE FROM sbt_tokens
		WHERE id = $1
		`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// update collection in database
func (m *NftsModel) UpdateCollection(collection *SBTCollection) (*SBTCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE sbt_collections
		SET raw_address = $1, friendly_address = $2, next_item_index = $3, content_uri = $4, raw_owner_address = $5, friendly_owner_address = $6, name = $7, description = $8, image = $9, content_json = $10, default_weight = $11, updated_at = $12, version = $13
		WHERE id = $14
		RETURNING *
		`

	row := m.DB.QueryRowContext(ctx, query, collection.RawAddress, collection.FriendlyAddress, collection.NextItemIndex, collection.ContentUri, collection.RawOwnerAddress, collection.FriendlyOwnerAddress, collection.Name, collection.Description, collection.Image, collection.ContentJson, collection.DefaultWeight, time.Now().Unix(), collection.Version+1, collection.ID)

	err := row.Scan(
		&collection.ID,
		&collection.RawAddress,
		&collection.FriendlyAddress,
		&collection.NextItemIndex,
		&collection.ContentUri,
		&collection.RawOwnerAddress,
		&collection.FriendlyOwnerAddress,
		&collection.Name,
		&collection.Description,
		&collection.Image,
		&collection.ContentJson,
		&collection.DefaultWeight,
		&collection.CreatedAt,
		&collection.UpdatedAt,
		&collection.Version,
	)

	if err != nil {
		return nil, err
	}

	return collection, nil
}

// update next item index in database
func (m *NftsModel) UpdateNextItemIndex(collectionAddr string, nextItem int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE sbt_collections
		SET next_item_index = $1
		WHERE friendly_address = $2
		`

	_, err := m.DB.ExecContext(ctx, query, nextItem, collectionAddr)
	if err != nil {
		return err
	}

	return nil
	


}

// update token in database
func (m *NftsModel) UpdateToken(token *SBTToken) (*SBTToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		UPDATE sbt_tokens
		SET raw_address = $1, friendly_address = $2, content_uri = $3, raw_owner_address = $4, friendly_owner_address = $5, name = $6, description = $7, image = $8, content_json = $9, weight = $10, updated_at = $11, version = $12
		WHERE id = $13
		RETURNING *
		`

	row := m.DB.QueryRowContext(ctx, query, token.RawAddress, token.FriendlyAddress, token.ContentUri, token.RawOwnerAddress, token.FriendlyOwnerAddress, token.Name, token.Description, token.Image, token.ContentJson, token.Weight, time.Now().Unix(), token.Version+1, token.ID)

	err := row.Scan(
		&token.ID,
		&token.RawAddress,
		&token.FriendlyAddress,
		&token.SBTCollectionID,
		&token.ContentUri,
		&token.RawOwnerAddress,
		&token.FriendlyOwnerAddress,
		&token.IsPinned,
		&token.Name,
		&token.Description,
		&token.Image,
		&token.ContentJson,
		&token.Weight,
		&token.Index,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.Version,
)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// get collection by id
func (m *NftsModel) GetCollectionById(id int64) (*SBTCollection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		SELECT *
		FROM sbt_collections
		WHERE id = $1
		`

	row := m.DB.QueryRowContext(ctx, query, id)

	var collection SBTCollection

	err := row.Scan(
		&collection.ID,
		&collection.RawAddress,
		&collection.FriendlyAddress,
		&collection.NextItemIndex,
		&collection.ContentUri,
		&collection.RawOwnerAddress,
		&collection.FriendlyOwnerAddress,
		&collection.Name,
		&collection.Description,
		&collection.Image,
		&collection.ContentJson,
		&collection.DefaultWeight,
		&collection.CreatedAt,
		&collection.UpdatedAt,
		&collection.Version,
	)

	if err != nil {
		return nil, err
	}

	return &collection, nil
}

