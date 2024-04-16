package budgethandler

import (
	"context"

	"github.com/Ali-Assar/CashWatch/db"
	pb "github.com/Ali-Assar/CashWatch/types"
	"github.com/go-playground/validator"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BudgetServiceServer struct {
	pb.UnimplementedBudgetServiceServer
	BudgetStore db.BudgetStorer
}

func (s *BudgetServiceServer) InsertBudget(ctx context.Context, req *pb.Budget) (*pb.Budget, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	budget := &pb.Budget{
		Title:    req.Title,
		Amount:   req.Amount,
		ExpireAt: req.ExpireAt,
		SetAt:    req.SetAt,
		UserId:   req.UserId,
	}

	insertedBudget, err := s.BudgetStore.InsertBudget(ctx, budget)
	if err != nil {
		return nil, err
	}

	return &pb.Budget{ID: insertedBudget.ID}, nil
}

func (s *BudgetServiceServer) GetBudget(ctx context.Context, req *pb.BudgetRequest) (*pb.Budget, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	fetchedBudget, err := s.BudgetStore.GetBudgetByID(ctx, req.GetID())
	if err != nil {
		return nil, err
	}

	return &pb.Budget{
		Title:    fetchedBudget.Title,
		Amount:   fetchedBudget.Amount,
		ExpireAt: fetchedBudget.ExpireAt,
		SetAt:    fetchedBudget.SetAt,
		UserId:   fetchedBudget.UserId,
	}, nil
}

func (s *BudgetServiceServer) DeleteBudgetByID(ctx context.Context, req pb.BudgetRequest) (*empty.Empty, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err := s.BudgetStore.DeleteBudgetByID(ctx, req.GetID())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *BudgetServiceServer) UpdateBudgetByEmail(ctx context.Context, reqId *pb.BudgetRequest, req *pb.Budget) (*empty.Empty, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	err := s.BudgetStore.UpdateBudgetByID(ctx, reqId.GetID(), req)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// ------- Income CRUD ------

func (s *BudgetServiceServer) InsertIncome(ctx context.Context, req *pb.Income) (*pb.Income, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	income := &pb.Income{
		Title:      req.Title,
		Amount:     req.Amount,
		ReceivedAt: req.ReceivedAt,
		UserId:     req.UserId,
	}

	insertedIncome, err := s.BudgetStore.InsertIncome(ctx, income)
	if err != nil {
		return nil, err
	}

	return &pb.Income{ID: insertedIncome.ID}, nil
}

func (s *BudgetServiceServer) GetIncome(ctx context.Context, req *pb.IncomeRequest) (*pb.Income, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	fetchedIncome, err := s.BudgetStore.GetIncomeByID(ctx, req.GetID())
	if err != nil {
		return nil, err
	}

	return &pb.Income{
		Title:      fetchedIncome.Title,
		Amount:     fetchedIncome.Amount,
		ReceivedAt: fetchedIncome.ReceivedAt,
		UserId:     fetchedIncome.UserId,
	}, nil
}

func (s *BudgetServiceServer) DeleteIncomeByID(ctx context.Context, req pb.IncomeRequest) (*empty.Empty, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err := s.BudgetStore.DeleteIncomeByID(ctx, req.GetID())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *BudgetServiceServer) UpdateIncomeByEmail(ctx context.Context, reqId *pb.IncomeRequest, req *pb.Income) (*empty.Empty, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	err := s.BudgetStore.UpdateIncomeByID(ctx, reqId.GetID(), req)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
