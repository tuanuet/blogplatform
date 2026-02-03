package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aiagent/internal/infrastructure/config"
)

// CreateVietQRRequest represents the request to create a VietQR code
type CreateVietQRRequest struct {
	AccountNo   string `json:"accountNo"`
	AccountName string `json:"accountName"`
	AcqId       string `json:"acqId"`
	Amount      int    `json:"amount"`
	AddInfo     string `json:"addInfo"`
	Format      string `json:"format"`
	Template    string `json:"template"`
}

// VietQRResponse represents the response from SePay for VietQR creation
type VietQRResponse struct {
	Status int         `json:"status"`
	Error  interface{} `json:"error"`
	Data   struct {
		QrCode    string `json:"qrCode"`
		QrDataURL string `json:"qrDataURL"`
	} `json:"data"`
}

// BankTransferInfo represents the bank account details for manual transfer
type BankTransferInfo struct {
	BankName    string `json:"bankName"`
	AccountNo   string `json:"accountNo"`
	AccountName string `json:"accountName"`
	Branch      string `json:"branch"`
}

// SePayTransaction represents a transaction record from SePay
type SePayTransaction struct {
	ID            string  `json:"id"`
	BankName      string  `json:"bankName"`
	AccountNo     string  `json:"accountNo"`
	Amount        float64 `json:"amount"`
	Content       string  `json:"content"`
	TransferDate  string  `json:"transferDate"`
	ReferenceCode string  `json:"referenceCode"`
}

// SePayTransactionResponse represents the response from SePay for transaction details
type SePayTransactionResponse struct {
	Status int              `json:"status"`
	Error  interface{}      `json:"error"`
	Data   SePayTransaction `json:"data"`
}

// SePayAdapter defines the interface for SePay payment gateway operations
type SePayAdapter interface {
	CreateVietQR(ctx context.Context, req CreateVietQRRequest) (*VietQRResponse, error)
	GetBankTransferInfo() BankTransferInfo
	VerifyWebhookSignature(payload map[string]interface{}, signature string) bool
	GetTransactionBySePayID(ctx context.Context, sepayID string) (*SePayTransaction, error)
}

type sePayAdapter struct {
	config  *config.SePayConfig
	client  *http.Client
	baseURL string
}

// NewSePayAdapter creates a new instance of SePayAdapter
func NewSePayAdapter(cfg *config.SePayConfig) SePayAdapter {
	return &sePayAdapter{
		config: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.sepay.vn/v1",
	}
}

// CreateVietQR generates a VietQR code for payment
func (a *sePayAdapter) CreateVietQR(ctx context.Context, req CreateVietQRRequest) (*VietQRResponse, error) {
	url := fmt.Sprintf("%s/qr/create", a.baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+a.config.APIKey)

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sepay api returned status: %d", resp.StatusCode)
	}

	var result VietQRResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetBankTransferInfo returns the bank account information from configuration
func (a *sePayAdapter) GetBankTransferInfo() BankTransferInfo {
	return BankTransferInfo{
		BankName:    a.config.BankName,
		AccountNo:   a.config.BankAccount,
		AccountName: a.config.BankOwner,
		Branch:      a.config.BankBranch,
	}
}

// VerifyWebhookSignature validates the authenticity of the webhook request
func (a *sePayAdapter) VerifyWebhookSignature(payload map[string]interface{}, signature string) bool {
	// SePay webhook token verification
	// SePay sends the token in the Authorization header as "Bearer <token>"
	// or in some cases as a query parameter. The caller should pass the token.
	return a.config.WebhookToken != "" && a.config.WebhookToken == signature
}

// GetTransactionBySePayID retrieves transaction details from SePay API
func (a *sePayAdapter) GetTransactionBySePayID(ctx context.Context, sepayID string) (*SePayTransaction, error) {
	url := fmt.Sprintf("%s/transactions/details/%s", a.baseURL, sepayID)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+a.config.APIKey)

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sepay api returned status: %d", resp.StatusCode)
	}

	var result SePayTransactionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Data, nil
}
