package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	_ "week11-assignment/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	// "golang.org/x/text/search"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	ISBN   string  `json:"isbn"`
	Year   int     `json:"year"`
	Price  float64 `json:"price"`

	// ฟิลด์ใหม่
	Category      string   `json:"category"`
	OriginalPrice *float64 `json:"original_price,omitempty"`
	Discount      int      `json:"discount"`
	CoverImage    string   `json:"cover_image"`
	Rating        float64  `json:"rating"`
	ReviewsCount  int      `json:"reviews_count"`
	IsNew         bool     `json:"is_new"`
	Pages         *int     `json:"pages,omitempty"`
	Language      string   `json:"language"`
	Publisher     string   `json:"publisher"`
	Description   string   `json:"description"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var db *sql.DB

func initDB() {
	var err error
	host := getEnv("DB_HOST", "")
	name := getEnv("DB_NAME", "")
	user := getEnv("DB_USER", "")
	password := getEnv("DB_PASSWORD", "")
	port := getEnv("DB_PORT", "")

	conSt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, name)
	// fmt.Println(conSt)
	db, err = sql.Open("postgres", conSt)
	if err != nil {
		log.Fatal("failed to open database")
	}
	// กำหนดจำนวน Connection สูงสุด
	db.SetMaxOpenConns(25)

	// กำหนดจำนวน Idle connection สูงสุด
	db.SetMaxIdleConns(20)

	// กำหนดอายุของ Connection
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		log.Fatal("failed to connect database")
	}

	log.Println("succesfully connected to database")
}

// @Summary Get all book
// @Description Get details of a book
// @Tags Books
// @Produce  json
// @Success 200  {array} Book
// @Failure 500 {object}  ErrorResponse
// @Router /books [get]
func getAllBooks(c *gin.Context) {
	var rows *sql.Rows
	var err error
	// ลูกค้าถาม "มีหนังสืออะไรบ้าง"
	rows, err = db.Query("SELECT id, title, author, isbn, year, price, created_at, updated_at FROM books")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close() // ต้องปิด rows เสมอ เพื่อคืน Connection กลับ pool

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			// handle error
		}
		books = append(books, book)
	}
	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get book by ID
// @Description Get details of a book by ID
// @Tags Books
// @Produce  json
// @Param   id   path      int     true  "Book ID"
// @Success 200  {object}  Book
// @Failure 404  {object}  ErrorResponse
// @Router /books/{id} [get]
func getBook(c *gin.Context) {
	id := c.Param("id")
	var book Book

	// QueryRow ใช้เมื่อคาดว่าจะได้ผลลัพธ์ 0 หรือ 1 แถว
	err := db.QueryRow("SELECT id, title, author FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

// @Summary Create a new book
// @Description Add a new book to the database
// @Tags Books
// @Accept  json
// @Produce  json
// @Param   book  body      Book  true  "Book data"
// @Success 201  {object}  Book
// @Failure 400  {object}  ErrorResponse
// @Failure 500  {object}  ErrorResponse
// @Router /books [post]
func createBook(c *gin.Context) {
	var newBook Book

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ใช้ RETURNING เพื่อดึงค่าที่ database generate (id, timestamps)
	var id int
	var createdAt, updatedAt time.Time

	err := db.QueryRow(
		`INSERT INTO books (title, author, isbn, year, price)
         VALUES ($1, $2, $3, $4, $5)
         RETURNING id, created_at, updated_at`,
		newBook.Title, newBook.Author, newBook.ISBN, newBook.Year, newBook.Price,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newBook.ID = id
	newBook.CreatedAt = createdAt
	newBook.UpdatedAt = updatedAt

	c.JSON(http.StatusCreated, newBook) // ใช้ 201 Created
}

// @Summary Update a book by ID
// @Description Update details of an existing book by its ID
// @Tags Books
// @Accept  json
// @Produce  json
// @Param   id    path      int    true  "Book ID"
// @Param   book  body      Book   true  "Updated book data"
// @Success 200   {object}  Book
// @Failure 400   {object}  ErrorResponse
// @Failure 404   {object}  ErrorResponse
// @Failure 500   {object}  ErrorResponse
// @Router /books/{id} [put]
func updateBook(c *gin.Context) {
	var u_id int
	id := c.Param("id")
	var updateBook Book

	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedAt time.Time
	err := db.QueryRow(
		`UPDATE books
         SET title = $1, author = $2, isbn = $3, year = $4, price = $5
         WHERE id = $6
         RETURNING id,updated_at`,
		updateBook.Title, updateBook.Author, updateBook.ISBN,
		updateBook.Year, updateBook.Price, id,
	).Scan(&id, &updatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updateBook.ID = u_id
	updateBook.UpdatedAt = updatedAt
	c.JSON(http.StatusOK, updateBook)
}

// @Summary Delete a book by ID
// @Description Delete a specific book from the database by its ID
// @Tags Books
// @Produce  json
// @Param   id  path      int     true  "Book ID"
// @Success 200  {object}  map[string]string  "Success message"
// @Failure 404  {object}  ErrorResponse
// @Failure 500  {object}  ErrorResponse
// @Router /books/{id} [delete]
func deleteBook(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}

func getNewBooks(c *gin.Context) {

	var rows *sql.Rows
	var err error

	// ลูกค้าถาม "มีหนังสืออะไรบ้าง"
	rows, err = db.Query("SELECT id, title, author, isbn, year, price, created_at, updated_at FROM books order by created_at desc limit 5")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close() // ต้องปิด rows เสมอ เพื่อคืน Connection กลับ pool

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			// handle error
		}
		books = append(books, book)
	}
	if books == nil {
		books = []Book{}
	}

	yearQuery := c.Query("year")
	if yearQuery != "" {
		filter := []Book{}
		for _, book := range books {
			if fmt.Sprint(book.Year) == yearQuery {
				filter = append(filter, book)
			}
		}
		c.JSON(http.StatusOK, filter)
		return
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get all categories
// @Description Get list of all book categories
// @Tags Categories
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} ErrorResponse
// @Router /categories [get]
func getCategories(c *gin.Context) {
	rows, err := db.Query("SELECT DISTINCT category FROM books WHERE category IS NOT NULL")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, category)
	}

	if categories == nil {
		categories = []string{}
	}

	c.JSON(http.StatusOK, categories)
}

// @Summary Search books by category
// @Description Get books filtered by category
// @Tags Books
// @Produce json
// @Param category path string true "Category name"
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /books/{category} [get]
func searchCategories(c *gin.Context) {
	category := c.Param("category")

	rows, err := db.Query("SELECT id, title, author, isbn, year, price, category, created_at, updated_at FROM books WHERE category = $1", category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.Category, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Search books by keyword
// @Description Search books by title or author
// @Tags Books
// @Produce json
// @Param q query string true "Search keyword"
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /books/search [get]
func searchBooks(c *gin.Context) {
	keyword := c.Query("q")

	rows, err := db.Query("SELECT id, title, author, isbn, year, price, created_at, updated_at FROM books WHERE title ILIKE $1 OR author ILIKE $1", "%"+keyword+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get featured books
// @Description Get list of featured books (high rating)
// @Tags Books
// @Produce json
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /books/featured [get]
func getFeaturedBooks(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, author, isbn, year, price, rating, created_at, updated_at FROM books WHERE rating >= 4.5 ORDER BY rating DESC LIMIT 10")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.Rating, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @Summary Get discounted books
// @Description Get list of books with discounts
// @Tags Books
// @Produce json
// @Success 200 {array} Book
// @Failure 500 {object} ErrorResponse
// @Router /books/discounted [get]
func getDiscountedBooks(c *gin.Context) {
	rows, err := db.Query("SELECT id, title, author, isbn, year, price, original_price, discount, created_at, updated_at FROM books WHERE discount > 0 ORDER BY discount DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Year, &book.Price, &book.OriginalPrice, &book.Discount, &book.CreatedAt, &book.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	if books == nil {
		books = []Book{}
	}

	c.JSON(http.StatusOK, books)
}

// @title           Bookstore API
// @version         1.0
// @description     This is a bookstore API with various endpoints for managing and searching books
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"message": "unhealthy", "error": err})
			return
		}
		c.JSON(200, gin.H{"message": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/books", getAllBooks)
		api.GET("/books/latest", getNewBooks)
		api.GET("/books/:id", getBook)
		api.POST("/books", createBook)
		api.PUT("/books/:id", updateBook)
		api.DELETE("/books/:id", deleteBook)

		api.GET("/categories", getCategories)
		api.GET("/books/search", searchBooks)
		api.GET("/books/featured", getFeaturedBooks)
		api.GET("/books/new", getNewBooks)
		api.GET("/books/discounted", getDiscountedBooks)
		api.GET("/books/:category", searchCategories)
	}

	r.Run(":8080")
}
