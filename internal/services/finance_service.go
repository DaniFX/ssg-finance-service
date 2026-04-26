// internal/services/finance_service.go
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DaniFX/ssg-finance-service/internal/models"
	"github.com/DaniFX/ssg-finance-service/internal/repository"
	"github.com/DaniFX/ssg-nexus-sdk/pkg/nexus"
)

type FinanceService struct {
	repo          repository.FinanceRepository
	nexusClient   *nexus.NexusClient
	docServiceURL string
}

func NewFinanceService(repo repository.FinanceRepository, nc *nexus.NexusClient, docURL string) *FinanceService {
	return &FinanceService{
		repo:          repo,
		nexusClient:   nc,
		docServiceURL: docURL,
	}
}

// UpdateInvoice aggiorna una fattura (utile per i DRAFT). Implementa il locking.
func (s *FinanceService) UpdateInvoice(ctx context.Context, inv models.Invoice) error {
	existing, err := s.repo.GetInvoice(ctx, inv.ID)
	if err != nil {
		return err
	}

	// Regola ERP: Se ISSUED, il documento è immutabile e non può essere aggiornato manualmente
	if existing.Status == models.StatusIssued || existing.Metadata["immutable"] == true {
		return errors.New("cannot update an issued invoice")
	}

	return s.repo.UpdateInvoice(ctx, inv)
}

// IssueInvoice implementa il Locking e la chiamata al Document Service per l'emissione
func (s *FinanceService) IssueInvoice(ctx context.Context, invoiceID string) error {
	invoice, err := s.repo.GetInvoice(ctx, invoiceID)
	if err != nil {
		return err
	}

	// 1. Regola ERP (Locking): Verifica stato e immutabilità
	if invoice.Status == models.StatusIssued || (invoice.Metadata != nil && invoice.Metadata["immutable"] == true) {
		return errors.New("il documento è già emesso e non può essere modificato")
	}

	// 2. Generazione numero sequenziale (Semplificata per l'esempio)
	invoice.ExternalID = fmt.Sprintf("SDI-%d-%s", time.Now().Year(), invoice.ID[:6])

	// 3. Chiamata HTTP al Document Service per generare/salvare il PDF
	docRef, err := s.generateDocument(ctx, invoice)
	if err != nil {
		return fmt.Errorf("errore generazione documento nel vault: %v", err)
	}

	// 4. Aggiornamento stato e metadata
	invoice.DocumentRef = docRef
	invoice.Status = models.StatusIssued
	if invoice.Metadata == nil {
		invoice.Metadata = make(map[string]any)
	}
	invoice.Metadata["immutable"] = true // Applicazione del Lock

	// 5. Salvataggio in Firestore
	return s.repo.UpdateInvoice(ctx, *invoice)
}

// RegisterPayment registra l'entrata nel ledger e applica la Riconciliazione (Ex Reconcile)
func (s *FinanceService) RegisterPayment(ctx context.Context, entry models.LedgerEntry) error {
	entry.Timestamp = time.Now()

	// 1. Salva la transazione finanziaria
	if err := s.repo.SaveLedgerEntry(ctx, entry); err != nil {
		return err
	}

	// 2. Logica di Riconciliazione (Sostituisce il vecchio metodo Reconcile)
	totalPaid, err := s.repo.GetTotalPaidForInvoice(ctx, entry.InvoiceID)
	if err != nil {
		return err
	}

	invoice, err := s.repo.GetInvoice(ctx, entry.InvoiceID)
	if err != nil {
		return err
	}

	// Se la somma dei pagamenti copre o supera il lordo, passa a PAID
	if totalPaid >= invoice.Totals.Gross && invoice.Status != models.StatusPaid {
		invoice.Status = models.StatusPaid
		now := time.Now()
		invoice.Dates.Paid = &now

		// Aggiorna lo stato della fattura
		return s.repo.UpdateInvoice(ctx, *invoice)
	}

	return nil
}

// generateDocument comunica con ssg-nexus-document-service
func (s *FinanceService) generateDocument(ctx context.Context, invoice *models.Invoice) (string, error) {
	payload, _ := json.Marshal(map[string]interface{}{
		"type": "INVOICE",
		"data": invoice,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", s.docServiceURL+"/api/v1/documents/generate", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.nexusClient.Do(ctx, req)
	if err != nil || resp.StatusCode >= 400 {
		return "", errors.New("fallita chiamata al Document Service")
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			DocumentID string `json:"documentId"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.DocumentID, nil
}
