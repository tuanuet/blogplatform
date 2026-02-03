package adapter

import (
	"testing"

	"context"
	"encoding/json"
	"github.com/aiagent/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
)

func TestSePayAdapter_CreateVietQR(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/qr/create", r.URL.Path)
		assert.Equal(t, "Bearer test_key", r.Header.Get("Authorization"))

		resp := VietQRResponse{
			Status: 200,
		}
		resp.Data.QrCode = "test_qr_code"
		resp.Data.QrDataURL = "test_url"

		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := &config.SePayConfig{
		APIKey: "test_key",
	}
	adapter := NewSePayAdapter(cfg).(*sePayAdapter)
	adapter.baseURL = server.URL

	req := CreateVietQRRequest{
		Amount: 10000,
	}
	res, err := adapter.CreateVietQR(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "test_qr_code", res.Data.QrCode)
}

func TestSePayAdapter_GetTransactionBySePayID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/transactions/details/123", r.URL.Path)
		assert.Equal(t, "Bearer test_key", r.Header.Get("Authorization"))

		resp := SePayTransactionResponse{
			Status: 200,
			Data: SePayTransaction{
				ID:     "123",
				Amount: 10000,
			},
		}

		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := &config.SePayConfig{
		APIKey: "test_key",
	}
	adapter := NewSePayAdapter(cfg).(*sePayAdapter)
	adapter.baseURL = server.URL

	res, err := adapter.GetTransactionBySePayID(context.Background(), "123")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "123", res.ID)
	assert.Equal(t, float64(10000), res.Amount)
}

func TestSePayAdapter_GetBankTransferInfo(t *testing.T) {
	cfg := &config.SePayConfig{
		BankName:    "Test Bank",
		BankAccount: "123456789",
		BankOwner:   "Test Owner",
		BankBranch:  "Test Branch",
	}
	adapter := NewSePayAdapter(cfg)

	info := adapter.GetBankTransferInfo()

	assert.Equal(t, cfg.BankName, info.BankName)
	assert.Equal(t, cfg.BankAccount, info.AccountNo)
	assert.Equal(t, cfg.BankOwner, info.AccountName)
	assert.Equal(t, cfg.BankBranch, info.Branch)
}

func TestSePayAdapter_VerifyWebhookSignature(t *testing.T) {
	cfg := &config.SePayConfig{
		WebhookToken: "test_token",
	}
	adapter := NewSePayAdapter(cfg)

	t.Run("Valid Signature", func(t *testing.T) {
		// SePay sends the token in the Authorization header as "Bearer <token>"
		// The signature parameter in VerifyWebhookSignature should be the token itself
		payload := map[string]interface{}{"id": "123"}
		signature := "test_token"

		assert.True(t, adapter.VerifyWebhookSignature(payload, signature))
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		payload := map[string]interface{}{"id": "123"}
		signature := "wrong_token"

		assert.False(t, adapter.VerifyWebhookSignature(payload, signature))
	})
}
