package service

import (
	"context"
	"fmt"
	"time"

	"github.com/raman-vhd/arvan-challenge/internal/model"
	"github.com/raman-vhd/arvan-challenge/internal/repository"
)

type IRateLimitService interface {
	CheckRateLimitPerMin(ctx context.Context, userID string) (bool, error)
	UseRateLimitPerMinToken(ctx context.Context, userID string) error
	CheckReqSizeLimitPerMon(ctx context.Context, userID string, size int) (bool, error)
	AddToCurrentReqSizePerMon(ctx context.Context, userID string, size int) error
}

type rateLimitService struct {
	repo repository.IRateLimitRepository
}

func NewRateLimit(
	repo repository.IRateLimitRepository,
) IRateLimitService {
	return rateLimitService{
		repo: repo,
	}
}

func (s rateLimitService) CheckRateLimitPerMin(ctx context.Context, userID string) (bool, error) {
	info, err := s.repo.GetRateLimitPerMinInfo(ctx, userID)
	if err != nil {
		return false, err
	}

	bucketCount := calculateRatePerMinBucketRefillIncrement(info) + info.Bucket
	bucket, err := s.repo.RefillRateLimitPerMinBucket(ctx, userID, bucketCount)
	if err != nil {
		return false, err
	}

	if bucket == 0 {
		return false, nil
	}

	return true, nil
}

// it calculates the refill increment for request per minute based on the elapsed time since the last refill and the maximum rate per minute
func calculateRatePerMinBucketRefillIncrement(info model.TokenBucketInfo) int {
	lastRefill := time.Unix(int64(info.LastRefill), 0)
	minSince := time.Since(lastRefill).Minutes()
	incr := int(minSince * float64(info.Rate))

	if incr+info.Bucket > info.Rate {
		return info.Rate - info.Bucket
	}

	return incr
}

func (s rateLimitService) CheckReqSizeLimitPerMon(ctx context.Context, userID string, size int) (bool, error) {
	info, err := s.repo.GetReqSizeLimitPerMonInfo(ctx, userID)
	if err != nil {
		return false, err
	}

    fmt.Println(info.Size, size, info.Max)
	if info.Size+size > info.Max {
		return false, nil
	}

	return true, nil
}

func (s rateLimitService) UseRateLimitPerMinToken(ctx context.Context, userID string) error {
	err := s.repo.UseRateLimitPerMinToken(ctx, userID)
	return err
}

func (s rateLimitService) AddToCurrentReqSizePerMon(ctx context.Context, userID string, size int) error {
    err := s.repo.AddToCurrentReqSizePerMon(ctx, userID, size)
    return err
}
