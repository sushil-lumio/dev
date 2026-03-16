package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sushil-lumio/user-api/db"
	"github.com/sushil-lumio/user-api/models"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	database, err := db.NewTestDB()
	if err != nil {
		t.Fatal("failed to create test db:", err)
	}
	DB = database
}

func registerTestUser(t *testing.T, name, email, password string) models.AuthResponse {
	t.Helper()
	body, _ := json.Marshal(models.User{Name: name, Email: email, Password: password})
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	RegisterHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: status %d, body: %s", w.Code, w.Body.String())
	}
	var resp models.AuthResponse
	json.NewDecoder(w.Body).Decode(&resp)
	return resp
}

// --- Auth Tests ---

func TestRegister(t *testing.T) {
	setupTestDB(t)

	resp := registerTestUser(t, "Alice", "alice@example.com", "secret123")

	if resp.Token == "" {
		t.Error("expected token to be set")
	}
	if resp.User.Name != "Alice" {
		t.Errorf("expected name Alice, got %s", resp.User.Name)
	}
	if resp.User.Email != "alice@example.com" {
		t.Errorf("expected email alice@example.com, got %s", resp.User.Email)
	}
	if resp.User.Password != "" {
		t.Error("password should not be in response")
	}
}

func TestRegisterMissingFields(t *testing.T) {
	setupTestDB(t)

	body := `{"name":"Alice","email":"alice@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	setupTestDB(t)

	registerTestUser(t, "Alice", "alice@example.com", "secret123")

	body := `{"name":"Alice2","email":"alice@example.com","password":"other"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	RegisterHandler(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d", w.Code)
	}
}

func TestLogin(t *testing.T) {
	setupTestDB(t)

	registerTestUser(t, "Bob", "bob@example.com", "password123")

	body := `{"email":"bob@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	LoginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp models.AuthResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Token == "" {
		t.Error("expected token")
	}
	if resp.User.Name != "Bob" {
		t.Errorf("expected name Bob, got %s", resp.User.Name)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	setupTestDB(t)

	registerTestUser(t, "Carol", "carol@example.com", "correct")

	body := `{"email":"carol@example.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	LoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestLoginUnknownEmail(t *testing.T) {
	setupTestDB(t)

	body := `{"email":"nobody@example.com","password":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	LoginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

// --- Middleware Tests ---

func TestAuthMiddlewareNoToken(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	AuthMiddleware(UserHandler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	setupTestDB(t)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	AuthMiddleware(UserHandler)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

// --- Protected User Endpoint Tests ---

func TestGetUsersWithAuth(t *testing.T) {
	setupTestDB(t)

	resp := registerTestUser(t, "Dave", "dave@example.com", "secret")

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+resp.Token)
	w := httptest.NewRecorder()
	AuthMiddleware(UserHandler)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var users []models.User
	json.NewDecoder(w.Body).Decode(&users)

	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestGetUserByIDWithAuth(t *testing.T) {
	setupTestDB(t)

	resp := registerTestUser(t, "Eve", "eve@example.com", "secret")

	req := httptest.NewRequest(http.MethodGet, "/users?id="+resp.User.ID, nil)
	req.Header.Set("Authorization", "Bearer "+resp.Token)
	w := httptest.NewRecorder()
	AuthMiddleware(UserHandler)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var user models.User
	json.NewDecoder(w.Body).Decode(&user)

	if user.Name != "Eve" {
		t.Errorf("expected name Eve, got %s", user.Name)
	}
}

func TestGetUserNotFound(t *testing.T) {
	setupTestDB(t)

	resp := registerTestUser(t, "Frank", "frank@example.com", "secret")

	req := httptest.NewRequest(http.MethodGet, "/users?id=999", nil)
	req.Header.Set("Authorization", "Bearer "+resp.Token)
	w := httptest.NewRecorder()
	AuthMiddleware(UserHandler)(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	setupTestDB(t)

	resp := registerTestUser(t, "Grace", "grace@example.com", "secret")

	req := httptest.NewRequest(http.MethodDelete, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+resp.Token)
	w := httptest.NewRecorder()
	AuthMiddleware(UserHandler)(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", w.Code)
	}
}
