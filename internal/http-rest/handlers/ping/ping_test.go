package ping_test

import (
	"io"
	"log/slog"
	"test4effectivemobile/internal/http-rest/handlers/ping"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	handlerPing := ping.New(log)

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	handlerPing(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(
		t,
		"{\"status\":\"OK\",\"message\":\"Pong\"}\n",
		string(data),
	)
}
