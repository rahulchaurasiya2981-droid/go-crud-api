package httprequest

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rahulchaurasiya2981-droid/go-crud-api/internal/user/dto"
)

func TestParseAndValidateJSONValidatesCreateUserFields(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{
			name: "valid request",
			body: `{"name":"Ada","email":"ada@example.com","age":0}`,
			want: true,
		},
		{
			name: "missing name",
			body: `{"email":"ada@example.com","age":37}`,
		},
		{
			name: "missing email",
			body: `{"name":"Ada","age":37}`,
		},
		{
			name: "missing age",
			body: `{"name":"Ada","email":"ada@example.com"}`,
		},
		{
			name: "invalid email",
			body: `{"name":"Ada","email":"not-an-email","age":37}`,
		},
		{
			name: "negative age",
			body: `{"name":"Ada","email":"ada@example.com","age":-1}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/users", strings.NewReader(test.body))
			response := httptest.NewRecorder()
			var destination dto.CreateUserRequest

			got := ParseAndValidateJSON(response, request, &destination, 1_000_000)
			if got != test.want {
				t.Fatalf("ParseAndValidateJSON() = %t, want %t", got, test.want)
			}
		})
	}
}
