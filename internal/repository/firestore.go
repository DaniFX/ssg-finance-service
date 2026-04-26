package repository

import (
	"context"

	"github.com/DaniFX/ssg-finance-service/internal/models"

	"cloud.google.com/go/firestore"
)

type FinanceRepository struct {
	client *firestore.Client
}

// NewFinanceRepository inizializza la connessione a Firestore e restituisce il repository
func NewFinanceRepository(projectID string) (FinanceRepository, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return FinanceRepository{}, err
	}
	return FinanceRepository{client: client}, nil
}

// GetInvoice recupera una fattura dal database tramite ID
func (r *FinanceRepository) GetInvoice(ctx context.Context, id string) (*models.Invoice, error) {
	doc, err := r.client.Collection("invoices").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var inv models.Invoice
	doc.DataTo(&inv)
	return &inv, nil
}

// UpdateInvoice salva o aggiorna un documento fattura
func (r *FinanceRepository) UpdateInvoice(ctx context.Context, inv models.Invoice) error {
	_, err := r.client.Collection("invoices").Doc(inv.ID).Set(ctx, inv)
	return err
}

// SaveLedgerEntry registra un nuovo pagamento nel libro giornale
func (r *FinanceRepository) SaveLedgerEntry(ctx context.Context, entry models.LedgerEntry) error {
	// Se l'ID non è impostato, ne generiamo uno usando Firestore
	if entry.ID == "" {
		docRef := r.client.Collection("ledger_entries").NewDoc()
		entry.ID = docRef.ID
	}
	_, err := r.client.Collection("ledger_entries").Doc(entry.ID).Set(ctx, entry)
	return err
}

// GetTotalPaidForInvoice somma tutti i pagamenti effettuati per una specifica fattura
func (r *FinanceRepository) GetTotalPaidForInvoice(ctx context.Context, invoiceID string) (float64, error) {
	iter := r.client.Collection("ledger_entries").Where("invoiceId", "==", invoiceID).Documents(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return 0, err
	}

	var totalPaid float64
	for _, doc := range docs {
		var entry models.LedgerEntry
		doc.DataTo(&entry)
		totalPaid += entry.Amount
	}

	return totalPaid, nil
}
