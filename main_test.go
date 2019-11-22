package todofun

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		body string
		want int
	}{
		{body: `{"title": ""}`, want: http.StatusNotFound},
		{body: `{"title": "book", "url": "example.com", "id": "1"}`, want: http.StatusOK},
	}

	for _, test := range tests {
		req := httptest.NewRequest("GET", "/", strings.NewReader(test.body))
		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		CreateTodo(rr, req)

		if got := rr.Code; got != test.want {
			t.Errorf("HTTP status %v, want %v", got, test.want)
		}
	}
}
