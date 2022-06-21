package payroll

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type PayrollAddUserRequest struct {
	Forename    string  `json:"foreName"`
	Surname     string  `json:"surname"`
	GrossSalary float64 `json:"grossSalary"`
	CountryCode int     `json:"counrtyCode"`
}

type PayrollHandler struct {
	logger  *log.Logger
	service *PayrollService
}

func NewPayrollHandler(service *PayrollService) *PayrollHandler {
	return &PayrollHandler{
		logger:  log.Default(),
		service: service,
	}
}

/**
 * I wasn't able to find away to register handlerFuncs for the same url but different HTTP methods
 */
func (handler *PayrollHandler) HandlePayrollRequest(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	switch method {
	case http.MethodGet:
		handler.GetPayrollForUser(w, r)

	case http.MethodPost:
		handler.AddUserToPayRoll(w, r)

	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Supported methods: GET | POST",
		})
		return
	}
}

/**
 * Get payroll for a user given an ID and a month
 */
func (handler *PayrollHandler) GetPayrollForUser(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	userid := urlQuery.Get("userId")
	month := urlQuery.Get("month")

	if userid == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "User id must be provided",
		})
		return
	}

	parsedId, err := strconv.Atoi(userid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid userid",
		})
	}

	parsedMonth, parsedMonthErr := strconv.Atoi(month)
	if parsedMonthErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid month",
		})
	}

	result, serviceError := handler.service.GetPayrollForUser(parsedId, parsedMonth)
	if serviceError != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to get users payroll %s", serviceError.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(result))
}

/*
 * Parse request and add user + salary into database
 */
func (handler *PayrollHandler) AddUserToPayRoll(w http.ResponseWriter, r *http.Request) {
	// Grab  forename, surname, country, gross salary
	var addUserRequest PayrollAddUserRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&addUserRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid add user request",
		})
		return
	}

	serviceError := handler.service.AddUserToPayRoll(&addUserRequest)
	if serviceError != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Failed to add user to payroll %s", serviceError.Error()),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}
