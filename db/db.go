package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	// Establish database connection
	db, err := sql.Open("postgres", "postgres://admin:admin@localhost:5432/database?sslmode=disable")
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(10)
	CreateUserTable(db)

	return db, nil
}

func CreateUserTable(db *sql.DB) {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		firstName VARCHAR(200) NOT NULL,
		lastName VARCHAR(200) NOT NULL,
		email VARCHAR(300) NOT NULL UNIQUE,
		encryptedPassword TEXT NOT NULL
	)
	`
	_, err := db.Exec(createUsersTable)
	if err != nil {
		log.Fatal("Could not create users table " + err.Error())
	}
}

func CreateCategoriesTable(db *sql.DB) {
	createCategoriesTable := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(200) NOT NULL,
		user_id INT,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`
	_, err := db.Exec(createCategoriesTable)
	if err != nil {
		log.Fatal("Could not create categories table " + err.Error())
	}
}

func CreateExpansesTable(db *sql.DB) {
	createExpansesTable := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		user_id INT,
		description TEXT,
		amount NUMERIC(10,2),
		FOREIGN KEY (category_id) REFERENCES Category(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`
	_, err := db.Exec(createExpansesTable)
	if err != nil {
		log.Fatal("Could not create expenses table " + err.Error())
	}
}

func CreateBudgetTable(db *sql.DB) {
	createBudgetsTable := `
	CREATE TABLE IF NOT EXISTS budgets (
		id serial PRIMARY KEY,
		user_id INT,
		title VARCHAR(200) NOT NULL,
		amount NUMERIC(10,2),
		expireAT TIMESTAMP,
		setAt TIMESTAMP,
		expireAt TIMESTAMP,
		
	)
	`
	_, err := db.Exec(createBudgetsTable)
	if err != nil {
		log.Fatal("Error creating budgets table:", err)
	}
}

func CreateIncomeTable(db *sql.DB) {
	createIncomeTable := `
	CREATE TABLE IF NOT EXISTS income (
		id serial PRIMARY KEY,
		user_id INT
		title VARCHAR(200) NOT NULL,
		amount NUMERIC(10,2),
		receivedAt TIMESTAMP,
		
	)
	`
	_, err := db.Exec(createIncomeTable)
	if err != nil {
		log.Fatal("Error creating income table:", err)
	}
}
