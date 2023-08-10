package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/flow"
	"github.com/hibiken/asynq"
	"github.com/ton-developer-program/internal/database"
	"github.com/ton-developer-program/internal/request"
	"github.com/ton-developer-program/internal/response"
	"github.com/ton-developer-program/internal/tonconnect"
	"github.com/ton-developer-program/internal/validator"
	"github.com/tonkeeper/tongo"
	"golang.org/x/oauth2"
)

func (app *application) manifestTonConnectHandler(w http.ResponseWriter, r *http.Request) {

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"url": "https://tdp.tonbuilders.com/",
		"name": "TON Developers Platform",
		"iconUrl": "https://images-tdp.s3.eu-central-1.amazonaws.com/favicon.png",
	})

}
func (app *application) proofHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get body
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}
	var tp tonconnect.TonProof

	err = json.Unmarshal(b, &tp)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// check payload
	err = tonconnect.CheckPayload(tp.Proof.Payload, app.config.Ton.SharedSecret)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	parsed, err := tonconnect.ConvertTonProofMessage(ctx, &tp)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	net := networks[tp.Network]

	if net == nil {
		app.serverError(w, r, errors.New("invalid network"))
		return
	}
	addr, err := tongo.ParseAccountID(tp.Address)
	if net == nil {
		app.serverError(w, r, errors.New("invalid address"))
		return
	}

	check, err := tonconnect.CheckProof(ctx, addr, net, parsed, app.config.App.DomainName)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}
	if !check {
		app.serverError(w, r, errors.New("proof verification failed"))
		return
	}

	var user *database.User
	// check if user exists
	user, err = app.sqlModels.Users.GetByAddress(tp.Address)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}


	if user == nil {

		// create user
		user = &database.User{
			RawAddress:      tp.Address,
			FriendlyAddress: tonconnect.ConvertToFriendlyAddr(tp.Address).String(),
		}

		user, err = app.sqlModels.Users.Insert(user)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}
	}

	// // get roles and permissions
	permissions, err := app.sqlModels.Users.GetAllPermissions(user.ID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	user.Permissions = permissions

	t, err := app.sqlModels.Tokens.New(user.ID, 24*7*time.Hour, database.ScopeAuthentication)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]interface{}{
		"token":   t.Plaintext,
		"expires": t.Expiry,
		"user":    user,
	})

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) getUsersHandler(w http.ResponseWriter, r *http.Request) {

	pagination, err := getPagination(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}


	nameLike := r.URL.Query().Get("name_like")

	// create filters
	filters := []string{}
	if nameLike != "" {
		filters = append(filters, fmt.Sprintf("first_name ILIKE '%%%s%%'", nameLike))
		filters = append(filters, fmt.Sprintf("last_name ILIKE '%%%s%%'", nameLike))
		filters = append(filters, fmt.Sprintf("username ILIKE '%%%s%%'", nameLike))
	}

	filter := ""
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s", strings.Join(filters, " OR "))
	}

	users, err := app.sqlModels.Users.GetMany(pagination, filter, "id ASC")
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	for i := range users {
		// get linked accounts
		accounts, err := app.sqlModels.Users.GetLinkedAccounts(users[i].ID)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}
		users[i].LinkedAccounts = accounts

		// get user roles
		role, err := app.sqlModels.Permissions.GetUserRoles(users[i].ID)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}
		users[i].Role = role
	}

	err = response.JSON(w, http.StatusOK, users)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}
}

// get user by id
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	user, err := app.sqlModels.Users.GetById(idInt64)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// get linkee accounts
	accounts, err := app.sqlModels.Users.GetLinkedAccounts(user.ID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	user.LinkedAccounts = accounts

	// get user roles
	role, err := app.sqlModels.Permissions.GetUserRoles(user.ID)

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	user.Role = role

	// get

	err = response.JSON(w, http.StatusOK, user)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}
}

