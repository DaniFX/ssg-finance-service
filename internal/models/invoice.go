package models

type InvoiceStatus string

const (
	StatusDraft     InvoiceStatus = "DRAFT"
	StatusIssued    InvoiceStatus = "ISSUED"
	StatusPaid      InvoiceStatus = "PAID"
	StatusCancelled InvoiceStatus = "CANCELLED"
)

type Invoice struct {
	ID          string         `json:"id" firestore:"id"`
	ExternalID  string         `json:"externalId" firestore:"externalId"`
	Type        string         `json:"type" firestore:"type"`
	Status      InvoiceStatus  `json:"status" firestore:"status"`
	Issuer      Entity         `json:"issuer" firestore:"issuer"`     // Copia da Registry
	Receiver    Entity         `json:"receiver" firestore:"receiver"` // Copia da Registry
	Totals      Totals         `json:"totals" firestore:"totals"`
	DocumentRef string         `json:"documentRef" firestore:"documentRef"` // Link a Document Service
	Metadata    map[string]any `json:"metadata" firestore:"metadata"`
}

type Totals struct {
	Gross    float64 `json:"gross" firestore:"gross"`
	Currency string  `json:"currency" firestore:"currency"`
}

type Entity struct {
	EntityID string `json:"entityId" firestore:"entityId"`
	Name     string `json:"name" firestore:"name"`
	VAT      string `json:"vat" firestore:"vat"`
}
