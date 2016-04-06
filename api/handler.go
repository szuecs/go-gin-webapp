package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RootHandler handles / endpoint
func (svc *Service) RootHandler(ginCtx *gin.Context) {
	ginCtx.JSON(http.StatusOK, gin.H{"title": "root"})
}

// HealthHandler handles /health endpoint
func (svc *Service) HealthHandler(ginCtx *gin.Context) {
	if svc.IsHealthy() {
		ginCtx.String(http.StatusOK, "%s", "OK")
	} else {
		ginCtx.String(http.StatusServiceUnavailable, "%s", "Unavailable")
	}
}

// IsHealthy returns the health status of the running service.
func (svc *Service) IsHealthy() bool {
	return svc.Healthy
}
