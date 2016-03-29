package frontend

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (svc *Service) rootHandler(ginCtx *gin.Context) {
	ginCtx.JSON(http.StatusOK, gin.H{"title": "root"})
}
