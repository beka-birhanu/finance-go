package expenses

import (
	"fmt"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/expenses/dto"
	"github.com/beka-birhanu/finance-go/api/utils"
	"github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/expense"
	"github.com/beka-birhanu/finance-go/application/expense/commands"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ExpensesHandler struct {
	addExpenseCommandHandler expense.IAddExpenseCommand
}

func NewHandler(addExpenseCommandHandler expense.IAddExpenseCommand) *ExpensesHandler {
	return &ExpensesHandler{
		addExpenseCommandHandler: addExpenseCommandHandler,
	}
}

func (h *ExpensesHandler) RegisterPublicRoutes(router *mux.Router) {}

func (h *ExpensesHandler) RegisterProtectedRoutes(router *mux.Router) {
	router.HandleFunc(
		"/expenses",
		h.handleAddExpense,
	).Methods(http.MethodPost)
}

func (h *ExpensesHandler) handleAddExpense(w http.ResponseWriter, r *http.Request) {
	var addExpenseRequest dto.AddExpenseRequest

	if err := utils.ParseJSON(r, &addExpenseRequest); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	if err := utils.Validate.Struct(addExpenseRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	addExpenseCommand := &commands.AddExpenseCommand{
		UserId:      addExpenseRequest.UserId,
		Description: addExpenseRequest.Description,
		Amount:      addExpenseRequest.Amount,
		Date:        addExpenseRequest.Date,
	}

	expense, err := h.addExpenseCommandHandler.Handle(addExpenseCommand)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	baseURL := utils.GetBaseURL(r)

	// Construct the resource location URL dynamically
	resourceLocation := fmt.Sprintf("%s%s/%s", baseURL, r.URL.Path, expense.ID().String())

	w.Header().Set("Location", resourceLocation)
	utils.WriteJSON(w, http.StatusCreated, nil)
}
