package services

import (
	"context"
	"errors"
	"finance-service/internal/models"
	"finance-service/internal/repository"
)

type FinanceService struct {
	repo repository.FinanceRepository
}

// UpdateInvoice implementa il locking richiesto
func (s *FinanceService) UpdateInvoice(ctx context.Context, inv models.Invoice) error {
	existing, err := s.repo.GetInvoice(ctx, inv.ID)
	if err != nil {
		return err
	}

	// Regola ERP: Se ISSUED, il documento è immutabile
	if existing.Status == models.StatusIssued || existing.Metadata["immutable"] == true {
		return errors.New("cannot update an issued invoice")
	}

	return s.repo.SaveInvoice(ctx, inv)
}

// Reconcile controlla se i pagamenti coprono il totale
func (s *FinanceService) Reconcile(ctx context.Context, invoiceID string) error {
	totalPaid, err := s.repo.GetTotalPaid(ctx, invoiceID)
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)

	if totalPaid >= invoice.Totals.Gross {
		invoice.Status = models.StatusPaid
		return s.repo.UpdateStatus(ctx, invoiceID, models.StatusPaid)
	}
	return nil
}
