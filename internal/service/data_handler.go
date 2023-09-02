package service

import (
	"context"

	"github.com/raman-vhd/arvan-challenge/internal/lib"
	"github.com/raman-vhd/arvan-challenge/internal/model"
)

type IDataHandlerService interface {
	ProcessData(ctx context.Context, data model.Data) error
}

type dataHandlerService struct {
	P lib.EventProducer
}

func NewDataHandler(
	p lib.EventProducer,
) IDataHandlerService {
	return dataHandlerService{
		P: p,
	}
}

func (s dataHandlerService) ProcessData(ctx context.Context, data model.Data) error {
	err := s.P.ProduceEvent(data, "data-processing")
	return err
}
