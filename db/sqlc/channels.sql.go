// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: channels.sql

package db

import (
	"context"
)

const createChannel = `-- name: CreateChannel :one
INSERT INTO channels (
  name
) VALUES (
  $1
)
RETURNING id, name, created_at
`

func (q *Queries) CreateChannel(ctx context.Context, name string) (*Channel, error) {
	row := q.db.QueryRow(ctx, createChannel, name)
	var i Channel
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return &i, err
}

const getChannelById = `-- name: GetChannelById :one
SELECT id, name, created_at
FROM channels
where id = $1
`

func (q *Queries) GetChannelById(ctx context.Context, id int64) (*Channel, error) {
	row := q.db.QueryRow(ctx, getChannelById, id)
	var i Channel
	err := row.Scan(&i.ID, &i.Name, &i.CreatedAt)
	return &i, err
}

const getChannels = `-- name: GetChannels :many
SELECT id, name, created_at
FROM channels
`

func (q *Queries) GetChannels(ctx context.Context) ([]*Channel, error) {
	rows, err := q.db.Query(ctx, getChannels)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Channel{}
	for rows.Next() {
		var i Channel
		if err := rows.Scan(&i.ID, &i.Name, &i.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
