package api

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"

	"github.com/nnkken/oracle-fetch/db"
)

var testConn *pgx.Conn

func TestMain(m *testing.M) {
	testConn = db.SetupTestConn()
	code := m.Run()
	os.Exit(code)
}

var _ http.ResponseWriter = &testResponseWriter{}

type testResponseWriter struct {
	HttpHeader http.Header
	Written    []byte
	Code       int
}

func NewTestResponseWriter() *testResponseWriter {
	return &testResponseWriter{
		HttpHeader: make(http.Header),
		Code:       200,
	}
}

func (w *testResponseWriter) Header() http.Header {
	return w.HttpHeader
}

func (w *testResponseWriter) Write(bz []byte) (int, error) {
	w.Written = append(w.Written, bz...)
	return len(bz), nil
}

func (w *testResponseWriter) WriteHeader(code int) {
	w.Code = code
}

func newTestContext(requestURL string) (*gin.Context, *testResponseWriter) {
	u, err := url.Parse(requestURL)
	if err != nil {
		panic(err)
	}
	writer := NewTestResponseWriter()
	ctx, _ := gin.CreateTestContext(writer)
	ctx.Request = &http.Request{
		URL: u,
	}
	ctx.Set("conn", testConn)
	ctx.Set("logger", zap.S())
	return ctx, writer
}
