package handlers

import (
	"github.com/DaniFX/ssg-finance-service/internal/models"
	"github.com/DaniFX/ssg-finance-service/internal/services"
	"github.com/DaniFX/ssg-nexus-sdk/pkg/nexus"
	"github.com/gin-gonic/gin"
)

// RegisterTransaction gestisce la POST /ledger
func RegisterTransaction(svc *services.FinanceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entry models.LedgerEntry

		if err := c.ShouldBindJSON(&entry); err != nil {
			nexus.Failure(c, 400, "INVALID_PAYLOAD", "Formato richiesta non valido", err.Error())
			return
		}

		// Valida dati minimi
		if entry.InvoiceID == "" || entry.Amount <= 0 {
			nexus.Failure(c, 400, "VALIDATION_ERROR", "InvoiceID o Amount mancanti/non validi", nil)
			return
		}

		err := svc.RegisterPayment(c.Request.Context(), entry)
		if err != nil {
			nexus.Failure(c, 500, "LEDGER_ERROR", "Impossibile registrare il pagamento e riconciliare", err.Error())
			return
		}

		nexus.Success(c, map[string]string{
			"message": "Pagamento registrato, riconciliazione effettuata",
		}, nil)
	}
}
