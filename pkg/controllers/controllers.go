package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/example/proofcraft/pkg/services"
)

type Controller struct {
	svc *services.Service
}

func NewController(db *gorm.DB) *Controller {
	svc := services.NewService(db)
	return &Controller{svc: svc}
}

// GenerateKey - generate a Schnorr keypair
func (c *Controller) GenerateKey(ctx *gin.Context) {
	var req struct {
		Label string `json:"label"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	res, err := c.svc.GenerateKey(req.Label)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// CreateProof - create a Schnorr proof using a provided private key
func (c *Controller) CreateProof(ctx *gin.Context) {
	var req struct {
		Message    string `json:"message"`
		PrivateHex string `json:"private_key"`
		Label      string `json:"label"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	res, err := c.svc.CreateProof(req.Message, req.PrivateHex, req.Label)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// VerifyProof - verify proof against provided public key
func (c *Controller) VerifyProof(ctx *gin.Context) {
	var req struct {
		Message string `json:"message"`
		Proof   struct {
			Rx  string `json:"rx"`
			Ry  string `json:"ry"`
			S   string `json:"s"`
		} `json:"proof"`
		PublicX string `json:"public_x"`
		PublicY string `json:"public_y"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	ok, err := c.svc.VerifyProof(req.Message, req.Proof.Rx, req.Proof.Ry, req.Proof.S, req.PublicX, req.PublicY)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"valid": ok})
}
