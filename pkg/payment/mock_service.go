package payment

import (
	"context"
	"math/rand"
	"time"

	"github.com/Ahmed1monm/backend-golang-task-2025/pkg/logger"
	"go.uber.org/zap"
)

// PaymentInfo represents payment details
type PaymentInfo struct {
	OrderID     uint    `json:"order_id"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	CardNumber  string  `json:"card_number"`
	CardExpiry  string  `json:"card_expiry"`
	CardCVC     string  `json:"card_cvc"`
	Description string  `json:"description"`
}

// PaymentResult represents the result of a payment processing attempt
type PaymentResult struct {
	Success       bool   `json:"success"`
	TransactionID string `json:"transaction_id,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

// Service defines the interface for payment processing
type Service interface {
	ProcessPayment(ctx context.Context, info PaymentInfo) (*PaymentResult, error)
}

type mockService struct {
	rng *rand.Rand
}

// NewMockService creates a new mock payment service
func NewMockService() Service {
	return &mockService{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// ProcessPayment simulates payment processing with a random delay and success rate
func (s *mockService) ProcessPayment(ctx context.Context, info PaymentInfo) (*PaymentResult, error) {
	logger.Info(ctx, "Processing payment",
		zap.Uint("order_id", info.OrderID),
		zap.Float64("amount", info.Amount),
		zap.String("currency", info.Currency),
	)

	// Simulate processing delay (up to 3 seconds)
	delay := time.Duration(s.rng.Intn(3000)) * time.Millisecond
	select {
	case <-ctx.Done():
		logger.Error(ctx, "Payment processing cancelled", zap.Error(ctx.Err()))
		return nil, ctx.Err()
	case <-time.After(delay):
		// Continue processing
	}

	// Simulate success rate (90% success)
	if s.rng.Float64() < 0.9 {
		result := &PaymentResult{
			Success:       true,
			TransactionID: generateTransactionID(s.rng),
		}
		logger.Info(ctx, "Payment processed successfully",
			zap.String("transaction_id", result.TransactionID),
		)
		return result, nil
	}

	// Simulate failure
	result := &PaymentResult{
		Success:      false,
		ErrorMessage: "Payment declined by issuer",
	}
	logger.Error(ctx, "Payment processing failed", zap.String("error", result.ErrorMessage))
	return result, nil
}

// generateTransactionID creates a random transaction ID
func generateTransactionID(rng *rand.Rand) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 12

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}
