package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/soulmate-dating/auth/internal/models"
)

type Repo struct {
	pool       *pgxpool.Pool
	mapUsers   func(row pgx.CollectableRow) (models.User, error)
	mapUserIDs func(row pgx.CollectableRow) (models.UserID, error)
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{
		pool:       pool,
		mapUsers:   pgx.RowToStructByName[models.User],
		mapUserIDs: pgx.RowToStructByName[models.UserID],
	}
}

func (r *Repo) CreateUser(ctx context.Context, p *models.User) (string, error) {
	var args []any
	args = append(args,
		p.ID, p.Email, p.Password,
	)
	rows, err := r.pool.Query(ctx, createUserQuery, args...)
	if err != nil {
		return "", fmt.Errorf("create profile: %w", err)
	}
	userID, err := pgx.CollectOneRow(rows, r.mapUserIDs)
	if err != nil {
		return "", fmt.Errorf("map id: %w", err)
	}
	return userID.ID, nil
}

func (r *Repo) GetUserByEmail(ctx context.Context, id string) (*models.User, error) {
	rows, err := r.pool.Query(ctx, getUserByEmailQuery, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	profile, err := pgx.CollectOneRow(rows, r.mapUsers)
	if err != nil {
		return nil, fmt.Errorf("map user: %w", err)
	}
	return &profile, nil
}

//
//func (r *Repo) UpdateUserLoginStatus(ctx context.Context, id string, loggedIn bool) error {
//	var args []any
//	args = append(args, id, loggedIn)
//
//	_, err := r.pool.Exec(ctx, updateUserLoginStatusQuery, args...)
//	if err != nil {
//		return fmt.Errorf("update user login status: %w", err)
//	}
//	return nil
//}
