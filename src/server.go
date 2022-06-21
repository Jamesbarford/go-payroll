package src

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Jamesbarford/go-payroll/src/payroll"
)

/* This is the best I can think of to obtain dependency injection.
 * https://github.com/uber-go/dig -> looks to solve this
 */
func registerApplicationHandlers(mux *http.ServeMux, dbConnection *sql.DB) {
	payrollRepository := payroll.NewPayrollRepository(dbConnection)
	payrollService := payroll.NewPayrollService(payrollRepository)
	payrollHandler := payroll.NewPayrollHandler(payrollService)

	mux.HandleFunc("/api/payroll", payrollHandler.HandlePayrollRequest)
}

func ServerMain(port string) {
	dbConnection, dbConnectionError := NewDbConnection(DbConfigFromEnvironment())

	if dbConnectionError != nil {
		panic(dbConnectionError)
	}

	mux := http.NewServeMux()

	registerApplicationHandlers(mux, dbConnection)

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	fmt.Printf("Server listening on port :: " + port + "\n")
	server.ListenAndServe()
}
