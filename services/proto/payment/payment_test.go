package payment_test

import (
	context "context"
	"testing"

	"github.com/gonzabosio/transaction/services/proto/payment"
	"github.com/gonzabosio/transaction/services/proto/payment/handlers"
	"github.com/stretchr/testify/assert"
)

func TestPaymentServiceCheckout(t *testing.T) {
	s := &handlers.PaymentService{ApiBaseUrl: "https://2kgy8.wiremockapi.cloud"}
	res, err := s.CheckoutOrder(context.Background(), &payment.CheckoutRequest{AccessToken: "abcd1234", OrderId: "OGI32JR23OAL"})
	if err != nil {
		t.Fatalf("checkout internal: %v", err)
	}
	assert.Equal(t, "COMPLETED", res.OrderStatus)
}
