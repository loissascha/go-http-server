package server

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerSetup(t *testing.T) {
	// test simple server (without any translation and server settings)
	s, err := NewServer()
	s.GET("/test", testRoute)
	s.POST("/test", testRoute)
	s.PUT("/test", testRoute)
	s.DELETE("/test", testRoute)
	assert.Equal(t, 0, len(s.Options))
	assert.False(t, s.TranslationsEnabled)
	assert.False(t, s.AutoDetectLanguageEnabled)
	assert.Equal(t, s.DefaultLanguage, "")
	assert.Equal(t, 5*time.Second, s.ReadHeaderTimeout)
	assert.Equal(t, 15*time.Second, s.ReadTimeout)
	assert.Equal(t, 15*time.Second, s.WriteTimeout)
	assert.Equal(t, 60*time.Second, s.IdleTimeout)
	assert.Equal(t, 1<<20, s.MaxHeaderBytes)
	assert.Equal(t, 1, len(s.Paths)) // has only one path because "/test" always refers to the same path
	assert.Nil(t, err)

	// test if language settings work
	s, err = NewServer(
		EnableTranslations(),
		EnableAutoDetectLanguage(),
		SetDefaultLanguage("en"),
		AddTranslationFile("en", "en_test.json"),
		AddTranslationFile("de", "de_test.json"),
	)
	s.GET("/test", testRoute)
	s.POST("/test", testRoute)
	s.PUT("/test", testRoute)
	s.DELETE("/test", testRoute)
	s.POSTI("/test/no/langs", testRoute)
	assert.Equal(t, 5, len(s.Options))
	assert.True(t, s.TranslationsEnabled)
	assert.True(t, s.AutoDetectLanguageEnabled)
	assert.Equal(t, s.DefaultLanguage, "en")
	assert.Equal(t, 4, len(s.Paths)) // 3 routes for "/test" ("/test", "/en/test", "/de/test") and one for "/test/no/langs"
	assert.Contains(t, s.Paths, "/test")
	assert.Contains(t, s.Paths, "/de/test")
	assert.Contains(t, s.Paths, "/en/test")
	assert.NotContains(t, s.Paths, "fr/test")
	assert.Nil(t, err)

	// test if the registered handlers actually work
	err = s.setupHandlers()
	assert.Nil(t, err)
	testserver := httptest.NewServer(s.mux)
	defer testserver.Close()

	resp, err := http.Get(testserver.URL + "/test")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotAcceptable)

	resp, err = http.Get(testserver.URL + "/en/test")
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotAcceptable)

	resp, err = http.Post(testserver.URL+"/test", "applicaton/json", bytes.NewBuffer([]byte("{'key':'value'}")))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotAcceptable)

	resp, err = http.Post(testserver.URL+"/en/test", "applicaton/json", bytes.NewBuffer([]byte("{'key':'value'}")))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotAcceptable)

	resp, err = http.Post(testserver.URL+"/de/test", "applicaton/json", bytes.NewBuffer([]byte("{'key':'value'}")))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotAcceptable)

	resp, err = http.Post(testserver.URL+"/fr/test", "applicaton/json", bytes.NewBuffer([]byte("{'key':'value'}")))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)

	resp, err = http.Post(testserver.URL+"/test/no/langs", "applicaton/json", bytes.NewBuffer([]byte("{'key':'value'}")))
	assert.Nil(t, err)
	body, err := io.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Test was successful\n", string(body))
	assert.Equal(t, resp.StatusCode, http.StatusNotAcceptable)

	resp, err = http.Post(testserver.URL+"/en/test/no/langs", "applicaton/json", bytes.NewBuffer([]byte("{'key':'value'}")))
	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func testRoute(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Test was successful", http.StatusNotAcceptable)
}

