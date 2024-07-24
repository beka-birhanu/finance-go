package expense

import (
	"fmt"
	"log"
	"net/http"

	apiError "github.com/beka-birhanu/finance-go/api/error"
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
		apiErr := apiError.NewErrBadRequest(fmt.Sprintf("invalid payload: %v", err))
		util.WriteError(w, apiErr)
		return
	}

	if err := util.Validate.Struct(addExpenseRequest); err != nil {
		errors := err.(validator.ValidationErrors)
		apiErr := apiError.NewErrValidation(fmt.Sprintf("invalid payload: %v", errors))
		util.WriteError(w, apiErr)
		return
	}

	userId, err := util.GetIdFromUrl(r, "userId")
	if err != nil {
		util.WriteError(w, err.(apiError.Error))
		return
	}

	addExpenseCommand := &expenseCommand.AddExpenseCommand{
		UserId:      userId,
		Description: addExpenseRequest.Description,
		Amount:      addExpenseRequest.Amount,
		Date:        addExpenseRequest.Date,
	}

	expense, err := h.addExpenseCommandHandler.Handle(addExpenseCommand)
	if err != nil {
		apiErr := apiError.NewErrBadRequest(err.Error())
		util.WriteError(w, apiErr)
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
		util.WriteError(w, err.(apiError.Error))
		return
	}

	expenseId, err := util.GetIdFromUrl(r, "expenseId")
	if err != nil {
		util.WriteError(w, err.(apiError.Error))
		return
	}

	log.Println(userId, expenseId)
}