func (app *application) updateAdminUserHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// get by id
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	user, err := app.sqlModels.Users.GetById(idInt64)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	var input struct {
		FirstName      *string   `json:"first_name"`
		LastName       *string   `json:"last_name"`
		Username       *string   `json:"username"`
		Job            *string   `json:"job"`
		Bio            *string   `json:"bio"`
		Languages      *[]string `db:"languages" json:"languages"`
		Certifications *[]string `db:"certifications" json:"certifications"`
		RoleId         *int64    `db:"role_id" json:"role_id"`
	}

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		user.LastName = *input.LastName
	}

	if input.Username != nil {
		user.Username = *input.Username
	}

	if input.Job != nil {
		user.Job = input.Job
	}

	if input.Bio != nil {
		user.Bio = input.Bio
	}

	if input.Languages != nil {
		user.Languages = *input.Languages
	}

	if input.Certifications != nil {
		user.Certifications = *input.Certifications
	}

	user, err = app.sqlModels.Users.Update(user)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	linkedAccounts, err := app.sqlModels.Users.GetLinkedAccounts(user.ID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	user.LinkedAccounts = linkedAccounts

	if input.RoleId != nil {
		// delete old role
		err = app.sqlModels.Permissions.DeleteUserRoles(user.ID)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}

		// add new role
		err = app.sqlModels.Permissions.InsertUserRole(user.ID, *input.RoleId)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}

	}

	err = response.JSON(w, http.StatusOK, user)

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) createMerch(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name string `json:"name"`
		Amount int64 `json:"amount"`
		Store string `json:"store"`
		UserId int64 `json:"user_id"`		
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	merch := &database.Merch{
		Name: input.Name,
		Amount: input.Amount,
		Store: input.Store,
		UserID: input.UserId,
	}

	// insert into merch
	err = app.sqlModels.Rewards.InsertMerch(merch)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, merch)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}
}

// get merch by id
func (app *application) getMerchHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	merch, err := app.sqlModels.Rewards.GetMerchByID(idInt64)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, merch)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}
}

// get all merchs
func (app *application) getMerchsHandler(w http.ResponseWriter, r *http.Request) {
	merchs, err := app.sqlModels.Rewards.GetAllMerch()
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, merchs)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}
}

// create user
func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName       *string `json:"first_name"`
		LastName        *string `json:"last_name"`
		Username        *string `json:"username"`
		Job             *string `json:"job"`
		Bio             *string `json:"bio"`
		FriendlyAddress *string `json:"friendly_address"`
		RawAddress      *string `json:"raw_address"`
		RoleId          *int64  `json:"role_id"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := &database.User{
		FirstName:       *input.FirstName,
		LastName:        *input.LastName,
		Username:        *input.Username,
		Job:             input.Job,
		Bio:             input.Bio,
		FriendlyAddress: *input.FriendlyAddress,
		RawAddress:      *input.RawAddress,
	}

	user, err = app.sqlModels.Users.Create(user)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	if input.RoleId != nil {
		// insert role
		err = app.sqlModels.Permissions.InsertUserRole(user.ID, *input.RoleId)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}
	}

	err = response.JSON(w, http.StatusCreated, user)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)

	}
}

// delete by id
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	err = app.sqlModels.Users.Delete(idInt64)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]string{
		"message": "user deleted",
	})

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(user.Version), 32) != r.Header.Get("X-Expected-Version") {
			app.editConclictResponse(w, r)
			return
		}
	}

	var input struct {
		FirstName      *string   `json:"first_name"`
		LastName       *string   `json:"last_name"`
		Username       *string   `json:"username"`
		Job            *string   `json:"job"`
		Bio            *string   `json:"bio"`
		Languages      *[]string `db:"languages" json:"languages"`
		Certifications *[]string `db:"certifications" json:"certifications"`
		AvatarUrl      *string   `db:"avatar_url" json:"avatar_url"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	if input.FirstName != nil {
		if yes := validator.MaxRunes(*input.FirstName, 50); !yes {
			app.badRequest(w, r, errors.New("first name must be less than 50 characters"))
			return
		}

		user.FirstName = *input.FirstName
	}

	if input.LastName != nil {

		if yes := validator.MaxRunes(*input.LastName, 50); !yes {
			app.badRequest(w, r, errors.New("last name must be less than 50 characters"))
			return
		}

		user.LastName = *input.LastName
	}

	if input.Username != nil {

		if yes := validator.MaxRunes(*input.Username, 50); !yes {
			app.badRequest(w, r, errors.New("username must be less than 10 characters"))
			return
		}

		user.Username = *input.Username
	}

	if input.Job != nil {
		if yes := validator.MaxRunes(*input.Job, 100); !yes {
			app.badRequest(w, r, errors.New("job must be less than 100 characters"))
			return
		}

		user.Job = input.Job
	}

	if input.Bio != nil {

		if yes := validator.MaxRunes(*input.Bio, 500); !yes {
			app.badRequest(w, r, errors.New("bio must be less than 500 characters"))
			return
		}
		user.Bio = input.Bio
	}

	if input.Languages != nil {
		user.Languages = *input.Languages
	}

	if input.Certifications != nil {
		user.Certifications = *input.Certifications
	}

	if input.AvatarUrl != nil {
		user.AvatarURL = input.AvatarUrl
	}

	user, err = app.sqlModels.Users.Update(user)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, user)

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) unlinkAccountHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	provider := flow.Param(r.Context(), "provider")

	err := app.sqlModels.Users.DeleteLinkedAccountByUserId(user.ID, provider)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// get linked accounts
	err = response.JSON(w, http.StatusOK, user)

}

