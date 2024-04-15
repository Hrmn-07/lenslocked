package main

import (
	"fmt"
	"lenslocked/context"
	"lenslocked/models"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

func main() {

	ctx := context.Background()

	user := models.User{
		Email: "jon@gmail.com",
	}

	ctx = context.WithUser(ctx, &user)

	retrievedUser := context.User(ctx)
	fmt.Println(retrievedUser.Email)
	// cfg := models.DefaultPostgresConfig()
	// db, err := models.Open(cfg)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("connected!")

	// us := models.UserService{
	// 	DB: db,
	// }

	// user, err := us.Create("bob@bob,com", "bob123")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(user)

	// (for creating new table)
	// _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	name TEXT,
	// 	email TEXT NOT NULL
	// );

	// 	CREATE TABLE IF NOT EXISTS orders (
	// 	id SERIAL PRIMARY KEY,
	// 	user_id INT NOT NULL,
	// 	amount INT,
	// 	description TEXT
	// );`)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Tables created.")

	// (For adding a user)
	// name := "New User"
	// email := "jon@calhoun.io"
	// row := db.QueryRow(`
	// INSERT INTO users(name, email)
	// VALUES($1, $2) RETURNING id;`, name, email)
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("User created. id =", id)

	// (For getting a user's info)
	// id := 1
	// row := db.QueryRow(
	// 	`SELECT name, email
	// 	FROM users
	// 	WHERE id=$1;`, id,
	// )
	// var name, email string
	// err = row.Scan(&name, &email)
	// if err == sql.ErrNoRows {
	// 	fmt.Println("Error: No Rows!")
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("User Information: name=%s, email=%s\n", name, email)

	// (For creating orders)
	// userID := 1
	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Fake order #%d", i)
	// 	_, err := db.Exec(`
	// 	INSERT INTO orders(user_id, amount, description)
	// 	VALUES($1, $2, $3);`, userID, amount, desc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fmt.Println("Created fake orders.")

	// (for getting order amount and desc)
	// type Order struct {
	// 	ID          int
	// 	UserID      int
	// 	Amount      int
	// 	Description string
	// }
	// var orders []Order

	// userID := 1
	// rows, err := db.Query(`
	// SELECT id, amount, description
	// FROM orders
	// WHERE user_id=$1;`, userID)
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var order Order
	// 	order.UserID = userID
	// 	err := rows.Scan(&order.ID, &order.Amount, &order.Description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	orders = append(orders, order)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Orders:", orders)
}
