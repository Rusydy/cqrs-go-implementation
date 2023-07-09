package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

// Command represents a write operation command.
type Command interface {
	Validate() error
}

// Query represents a read operation query.
type Query interface {
	Validate() error
}

// ArticleWriteModel represents the write model for an article.
type ArticleWriteModel struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

// ArticleReadModel represents the read model for an article.
type ArticleReadModel struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Created string `json:"created"`
}

// CreateArticleCommand represents the command for creating an article.
type CreateArticleCommand struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// Validate validates the create article command.
func (c *CreateArticleCommand) Validate() error {
	if c.Author == "" {
		return errors.New("author is required")
	}
	if c.Title == "" {
		return errors.New("title is required")
	}
	if c.Body == "" {
		return errors.New("body is required")
	}
	return nil
}

// GetArticlesQuery represents the query for retrieving articles.
type GetArticlesQuery struct {
	Query  string `json:"query"`
	Author string `json:"author"`
}

// Validate validates the get articles query.
func (q *GetArticlesQuery) Validate() error {
	return nil
}

// ArticleWriteRepository represents the repository for write operations on articles.
type ArticleWriteRepository interface {
	Create(article *ArticleWriteModel) error
}

// ArticleReadRepository represents the repository for read operations on articles.
type ArticleReadRepository interface {
	GetAll() ([]ArticleReadModel, error)
}

// ArticleWriteService represents the service for write operations on articles.
type ArticleWriteService struct {
	repo ArticleWriteRepository
}

// CreateArticle creates a new article.
func (s *ArticleWriteService) CreateArticle(command *CreateArticleCommand) error {
	err := command.Validate()
	if err != nil {
		return err
	}

	article := &ArticleWriteModel{
		Author:  command.Author,
		Title:   command.Title,
		Body:    command.Body,
		Created: time.Now(),
	}

	err = s.repo.Create(article)
	if err != nil {
		return err
	}

	return nil
}

// ArticleReadService represents the service for read operations on articles.
type ArticleReadService struct {
	repo ArticleReadRepository
}

// GetArticles retrieves a list of articles.
func (s *ArticleReadService) GetArticles(query *GetArticlesQuery) ([]ArticleReadModel, error) {
	err := query.Validate()
	if err != nil {
		return nil, err
	}

	articles, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// ArticleWriteRepositoryImpl represents the PostgreSQL repository for write operations on articles.
type ArticleWriteRepositoryImpl struct {
	db *sql.DB
}

// Create creates a new article in the PostgreSQL database.
func (r *ArticleWriteRepositoryImpl) Create(article *ArticleWriteModel) error {
	// Perform the create operation on the database
	_, err := r.db.Exec("INSERT INTO articles(author, title, body, created) VALUES($1, $2, $3, $4)",
		article.Author, article.Title, article.Body, article.Created)
	if err != nil {
		return err
	}

	return nil
}

// ArticleReadRepositoryImpl represents the PostgreSQL repository for read operations on articles.
type ArticleReadRepositoryImpl struct {
	db *sql.DB
}

// GetAll retrieves all articles from the PostgreSQL database.
func (r *ArticleReadRepositoryImpl) GetAll() ([]ArticleReadModel, error) {
	// Perform the query operation on the database
	rows, err := r.db.Query("SELECT id, author, title, body, created FROM articles ORDER BY created DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]ArticleReadModel, 0)
	for rows.Next() {
		var article ArticleReadModel
		err := rows.Scan(&article.ID, &article.Author, &article.Title, &article.Body, &article.Created)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// ArticleHandler represents the HTTP handler for articles.
type ArticleHandler struct {
	writeService *ArticleWriteService
	readService  *ArticleReadService
}

// CreateArticle handles the creation of a new article.
func (h *ArticleHandler) CreateArticle(c echo.Context) error {
	command := new(CreateArticleCommand)
	if err := c.Bind(command); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	err := h.writeService.CreateArticle(command)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create article")
	}

	return c.JSON(http.StatusCreated, "Article created successfully")
}

// GetArticles handles the retrieval of articles.
func (h *ArticleHandler) GetArticles(c echo.Context) error {
	query := new(GetArticlesQuery)
	if err := c.Bind(query); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request query")
	}

	articles, err := h.readService.GetArticles(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve articles")
	}

	return c.JSON(http.StatusOK, articles)
}

func main() {
	// Retrieve PostgreSQL connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	// Construct the connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	// Initialize the PostgreSQL database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create instances of the repositories
	writeRepo := &ArticleWriteRepositoryImpl{db: db}
	readRepo := &ArticleReadRepositoryImpl{db: db}

	// Create instances of the services
	writeService := &ArticleWriteService{repo: writeRepo}
	readService := &ArticleReadService{repo: readRepo}

	// Create an instance of the article handler
	articleHandler := &ArticleHandler{
		writeService: writeService,
		readService:  readService,
	}

	// Initialize the Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define the HTTP routes
	e.POST("/articles", articleHandler.CreateArticle)
	e.GET("/articles", articleHandler.GetArticles)

	// Error handling middleware
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		// Check if it's an Echo HTTPError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		}

		// Send error response and handle the error
		if err := c.JSON(code, map[string]interface{}{
			"error": message,
		}); err != nil {
			log.Println("Failed to send error response:", err)
		}
	}

	// Start the server
	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
