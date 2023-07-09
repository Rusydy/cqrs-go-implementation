package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// MockWriteRepository represents a mock repository for write operations on articles.
type MockWriteRepository struct {
	articles []*ArticleWriteModel
}

// Create adds a new article to the repository.
func (r *MockWriteRepository) Create(article *ArticleWriteModel) error {
	r.articles = append(r.articles, article)
	return nil
}

// MockReadRepository represents a mock repository for read operations on articles.
type MockReadRepository struct {
	articles []*ArticleReadModel
}

// GetAll retrieves all articles from the repository.
func (r *MockReadRepository) GetAll() ([]ArticleReadModel, error) {
	articles := make([]ArticleReadModel, len(r.articles))
	for i, article := range r.articles {
		articles[i] = *article
	}
	return articles, nil
}

func TestCreateArticleHandler(t *testing.T) {
	// Create a mock write repository
	mockWriteRepo := &MockWriteRepository{}

	// Create an article handler with the mock repository
	articleHandler := &ArticleHandler{
		writeService: &ArticleWriteService{repo: mockWriteRepo},
	}

	// Create a new Echo instance
	e := echo.New()

	// Create a POST request with a JSON payload
	req := httptest.NewRequest(http.MethodPost, "/articles", strings.NewReader(`{
		"author": "John Doe",
		"title": "Hello World",
		"body": "This is the article body"
	}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Invoke the CreateArticle handler
	err := articleHandler.CreateArticle(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the response code is HTTP 201 Created
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Assert that the article was added to the repository
	assert.Equal(t, 1, len(mockWriteRepo.articles))
	assert.Equal(t, "John Doe", mockWriteRepo.articles[0].Author)
	assert.Equal(t, "Hello World", mockWriteRepo.articles[0].Title)
	assert.Equal(t, "This is the article body", mockWriteRepo.articles[0].Body)
}

func TestGetArticlesHandler(t *testing.T) {
	// Create a mock read repository with some articles
	mockReadRepo := &MockReadRepository{
		articles: []*ArticleReadModel{
			{
				ID:     1,
				Author: "John Doe",
				Title:  "Hello World",
				Body:   "This is the article body",
				// convert to string to match the JSON format
				Created: time.Now().String(),
			},
			{
				ID:      2,
				Author:  "Jane Smith",
				Title:   "Greetings",
				Body:    "Welcome to the world",
				Created: time.Now().String(),
			},
		},
	}

	// Create an article handler with the mock repository
	articleHandler := &ArticleHandler{
		readService: &ArticleReadService{repo: mockReadRepo},
	}

	// Create a new Echo instance
	e := echo.New()

	// Create a GET request
	req := httptest.NewRequest(http.MethodGet, "/articles", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Invoke the GetArticles handler
	err := articleHandler.GetArticles(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert that the response code is HTTP 200 OK
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse the response body
	var response []ArticleReadModel
	err = json.Unmarshal(rec.Body.Bytes(), &response)

	// Assert that the response body was parsed successfully
	assert.NoError(t, err)

	// Assert that the correct number of articles were returned
	assert.Equal(t, len(mockReadRepo.articles), len(response))

	// Assert the contents of the first article
	assert.Equal(t, mockReadRepo.articles[0].ID, response[0].ID)
	assert.Equal(t, mockReadRepo.articles[0].Author, response[0].Author)
	assert.Equal(t, mockReadRepo.articles[0].Title, response[0].Title)
	assert.Equal(t, mockReadRepo.articles[0].Body, response[0].Body)
}
