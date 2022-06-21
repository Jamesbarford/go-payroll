/**
 * This is fast though hack and slash. An ORM would give better safety for objects.
 * I am unfamiliar with any in Go, this looks promising: https://gorm.io/index.html
 */
package payroll

import (
	"database/sql"
	"errors"
	"log"
)

type PayrollRepository struct {
	logger *log.Logger
	db     *sql.DB
}

func NewPayrollRepository(db *sql.DB) *PayrollRepository {
	return &PayrollRepository{
		logger: log.Default(),
		db:     db,
	}
}

/* Prepare SQL statement to prevent sql injection */
func prepareStmt(
	sqlQuery string,
	db *sql.DB,
	args ...interface{},
) (*sql.Rows, error) {
	stmt, prepareErr := db.Prepare(sqlQuery)
	if prepareErr != nil {
		return nil, prepareErr
	}

	rows, queryErr := stmt.Query(args...)
	if queryErr != nil {
		return nil, queryErr
	}

	return rows, nil
}

/**
 * As we are using postgres to generate the JSON we only need the first row
 * and do not need to parse datatypes. This should be fairly safe
 */
func execGetJsonString(
	sqlQuery string,
	db *sql.DB,
	args ...interface{},
) (string, error) {
	rows, rowsError := prepareStmt(sqlQuery, db, args...)
	if rowsError != nil {
		return "", rowsError
	}

	rowResult := ""
	defer rows.Close()
	rows.Next()
	rows.Scan(&rowResult)

	if rowResult == "" {
		rowResult = "[]"
	}

	return rowResult, nil
}

/* Get a payroll for a user if they exist or error */
func (repository *PayrollRepository) GetPayrollForUser(userId int, month int) (string, error) {
	rows, rowsError := prepareStmt(`SELECT id FROM users WHERE id = $1;`, repository.db, userId)

	if rowsError != nil {
		repository.logger.Printf("Failed to get payroll for user: %s\n", rowsError)
		return "", errors.New("No user found")
	}

	/* Ensure we have got one user back */
	userCount := 0
	for rows.Next() {
		userCount++
		break
	}

	/* Throw an error */
	if userCount == 0 {
		return "", errors.New("No user found")
	}

	/*
	 * There is a fair amount of magic.
	 * - To calculate the percentage product you can do this using the Exponent of the sum of the logarthim
	 * - We use this on the bonuses to find the total gross payable salary
	 * - Deduct tax from total gross payable salary
	 * - Use the product of the deductions
	 * - calculate the final amount and prepare as JSON
	 *
	 * Joining on NULL is a bit hacky but works for when we only need a single value and have
	 * nothing to join on.
	 */
	return execGetJsonString(
		`
		WITH bonus AS (
			SELECT
				COALESCE(EXP(SUM(LN((country_bonuses.percentage / 100)))), 0) AS percentage
			FROM
				country_bonuses
			WHERE country_bonuses.month_payable = $1
		), tax_deducted AS (
			SELECT
				payroll.gross_salary + (payroll.gross_salary * bonus.percentage) AS gross_salary,
				bonus.percentage AS bonus,
				users.country_code,
				users.id AS user_id,
				users.forename,
				users.surname,
				(payroll.gross_salary + (payroll.gross_salary * bonus.percentage)) * (1 - (taxes.percent / 100)) AS salary_minus_tax
			FROM
				payroll
			LEFT JOIN users ON users.id = payroll.user_id
			LEFT JOIN taxes ON taxes.country_code = users.country_code
			LEFT JOIN bonus ON bonus.percentage IS NOT NULL
			WHERE users.id = $2
		), cumulative_percentage_deductions AS (
			SELECT
				COALESCE(EXP(SUM(LN(1-(deductions.percentage/100)))), 0) AS product
			FROM
				deductions
			LEFT JOIN tax_deducted ON tax_deducted.country_code = deductions.country_code
		)
		SELECT
			json_build_object(
				'forename', tax_deducted.forename,
				'surname', tax_deducted.surname,
				'grossSalary', tax_deducted.gross_salary,
				'bonus', tax_deducted.bonus,
				'netSalary', tax_deducted.salary_minus_tax * cumulative_percentage_deductions.product
			) AS result
		FROM
			tax_deducted
		LEFT JOIN 
			cumulative_percentage_deductions ON cumulative_percentage_deductions.product IS NOT NULL;
		`,
		repository.db,
		month,
		userId,
	)
}

/* Add a user to the payroll table and user table */
func (repository *PayrollRepository) AddUserToPayRoll(addUser *PayrollAddUserRequest) error {
	insertRow, rowsError := prepareStmt(
		`INSERT INTO users(forename, surname, country_code) VALUES ($1, $2, $3) RETURNING id;`,
		repository.db,
		addUser.Forename, addUser.Surname, addUser.CountryCode,
	)

	if rowsError != nil {
		repository.logger.Printf("Failed to prepare statement for insertion of user: %s\n", rowsError.Error())
		return errors.New("Invalid user paramaters")
	}

	var userId int
	for insertRow.Next() {
		scanErr := insertRow.Scan(&userId)

		if scanErr != nil {
			repository.logger.Printf("Failed to return userid: %s\n", scanErr)
			return errors.New("No id created for user")
		}

		/* We should have our id now */
		break
	}

	insertRow.Close()

	payrollRow, payrollError := prepareStmt(
		`INSERT INTO payroll (gross_salary, user_id) VALUES ($1, $2);`,
		repository.db,
		addUser.GrossSalary, userId,
	)

	defer payrollRow.Close()

	if payrollError != nil {
		repository.logger.Printf("Failed to perpare statement for payroll: %s\n", payrollError.Error())
		return errors.New("Invalid payroll parameters")
	}

	return nil
}
