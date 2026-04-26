// internal/models/ledger.go
package models

import "time"

type LedgerEntry struct {
	ID        string    `json:"id" firestore:"id"`
	EntityID  string    `json:"entityId" firestore:"entityId"`
	InvoiceID string    `json:"invoiceId" firestore:"invoiceId"`
	Amount    float64   `json:"amount" firestore:"amount"`
	Type      string    `json:"type" firestore:"type"`     // DEBIT | CREDIT
	Method    string    `json:"method" firestore:"method"` // STRIPE | BANK_TRANSFER
	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
}
