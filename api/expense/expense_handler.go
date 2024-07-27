package expense

import (
	"fmt"
	"log"
	"net/http"

	baseapi "github.com/beka-birhanu/finance-go/api/base_handler"
	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/expense/dto"
	httputil "github.com/beka-birhanu/finance-go/api/http_util"
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	itimeservice "github.com/beka-birhanu/finance-go/application/common/interface/time_service"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/gorilla/mux"
)

type ExpensesHandler struct {
	baseapi.BaseHandler
	addHandler icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
}

func NewHandler(
	addHandler icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense],
	timeServie itimeservice.IService,
) *ExpensesHandler {
	return &ExpensesHandler{
		addHandler: addHandler,
	}
}

func (h *ExpensesHandler) RegisterPublicRoutes(router *mux.Router) {}

func (h *ExpensesHandler) RegisterProtectedRoutes(router *mux.Router) {
	router.HandleFunc(
		"/users/{userId}/expenses",
		h.handleAdd,
	).Methods(http.MethodPost)

	router.HandleFunc(
		"/users/{userId}/expenses/{expenseId}",
		h.handleById,
	).Methods(http.MethodGet)
}

func (h *ExpensesHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var addExpenseRequest dto.AddExpenseRequest

	// Populate addExpenseRequest from request body
	if err := h.ValidatedBody(r, &addExpenseRequest); err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	userId, err := httputil.UUIDParam(r, "userId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	addExpenseCommand := &expensecmd.AddCommand{
		UserId:      userId,
		Description: addExpenseRequest.Description,
		Amount:      addExpenseRequest.Amount,
		Date:        addExpenseRequest.Date,
	}

	expense, err := h.addHandler.Handle(addExpenseCommand)
	if err != nil {
		apiErr := errapi.NewBadRequest(err.Error())
		h.Problem(w, apiErr)
		return
	}

	baseURL := httputil.BaseURL(r)

	// Construct the resource location URL
	resourceLocation := fmt.Sprintf("%s%s/%s", baseURL, r.URL.Path, expense.ID().String())
	httputil.ResondWithLocation(w, http.StatusCreated, nil, resourceLocation)
}

func (h *ExpensesHandler) handleById(w http.ResponseWriter, r *http.Request) {
	userId, err := httputil.UUIDParam(r, "userId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	expenseId, err := httputil.UUIDParam(r, "expenseId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	log.Println(userId, expenseId)
}
