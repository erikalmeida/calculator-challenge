package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sezzle/backend/internal/adapters/out/strategyregistry"
	"sezzle/backend/internal/application/calculator"
	"sezzle/backend/internal/domain/calculation"
)

type failingCalculateUseCase struct {
	err         error
	shouldPanic bool
}

func (u failingCalculateUseCase) Calculate(_ string, _ []float64) (float64, error) {
	if u.shouldPanic {
		panic("simulated unexpected failure")
	}
	return 0, u.err
}

func TestCalculateEndpoint(t *testing.T) {
	handler := NewCalculatorHandler(calculator.NewService(strategyregistry.NewInMemoryCalculationStrategyRegistry()))

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantCode   string
		wantResult *float64
		wantAllow  string
	}{
		{name: "successful calculation", method: http.MethodPost, body: `{"operation":"multiply","operands":[2,3,4]}`, wantStatus: http.StatusOK, wantResult: float64Pointer(24)},
		{name: "invalid JSON", method: http.MethodPost, body: `{`, wantStatus: http.StatusBadRequest, wantCode: "invalid_json"},
		{name: "unknown field", method: http.MethodPost, body: `{"operation":"add","operands":[2,3],"extra":true}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_json"},
		{name: "multiple JSON values", method: http.MethodPost, body: `{"operation":"add","operands":[2,3]} {}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_json"},
		{name: "missing operation", method: http.MethodPost, body: `{"operands":[2,3]}`, wantStatus: http.StatusBadRequest, wantCode: calculation.ErrorCodeValidation},
		{name: "unknown operation", method: http.MethodPost, body: `{"operation":"modulo","operands":[2,3]}`, wantStatus: http.StatusBadRequest, wantCode: calculation.ErrorCodeUnknownOperation},
		{name: "invalid arity", method: http.MethodPost, body: `{"operation":"divide","operands":[8,2,1]}`, wantStatus: http.StatusBadRequest, wantCode: calculation.ErrorCodeInvalidArity},
		{name: "division by zero", method: http.MethodPost, body: `{"operation":"divide","operands":[8,0]}`, wantStatus: http.StatusBadRequest, wantCode: calculation.ErrorCodeDivisionByZero},
		{name: "null operand", method: http.MethodPost, body: `{"operation":"add","operands":[null,2]}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_number_format"},
		{name: "comma decimal separator", method: http.MethodPost, body: `{"operation":"add","operands":["2,5",3]}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_number_format"},
		{name: "scientific notation", method: http.MethodPost, body: `{"operation":"multiply","operands":[1e3,2]}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_number_format"},
		{name: "too many decimal places", method: http.MethodPost, body: `{"operation":"add","operands":[0.00000000001,2]}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_number_format"},
		{name: "operand above range", method: http.MethodPost, body: `{"operation":"add","operands":[100000000000001,2]}`, wantStatus: http.StatusBadRequest, wantCode: calculation.ErrorCodeOperandOutOfRange},
		{name: "too many operands", method: http.MethodPost, body: calculationRequestBody("add", 21), wantStatus: http.StatusBadRequest, wantCode: calculation.ErrorCodeTooManyOperands},
		{name: "wrong method", method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed, wantCode: "method_not_allowed", wantAllow: http.MethodPost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBufferString(tt.body))
			if tt.method == http.MethodPost {
				request.Header.Set("Content-Type", "application/json")
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.wantStatus, response.Body.String())
			}
			if contentType := response.Header().Get("Content-Type"); contentType != "application/json; charset=utf-8" {
				t.Fatalf("Content-Type = %q, want JSON", contentType)
			}
			if requestID := response.Header().Get(requestIDHeader); requestID == "" {
				t.Fatal("X-Request-ID is empty")
			}
			if tt.wantAllow != "" && response.Header().Get("Allow") != tt.wantAllow {
				t.Fatalf("Allow = %q, want %q", response.Header().Get("Allow"), tt.wantAllow)
			}

			if tt.wantResult != nil {
				var body calculationResponse
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}
				if body.Result != *tt.wantResult {
					t.Fatalf("result = %v, want %v", body.Result, *tt.wantResult)
				}
				return
			}

			var body errorResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatalf("unmarshal response: %v", err)
			}
			if body.Error.Code != tt.wantCode {
				t.Fatalf("error code = %q, want %q", body.Error.Code, tt.wantCode)
			}
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	handler := NewCalculatorHandler(calculator.NewService(strategyregistry.NewInMemoryCalculationStrategyRegistry()))

	t.Run("returns an OK status", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/health", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
		}
		if response.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			t.Fatalf("Content-Type = %q, want JSON", response.Header().Get("Content-Type"))
		}
		if response.Header().Get(requestIDHeader) == "" {
			t.Fatal("X-Request-ID is empty")
		}

		var body healthResponse
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
		if body.Status != "ok" {
			t.Fatalf("status body = %q, want ok", body.Status)
		}
	})

	t.Run("rejects unsupported methods", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodPost, "/health", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertErrorResponse(t, response, http.StatusMethodNotAllowed, "method_not_allowed")
		if response.Header().Get("Allow") != http.MethodGet {
			t.Fatalf("Allow = %q, want %q", response.Header().Get("Allow"), http.MethodGet)
		}
	})
}

