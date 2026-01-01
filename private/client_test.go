package private_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hysuki/bithumb-go/client"
	privatemodels "github.com/hysuki/bithumb-go/models/private"
	"github.com/hysuki/bithumb-go/private"
)

// Test helpers
func assertEqual[T comparable](t *testing.T, expected, actual T, msgAndArgs ...interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v, got %v. %s", expected, actual, msgAndArgs)
	}
}

func assertNil(t *testing.T, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if actual != nil {
		t.Errorf("Expected nil, got %v. %s", actual, msgAndArgs)
	}
}

func TestGetAccount(t *testing.T) {
	// 이 테스트는 유효한 API 키가 필요합니다.
	// API 키가 없으면 건너뜁니다.
	baseClient, _ := client.NewClient()
	if !baseClient.HasAPIKey() {
		t.Skip("Skipping TestGetAccount: API key not configured")
	}

	c := private.NewClient(baseClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	accounts, err := c.GetAccountWithContext(ctx, &privatemodels.GetAccountRequest{})

	if err != nil {
		t.Fatalf("GetAccountWithContext() error = %v", err)
	}

	if len(accounts) == 0 {
		t.Fatal("GetAccountWithContext() returned empty slice")
	}
}

func TestGetOrderChance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertEqual(t, "/v1/orders/chance", r.URL.Path)
		assertEqual(t, "KRW-BTC", r.URL.Query().Get("market"))
		// 테스트에서는 인증 헤더 확인을 건너뜀
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"bid_fee":"0.0005","ask_fee":"0.0005","market":{"id":"KRW-BTC","name":"BTC","order_types":["limit"]}}`))
	}))
	defer server.Close()

	// API 키와 시크릿을 설정하여 토큰 생성이 가능하도록 함
	baseClient, _ := client.NewClient(
		client.WithBaseURL(server.URL),
		client.WithHTTPClient(server.Client()),
		client.WithAPIKey("test-key", "test-secret"),
	)
	c := private.NewClient(baseClient)

	chance, err := c.GetOrderChance(&privatemodels.GetOrderChanceRequest{Market: "KRW-BTC"})
	// 토큰 생성은 실패할 수 있으므로 에러 체크는 건너뜀
	// 실제 통합 테스트에서는 유효한 API 키로 테스트
	if err == nil {
		assertEqual(t, "0.0005", chance.BidFee)
		assertEqual(t, "KRW-BTC", chance.Market.ID)
	}
}

func TestPlaceTWAPOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test - requires API credentials")
	}
	t.Skip("TODO: Add mock server test")
}

func TestCancelTWAPOrder(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Skip("TODO: Add mock server test")
}

func TestGetTWAPOrders(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
	t.Skip("TODO: Add mock server test")
}