func (app *application) pinNftHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	id := flow.Param(r.Context(), "id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}
	// check if nft belongs to user
	nft, err := app.sqlModels.Nfts.GetTokenByID(idInt64)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	if nft.FriendlyOwnerAddress != user.FriendlyAddress {
		app.badRequest(w, r, errors.New("nft does not belong to user"))
		return
	}

	var pinned bool

	if nft.IsPinned == true {
		pinned = false
	} else {
		pinned = true
	}

	err = app.sqlModels.Nfts.PinNFT(idInt64, pinned)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// success
	err = response.JSON(w, http.StatusOK, map[string]string{
		"message": "nft pinned",
	})

}



func (app *application) getIncomingAchievementsHandler(w http.ResponseWriter, r *http.Request) {
	pagination, err := getPagination(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}


	user := app.contextGetUser(r)

	achievements, err := app.sqlModels.Rewards.GetStoredRewardsByUserID(user.FriendlyAddress, pagination)
	if err != nil {
		if err == sql.ErrNoRows {
			// return nil for achievements
			err = response.JSON(w, http.StatusOK, map[string]interface{}{
				"achievements": nil,
			})
			if err != nil {
				app.serverError(w, r, err)
				app.logger.Error(err, nil)
			}
			return
		}
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	var outputs []interface{}

	for _, achievement := range achievements {
		// get by achievement base 64 nft metadata
		nftMetaData, err := app.sqlModels.Nfts.GetNFTMetadataByBase64(achievement.Base64Metadata)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}

		var output struct {
			AchievementID  int64  `json:"achievement_id"`
			ImageUrl       string `json:"image_url"`
			Name           string `json:"name"`
			Description    string `json:"description"`
			Weight         int64  `json:"weight"`
			ApprovedByUser bool   `json:"approved"`
			Processed      bool   `json:"processed"`
		}

		output.AchievementID = achievement.ID
		output.ImageUrl = nftMetaData.Image
		output.Name = nftMetaData.Name
		output.Description = nftMetaData.Description
		output.Weight = 1000
		output.ApprovedByUser = achievement.ApprovedByUser
		output.Processed = achievement.Processed

		outputs = append(outputs, output)
	}

	// count
	count, err := app.sqlModels.Rewards.CountStoredRewardsByUserID(user.FriendlyAddress)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]interface{}{
		"achievements": outputs,
		"count":        count,
	})

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}
}

