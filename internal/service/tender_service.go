package service

import (
	"errors"
	repository "tender-managment/internal/db/repo"
	"tender-managment/internal/model"
)

type TenderService struct {
	repo *repository.TenderRepository
}

func NewTenderService(repo *repository.TenderRepository) *TenderService {
	return &TenderService{repo: repo}
}

func (s *TenderService) CreateTender(tender *model.Tender) (*model.Tender, error) {
	return s.repo.CreateTender(tender)
}

func (s *TenderService) GetTendersByClient(clientID int) ([]model.GetTender, error) {
	return s.repo.GetTendersByClientID(clientID)
}

func (s *TenderService) ListTenders(clientID int) ([]model.Tender, error) {
	return s.repo.ListTendersByClientID(clientID)
}

func (s *TenderService) UpdateTenderStatus(clientID int, tenderID int, status string) error {
	if status != "open" && status != "closed" && status != "awarded" {
		return errors.New("invalid status")
	}

	tender, err := s.repo.GetTenderByID(tenderID)
	if err != nil {
		return errors.New("tender not found")
	}

	if tender.ClientID != clientID {
		return errors.New("tender not found")
	}

	return s.repo.UpdateTenderStatus(tenderID, status)
}

func (s *TenderService) DeleteTender(clientID, tenderID int) error {
	tender, err := s.repo.GetTenderByID(tenderID)
	if err != nil {
		return errors.New("tender not found")
	}

	if tender.ClientID != clientID {
		return errors.New("tender not found")
	}

	return s.repo.DeleteTender(tenderID)
}
