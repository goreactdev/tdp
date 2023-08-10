package database

import "github.com/jmoiron/sqlx"

type Models struct {
	Users 	   UserModel
	Nfts  NftsModel
	Tokens TokensModel
	Activities ActivitiesModel
	Permissions PermissionModel
	Rewards RewardModel
}

func NewModels(db *sqlx.DB) Models {
	return Models{
		Users: 	   UserModel{DB: db},
		Nfts: NftsModel{DB: db},
		Tokens: TokensModel{DB: db},
		Activities: ActivitiesModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Rewards: RewardModel{DB: db},
	}
}
