package handlers

import (
	"finance-service/internal/services"
	"ssg-nexus-sdk/pkg/nexus"

	"github.com/gin-gonic/gin"
)

// IssueInvoice gestisce la richiesta PATCH /invoices/:id/issue
func IssueInvoice(svc *services.FinanceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceID := c.Param("id")

		// Recupera l'identità dell'utente dal contesto iniettata dal NexusGuard (opzionale per audit log)
		// identity := nexus.FromContext(c.Request.Context())

		err := svc.IssueInvoice(c.Request.Context(), invoiceID)
		if err != nil {
			// Utilizza lo standard responder del Nexus SDK
			nexus.ErrorResponse(c, 400, "ERP_LOCK_ERROR", err.Error())
			return
		}

		nexus.SuccessResponse(c, map[string]string{
			"status":  "ISSUED",
			"message": "Fattura emessa e bloccata con successo",
		})
	}
}
