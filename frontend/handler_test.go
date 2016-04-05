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
	gin.SetMode(gin.ReleaseMode)
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
	gin.SetMode(gin.ReleaseMode)
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

func Benchmark_IsHealthy(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	svc := getDefaultService()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		svc.IsHealthy()
	}
}

func Benchmark_RootHandler(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	svc := getDefaultService()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ctx, _, _ := gin.CreateTestContext()
		svc.RootHandler(ctx)
	}
}

func Benchmark_HealthHandler(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)
	svc := getDefaultService()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		ctx, _, _ := gin.CreateTestContext()
		b.StartTimer()
		svc.HealthHandler(ctx)
	}
}
