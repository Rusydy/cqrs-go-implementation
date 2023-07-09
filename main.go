package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

// Article represents an article.
type Article struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

// ArticleQuery represents a query to filter articles.
type ArticleQuery struct {
	Query  string `json:"query"`
	Author string `json:"author"`
}

// CreateArticleCommand represents a command to create a new article.
type CreateArticleCommand struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

// ArticleRepository provides an interface to interact with the article data store.
type ArticleRepository interface {
	Create(article *Article) error
	List(query *ArticleQuery) ([]Article, error)
}

// SQLArticleRepository is an implementation of ArticleRepository using PostgreSQL.
type SQLArticleRepository struct {
	db *sql.DB
}

// Create creates a new article in the database.
func (repo *SQLArticleRepository) Create(article *Article) error {
	query := `
		INSERT INTO articles (author, title, body, created)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := repo.db.QueryRow(query, article.Author, article.Title, article.Body, article.Created).Scan(&article.ID)
	if err != nil {
		return err
	}

	return nil
}

// List retrieves a list of articles from the database based on the query parameters.
func (repo *SQLArticleRepository) List(query *ArticleQuery) ([]Article, error) {
	whereClauses := make([]string, 0)
	params := make([]interface{}, 0)

	if query.Query != "" {
		whereClauses = append(whereClauses, "(LOWER(title) LIKE $1 OR LOWER(body) LIKE $1)")
		params = append(params, "%"+strings.ToLower(query.Query)+"%")
	}

	if query.Author != "" {
		whereClauses = append(whereClauses, "LOWER(author) = $2")
		params = append(params, strings.ToLower(query.Author))
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	queryStr := `
		SELECT id, author, title, body, created
		FROM articles
		` + whereClause + `
		ORDER BY created DESC
	`

	rows, err := repo.db.Query(queryStr, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Author, &article.Title, &article.Body, &article.Created)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// ArticleService provides an interface to interact with articles.
type ArticleService struct {
	repo ArticleRepository
}

// CreateArticle creates a new article.
func (s *ArticleService) CreateArticle(command *CreateArticleCommand) error {
	article := &Article{
		Author:  command.Author,
		Title:   command.Title,
		Body:    command.Body,
		Created: time.Now(),
	}

	err := s.repo.Create(article)
	if err != nil {
		return err
	}

	return nil
}

// GetArticles returns a list of articles based on the query parameters.
func (s *ArticleService) GetArticles(query *ArticleQuery) ([]Article, error) {
	articles, err := s.repo.List(query)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

// ErrorHandler handles errors and sends appropriate HTTP responses.
func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := http.StatusText(code)

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = he.Message.(string)
	}

	c.JSON(code, map[string]interface{}{
		"error": message,
	})
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

	// Initialize the article repository
	repo := &SQLArticleRepository{db: db}

	// Initialize the article service
	service := &ArticleService{repo: repo}

	// Initialize the Echo router
	e := echo.New()

	// Define the HTTP endpoints
	e.POST("/articles", func(c echo.Context) error {
		var command CreateArticleCommand
		if err := c.Bind(&command); err != nil {
			return err
		}

		if err := service.CreateArticle(&command); err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, command)
	})

	e.GET("/articles", func(c echo.Context) error {
		query := &ArticleQuery{
			Query:  c.QueryParam("query"),
			Author: c.QueryParam("author"),
		}

		articles, err := service.GetArticles(query)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, articles)
	})

	// Register the error handler
	e.HTTPErrorHandler = ErrorHandler

	// Start the HTTP server
	fmt.Println("Server listening on port 8000")
	log.Fatal(e.Start(":8000"))
}
