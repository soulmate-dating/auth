package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/soulmate-dating/auth/internal/domain"
)

type Repo struct {
	pool       ConnPool
	mapUsers   func(row pgx.CollectableRow) (domain.User, error)
	mapUserIDs func(row pgx.CollectableRow) (domain.UserID, error)
}

func NewRepo(pool ConnPool) *Repo {
	return &Repo{
		pool:       pool,
		mapUsers:   pgx.RowToStructByName[domain.User],
		mapUserIDs: pgx.RowToStructByName[domain.UserID],
	}
}

func (r *Repo) CreateUser(ctx context.Context, p *domain.User) (uuid.UUID, error) {
	var args []any
	args = append(args,
		p.ID, p.Email, p.Password,
	)
	rows, err := r.pool.GetTx(ctx).Query(ctx, createUserQuery, args...)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("create profile: %w", err)
	}
	userID, err := pgx.CollectOneRow(rows, r.mapUserIDs)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("map id: %w", err)
	}
	return userID.ID, nil
}

func (r *Repo) GetUserByEmail(ctx context.Context, id string) (*domain.User, error) {
	rows, err := r.pool.GetTx(ctx).Query(ctx, getUserByEmailQuery, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	profile, err := pgx.CollectOneRow(rows, r.mapUsers)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("map user: %w", err)
	}
	return &profile, nil
}
