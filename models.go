package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

// getBooks retrieves all books
func GetBooks(db *gorm.DB, c *fiber.Ctx) error {
	var books []Book
	db.Find(&books)
	return c.JSON(books)
}

// Search by id
func GetBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var book Book
	db.First(&book, id)
	return c.JSON(book)
}

func CreateBook(db *gorm.DB, c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return err
	}
	db.Create(&book)
	return c.JSON(book)
}

func UpdateBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	book := new(Book)
	db.First(&book, id)

	if err := c.BodyParser(book); err != nil {
		return err
	}
	db.Save(&book)
	return c.JSON(book)
}

func DeleteBook(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	db.Delete(&Book{}, id)
	return c.SendString("Book successfully deleted")
}
