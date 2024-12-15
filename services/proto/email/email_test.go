package email_test

import (
	"context"
	"os"
	"testing"

	pb "github.com/gonzabosio/transaction/services/proto/email"
	"github.com/gonzabosio/transaction/services/proto/email/handlers"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestEmailService(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal(err)
	}
	s := handlers.EmailService{}
	res, err := s.SendEmail(context.Background(), &pb.EmailRequest{
		Subject:    "Testing Email",
		BodyText:   "This is a testing body text",
		PayerEmail: os.Getenv("PAYER_TEST_EMAIL"),
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Email was sent successfully", res.Message)
}