// put /incoming-achievements/:id
func (app *application) updateIncomingAchievementHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)

	id := flow.Param(r.Context(), "id")

	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	var input struct {
		ApprovedByUser *bool `json:"approved_by_user"`
	}

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	achievement, err := app.sqlModels.Rewards.GetStoredRewardByID(idInt64)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	if achievement.UserAddress != user.FriendlyAddress {
		app.badRequest(w, r, errors.New("achievement does not belong to user"))
		return
	}

	if input.ApprovedByUser != nil {
		achievement.ApprovedByUser = *input.ApprovedByUser
	}

	err = app.sqlModels.Rewards.UpdateStoredRewardApprovedByUser(achievement.ID, true)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	runGetStoredRewards := asynq.NewTask(database.TYPE_MINT_STORED_REWARDS, nil)

	info, err := app.asynqClient.Enqueue(runGetStoredRewards, asynq.TaskID("MINT_STORED_REWARDS"), asynq.ProcessIn(10*time.Second), asynq.MaxRetry(5),asynq.Retention(30 * time.Second), asynq.Queue(database.PRIORITY_URGENT))
	if err != nil {
		app.logger.Error(fmt.Errorf("error: %s", err), nil)
		return
	}

	app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))

	err = response.JSON(w, http.StatusOK, map[string]string{
		"message": "achievement updated",
	})

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) getTopUsersHandler(w http.ResponseWriter, r *http.Request) {

	pagination, err := getPagination(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	users, err := app.sqlModels.Users.GetTopUsers(pagination)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// get linked accounts
	for _, user := range users {
		linkedAccounts, err := app.sqlModels.Users.GetLinkedAccounts(user.ID)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}
		user.LinkedAccounts = linkedAccounts
	}

	err = response.JSON(w, http.StatusOK, users)

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) getUserByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := flow.Param(ctx, "username")

	user, err := app.sqlModels.Users.GetByUsername(username)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// get linked accounts

	linkedAccounts, err := app.sqlModels.Users.GetLinkedAccounts(user.ID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	user.LinkedAccounts = linkedAccounts

	// get linked accounts
	err = response.JSON(w, http.StatusOK, map[string]interface{}{
		"user": user,
	})
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}
}

// get nfts by user id
func (app *application) getNftsByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username := flow.Param(ctx, "username")

	pagination, err := getPagination(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user, err := app.sqlModels.Users.GetByUsername(username)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	nfts, err := app.sqlModels.Nfts.GetTokensByOwnerAddress(user.FriendlyAddress, pagination)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// get total count
	count, err := app.sqlModels.Nfts.GetTotalTokensByOwnerAddress(user.FriendlyAddress)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// get linked accounts
	err = response.JSON(w, http.StatusOK, map[string]interface{}{
		"nfts":  nfts,
		"count": count,
	})
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}
}

func (app *application) getMyAccountHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	// get linked accounts

	linkedAccounts, err := app.sqlModels.Users.GetLinkedAccounts(user.ID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	user.LinkedAccounts = linkedAccounts

	// get linked accounts
	err = response.JSON(w, http.StatusOK, user)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}
}

