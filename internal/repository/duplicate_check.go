package repository

import (
	"context"

	"github.com/raman-vhd/arvan-challenge/internal/lib"
	"github.com/redis/go-redis/v9"
)

type IDuplicateCheckRepository interface {
	// check if data is duplicate for the user or not
	// returns true if duplicate
	IsDuplicate(ctx context.Context, userID string, dataID string) (bool, error)
}

type duplicateCheckRepository struct {
	redis *redis.Client
}

func NewDuplicateCheck(db lib.Database) IDuplicateCheckRepository {
	return duplicateCheckRepository{
		redis: db.RedisClient,
	}
}

func (r duplicateCheckRepository) IsDuplicate(ctx context.Context, userID string, dataID string) (bool, error) {
	key := "data." + userID
	exist, err := r.redis.SIsMember(ctx, key, dataID).Result()
	if err != nil {
		return false, err
	}

    return exist, nil
}
