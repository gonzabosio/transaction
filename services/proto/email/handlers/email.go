package handlers

import (
	"context"
	"fmt"
	"net/smtp"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/email"
	"github.com/jordan-wright/email"
)

type EmailService struct {
	pb.UnimplementedEmailServiceServer
}

func (e *EmailService) SendEmail(ctx context.Context, req *pb.EmailRequest) (*pb.Result, error) {
	bsnEmail := os.Getenv("BUSINESS_EMAIL")
	em := email.Email{
		From:    fmt.Sprintf("Fictitious Bussiness <%v>", bsnEmail),
		To:      []string{req.PayerEmail},
		Subject: req.Subject,
		Text:    []byte(req.BodyText),
	}

	if err := em.Send("smtp.gmail.com:587", smtp.PlainAuth("", bsnEmail, os.Getenv("BUSINESS_EMAIL_PW"), "smtp.gmail.com")); err != nil {
		return nil, err
	}
	return &pb.Result{Message: "Email was sent successfully"}, nil
}
