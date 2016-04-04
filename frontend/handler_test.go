package frontend

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func getDefaultService() *Service {
	return &Service{
		Healthy: true,
	}
}

func TestService_HealthHandler(t *testing.T) {
	ctx, _, _ := gin.CreateTestContext()
	svc := getDefaultService()

	svc.HealthHandler(ctx)

	if ctx.Writer.Status() != 200 {
		t.Fatal("Wrong status code")
	}

	svc.Healthy = false
	svc.HealthHandler(ctx)
	if ctx.Writer.Status() == 200 {
		t.Fatal("Wrong status code")
	}
}

func TestService_RootHandler(t *testing.T) {
	ctx, w, router := gin.CreateTestContext()
	svc := getDefaultService()

	svc.RootHandler(ctx)

	if ctx.Writer.Status() != 200 {
		t.Fatal("Wrong status code")
	}
	_ = w
	_ = router
}

func TestService_IsHealthy(t *testing.T) {
	svc := getDefaultService()
	if !svc.IsHealthy() {
		t.Fatal("err: Service is not healthy")
	}
}
