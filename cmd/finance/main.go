package main

import (
	"log"
	"os"

	"github.com/DaniFX/ssg-nexus-sdk/pkg/nexus"
	"github.com/gin-gonic/gin"

	"github.com/DaniFX/ssg-finance-service/internal/handlers"
	"github.com/DaniFX/ssg-finance-service/internal/repository"
	"github.com/DaniFX/ssg-finance-service/internal/services"
)

func main() {
	// 1. Configurazione Ambiente
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Println("WARNING: GOOGLE_CLOUD_PROJECT non settato. Usando default.")
		projectID = "ssg-nexus-dev"
	}

	// 2. Inizializzazione Repository e Service
	repo, err := repository.NewFinanceRepository(projectID)
	if err != nil {
		log.Fatalf("Errore avvio Firestore: %v", err)
	}

	// Inizializziamo il NexusClient dell'SDK per le chiamate inter-service
	nexusClient := &nexus.NexusClient{}
	docServiceURL := os.Getenv("DOCUMENT_SERVICE_URL")

	financeSvc := services.NewFinanceService(repo, nexusClient, docServiceURL)

	r := gin.Default()

	// 3. Service Definition e Handshake (Novità)
	// Definiamo il contratto del servizio per il Gateway
	def := nexus.ServiceDefinition{
		ServiceName: "finance-service",
		Version:     "1.0.0",
		Endpoints: []nexus.Endpoint{
			{
				Path:         "/api/v1/finance/invoices/:id/issue",
				Method:       "PATCH",
				AuthRequired: true,
				Summary:      "Emette una fattura e la rende immutabile",
			},
			{
				Path:         "/api/v1/finance/ledger",
				Method:       "POST",
				AuthRequired: true,
				Summary:      "Registra un pagamento nel libro giornale",
			},
		},
	}

	// Registra automaticamente l'endpoint GET /_discover
	nexus.RegisterDiscovery(r, def)

	// Avvia l'handshake (PUSH) verso il Gateway per la registrazione dinamica
	nexus.StartGatewayHandshake(def)

	// 4. Rotte di Business (PROTETTE DAL GUARD)
	api := r.Group("/api/v1/finance")

	// Applichiamo il middleware di sicurezza dell'SDK
	api.Use(nexus.Guard())
	{
		api.PATCH("/invoices/:id/issue", handlers.IssueInvoice(financeSvc))
		api.POST("/ledger", handlers.RegisterTransaction(financeSvc))
	}

	log.Printf("Finance Service avviato sulla porta %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Errore critico: %v", err)
	}
}
