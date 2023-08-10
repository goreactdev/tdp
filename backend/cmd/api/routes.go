package main

import (
	"net/http"

	"github.com/alexedwards/flow"
)

func (app *application) routes() http.Handler {
	mux := flow.New()

	mux.NotFound = http.HandlerFunc(app.notFound)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)

	mux.Use(app.recoverPanic)
	mux.Use(app.authenticate)
	mux.Use(app.enableCORS)

	mux.HandleFunc("/v1/status", app.status, "GET")

	mux.HandleFunc("/v1/ton-connect/generate-payload", app.payloadHandler, "GET")
	mux.HandleFunc("/v1/ton-connect/check-proof", app.proofHandler, "POST")
	mux.HandleFunc("/v1/manifest-ton-connect", app.manifestTonConnectHandler, "GET")

	// auth
	mux.HandleFunc("/v1/github/callback", app.githubCallbackHandler, "GET")
	mux.HandleFunc("/v1/github/login", app.githubLoginHandler, "GET")


	mux.HandleFunc("/v1/deployed-nft/n/:base64/meta.json", app.getMetaJsonNft, "GET")
	mux.HandleFunc("/v1/deployed-nft/c/:base64/meta.json", app.getMetaJsonCollection, "GET")
	
	mux.HandleFunc("/v1/users", app.getTopUsersHandler, "GET")
	mux.HandleFunc("/v1/users/:username", app.getUserByUsernameHandler, "GET")
	mux.HandleFunc("/v1/nfts/:username", app.getNftsByUserIdHandler, "GET")
	mux.HandleFunc("/v1/admin/csv/upload", app.uploadCSVHandler, "POST")
	mux.HandleFunc("/v1/admin/media/upload", app.uploadImageHandler, "POST")


	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requireAuthenticatedUser)
		mux.HandleFunc("/v1/telegram/check_authorization", app.checkTelegramAuthorization, "POST")

		mux.HandleFunc("/v1/my-account", app.getMyAccountHandler, "GET")

		// deploy and upload media
		mux.HandleFunc("/v1/update/users", app.updateUserHandler, "PATCH")
		

		mux.HandleFunc("/v1/nft/:id/pin", app.pinNftHandler, "PUT")
		mux.HandleFunc("/v1/unlink/:provider", app.unlinkAccountHandler, "DELETE")
	
		mux.HandleFunc("/v1/incoming-achievements", app.getIncomingAchievementsHandler, "GET")
		mux.HandleFunc("/v1/incoming-achievements/:id", app.updateIncomingAchievementHandler, "PUT")

	})

	mux.Group(func(mux *flow.Mux) {
		mux.Use(func(next http.Handler) http.Handler {
			return app.requirePermission(next)
		})

		mux.Use(app.requireAuthenticatedUser)

		mux.HandleFunc("/v1/admin/existing-collection", app.insertExistingCollectionHandler, "POST")

		mux.HandleFunc("/v1/admin/merch", app.createMerch, "POST")
		mux.HandleFunc("/v1/admin/merch/:id", app.getMerchHandler, "GET")
		mux.HandleFunc("/v1/admin/merch", app.getMerchsHandler, "GET")


		mux.HandleFunc("/v1/admin/rewards", app.getRewardsHandler, "GET")
		mux.HandleFunc("/v1/admin/rewards/:id", app.getRewardHandler, "GET")
		mux.HandleFunc("/v1/admin/rewards/:id", app.deleteRewardHandler, "DELETE")

		mux.HandleFunc("/v1/admin/users", app.getUsersHandler, "GET")
		mux.HandleFunc("/v1/admin/users/:id", app.getUserHandler, "GET")

		mux.HandleFunc("/v1/admin/users", app.createUserHandler, "POST")
		mux.HandleFunc("/v1/admin/users/:id", app.deleteUserHandler, "DELETE")
		mux.HandleFunc("/v1/admin/users/:id", app.updateAdminUserHandler, "PATCH")


		mux.HandleFunc("/v1/admin/collections", app.getCollectionsHandler, "GET")
		mux.HandleFunc("/v1/admin/collections/:id", app.getCollectionHandler, "GET")

		mux.HandleFunc("/v1/admin/collections", app.insertCollectionHandler, "POST")
		mux.HandleFunc("/v1/admin/collections/:id", app.deleteCollectionHandler, "DELETE")
		mux.HandleFunc("/v1/admin/collections/:id", app.updateCollectionHandler, "PATCH")

		mux.HandleFunc("/v1/admin/prototype-nfts", app.getPrototypeTokensHandler, "GET")
		mux.HandleFunc("/v1/admin/prototype-nfts/:id", app.getMetadataNftsHandler, "GET")

		mux.HandleFunc("/v1/admin/prototype-nfts", app.insertPrototypeNft, "POST")
		mux.HandleFunc("/v1/admin/prototype-nfts/:id", app.deletePrototypeHandler, "DELETE")
		mux.HandleFunc("/v1/admin/prototype-nfts/:id", app.updatePrototypeHandler, "PATCH")
		

		mux.HandleFunc("/v1/admin/minted-nfts", app.getTokensHandler, "GET")
		mux.HandleFunc("/v1/admin/minted-nfts/:id", app.getTokenHandler, "GET")

		mux.HandleFunc("/v1/admin/minted-nfts/:id", app.deleteTokenHandler, "DELETE")
		mux.HandleFunc("/v1/admin/minted-nfts", app.mintHandler, "POST")
		mux.HandleFunc("/v1/admin/minted-nfts/:id", app.updateTokenHandler, "PATCH")

		mux.HandleFunc("/v1/admin/activities", app.getActivitiesHandler, "GET")
		mux.HandleFunc("/v1/admin/activities/:id", app.getActivityHandler, "GET")

		mux.HandleFunc("/v1/admin/activities", app.InsertActivityHandler, "POST")
		mux.HandleFunc("/v1/admin/activities/:id", app.deleteActivityHandler, "DELETE")
		mux.HandleFunc("/v1/admin/activities/:id", app.updateActivityHandler, "PATCH")

		mux.HandleFunc("/v1/admin/permissions", app.getPermissionsHandler, "GET")

		mux.HandleFunc("/v1/admin/roles", app.getRolesHandler, "GET")
		mux.HandleFunc("/v1/admin/roles/:id", app.getRoleHandler, "GET")

		mux.HandleFunc("/v1/admin/roles", app.insertRoleHandler, "POST")
		mux.HandleFunc("/v1/admin/roles/:id", app.deleteRolesHandler, "DELETE")
		mux.HandleFunc("/v1/admin/roles/:id", app.updateRoleHandler, "PATCH")

	})

	return mux
}
