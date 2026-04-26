package handlers

import "github.com/gin-gonic/gin"

// GetDiscovery espone le rotte al Gateway (Standard SSG Nexus)
func GetDiscovery(c *gin.Context) {
	discovery := map[string]interface{}{
		"serviceName": "finance-service",
		"version":     "1.0.0",
		"endpoints": []map[string]interface{}{
			{
				"path":         "/api/v1/invoices/:id/issue",
				"method":       "PATCH",
				"summary":      "Emette una fattura e la blocca",
				"authRequired": true,
			},
			{
				"path":         "/api/v1/ledger",
				"method":       "POST",
				"summary":      "Registra un pagamento e applica la riconciliazione",
				"authRequired": true,
			},
		},
	}
	c.JSON(200, discovery)
}
