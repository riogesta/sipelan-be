package repository

import (
	"sipelan/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *models.Transaction) error
	Update(tx *models.Transaction) error
	Delete(id uint, personID uint) error
	GetByID(id uint, personID uint) (*models.Transaction, error)
	GetAll(personID uint, filters map[string]string, limit, offset int) ([]models.Transaction, int64, error)
	GetCategorySummary(personID uint, month, year string) ([]models.CategorySummary, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) Update(tx *models.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *transactionRepository) Delete(id uint, personID uint) error {
	return r.db.Where("id = ? AND person_id = ?", id, personID).Delete(&models.Transaction{}).Error
}

func (r *transactionRepository) GetByID(id uint, personID uint) (*models.Transaction, error) {
	var tx models.Transaction
	err := r.db.Preload("Category").Where("id = ? AND person_id = ?", id, personID).First(&tx).Error
	return &tx, err
}

func (r *transactionRepository) GetAll(personID uint, filters map[string]string, limit, offset int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	query := r.db.Model(&models.Transaction{}).Where("person_id = ?", personID)

	if val, ok := filters["type"]; ok && val != "" {
		query = query.Where("type = ?", val)
	}
	if val, ok := filters["start_date"]; ok && val != "" {
		query = query.Where("date >= ?", val)
	}
	if val, ok := filters["end_date"]; ok && val != "" {
		query = query.Where("date <= ?", val)
	}
	if val, ok := filters["search"]; ok && val != "" {
		query = query.Where("description LIKE ?", "%"+val+"%")
	}
	if val, ok := filters["category_id"]; ok && val != "" {
		query = query.Where("category_id = ?", val)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Category", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Limit(limit).Offset(offset).Order("date DESC").Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) GetCategorySummary(personID uint, month, year string) ([]models.CategorySummary, error) {
	var results []models.CategorySummary

	err := r.db.Table("transactions").
		Select("categories.name as name, SUM(transactions.total) as value").
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.person_id = ? AND transactions.type = ? AND strftime('%m', transactions.date, 'localtime') = ? AND strftime('%Y', transactions.date, 'localtime') = ? AND transactions.deleted_at IS NULL", personID, "pengeluaran", month, year).
		Group("transactions.category_id").
		Order("value DESC").
		Scan(&results).Error

	return results, err
}
