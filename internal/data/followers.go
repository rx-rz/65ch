package data

import (
	"context"
	"database/sql"
	"time"
)

type Follower struct {
	FollowerID string    `json:"follower_id"`
	FollowedID string    `json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowerModel struct {
	DB *sql.DB
}

func (m *FollowerModel) FollowUser(ctx context.Context, followerID, followedID string) (*Follower, error) {

	const query = `
	INSERT INTO followers (follower_id, followed_id)
	VALUES ($1, $2)
	RETURNING follower_id, followed_id
	`
	newFollow := &Follower{
		FollowerID: followerID,
		FollowedID: followedID,
	}
	_, err := m.DB.ExecContext(
		ctx,
		query,
		followerID,
		followedID,
	)

	if err != nil {
		return nil, DetermineDBError(err, "follower_followuser")
	}
	return newFollow, nil
}

func (m *FollowerModel) UnfollowUser(ctx context.Context, followerID, followedID string) error {
	const query = `
	DELETE FROM followers 
	WHERE follower_id = $1
	AND followed_id = $2
	`
	_, err := m.DB.ExecContext(
		ctx,
		query,
		followerID,
		followedID,
	)
	if err != nil {
		return DetermineDBError(err, "follower_unfollowuser")
	}
	return nil
}

func (m *FollowerModel) ListFollowers(ctx context.Context, userID string) ([]*Follower, error) {
	const query = `
	SELECT (follower_id, followed_id, created_at) FROM followers
	WHERE followed_id = $1
	`
	rows, err := m.DB.QueryContext(
		ctx,
		query,
		userID)
	if err != nil {
		return nil, DetermineDBError(err, "followers_listfollowers")
	}
	defer rows.Close()
	var followers []*Follower
	for rows.Next() {
		var follower *Follower
		err = rows.Scan(
			&follower.FollowerID,
			&follower.FollowedID,
			&follower.CreatedAt,
		)
		if err != nil {
			return nil, DetermineDBError(err, "followers_listfollowers")
		}
		followers = append(followers, follower)
	}
	if err = rows.Err(); err != nil {
		return nil, DetermineDBError(err, "followers_listfollowers")
	}
	return followers, nil
}

func (m *FollowerModel) ListUsersFollowed(ctx context.Context, userID string) ([]*Follower, error) {
	const query = `
	SELECT (follower_id, followed_id, created_at) FROM followers
	WHERE follower_id = $1
	`
	rows, err := m.DB.QueryContext(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, DetermineDBError(err, "followers_listfollowers")
	}
	defer rows.Close()
	var usersFollowed []*Follower
	for rows.Next() {
		var userFollowed *Follower
		err = rows.Scan(
			&userFollowed.FollowerID,
			&userFollowed.FollowedID,
			&userFollowed.CreatedAt,
		)
		if err != nil {
			return nil, DetermineDBError(err, "followers_listfollowers")
		}
		usersFollowed = append(usersFollowed, userFollowed)
	}
	if err = rows.Err(); err != nil {
		return nil, DetermineDBError(err, "followers_listfollowers")
	}
	return usersFollowed, nil
}
