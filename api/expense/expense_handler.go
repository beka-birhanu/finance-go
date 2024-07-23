package expense

import (
	"fmt"
	"log"
	"net/http"

	"github.com/beka-birhanu/finance-go/api/expense/dto"
	"github.com/beka-birhanu/finance-go/api/util"
	handlerInterface "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	expenseCommand "github.com/beka-birhanu/finance-go/application/expense/command"
	"github.com/beka-birhanu/finance-go/domain/model"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type ExpensesHandler struct {
	addExpenseCommandHandler handlerInterface.ICommandHandler[*expenseCommand.AddExpenseCommand, *model.Expense]
}

func NewHandler(
	addExpenseCommandHandler handlerInterface.ICommandHandler[*expenseCommand.AddExpenseCommand, *model.Expense],
) *ExpensesHandler {
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

	router.HandleFunc(
		"/users/{userId}/expenses/{expenseId}",
		h.handleGetSingleExpenseById,
	).Methods(http.MethodGet)
}

func (h *ExpensesHandler) handleAddExpense(w http.ResponseWriter, r *http.Request) {
	var addExpenseRequest dto.AddExpenseRequest

	if err := util.ParseJSON(r, &addExpenseRequest); err != nil {
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	if err := util.Validate.Struct(addExpenseRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		util.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	userId, err := util.GetIdFromUrl(r, "userId")
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
	}

	addExpenseCommand := &expenseCommand.AddExpenseCommand{
		UserId:      userId,
		Description: addExpenseRequest.Description,
		Amount:      addExpenseRequest.Amount,
		Date:        addExpenseRequest.Date,
	}

	expense, err := h.addExpenseCommandHandler.Handle(addExpenseCommand)
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
		return
	}

	baseURL := util.GetBaseURL(r)

	// Construct the resource location URL dynamically
	resourceLocation := fmt.Sprintf("%s%s/%s", baseURL, r.URL.Path, expense.ID().String())

	w.Header().Set("Location", resourceLocation)
	util.WriteJSON(w, http.StatusCreated, nil)
}

func (h *ExpensesHandler) handleGetSingleExpenseById(w http.ResponseWriter, r *http.Request) {
	userId, err := util.GetIdFromUrl(r, "userId")
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
	}

	expenseId, err := util.GetIdFromUrl(r, "expenseId")
	if err != nil {
		util.WriteError(w, http.StatusBadRequest, err)
	}

	log.Println(userId, expenseId)
}
