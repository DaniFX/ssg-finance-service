package main

import (
	"finance-service/internal/handlers"
	"finance-service/internal/repository"
	"log"
	"ssg-nexus-sdk/pkg/nexus"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Inizializzazione Repository Firestore (Seguendo il pattern ssg-db)
	repo, err := repository.NewFirestoreRepository("project-id")
	if err != nil {
		log.Fatalf("Errore Firestore: %v", err)
	}

	// Middleware globale di sicurezza dall'SDK
	r.Use(nexus.NexusGuard())

	// Gruppo API V1
	v1 := r.Group("/api/v1")
	{
		// Invoices
		v1.POST("/invoices", handlers.CreateInvoice(repo))
		v1.PATCH("/invoices/:id/issue", handlers.IssueInvoice(repo))

		// Ledger (Libro Giornale)
		v1.POST("/ledger", handlers.RegisterTransaction(repo))
	}

	// Endpoint richiesto per il Service Discovery del Gateway
	r.GET("/_discover", handlers.GetDiscovery)

	log.Println("Finance Service in ascolto sulla porta 8080...")
	r.Run(":8080")
}
