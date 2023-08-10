package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ActivitiesModel struct {
	DB *sqlx.DB
}

type Activity struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	TokenThreshold int64 `db:"token_threshold" json:"token_threshold"`
	SBTPrototypeID int64 `db:"sbt_prototype_id" json:"sbt_prototype_id"`
}

type UserActivity struct {
	ID         int64 `db:"id" json:"id"`
	UserID     int64 `db:"user_id" json:"user_id"`
	ActivityID int64 `db:"activity_id" json:"activity_id"`
	SBTTokenID int64 `db:"sbt_token_id" json:"sbt_token_id"`
	PointsEarned int64 `db:"points_earned" json:"points_earned"`
	CreatedAt  int64 `db:"created_at" json:"created_at"`
}

// insert user activity
func (m *ActivitiesModel) InsertUserActivity(userID, activityID, sbtTokenID, pointsEarned int64) (*UserActivity, error) {
	stmt := `INSERT INTO user_activities (user_id, activity_id, sbt_token_id, points_earned) VALUES ($1, $2, $3, $4) RETURNING *`

	row := m.DB.QueryRow(stmt, userID, activityID, sbtTokenID, pointsEarned)

	var userActivity UserActivity

	err := row.Scan(
		&userActivity.ID,
		&userActivity.UserID,
		&userActivity.ActivityID,
		&userActivity.SBTTokenID,
		&userActivity.PointsEarned,
		&userActivity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &userActivity, nil
}

// get user activities
func (m *ActivitiesModel) GetUserActivities(userID int64) ([]*UserActivity, error) {
	query := `SELECT * FROM user_activities WHERE user_id=$1 ORDER BY id`

	userActivities := []*UserActivity{}

	err := m.DB.Select(&userActivities, query, userID)
	if err != nil {
		return nil, err
	}

	return userActivities, nil
}

// delete user activity
func (m *ActivitiesModel) DeleteUserActivity(userID, activityID int64) error {
	stmt := `DELETE FROM user_activities WHERE user_id=$1 AND activity_id=$2`

	_, err := m.DB.Exec(stmt, userID, activityID)
	if err != nil {
		return err
	}

	return nil
}

// get user activity
func (m *ActivitiesModel) GetUserActivity(userID, activityID int64) (*UserActivity, error) {
	stmt := `SELECT * FROM user_activities WHERE user_id=$1 AND activity_id=$2`

	userActivity := &UserActivity{}

	err := m.DB.Get(userActivity, stmt, userID, activityID)
	if err != nil {
		return nil, err
	}

	return userActivity, nil
}

// update user activity
func (m *ActivitiesModel) UpdateUserActivity(userID, activityID, sbtTokenID, pointsEarned int64) (*UserActivity, error) {
	stmt := `UPDATE user_activities SET sbt_token_id=$1, points_earned=$2 WHERE user_id=$3 AND activity_id=$4 RETURNING *`

	row := m.DB.QueryRow(stmt, sbtTokenID, pointsEarned, userID, activityID)

	var userActivity UserActivity

	err := row.Scan(
		&userActivity.ID,
		&userActivity.UserID,
		&userActivity.ActivityID,
		&userActivity.SBTTokenID,
		&userActivity.PointsEarned,
		&userActivity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &userActivity, nil
}



// // insert see in backend\internal\database\users.go

func (m *ActivitiesModel) Insert(name, description string, points int64, sbtId int64) (*Activity, error) {
	
	stmt := `INSERT INTO activities (name, description, token_threshold, sbt_prototype_id) VALUES ($1, $2, $3, $4) RETURNING *`

	row := m.DB.QueryRow(stmt, name, description, points, sbtId)

	var activity Activity

	err := row.Scan(
		&activity.ID,
		&activity.Name,
		&activity.Description,
		&activity.TokenThreshold,
		&activity.SBTPrototypeID,
	)

	if err != nil {
		return nil, err
	}

	return &activity, nil

}

// get all
func (m *ActivitiesModel) GetAll(start, end int, filter string) ([]*Activity, error) {
	query := fmt.Sprintf(`SELECT * FROM activities %s ORDER BY id LIMIT $1 OFFSET $2`, filter)

	activities := []*Activity{}

	err := m.DB.Select(&activities, query, end, start)
	if err != nil {
		return nil, err
	}

	return activities, nil

}

// delete
func (m *ActivitiesModel) Delete(id int64) error {
	stmt := `DELETE FROM activities WHERE id = $1`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// get by id
func (m *ActivitiesModel) GetByID(id int64) (*Activity, error) {
	stmt := `SELECT * FROM activities WHERE id = $1`

	activity := &Activity{}

	err := m.DB.Get(activity, stmt, id)
	if err != nil {
		return nil, err
	}

	return activity, nil
}

// get by name
func (m *ActivitiesModel) GetByName(name string) (*Activity, error) {
	stmt := `SELECT * FROM activities WHERE name = $1`

	activity := &Activity{}

	err := m.DB.Get(activity, stmt, name)

	if err != nil {
		return nil, err
	}

	return activity, nil
}

// update
func (m *ActivitiesModel) Update(id int64, name, description string, points int64) (*Activity, error) {
	stmt := `UPDATE activities SET name = $1, description = $2, token_threshold = $3 WHERE id = $4 RETURNING *`

	activity := &Activity{}

	err := m.DB.Get(activity, stmt, name, description, points, id)
	if err != nil {
		return nil, err
	}

	return activity, nil
}



