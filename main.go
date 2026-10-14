package main

import (
	"log"
	"github.com/example/proofcraft/pkg/db"
	"github.com/example/proofcraft/pkg/controllers"
	"github.com/example/proofcraft/pkg/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize DB
	dbConn, err := db.Init("proofcraft.db")
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	// Auto migrate models
	err = dbConn.AutoMigrate(&models.Key{}, &models.Proof{})
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}

	// Setup router
	r := gin.Default()

	// Controllers
	ctrl := controllers.NewController(dbConn)

	r.POST("/keys/generate", ctrl.GenerateKey)
	r.POST("/proofs/prove", ctrl.CreateProof)
	r.POST("/proofs/verify", ctrl.VerifyProof)

	log.Println("ProofCraft server starting on :8080")
	r.Run(":8080")
}
