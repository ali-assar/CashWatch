package expensehandler

import (
	"context"

	"github.com/Ali-Assar/CashWatch/db"
	pb "github.com/Ali-Assar/CashWatch/types"
	"github.com/go-playground/validator"
	"github.com/golang/protobuf/ptypes/empty"
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

func (s *ExpenseServiceServer) GetCategory(ctx context.Context, req *pb.CategoryRequest) (*pb.Category, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	fetchedCategory, err := s.ExpenseStore.GetCategoryById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.Category{
		ID:     fetchedCategory.ID,
		Name:   fetchedCategory.Name,
		UserId: fetchedCategory.UserId,
	}, nil
}

func (s *ExpenseServiceServer) DeleteCategoryByID(ctx context.Context, req pb.CategoryRequest) (*empty.Empty, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err := s.ExpenseStore.DeleteCategoryById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *ExpenseServiceServer) UpdateCategoryByEmail(ctx context.Context, reqId *pb.CategoryRequest, req *pb.Category) (*empty.Empty, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	err := s.ExpenseStore.UpdateCategoryById(ctx, reqId.GetId(), req)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// ------- Expense CRUD ------

func (s *ExpenseServiceServer) InsertExpense(ctx context.Context, req *pb.Expense) (*pb.Expense, error) {
	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	expense := &pb.Expense{
		ID:          req.ID,
		Description: req.Description,
		Amount:      req.Amount,
		CategoryId:  req.CategoryId,
		UserId:      req.UserId,
	}

	insertedExpense, err := s.ExpenseStore.InsertExpense(ctx, expense)
	if err != nil {
		return nil, err
	}

	return &pb.Expense{ID: insertedExpense.ID}, nil
}

func (s *ExpenseServiceServer) GetExpanse(ctx context.Context, req *pb.ExpenseRequest) (*pb.Expense, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	fetchedExpanse, err := s.ExpenseStore.GetExpenseById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.Expense{
		ID:          fetchedExpanse.ID,
		Description: fetchedExpanse.Description,
		Amount:      fetchedExpanse.Amount,
		CategoryId:  fetchedExpanse.CategoryId,
		UserId:      fetchedExpanse.UserId,
	}, nil
}

func (s *ExpenseServiceServer) DeleteExpanseByID(ctx context.Context, req pb.ExpenseRequest) (*empty.Empty, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation error: %v", err)
	}

	err := s.ExpenseStore.DeleteExpenseById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *ExpenseServiceServer) UpdateExpenseByEmail(ctx context.Context, reqId *pb.ExpenseRequest, req *pb.Expense) (*empty.Empty, error) {
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Validation error: %v", err)
	}

	err := s.ExpenseStore.UpdateExpenseById(ctx, reqId.GetId(), req)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
