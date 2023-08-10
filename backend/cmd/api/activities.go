package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexedwards/flow"
	"github.com/ton-developer-program/internal/request"
	"github.com/ton-developer-program/internal/response"
)


func (app *application) InsertActivityHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		TokenThreshold *int64 `json:"token_threshold"`
		SBTMetadata *string `json:"sbt_token_metadata"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// get sbtMetadata by base64
	sbtMetadata, err := app.sqlModels.Nfts.GetNFTMetadataByBase64(*input.SBTMetadata)

	var tokenThreshold int64 = 0

	if input.TokenThreshold != nil {
		tokenThreshold = *input.TokenThreshold
	}

	// Assuming you have a method in your `sqlModels` to insert the activity
	activity, err := app.sqlModels.Activities.Insert(*input.Name, *input.Description, tokenThreshold, sbtMetadata.ID)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusCreated, activity)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
	}
}


func (app *application) getActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	end := r.URL.Query().Get("_end")
	

	start := r.URL.Query().Get("_start")

	endInt, err := strconv.Atoi(end)
	if err != nil {
		app.badRequest(w, r, errors.New("current must be an integer"))
		return
	}

	startInt, err := strconv.Atoi(start)
	if err != nil {
		app.badRequest(w, r, errors.New("page_size must be an integer"))
		return
	}

	nameLike := r.URL.Query().Get("name_like")

	// create filters
	filters := []string{}
	if nameLike != "" {
		filters = append(filters, fmt.Sprintf("name ILIKE '%%%s%%'", nameLike))
		filters = append(filters, fmt.Sprintf("description ILIKE '%%%s%%'", nameLike))
	}

	filter := ""
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s", strings.Join(filters, " OR "))
	}
	
	activities, err := app.sqlModels.Activities.GetAll(startInt, endInt, filter)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, activities)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
	}
}

func (app *application) deleteActivityHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the activity ID from the URL path or request query parameters
	activityID := flow.Param(r.Context(), "id")

	// convert to int64
	activityIDInt, err := strconv.ParseInt(activityID, 10, 64)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	err = app.sqlModels.Activities.Delete(activityIDInt)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]string{"message": "Activity deleted successfully"})
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
	}
}


func (app *application) updateActivityHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// Parse the activity ID from the URL parameter
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	// Retrieve the activity by ID
	activity, err := app.sqlModels.Activities.GetByID(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	var input struct {
		Id int64 `json:"id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		TokenThreshold *int64 `json:"token_threshold"`
	}

	err = request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// Check if input values are present and update the activity fields accordingly
	if input.Name != nil {
		activity.Name = *input.Name
	}

	if input.Description != nil {
		activity.Description = *input.Description
	}

	if input.TokenThreshold != nil {
		activity.TokenThreshold = *input.TokenThreshold
	}


	// Update the activity in the database
	updatedActivity, err := app.sqlModels.Activities.Update(activity.ID, activity.Name, activity.Description, activity.TokenThreshold)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, updatedActivity)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
	}
}

func (app *application) getActivityHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// Parse the activity ID from the URL parameter
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	// Retrieve the activity by ID
	activity, err := app.sqlModels.Activities.GetByID(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, activity)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
	}
}

type Reward struct {
	UserID int64 `json:"user_id"`
	ActivityID int64 `json:"activity_id"`
	SBTTokenID int64 `json:"sbt_token_id"`
	PointsEarned int64 `json:"points_earned"`
}

func (app *application) getRewardsHandler(w http.ResponseWriter, r *http.Request) {
	// get rewards from db
	end := r.URL.Query().Get("_end")
	

	start := r.URL.Query().Get("_start")

	endInt, err := strconv.Atoi(end)
	if err != nil {
		app.badRequest(w, r, errors.New("current must be an integer"))
		return
	}

	startInt, err := strconv.Atoi(start)
	if err != nil {
		app.badRequest(w, r, errors.New("page_size must be an integer"))
		return
	}

	nameLike := r.URL.Query().Get("name_like")

	// create filters
	filters := []string{}
	if nameLike != "" {
		filters = append(filters, fmt.Sprintf("name ILIKE '%%%s%%'", nameLike))
		filters = append(filters, fmt.Sprintf("description ILIKE '%%%s%%'", nameLike))
	}

	filter := ""
	if len(filters) > 0 {
		filter = fmt.Sprintf("WHERE %s", strings.Join(filters, " OR "))
	}
	
	rewards, err := app.sqlModels.Rewards.GetAll(startInt, endInt, filter)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	response.JSON(w, http.StatusOK,  rewards)
}


func (app *application) deleteRewardHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// Parse the reward ID from the URL parameter
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	// Delete the reward from the database
	err = app.sqlModels.Rewards.Delete(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, map[string]string{"message": "Reward deleted successfully"})
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
	}



}

func (app *application) getRewardHandler(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	// Parse the reward ID from the URL parameter
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.badRequest(w, r, errors.New("id must be an integer"))
		return
	}

	// Retrieve the reward by ID
	reward, err := app.sqlModels.Rewards.GetById(idInt64)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
		return
	}

	err = response.JSON(w, http.StatusOK, reward)
	if err != nil {
		app.serverError(w, r, err) 
		app.logger.Error(err,nil)
	}
}


