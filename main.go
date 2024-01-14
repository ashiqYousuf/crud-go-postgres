package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=testdb host=localhost port=5432")
	if err != nil {
		return nil, err
	}
	return db, nil
}

type User struct {
	UserID   int
	UserName string
	Email    string
}

func (u User) String() string {
	// Better use Env Variables
	return fmt.Sprintf("{UserId:%2v | UserName:%15v | Email:%15v}", u.UserID, u.UserName, u.Email)
}

func readUser(userID int) (User, error) {
	user := User{}
	db, err := connectDB()

	if err != nil {
		return user, err
	}
	defer db.Close()

	query := "SELECT id, username, email FROM users WHERE id=$1"
	err = db.QueryRow(query, userID).Scan(&user.UserID, &user.UserName, &user.Email)

	if err != nil {
		return user, err
	}
	return user, nil
}

func createUser(username, email string) error {
	db, err := connectDB()
	if err != nil {
		return err
	}

	defer db.Close()

	query := "INSERT INTO users (username, email) VALUES ($1, $2)"
	_, err = db.Exec(query, username, email)

	if err != nil {
		return err
	}
	return nil
}

func updateUser(user User) error {
	db, err := connectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := "UPDATE users SET username=$1, email=$2 WHERE  id=$3"
	_, err = db.Exec(query, user.UserName, user.Email, user.UserID)

	if err != nil {
		return err
	}
	return nil
}

func deleteUser(userID int) error {
	db, err := connectDB()

	if err != nil {
		return err
	}
	defer db.Close()

	query := "DELETE FROM users WHERE id=$1"
	_, err = db.Exec(query, userID)

	if err != nil {
		return err
	}
	return nil
}

func readAllUsers() ([]User, error) {
	users := []User{}
	db, err := connectDB()

	if err != nil {
		return users, err
	}
	defer db.Close()

	query := "SELECT * FROM users"
	rows, err := db.Query(query)

	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		user := User{}
		rows.Scan(&user.UserID, &user.UserName, &user.Email)
		users = append(users, user)
	}

	return users, nil
}

func main() {
	users, err := readAllUsers()

	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		fmt.Printf("%v\n", user)
	}
}
