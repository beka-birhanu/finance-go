package expense

import (
	"fmt"
	"net/http"

	baseapi "github.com/beka-birhanu/finance-go/api/base_handler"
	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/expense/dto"
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
	expensqry "github.com/beka-birhanu/finance-go/application/expense/query"
	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/gorilla/mux"
)

type ExpensesHandler struct {
	baseapi.BaseHandler
	addHandler         icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
	getHandler         iquery.IHandler[*expensqry.GetQuery, *expensemodel.Expense]
	getMultipleHandler iquery.IHandler[*expensqry.GetMultipleQuery, []*expensemodel.Expense]
	patchHandler       iquery.IHandler[*expensecmd.PatchCommand, *expensemodel.Expense]
}

type Config struct {
	AddHandler         icmd.IHandler[*expensecmd.AddCommand, *expensemodel.Expense]
	GetHandler         iquery.IHandler[*expensqry.GetQuery, *expensemodel.Expense]
	PatchHandler       iquery.IHandler[*expensecmd.PatchCommand, *expensemodel.Expense]
	GetMultipleHandler iquery.IHandler[*expensqry.GetMultipleQuery, []*expensemodel.Expense]
}

func NewHandler(config Config) *ExpensesHandler {
	return &ExpensesHandler{
		addHandler:         config.AddHandler,
		getHandler:         config.GetHandler,
		patchHandler:       config.PatchHandler,
		getMultipleHandler: config.GetMultipleHandler,
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

	router.HandleFunc(
		"/users/{userId}/expenses",
		h.handleByUserId,
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/users/{userId}/expenses/{expenseId}",
		h.handlePatch,
	).Methods(http.MethodPatch)
}

func (h *ExpensesHandler) handleAdd(w http.ResponseWriter, r *http.Request) {
	var addExpenseRequest dto.AddExpenseRequest

	// Populate addExpenseRequest from request body
	if err := h.ValidatedBody(r, &addExpenseRequest); err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	userId, err := h.UUIDParam(r, "userId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	// TODO: match the id in the ctx
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

	baseURL := h.BaseURL(r)

	// Construct the resource location URL
	resourceLocation := fmt.Sprintf("%s%s/%s", baseURL, r.URL.Path, expense.ID().String())
	response := dto.FromExpenseModel(expense)
	h.RespondWithLocation(w, http.StatusCreated, response, resourceLocation)
}

func (h *ExpensesHandler) handleById(w http.ResponseWriter, r *http.Request) {
	userId, err := h.UUIDParam(r, "userId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	expenseId, err := h.UUIDParam(r, "expenseId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	expense, err := h.getHandler.Handle(&expensqry.GetQuery{UserId: userId, ExpenseId: expenseId})
	if err != nil {
		h.Problem(w, errapi.NewBadRequest(err.Error()))
		return
	}
	response := dto.FromExpenseModel(expense)
	h.Respond(w, http.StatusOK, response)
}

func (h *ExpensesHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	var patchRequest dto.PatchRequest
	userId, err := h.UUIDParam(r, "userId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	expenseId, err := h.UUIDParam(r, "expenseId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	// Extract userId for context and match with the userId form URL.
	err = h.MatchPathUserIdctxUserId(r, userId)
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	err = h.BaseHandler.ValidatedBody(r, &patchRequest)
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	expense, err := h.patchHandler.Handle(&expensecmd.PatchCommand{
		Description: patchRequest.Description,
		Amount:      patchRequest.Amount,
		Date:        patchRequest.Date,
		Id:          expenseId,
		UserId:      userId,
	})

	if err != nil {
		switch err.(*errdmn.Error).Type() {
		case errdmn.NotFound:
			h.Problem(w, errapi.NewNotFound(err.Error()))
		case errdmn.Validation:
			h.Problem(w, errapi.NewBadRequest(err.Error()))
		default:
			h.Problem(w, errapi.NewServerError("unknown error occured while patching expense"))
		}
		return
	}
	response := dto.FromExpenseModel(expense)
	h.Respond(w, http.StatusOK, response)
}

func (h *ExpensesHandler) handleByUserId(w http.ResponseWriter, r *http.Request) {
	userId, err := h.UUIDParam(r, "userId")
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	// Extract userId for context and match with the userId form URL.
	err = h.MatchPathUserIdctxUserId(r, userId)
	if err != nil {
		h.Problem(w, err.(errapi.Error))
		return
	}

	expenses, err := h.getMultipleHandler.Handle(&expensqry.GetMultipleQuery{UserId: userId})
	expenses, err := h.getMultipleHandler.Handle(&expensqry.GetMultipleQuery{UserID: userId})
	if err != nil {
		h.Problem(w, errapi.NewServerError(err.Error()))
	}

	response := make([]*dto.GetExpenseResponse, 0)
	for _, expense := range expenses {
		if expense != nil {
			response = append(response, dto.FromExpenseModel(expense))
		}
	}

	h.Respond(w, http.StatusOK, response)
}
