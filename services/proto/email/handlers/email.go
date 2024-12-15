package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"

	pb "github.com/gonzabosio/transaction/services/proto/email"
	"github.com/jordan-wright/email"
	"google.golang.org/grpc"
)

type EmailService struct {
	pb.UnimplementedEmailServiceServer
}

func StartEmailServiceServer() {
	lis, err := net.Listen("tcp", os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterEmailServiceServer(grpcServer, &EmailService{})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
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
