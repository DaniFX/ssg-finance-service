package handlers

import (
	"github.com/DaniFX/ssg-finance-service/internal/services"
	"github.com/DaniFX/ssg-nexus-sdk/pkg/nexus"
	"github.com/gin-gonic/gin"
)

// IssueInvoice gestisce la richiesta PATCH /invoices/:id/issue
func IssueInvoice(svc *services.FinanceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceID := c.Param("id")

		err := svc.IssueInvoice(c.Request.Context(), invoiceID)
		if err != nil {
			// Usiamo nexus.Failure passando 'nil' per i details se non ne abbiamo
			nexus.Failure(c, 400, "ERP_LOCK_ERROR", err.Error(), nil)
			return
		}

		// Usiamo nexus.Success passando 'nil' per i meta
		nexus.Success(c, map[string]string{
			"status":  "ISSUED",
			"message": "Fattura emessa e bloccata con successo",
		}, nil)
	}
}
