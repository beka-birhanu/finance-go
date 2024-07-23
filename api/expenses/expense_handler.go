package expenses

import (
	"fmt"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/expenses/dto"
	"github.com/beka-birhanu/finance-go/api/utils"
	"github.com/beka-birhanu/finance-go/application/common/cqrs/i_commands/expense"
	"github.com/beka-birhanu/finance-go/application/expense/commands"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
		"/users/{userId}/expenses",
		h.handleAddExpense,
	).Methods(http.MethodPost)
}

func (h *ExpensesHandler) handleAddExpense(w http.ResponseWriter, r *http.Request) {
	var addExpenseRequest dto.AddExpenseRequest

	// Extract user ID from URL path
	vars := mux.Vars(r)
	userIdStr, ok := vars["userId"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user ID in path"))
		return
	}

	// Parse user ID to UUID
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user ID format: %v", err))
		return
	}

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
		UserId:      userId,
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
