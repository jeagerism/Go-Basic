package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "myuser"
	password = "mypassword"
	dbname   = "gorm"
)

func AuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	var jwtSecret = []byte("your_jwt_secret")
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})
	if err != nil {
		panic("failed to connect to database")
	}
	//Migrate the schema
	db.AutoMigrate(&Book{}, &User{})
	fmt.Println("Database migration completed!")

	app := fiber.New()
	app.Use("/book", AuthRequired)

	//CRUD routes
	app.Get("/books", func(c *fiber.Ctx) error {
		return GetBooks(db, c)
	})
	app.Get("/book/:id", func(c *fiber.Ctx) error {
		return GetBook(db, c)
	})
	app.Post("/book", func(c *fiber.Ctx) error {
		return CreateBook(db, c)
	})
	app.Put("/book/:id", func(c *fiber.Ctx) error {
		return UpdateBook(db, c)
	})
	app.Delete("/book/:id", func(c *fiber.Ctx) error {
		return DeleteBook(db, c)
	})

	//AUTH Routes
	app.Post("/register", func(c *fiber.Ctx) error {
		return CreateUser(db, c)
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return LoginUser(db, c)
	})

	//Start server
	log.Fatal(app.Listen(":8000"))
}
