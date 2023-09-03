package service

import (
	"context"

	"github.com/raman-vhd/arvan-challenge/internal/repository"
)

type IDuplicateCheckService interface {
	// checks if the user already has the data with the given dataID
	// returns true if it's duplicated
	IsDuplicate(ctx context.Context, userID string, dataID string) (bool, error)
}

type duplicateCheckService struct {
	repo repository.IDuplicateCheckRepository
}

func NewDuplicateCheck(
	repo repository.IDuplicateCheckRepository,
) IDuplicateCheckService {
	return duplicateCheckService{
		repo: repo,
	}
}

func (s duplicateCheckService) IsDuplicate(ctx context.Context, userID string, dataID string) (bool, error) {
	d, err := s.repo.IsDuplicate(ctx, userID, dataID)
	if err != nil {
		return false, err
	}

	return d, nil
}
