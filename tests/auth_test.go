package tests

import (
	"net/http/httptest"
	"strings"
	"testing"
	"wallet-api/handlers"
	"wallet-api/mocks"
	"wallet-api/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUpSuccess(t *testing.T) {
	testRepo := new(mocks.MockUserRepository)

	testRepo.On("CreateUserWithWallet", mock.AnythingOfType("User")).Return(models.User{ID: 1, Username: "yasmin", Role: "user"}, nil)

	body := strings.NewReader(`{"username": "yasmin", "role":"user", "password":"test123"}`)
	req := httptest.NewRequest("POST", "/signup", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))
	c.Request = req

	authHandler := handlers.NewAuthHandler(testRepo)
	authHandler.SignUp(c)
	assert.Equal(t, 201, w.Code)
	assert.Contains(t, w.Body.String(), "yasmin")
}
