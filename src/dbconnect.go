/**
 * Create database connection, Idea being sometimes we might want this from environment variables
 * other times it could be a connection string retrieved from elsewhere etc...
 */
package src

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

type DbConfig struct {
	host       string
	port       int
	dbname     string
	dbuser     string
	dbpassword string
}

/* Create database configuration from environment variables or panic */
func DbConfigFromEnvironment() *DbConfig {
	dbUser, userOk := os.LookupEnv("DB_USER")
	if !userOk {
		panic("Environment variable DB_USER must be set to connect to database")
	}

	dbPassword, passwordOk := os.LookupEnv("DB_PASSWORD")
	if !passwordOk {
		panic("Environment variable DB_PASSWORD must be set to connect to database")
	}

	dbPort, portOk := os.LookupEnv("DB_PORT")
	if !portOk {
		panic("Environment variable DB_PORT must be set to connect to database")
	}

	dbHost, hostOk := os.LookupEnv("DB_HOST")
	if !hostOk {
		panic("Environment variable DB_HOST must be set to connect to database")
	}

	dbPortInt, portIntOk := strconv.Atoi(dbPort)
	if portIntOk != nil {
		fmt.Printf("Port: %s\n", dbPort)
		panic("Failed to convert port number to int: " + dbPort)
	}

	dbName, nameOk := os.LookupEnv("DB_NAME")
	if !nameOk {
		panic("Environment variable DB_NAME must be set to connect to database")
	}

	return &DbConfig{
		host:       dbHost,
		port:       dbPortInt,
		dbname:     dbName,
		dbpassword: dbPassword,
		dbuser:     dbUser,
	}
}

/* Create connection to database from configuration */
func NewDbConnection(config *DbConfig) (*sql.DB, error) {
	pgConnectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.dbuser,
		config.dbpassword,
		config.host,
		config.port,
		config.dbname)

	dbConnection, dbConnectionError := sql.Open("postgres", pgConnectionString)

	if dbConnectionError != nil {
		return nil, dbConnectionError
	}

	/* Ping postgres to make sure it is up and running as it should be */
	if dbPingError := dbConnection.Ping(); dbPingError != nil {
		return nil, dbPingError
	}

	fmt.Printf("Connection to database on port :: %d\n", config.port)

	return dbConnection, nil
}
