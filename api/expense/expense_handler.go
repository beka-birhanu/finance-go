package expense

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	baseapi "github.com/beka-birhanu/finance-go/api/base_handler"
	errapi "github.com/beka-birhanu/finance-go/api/error"
	"github.com/beka-birhanu/finance-go/api/expense/dto"
	icmd "github.com/beka-birhanu/finance-go/application/common/cqrs/command"
	iquery "github.com/beka-birhanu/finance-go/application/common/cqrs/query"
	expensecmd "github.com/beka-birhanu/finance-go/application/expense/command"
	expensqry "github.com/beka-birhanu/finance-go/application/expense/query"
	errdmn "github.com/beka-birhanu/finance-go/domain/error/common"
	expensemodel "github.com/beka-birhanu/finance-go/domain/model/expense"
	"github.com/google/uuid"
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

	// Extract userId for context and match with the userId form URL.
	err = h.MatchPathUserIdctxUserId(r, userId)
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

	// Extract userId for context and match with the userId form URL.
	err = h.MatchPathUserIdctxUserId(r, userId)
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

// handleByUserId handles the request to get expenses by user ID
func (h *ExpensesHandler) handleByUserId(w http.ResponseWriter, r *http.Request) {
	userId, err := h.UUIDParam(r, "userId")
	if err != nil {
		h.Problem(w, errapi.NewBadRequest(err.Error()))
		return
	}

	err = h.MatchPathUserIdctxUserId(r, userId)
	if err != nil {
		h.Problem(w, errapi.NewBadRequest(err.Error()))
		return
	}

	cursor, limit, sortField, sortOrder, err := h.extractAndValidateParams(r)
	if err != nil {
		h.Problem(w, errapi.NewBadRequest(err.Error()))
		return
	}

	queryParams, err := h.constructQueryParams(userId, cursor, limit, sortField, sortOrder)
	if err != nil {
		h.Problem(w, errapi.NewBadRequest(err.Error()))
		return
	}

	expenses, err := h.getMultipleHandler.Handle(queryParams)
	if err != nil {
		h.Problem(w, errapi.NewServerError(err.Error()))
		return
	}

	expensesResponse := make([]*dto.GetExpenseResponse, 0)
	for _, expense := range expenses {
		if expense != nil {
			expensesResponse = append(expensesResponse, dto.FromExpenseModel(expense))
		}
	}

	nextCursor := ""
	if len(expenses) > 0 {
		nextCursor = h.buildCursor(expenses[len(expenses)-1], sortField)
	}

	response := dto.GetMultipleResponse{
		Expenses: expensesResponse,
		Cursor:   nextCursor,
	}

	h.Respond(w, http.StatusOK, response)
}

// extractAndValidateParams extracts and validates the query parameters from the request
func (h *ExpensesHandler) extractAndValidateParams(r *http.Request) (string, int, string, string, error) {
	cursor := h.StringQueryParam(r, "cursor")

	limit, err := h.IntQueryParam(r, "limit")
	if err != nil {
		return "", 0, "", "", err
	}

	sortBy := h.StringQueryParam(r, "sortBy")
	sortField := "createdAt"
	sortOrder := "desc"
	if sortBy != "" {
		parts := strings.Split(sortBy, ".")
		if len(parts) != 2 {
			return "", 0, "", "", errapi.NewBadRequest("invalid sortBy format")
		}
		sortField = parts[0]
		sortOrder = parts[1]
		if sortField != "createdAt" && sortField != "amount" {
			return "", 0, "", "", errapi.NewBadRequest(fmt.Sprintf("invalid sortBy field: %s", sortField))
		}
		if sortOrder != "asc" && sortOrder != "desc" {
			return "", 0, "", "", errapi.NewBadRequest(fmt.Sprintf("invalid sortBy order: %s", sortOrder))
		}
	}

	return cursor, limit, sortField, sortOrder, nil
}

// constructQueryParams constructs the appropriate query parameters struct
func (h *ExpensesHandler) constructQueryParams(userId uuid.UUID, cursor string, limit int, sortField string, sortOrder string) (*expensqry.GetMultipleQuery, error) {
	var lastSeenID uuid.UUID
	var lastSeenTime time.Time
	var lastSeenAmt float64
	var ascending bool

	if cursor != "" {
		cursorByte, err := base64.StdEncoding.DecodeString(cursor)
		if err != nil {
			return &expensqry.GetMultipleQuery{}, errapi.NewBadRequest("invalid cursor format1")
		}

		cursor = string(cursorByte)
		cursorParts := strings.Split(cursor, ",")
		if len(cursorParts) != 2 {
			return &expensqry.GetMultipleQuery{}, errapi.NewBadRequest("invalid cursor format1")
		}
		lastSeenID, err = uuid.Parse(cursorParts[0])
		if err != nil {
			return &expensqry.GetMultipleQuery{}, errapi.NewBadRequest("invalid cursor format2")
		}

		if sortField == "createdAt" {
			lastSeenTime, err = time.Parse(time.RFC3339, cursorParts[1])
			if err != nil {
				return &expensqry.GetMultipleQuery{}, fmt.Errorf("invalid cursor format for createdAt: %v", err)
			}
		} else if sortField == "amount" {
			lastSeenAmt, err = strconv.ParseFloat(cursorParts[1], 64)
			if err != nil {
				return &expensqry.GetMultipleQuery{}, fmt.Errorf("invalid cursor format for amount")
			}
		}
	}

	if sortOrder == "asc" {
		ascending = true
	}

	return &expensqry.GetMultipleQuery{
		UserID:       userId,
		Limit:        limit,
		By:           sortField,
		LastSeenID:   &lastSeenID,
		LastSeenTime: &lastSeenTime,
		LastSeenAmt:  lastSeenAmt,
		Ascending:    ascending,
	}, nil
}

// respondWithExpenses constructs and sends the response with the expenses and the next cursor
func (h *ExpensesHandler) buildCursor(lastExpense *expensemodel.Expense, field string) string {
	nextCursor := ""
	if lastExpense != nil {
		if field == "amount" {
			nextCursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%f", lastExpense.ID(), lastExpense.Amount())))
		} else {
			createdAt := lastExpense.CreatedAt().Format(time.RFC3339)
			nextCursor = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s,%v", lastExpense.ID(), createdAt)))
		}
	}

	return nextCursor
}
