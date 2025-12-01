package query

import (
	"context"
	"log/slog"
	"mini-erp-backend/api/repository"
	"mini-erp-backend/model"

	"gorm.io/gorm"
)

type SearchSuppliers struct {
	logger       *slog.Logger
	db           *gorm.DB
	SupplierRepo repository.SupplierRepository
}

type SearchSuppliersRequest struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	OrderBy string `json:"order_by"`
}

type SearchSuppliersResult struct {
	Suppliers []*model.Supplier `json:"suppliers"`
}

func NewSearchSuppliersHandler(logger *slog.Logger, db *gorm.DB, repo repository.SupplierRepository) *SearchSuppliers {
	return &SearchSuppliers{
		logger:       logger,
		db:           db,
		SupplierRepo: repo,
	}
}

func (h *SearchSuppliers) Handle(ctx context.Context, req *SearchSuppliersRequest) (interface{}, error) {
	// Build search conditions
	conditions := make(map[string]interface{})
	
	if req.Email != "" {
		conditions["email"] = req.Email
	}
	if req.Name != "" {
		conditions["name"] = req.Name
	}

	// Set default order by
	orderBy := req.OrderBy
	if orderBy == "" {
		orderBy = "name ASC"
	}

	// Search suppliers from database
	suppliers, err := h.SupplierRepo.Searches(h.db, conditions, orderBy)
	if err != nil {
		h.logger.Error("Failed to search suppliers", "error", err)
		return nil, err
	}

	return SearchSuppliersResult{Suppliers: suppliers}, nil
}