func TestCalculateEndpointRejectsUnsupportedMediaType(t *testing.T) {
	handler := NewCalculatorHandler(calculator.NewService(strategyregistry.NewInMemoryCalculationStrategyRegistry()))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"operation":"add","operands":[2,3]}`))
	request.Header.Set("Content-Type", "text/plain")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusUnsupportedMediaType, "unsupported_media_type")
}

func TestCalculateEndpointRejectsOversizedRequest(t *testing.T) {
	handler := NewCalculatorHandler(calculator.NewService(strategyregistry.NewInMemoryCalculationStrategyRegistry()))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(strings.Repeat(" ", maxRequestBodyBytes+1)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusRequestEntityTooLarge, "request_too_large")
}

func TestCalculateEndpointHandlesUnexpectedError(t *testing.T) {
	handler := NewCalculatorHandler(failingCalculateUseCase{err: errors.New("storage unavailable")})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"operation":"add","operands":[2,3]}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusInternalServerError, "internal_error")
}

func TestCalculateEndpointRecoversFromPanic(t *testing.T) {
	handler := NewCalculatorHandler(failingCalculateUseCase{shouldPanic: true})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(`{"operation":"add","operands":[2,3]}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	assertErrorResponse(t, response, http.StatusInternalServerError, "internal_error")
	if response.Header().Get(requestIDHeader) == "" {
		t.Fatal("panic response must include X-Request-ID")
	}
}

func TestAPIDocumentationEndpoints(t *testing.T) {
	handler := NewCalculatorHandler(calculator.NewService(strategyregistry.NewInMemoryCalculationStrategyRegistry()))

	t.Run("serves the OpenAPI specification", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/api/v1/openapi.yaml", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
		}
		if response.Header().Get("Content-Type") != "application/yaml; charset=utf-8" {
			t.Fatalf("Content-Type = %q, want YAML", response.Header().Get("Content-Type"))
		}
		if !strings.Contains(response.Body.String(), "openapi: 3.1.0") {
			t.Fatal("OpenAPI specification is missing its version")
		}
	})

	t.Run("serves Swagger UI", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, "/docs", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
		}
		if !strings.Contains(response.Body.String(), "SwaggerUIBundle") {
			t.Fatal("Swagger UI HTML is missing the SwaggerUIBundle initialization")
		}
	})
}

func assertErrorResponse(t *testing.T, response *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()
	if response.Code != wantStatus {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, wantStatus, response.Body.String())
	}
	var body errorResponse
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if body.Error.Code != wantCode {
		t.Fatalf("error code = %q, want %q", body.Error.Code, wantCode)
	}
}

func float64Pointer(value float64) *float64 { return &value }

func calculationRequestBody(operation string, operandCount int) string {
	operands := strings.TrimSuffix(strings.Repeat("1,", operandCount), ",")
	return `{"operation":"` + operation + `","operands":[` + operands + `]}`
}
