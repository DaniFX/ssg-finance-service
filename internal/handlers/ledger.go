package handlers

import (
	"finance-service/internal/models"
	"finance-service/internal/services"
	"ssg-nexus-sdk/pkg/nexus"

	"github.com/gin-gonic/gin"
)

// RegisterTransaction gestisce la POST /ledger
func RegisterTransaction(svc *services.FinanceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entry models.LedgerEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			nexus.ErrorResponse(c, 400, "INVALID_PAYLOAD", "Formato richiesta non valido")
			return
		}

		// Valida dati minimi
		if entry.InvoiceID == "" || entry.Amount <= 0 {
			nexus.ErrorResponse(c, 400, "VALIDATION_ERROR", "InvoiceID o Amount mancanti/non validi")
			return
		}

		err := svc.RegisterPayment(c.Request.Context(), entry)
		if err != nil {
			nexus.ErrorResponse(c, 500, "LEDGER_ERROR", "Impossibile registrare il pagamento e riconciliare")
			return
		}

		nexus.SuccessResponse(c, map[string]string{
			"message": "Pagamento registrato, riconciliazione effettuata",
		})
	}
}
