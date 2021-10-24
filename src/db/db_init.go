package db

import (
	"database/sql"
	"fmt"

	"tira.com/src/helper"
)

const (
	PG_DB_USER     = "YOUR USER POSTGRE"
	PG_DB_PASSWORD = "YOUR PASSWORD POSTGRE"
	PG_DB_NAME     = "YOUR DB POSTGRE"
)

const (
	MSSQL_DB_USER     = "YOUR USER MSSQL"
	MSSQL_DB_PASSWORD = "YOUR PASSWORD MSSQL"
	MSSQL_DB_HOST     = "YOUR HOST MSSQL"
	MSSQL_DB_PORT     = "YOUR PORT MSSQL"
	MSSQL_DB_NAME     = "YOUR DB MSSQL"
)

// Postgres DB set up (Exported upper-case)
func SetupDBPostgres() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", PG_DB_USER, PG_DB_PASSWORD, PG_DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	helper.CheckErr(err)

	return db
}

// MSSQL Server DB set up (Exported upper-case)
func SetupDBMssql() *sql.DB {
	dbinfo := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&connection+timeout=30",
		MSSQL_DB_USER, MSSQL_DB_PASSWORD, MSSQL_DB_HOST, MSSQL_DB_PORT, MSSQL_DB_NAME)

	db, err := sql.Open("sqlserver", dbinfo)

	helper.CheckErr(err)

	return db
}