func (app *application) githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	// query username
	username := r.URL.Query().Get("username")

	url := app.githubOauthConfig.AuthCodeURL(username, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) githubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	user, err := app.sqlModels.Users.GetByUsername(state)
	if err != nil {
		fmt.Printf("Failed to get user: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// if username is not matched to the state
	if user.Username != state {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", user.Username, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := app.githubOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("githubOauthConfig.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		fmt.Printf("Failed to create request: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	req.Header.Set("Authorization", "token "+token.AccessToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed to get user info: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var githubUser database.GithubUser

	err = json.Unmarshal(content, &githubUser)
	if err != nil {
		fmt.Printf("Failed to unmarshal response body: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// insert linked account
	linkedAccount := &database.LinkedAccount{
		UserID:      user.ID,
		Provider:    database.ProviderGithub,
		AvatarURL:   githubUser.AvatarURL,
		Login:       githubUser.Login,
		AccessToken: token.AccessToken,
		CreatedAt:   uint64(time.Now().Unix()),
		UpdatedAt:   uint64(time.Now().Unix()),
		Version:     1,
	}

	err = app.sqlModels.Users.InsertLinkedAccount(linkedAccount)
	if err != nil {
		fmt.Printf("Failed to insert linked account: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// check if user has two linked accounts
	hasTwoAccounts, err := app.sqlModels.Users.HasTwoLinkedAccounts(user.ID)
	if err != nil {
		fmt.Printf("Failed to check if user has two linked accounts: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// get auth nft by metadata id
	authNft, err := app.sqlModels.Nfts.GetNFTMetadataByID(app.config.App.AuthMetadataID)
	if err != nil {
		fmt.Printf("Failed to get auth nft: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// check if user has auth nft
	hasAuthNft, err := app.sqlModels.Nfts.HasAuthNFT(user.FriendlyAddress, authNft.Base64)
	if err != nil {
		fmt.Printf("Failed to check if user has auth nft: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// if user has two linked accounts and has auth nft
	if hasTwoAccounts && !hasAuthNft {
		// update user status

		payload, err := json.Marshal(user.ID)

		if err != nil {
			app.logger.Error(err, nil)
			return
		}

		runGetReward := asynq.NewTask(database.TYPE_REWARD_FOR_LINKED_ACCOUNT, payload)

		info, err := app.asynqClient.Enqueue(runGetReward, asynq.TaskID(fmt.Sprint("reward_auth", user.ID)), asynq.ProcessIn(10*time.Second), asynq.MaxRetry(5), asynq.ProcessIn(5*time.Second), asynq.Retention(10*time.Minute), asynq.Queue(database.PRIORITY_URGENT))
		if err != nil {
			app.logger.Error(fmt.Errorf("error: %s", err), nil)
			return
		}

		app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))

	}

	http.Redirect(w, r, app.config.App.BaseUrl+"/settings", http.StatusTemporaryRedirect)

	return

}

type AuthData struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
	Hash      string `json:"hash"`
}



func (app *application) checkTelegramAuthorization(w http.ResponseWriter, r *http.Request) {
	

	user := app.contextGetUser(r)
	// parse body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	var jsonBody struct {
		AuthObj string `json:"auth_obj"`
	}

	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// base64 decode
	decoded, err := base64.RawStdEncoding.DecodeString(jsonBody.AuthObj)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	var authData AuthData

	// parse json
	err = json.Unmarshal(decoded, &authData)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}


	checkHash := authData.Hash
	authData.Hash = ""

	dataCheckArr := []string{}
	dataCheckArr = append(dataCheckArr, fmt.Sprintf("auth_date=%d", authData.AuthDate))

	if authData.FirstName != "" {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("first_name=%s", authData.FirstName))
		user.FirstName = authData.FirstName
	}

	if authData.LastName != "" {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("last_name=%s", authData.LastName))
		user.LastName = authData.LastName
	}
	
	dataCheckArr = append(dataCheckArr, fmt.Sprintf("id=%d", authData.ID))

	if authData.PhotoURL != "" {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("photo_url=%s", authData.PhotoURL))

		if user.AvatarURL == nil {
			user.AvatarURL = &authData.PhotoURL
		}
	}

	if authData.Username != "" {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("username=%s", authData.Username))

		// check if username exists
		_, err := app.sqlModels.Users.GetByUsername(authData.Username)

		if err != nil {
			if err == sql.ErrNoRows {
				user.Username = authData.Username
			}
		}

	}

	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	secretKey := sha256.Sum256([]byte(app.config.Auth.TelegramBotToken)) // replace with your bot token

	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))

	hash := hex.EncodeToString(h.Sum(nil))

	if hash != checkHash {
		app.logger.Error(fmt.Errorf("hashes do not match: %s != %s", hash, checkHash), nil)
		app.serverError(w, r, errors.New("Data is NOT from Telegram"))
		return
	}

	if time.Now().Unix()-authData.AuthDate > 86400 {
		app.logger.Error(fmt.Errorf("auth data is outdated: %d", time.Now().Unix()-authData.AuthDate), nil)
		app.serverError(w, r, errors.New("Data is outdated"))
		return
	}

	// insert linked account
	linkedAccount := &database.LinkedAccount{
		UserID:         user.ID,
		TelegramUserID: &authData.ID,
		Provider:       database.ProviderTelegram,
		AvatarURL:      authData.PhotoURL,
		Login:          authData.Username,
		AccessToken:    "",
		CreatedAt:      uint64(time.Now().Unix()),
		UpdatedAt:      uint64(time.Now().Unix()),
		Version:        1,
	}

	err = app.sqlModels.Users.InsertLinkedAccount(linkedAccount)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// update user
	_, err = app.sqlModels.Users.Update(user)
	if err != nil {
		app.logger.Error(err, nil)
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// check if user has two linked accounts
	hasTwoAccounts, err := app.sqlModels.Users.HasTwoLinkedAccounts(user.ID)
	if err != nil {
		fmt.Printf("Failed to check if user has two linked accounts: %s\n", err.Error())
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// get auth nft by metadata id
	authNft, err := app.sqlModels.Nfts.GetNFTMetadataByID(app.config.App.AuthMetadataID)
	if err != nil {
		fmt.Printf("Failed to get auth nft: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// check if user has auth nft
	hasAuthNft, err := app.sqlModels.Nfts.HasAuthNFT(user.FriendlyAddress, authNft.Base64)
	if err != nil {
		fmt.Printf("Failed to check if user has auth nft: %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// if user has two linked accounts and has auth nft
	if hasTwoAccounts && !hasAuthNft {
		// update user status

		payload, err := json.Marshal(user.ID)

		if err != nil {
			app.logger.Error(err, nil)
			return
		}

		runGetReward := asynq.NewTask(database.TYPE_REWARD_FOR_LINKED_ACCOUNT, payload)

		info, err := app.asynqClient.Enqueue(runGetReward, asynq.TaskID(fmt.Sprint("reward_auth", user.ID)), asynq.ProcessIn(10*time.Second), asynq.MaxRetry(5), asynq.ProcessIn(5*time.Second), asynq.Retention(10*time.Minute), asynq.Queue(database.PRIORITY_URGENT))
		if err != nil {
			app.logger.Error(fmt.Errorf("error: %s", err), nil)
			return
		}

		app.logger.Info(fmt.Sprintf("enqueued task with id %s", info.ID))

	}

	response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// payloadHandler handles payload generation for TonConnect.
// It generates a payload, with the time to live based on the configured value.
// It then returns a JSON response containing the payload and time to live.

func (app *application) payloadHandler(w http.ResponseWriter, r *http.Request) {

	// seconds proofLifetime time.second

	ttl := time.Duration(app.config.Ton.ProfLifeTimeSec) * time.Second

	payload, err := tonconnect.GeneratePayload(app.config.Ton.SharedSecret, ttl)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]string{
		"payload": payload,
		"ttl":     fmt.Sprintf("%d", app.config.Ton.ProfLifeTimeSec),
	})
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}
func (app *application) getRolesHandler(w http.ResponseWriter, r *http.Request) {

	roles, err := app.sqlModels.Permissions.GetAllRoles()
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// get all permissions for these roles
	for i, role := range roles {
		permissions, err := app.sqlModels.Permissions.GetRolePermissions(role.ID)
		if err != nil {
			app.serverError(w, r, err)
			app.logger.Error(err, nil)
			return
		}
		roles[i].Permissions = permissions

	
	}

	err = response.JSON(w, http.StatusOK, roles)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

// get all types of permissions
func (app *application) getPermissionsHandler(w http.ResponseWriter, r *http.Request) {

	permissions, err := app.sqlModels.Permissions.GetAllPermissions()
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, permissions)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) getRoleHandler(w http.ResponseWriter, r *http.Request) {

	roleID, err := strconv.ParseInt(flow.Param(r.Context(), "id"), 10, 64)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	role, err := app.sqlModels.Permissions.GetRoleByID(roleID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	permissions, err := app.sqlModels.Permissions.GetRolePermissions(role.ID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	
	var tags []int64

	for _, permission := range permissions {
		tags = append(tags, permission.ID)
	}


	role.Tags = tags
	
	role.Permissions = permissions

	err = response.JSON(w, http.StatusOK, role)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) insertRoleHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Permissions []int64 `json:"permissions"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	role := database.Role{
		Name:        input.Name,
		Description: input.Description,
	}

	roleID, err := app.sqlModels.Permissions.InsertRole(role)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = app.sqlModels.Permissions.InsertRolePermissionsByIDs(roleID, input.Permissions)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]int64{
		"id": roleID,
	})
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) deleteRolesHandler(w http.ResponseWriter, r *http.Request) {

	id := flow.Param(r.Context(), "id")

	roleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = app.sqlModels.Permissions.DeleteRole(roleID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}

func (app *application) updateRoleHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Permissions []int64 `json:"tags"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// delete old permissions for this role and get role by name
	role, err := app.sqlModels.Permissions.GetRoleByName(input.Name)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = app.sqlModels.Permissions.DeleteRolePermissions(role.ID)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// update role
	role.Name = input.Name
	role.Description = input.Description

	err = app.sqlModels.Permissions.UpdateRole(role)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	// insert new permissions
	err = app.sqlModels.Permissions.InsertRolePermissionsByIDs(role.ID, input.Permissions)
	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})

	if err != nil {
		app.serverError(w, r, err)
		app.logger.Error(err, nil)
	}

}
