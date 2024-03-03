package expensehandler

import (
	"github.com/Ali-Assar/CashWatch/db"
	pb "github.com/Ali-Assar/CashWatch/types"
)

type ExpenseServiceServer struct {
	pb.UnimplementedExpenseServiceServer
	ExpenseStore db.ExpenseStorer
}
