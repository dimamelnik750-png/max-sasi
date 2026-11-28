package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"maxsasi/internal/cache"
	"maxsasi/internal/config"
	"maxsasi/internal/database"
	"maxsasi/internal/handler"
	"maxsasi/internal/repository"
	"maxsasi/internal/service"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

var (
	testDB             *sql.DB
	testPGContainer    *tcpostgres.PostgresContainer
	testRedisContainer *tcredis.RedisContainer
	testRedisHost      string
	testRedisPort      int
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16",
		tcpostgres.WithDatabase("todos_test"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("1234"),
		tcpostgres.WithInitScripts(),
	)
	if err != nil {
		panic(err)
	}
	testPGContainer = pgContainer

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	deadline := time.Now().Add(10 * time.Second)
	for testDB.PingContext(ctx) != nil {
		if time.Now().After(deadline) {
			panic("PostgreSQL did not become ready")
		}
		time.Sleep(100 * time.Millisecond)
	}

	redisContainer, err := tcredis.Run(ctx, "redis:7")
	if err != nil {
		panic(err)
	}
	testRedisContainer = redisContainer

	redisEndpoint, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		panic(err)
	}

	if err := database.RunMigrations(testDB); err != nil {
		panic(err)
	}

	host, portString, err := net.SplitHostPort(redisEndpoint)
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic(err)
	}

	testRedisHost = host
	testRedisPort = port

	code := m.Run()

	_ = testDB.Close()
	_ = testPGContainer.Terminate(ctx)
	_ = testRedisContainer.Terminate(ctx)
	os.Exit(code)
}

func setupTest(t *testing.T) http.Handler {
	t.Helper()

	cleanDatabase(t)
	t.Cleanup(func() {
		cleanDatabase(t)
	})

	cfg := config.Config{JWTSecret: "test-secret"}
	todoRepo := repository.NewPostgresTodoRepository(testDB)
	userRepo := repository.NewPostgresUserRepository(testDB)
	redisCache := cache.NewRedisCache(testRedisHost, testRedisPort)
	if err := redisCache.DeleteByPrefix(context.Background(), "todos:"); err != nil {
		t.Fatalf("clean Redis cache: %v", err)
	}
	todoService := service.NewTodoService(todoRepo)
	cachedTodoService := service.NewCachedTodoService(todoService, redisCache)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	httpHandler := handler.New(cachedTodoService)
	authHandler := handler.NewAuthHandler(authService)
	mux := handler.NewRouter(httpHandler, authHandler)

	return handler.CORS(
		handler.Recovery(
			handler.JWTAuth(cfg.JWTSecret)(
				handler.Gzip(mux),
			),
		),
	)
}

func cleanDatabase(t *testing.T) {
	t.Helper()
	if _, err := testDB.Exec("DELETE FROM todos"); err != nil {
		t.Fatalf("clean todos: %v", err)
	}
	if _, err := testDB.Exec("DELETE FROM users"); err != nil {
		t.Fatalf("clean users: %v", err)
	}
}

func doRequest(server http.Handler, method, path, token string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)
	return rr
}

func registerAndLogin(server http.Handler) string {
	doRequest(server, http.MethodPost, "/auth/register", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	rr := doRequest(server, http.MethodPost, "/auth/login", "", map[string]string{
		"username": "testuser",
		"password": "testpass123",
	})

	var tokens map[string]string
	_ = json.NewDecoder(rr.Body).Decode(&tokens)
	return tokens["access_token"]
}
