package handler

import (
	"net/http"
	"virtpanel/internal/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) ListBridges(c *gin.Context) {
	bridges, err := h.svc.ListBridges()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bridges)
}

func (h *Handler) CreateBridge(c *gin.Context) {
	var req model.CreateBridgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateBridge(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

func (h *Handler) DeleteBridge(c *gin.Context) {
	if err := h.svc.DeleteBridge(c.Param("name")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
