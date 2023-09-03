package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/raman-vhd/arvan-challenge/internal/lib"
	"github.com/raman-vhd/arvan-challenge/internal/model"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRateLimitRepository interface {
	// gets the information about rate limit per minute by userID and returns max-rate and current bucket state
	// first it requests redis, if data exists it'll return. if not, it'll get data from main DB, cache it to redis and then return data
	GetRateLimitPerMinInfo(ctx context.Context, userID string) (model.TokenBucketInfo, error)

	// refills the token bucket by size for user with userID. and returns token bucket state
	RefillRateLimitPerMinBucket(ctx context.Context, userID string, size int) (int, error)

	// decreases the token bucket state by 1
	UseRateLimitPerMinToken(ctx context.Context, userID string) error

	// gets the information about request size limit per month by userID and returns max-size and current month state
	// first it requests redis, if data exists it'll return. if not, it'll get data from main DB, cache it to redis and then return data
	GetReqSizeLimitPerMonInfo(ctx context.Context, userID string) (model.ReqSizePerMonInfo, error)

	// adds the given request size to the current used size
	AddToCurrentReqSizePerMon(ctx context.Context, userID string, size int) error
}

type rateLimitRepository struct {
	Mongo *mongo.Collection
	Redis *redis.Client
}

func NewRateLimit(db lib.Database) IRateLimitRepository {
	return rateLimitRepository{
		Mongo: db.GetCollection("rate_limit"),
		Redis: db.RedisClient,
	}
}

func (r rateLimitRepository) GetRateLimitPerMinInfo(ctx context.Context, userID string) (model.TokenBucketInfo, error) {
	key := "RatePerMin." + userID
	exist, err := r.Redis.Exists(ctx, key).Result()
	if err != nil {
		return model.TokenBucketInfo{}, err
	}

	if exist == 1 {
		var info model.TokenBucketInfo
		i := r.Redis.HGetAll(ctx, key)
		err := i.Scan(&info)
		if err != nil {
			return model.TokenBucketInfo{}, err
		}
		return info, nil
	}

	var rateLimitInfo model.RateLimit
	err = r.Mongo.FindOne(ctx,
		bson.M{
			"userid": userID,
		}).Decode(&rateLimitInfo)
	if err != nil {
		return model.TokenBucketInfo{}, err
	}

	rateInt, err := strconv.Atoi(rateLimitInfo.ReqPerMin)
	if err != nil {
		return model.TokenBucketInfo{}, err
	}

	tokenBucketInfo := model.TokenBucketInfo{
		Rate:       rateInt,
		LastRefill: int(time.Now().Unix()),
		Bucket:     rateInt,
	}

	err = r.Redis.HSet(ctx, key,
		"rate", tokenBucketInfo.Rate,
		"lastRefill", tokenBucketInfo.LastRefill,
		"bucket", tokenBucketInfo.Bucket,
	).Err()
	if err != nil {
		return model.TokenBucketInfo{}, err
	}

	return tokenBucketInfo, nil
}

func (r rateLimitRepository) RefillRateLimitPerMinBucket(ctx context.Context, userID string, size int) (int, error) {
	key := "RatePerMin." + userID
	err := r.Redis.HSet(ctx, key,
		"bucket", size,
		"lastRefill", time.Now().Unix(),
	).Err()
	if err != nil {
		return 0, err
	}
	return size, nil
}

func (r rateLimitRepository) UseRateLimitPerMinToken(ctx context.Context, userID string) error {
	key := "RatePerMin." + userID
	var size int
	err := r.Redis.HGet(ctx, key, "bucket").Scan(&size)
	if err != nil {
		return err
	}

	err = r.Redis.HSet(ctx, key, "bucket", size-1).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r rateLimitRepository) GetReqSizeLimitPerMonInfo(ctx context.Context, userID string) (model.ReqSizePerMonInfo, error) {
	key := "ReqSizePerMon." + userID
	exist, err := r.Redis.Exists(ctx, key).Result()
	if err != nil {
		return model.ReqSizePerMonInfo{}, err
	}

	if exist == 1 {
		var info model.ReqSizePerMonInfo
		i := r.Redis.HGetAll(ctx, key)
		err := i.Scan(&info)
		if err != nil {
			return model.ReqSizePerMonInfo{}, err
		}
		return info, nil
	}

	var rateLimitInfo model.RateLimit
	err = r.Mongo.FindOne(ctx,
		bson.M{
			"userid": userID,
		}).Decode(&rateLimitInfo)
	if err != nil {
		return model.ReqSizePerMonInfo{}, err
	}

	rateInt, err := strconv.Atoi(rateLimitInfo.ReqPerMon)
	if err != nil {
		return model.ReqSizePerMonInfo{}, err
	}

	reqSizePerMonInfo := model.ReqSizePerMonInfo{
		Max:  rateInt,
		Size: 0,
		Date: int(time.Now().Unix()),
	}

	err = r.Redis.HSet(ctx, key,
		"max", reqSizePerMonInfo.Max,
		"size", 0,
		"date", reqSizePerMonInfo.Date,
	).Err()
	if err != nil {
		return model.ReqSizePerMonInfo{}, err
	}

	return reqSizePerMonInfo, nil
}

func (r rateLimitRepository) AddToCurrentReqSizePerMon(ctx context.Context, userID string, size int) error {
	key := "ReqSizePerMon." + userID
	var currentSize int
	err := r.Redis.HGet(ctx, key, "size").Scan(&currentSize)
	if err != nil {
		return err
	}

	err = r.Redis.HSet(ctx, key, "size", currentSize+size).Err()
	if err != nil {
		return err
	}
	return nil
}
