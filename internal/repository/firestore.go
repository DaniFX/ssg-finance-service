package repository

import (
	"context"
	"finance-service/internal/models"

	"cloud.google.com/go/firestore"
)

type FinanceRepository struct {
	client *firestore.Client
}

// Crea una fattura in stato DRAFT
func (r *FinanceRepository) CreateInvoice(ctx context.Context, inv models.Invoice) error {
	_, err := r.client.Collection("invoices").Doc(inv.ID).Set(ctx, inv)
	return err
}

// Recupera una fattura per ID
func (r *FinanceRepository) GetInvoice(ctx context.Context, id string) (*models.Invoice, error) {
	doc, err := r.client.Collection("invoices").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var inv models.Invoice
	doc.DataTo(&inv)
	return &inv, nil
}
