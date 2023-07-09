package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockArticleRepository is a mock implementation of the ArticleRepository interface.
type MockArticleRepository struct {
	createCalled   bool
	listCalled     bool
	createError    error
	listError      error
	listResult     []Article
	expectedCreate *Article
	expectedQuery  *ArticleQuery
}

func (m *MockArticleRepository) Create(article *Article) error {
	m.createCalled = true
	m.expectedCreate = article
	return m.createError
}

func (m *MockArticleRepository) List(query *ArticleQuery) ([]Article, error) {
	m.listCalled = true
	m.expectedQuery = query
	return m.listResult, m.listError
}

func TestCreateArticle(t *testing.T) {
	repo := &MockArticleRepository{}
	service := &ArticleService{repo: repo}

	// Define the test case
	command := &CreateArticleCommand{
		Author: "John Doe",
		Title:  "Test Article",
		Body:   "This is a test article",
	}

	// Execute the method
	err := service.CreateArticle(command)

	// Assert the results
	assert.NoError(t, err)
	assert.True(t, repo.createCalled)
	assert.Equal(t, command.Author, repo.expectedCreate.Author)
	assert.Equal(t, command.Title, repo.expectedCreate.Title)
	assert.Equal(t, command.Body, repo.expectedCreate.Body)
}

func TestGetArticles(t *testing.T) {
	repo := &MockArticleRepository{
		listResult: []Article{
			{ID: 1, Author: "John Doe", Title: "Article 1", Body: "Body 1"},
			{ID: 2, Author: "Jane Smith", Title: "Article 2", Body: "Body 2"},
		},
	}
	service := &ArticleService{repo: repo}

	// Define the test case
	query := &ArticleQuery{
		Query:  "test",
		Author: "John Doe",
	}

	// Execute the method
	articles, err := service.GetArticles(query)

	// Assert the results
	assert.NoError(t, err)
	assert.True(t, repo.listCalled)
	assert.Equal(t, query, repo.expectedQuery)
	assert.Equal(t, repo.listResult, articles)
}

func TestGetArticles_Error(t *testing.T) {
	repo := &MockArticleRepository{
		listError: errors.New("list error"),
	}
	service := &ArticleService{repo: repo}

	// Define the test case
	query := &ArticleQuery{
		Query:  "test",
		Author: "John Doe",
	}

	// Execute the method
	articles, err := service.GetArticles(query)

	// Assert the error
	assert.Error(t, err)
	assert.Nil(t, articles)
	assert.True(t, repo.listCalled)
	assert.Equal(t, query, repo.expectedQuery)
}