func TestOPTIONSPreflightOnPostRouteReturnsNoContentAndCORSHeaders(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "*")
	t.Setenv("APP_ENV", "development")
	s, err := NewServer()
	assert.NoError(t, err)
	s.POSTI("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	err = s.setupHandlers()
	assert.NoError(t, err)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodOptions, ts.URL+"/auth/register", nil)
	assert.NoError(t, err)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, "http://localhost:5173", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, resp.Header.Get("Access-Control-Allow-Headers"), "Content-Type")
}

func TestUnsupportedMethodReturnsMethodNotAllowed(t *testing.T) {
	s, err := NewServer()
	assert.NoError(t, err)
	s.GETI("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	err = s.setupHandlers()
	assert.NoError(t, err)
	ts := httptest.NewServer(s.mux)
	defer ts.Close()
	req, err := http.NewRequest(http.MethodPatch, ts.URL+"/health", nil)
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestServerTimeoutOptions(t *testing.T) {
	s, err := NewServer(
		SetReadHeaderTimeout(2*time.Second),
		SetReadTimeout(10*time.Second),
		SetWriteTimeout(20*time.Second),
		SetIdleTimeout(90*time.Second),
		SetMaxHeaderBytes(2048),
	)
	assert.NoError(t, err)
	assert.Equal(t, 2*time.Second, s.ReadHeaderTimeout)
	assert.Equal(t, 10*time.Second, s.ReadTimeout)
	assert.Equal(t, 20*time.Second, s.WriteTimeout)
	assert.Equal(t, 90*time.Second, s.IdleTimeout)
	assert.Equal(t, 2048, s.MaxHeaderBytes)
}

func TestInvalidServerTimeoutOptionReturnsError(t *testing.T) {
	s, err := NewServer(ServerOption{Name: READ_TIMEOUT, Value: "not-a-duration"})
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid read timeout")
}

func TestShutdownBeforeServeReturnsNil(t *testing.T) {
	s, err := NewServer()
	assert.NoError(t, err)

	err = s.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestServeReturnsNilAfterShutdown(t *testing.T) {
	s, err := NewServer()
	assert.NoError(t, err)
	s.GETI("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := getFreeAddr(t)
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Serve(addr)
	}()

	waitForServer(t, addr)

	err = s.Shutdown(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, <-errCh)

	_, err = http.Get("http://" + addr + "/health")
	assert.Error(t, err)
}

func TestServeReturnsErrorWhenAlreadyRunning(t *testing.T) {
	s, err := NewServer()
	assert.NoError(t, err)
	s.GETI("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	addr := getFreeAddr(t)
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Serve(addr)
	}()

	waitForServer(t, addr)

	err = s.Serve(addr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "server already running")

	err = s.Shutdown(context.Background())
	assert.NoError(t, err)
	assert.NoError(t, <-errCh)
}

func TestShutdownWaitsForInFlightRequest(t *testing.T) {
	s, err := NewServer()
	assert.NoError(t, err)

	started := make(chan struct{}, 1)
	finished := make(chan struct{}, 1)
	s.GETI("/slow", func(w http.ResponseWriter, r *http.Request) {
		started <- struct{}{}
		time.Sleep(150 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		finished <- struct{}{}
	})

	addr := getFreeAddr(t)
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Serve(addr)
	}()

	waitForServer(t, addr)

	respCh := make(chan error, 1)
	go func() {
		resp, err := http.Get("http://" + addr + "/slow")
		if err != nil {
			respCh <- err
			return
		}
		defer resp.Body.Close()
		respCh <- nil
	}()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not start")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Shutdown(shutdownCtx)
	assert.NoError(t, err)

	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("in-flight request did not finish before shutdown")
	}

	assert.NoError(t, <-respCh)
	assert.NoError(t, <-errCh)
}

func getFreeAddr(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate free address: %v", err)
	}
	defer listener.Close()

	return listener.Addr().String()
}

func waitForServer(t *testing.T, addr string) {
	t.Helper()

	client := &http.Client{Timeout: 100 * time.Millisecond}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		resp, err := client.Get("http://" + addr + "/health")
		if err == nil {
			resp.Body.Close()
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("server at %s did not start in time", addr)
}
