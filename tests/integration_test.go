package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"wallet-api/db"
	"wallet-api/handlers"
	"wallet-api/middleware"
	"wallet-api/repositories"
	"wallet-api/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var testRouter *gin.Engine

func setupTestDB(t *testing.T) {
	godotenv.Load("../.env")
	os.Setenv("DB_NAME", os.Getenv("TEST_DB_NAME"))

	database, err := db.Connect()
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}

	// clean up before each test
	database.Exec("DELETE FROM transactions")
	database.Exec("DELETE FROM wallets")
	database.Exec("DELETE FROM users")

	userRepo := repositories.NewUserRepository(database)
	walletRepo := repositories.NewWalletRepository(database)
	transactionRepo := repositories.NewTransactionRepository(database)
	walletService := services.NewWalletService(database, walletRepo, transactionRepo, userRepo)

	authHandler := handlers.NewAuthHandler(userRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	gin.SetMode(gin.TestMode)
	testRouter = gin.New()

	testRouter.POST("/signup", authHandler.SignUp)
	testRouter.POST("/login", authHandler.LogIn)

	protected := testRouter.Group("/wallet")
	protected.Use(middleware.RequireAuth())
	{
		protected.GET("", walletHandler.GetWallet)
		protected.POST("/deposit", walletHandler.Deposit)
		protected.POST("/withdraw", walletHandler.Withdraw)
		protected.POST("/transfer", walletHandler.Transfer)
	}
}

func createUserAndGetToken(t *testing.T, username, password string) string {
	// signup
	body := fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("signup failed: %s", w.Body.String())
	}

	// login
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["token"]
}

func depositAmount(t *testing.T, token string, amount int) {
	body := fmt.Sprintf(`{"amount":%d,"category":"test","note":"test deposit"}`, amount)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/wallet/deposit", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("deposit failed: %s", w.Body.String())
	}
}

func TestTransferHappyPath(t *testing.T) {
	setupTestDB(t)

	token1 := createUserAndGetToken(t, "sender", "password")
	createUserAndGetToken(t, "receiver", "password")
	depositAmount(t, token1, 200)

	body := `{"toUsername":"receiver","amount":50,"category":"gift","note":"here"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/wallet/transfer", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 got %d: %s", w.Code, w.Body.String())
	}
}

func TestTransferInsufficientFunds(t *testing.T) {
	setupTestDB(t)

	token1 := createUserAndGetToken(t, "sender", "password")
	createUserAndGetToken(t, "receiver", "password")
	depositAmount(t, token1, 30)

	body := `{"toUsername":"receiver","amount":100,"category":"gift","note":"too much"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/wallet/transfer", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d: %s", w.Code, w.Body.String())
	}
}

func TestTransferNonexistentUser(t *testing.T) {
	setupTestDB(t)

	token1 := createUserAndGetToken(t, "sender", "password")
	depositAmount(t, token1, 200)

	body := `{"toUsername":"ghost","amount":50,"category":"gift","note":"nobody"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/wallet/transfer", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 got %d: %s", w.Code, w.Body.String())
	}
}

func TestTransferSelf(t *testing.T) {
	setupTestDB(t)

	token1 := createUserAndGetToken(t, "sender", "password")
	depositAmount(t, token1, 200)

	body := `{"toUsername":"sender","amount":50,"category":"gift","note":"self"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/wallet/transfer", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token1)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 got %d: %s", w.Code, w.Body.String())
	}
}

func TestConcurrentWithdrawals(t *testing.T) {
	setupTestDB(t)

	token1 := createUserAndGetToken(t, "sender", "password")
	depositAmount(t, token1, 100)

	var wg sync.WaitGroup
	results := make([]int, 2)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			body := `{"amount":100,"category":"test","note":"concurrent"}`
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/wallet/withdraw", bytes.NewBufferString(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token1)
			testRouter.ServeHTTP(w, req)
			results[index] = w.Code
		}(i)
	}

	wg.Wait()

	successCount := 0
	for _, code := range results {
		if code == http.StatusOK {
			successCount++
		}
	}

	if successCount != 1 {
		t.Errorf("expected exactly 1 success, got %d", successCount)
	}
}
