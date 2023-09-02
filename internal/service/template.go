package service

import (
	"github.com/raman-vhd/arvan-challenge/internal/repository"
)

type ITemplateService interface {
	Action()
}

type templateService struct {
	repo repository.ITemplateRepository
}

func NewTemplate(
	repo repository.ITemplateRepository,
) ITemplateService {
	return templateService{
		repo: repo,
	}
}

func (s templateService) Action() {
	return
}
