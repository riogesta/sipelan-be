package service

import (
	"sipelan/internal/repository"
	"sipelan/models"
)

type TransactionService interface {
	CreateTransaction(tx *models.Transaction) error
	UpdateTransaction(id uint, personID uint, input *models.Transaction) (*models.Transaction, error)
	DeleteTransaction(id uint, personID uint) error
	GetTransactionByID(id uint, personID uint) (*models.Transaction, error)
	GetAllTransactions(personID uint, filters map[string]string, limit, offset int) ([]models.Transaction, int64, error)
	GetCategorySummary(personID uint, month, year string) ([]models.CategorySummary, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo}
}

func (s *transactionService) CreateTransaction(tx *models.Transaction) error {
	// Business logic: check if category exists and belongs to user (optional repo check)
	return s.repo.Create(tx)
}

func (s *transactionService) UpdateTransaction(id uint, personID uint, input *models.Transaction) (*models.Transaction, error) {
	tx, err := s.repo.GetByID(id, personID)
	if err != nil {
		return nil, err
	}

	tx.Date = input.Date
	tx.CategoryID = input.CategoryID
	tx.Description = input.Description
	tx.Total = input.Total
	tx.Type = input.Type
	tx.Attachment = input.Attachment

	err = s.repo.Update(tx)
	return tx, err
}

func (s *transactionService) DeleteTransaction(id uint, personID uint) error {
	return s.repo.Delete(id, personID)
}

func (s *transactionService) GetTransactionByID(id uint, personID uint) (*models.Transaction, error) {
	return s.repo.GetByID(id, personID)
}

func (s *transactionService) GetAllTransactions(personID uint, filters map[string]string, limit, offset int) ([]models.Transaction, int64, error) {
	return s.repo.GetAll(personID, filters, limit, offset)
}

func (s *transactionService) GetCategorySummary(personID uint, month, year string) ([]models.CategorySummary, error) {
	return s.repo.GetCategorySummary(personID, month, year)
}
