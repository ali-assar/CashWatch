package expensehandler

import (
	"context"

	"github.com/Ali-Assar/CashWatch/db"
	pb "github.com/Ali-Assar/CashWatch/types"
	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExpenseServiceServer struct {
	pb.UnimplementedExpenseServiceServer
	ExpenseStore db.ExpenseStorer
}

func (s *ExpenseServiceServer) InsertCategory(ctx context.Context, req *pb.Category) (*pb.Category, error) {
	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	category := &pb.Category{
		Name: req.Name,
	}

	insertedCategory, err := s.ExpenseStore.InsertCategory(ctx, category)
	if err != nil {
		return nil, err
	}

	return &pb.Category{ID: insertedCategory.ID}, nil
}
