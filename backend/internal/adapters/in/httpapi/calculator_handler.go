// Package httpapi is the inbound HTTP adapter for the calculator use case.
package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"sezzle/backend/internal/application/calculator"
	"sezzle/backend/internal/domain/calculation"
)

const maxRequestBodyBytes = 1 << 20 // 1 MiB

var decimalNumberPattern = regexp.MustCompile(`^-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?$`)

type calculationRequest struct {
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
}

type calculationRequestPayload struct {
	Operation string          `json:"operation"`
	Operands  json.RawMessage `json:"operands"`
}

type requestValidationError struct {
	code    string
	message string
}

func (e *requestValidationError) Error() string {
	return e.message
}

type calculationResponse struct {
	Operation string    `json:"operation"`
	Operands  []float64 `json:"operands"`
	Result    float64   `json:"result"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewCalculatorHandler creates the API router for the calculator endpoint.
func NewCalculatorHandler(calculateUseCase calculator.CalculateUseCase) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealthCheck)
	mux.HandleFunc("/api/v1/calculate", handleCalculateRequest(calculateUseCase))
	mux.HandleFunc("/api/v1/openapi.yaml", serveOpenAPISpecification)
	mux.HandleFunc("/docs", serveSwaggerUI)
	return withRequestID(withRequestLogging(withRecovery(mux)))
}

func handleCalculateRequest(calculateUseCase calculator.CalculateUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			logRequestRejection(r, "method_not_allowed", "", 0)
			writeErrorResponse(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only the POST method is allowed.")
			return
		}
		if !isJSONContentType(r.Header.Get("Content-Type")) {
			logRequestRejection(r, "unsupported_media_type", "", 0)
			writeErrorResponse(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json.")
			return
		}

		request, err := decodeCalculationRequest(w, r)
		if err != nil {
			logRequestRejection(r, requestErrorCode(err), "", 0)
			writeRequestErrorResponse(w, err)
			return
		}
		if strings.TrimSpace(request.Operation) == "" {
			logRequestRejection(r, calculation.ErrorCodeValidation, "", len(request.Operands))
			writeErrorResponse(w, http.StatusBadRequest, calculation.ErrorCodeValidation, "An operation is required.")
			return
		}

		result, err := calculateUseCase.Calculate(request.Operation, request.Operands)
		if err != nil {
			logCalculationError(r, err, request.Operation, len(request.Operands))
			writeCalculationErrorResponse(w, err)
			return
		}

		writeJSONResponse(w, http.StatusOK, calculationResponse{Operation: request.Operation, Operands: request.Operands, Result: result})
	}
}

func isJSONContentType(contentType string) bool {
	mediaType, _, err := mime.ParseMediaType(contentType)
	return err == nil && mediaType == "application/json"
}

func decodeCalculationRequest(w http.ResponseWriter, r *http.Request) (calculationRequest, error) {
	var payload calculationRequestPayload
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		return calculationRequest{}, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return calculationRequest{}, errors.New("request body must contain one JSON object")
	}

	operands, err := decodeOperands(payload.Operands)
	if err != nil {
		return calculationRequest{}, err
	}
	return calculationRequest{Operation: payload.Operation, Operands: operands}, nil
}

func decodeOperands(rawOperands json.RawMessage) ([]float64, error) {
	trimmedOperands := bytes.TrimSpace(rawOperands)
	if len(trimmedOperands) == 0 || bytes.Equal(trimmedOperands, []byte("null")) {
		return nil, newRequestValidationError("invalid_number_format", "Operands must be a JSON array of decimal numbers.")
	}

	var rawValues []json.RawMessage
	if err := json.Unmarshal(trimmedOperands, &rawValues); err != nil {
		return nil, newRequestValidationError("invalid_number_format", "Operands must be a JSON array of decimal numbers.")
	}

	operands := make([]float64, 0, len(rawValues))
	for _, rawValue := range rawValues {
		operand, err := parseOperand(rawValue)
		if err != nil {
			return nil, err
		}
		operands = append(operands, operand)
	}
	return operands, nil
}

func parseOperand(rawValue json.RawMessage) (float64, error) {
	value := string(bytes.TrimSpace(rawValue))
	if len(value) > calculation.MaxOperandInputLength {
		return 0, newRequestValidationError("invalid_number_format", "Each operand exceeds the maximum supported input length.")
	}
	if !decimalNumberPattern.MatchString(value) {
		return 0, newRequestValidationError("invalid_number_format", "Operands must use digits, an optional leading minus sign, and a period decimal separator.")
	}

	if decimalIndex := strings.IndexByte(value, '.'); decimalIndex >= 0 && len(value)-decimalIndex-1 > calculation.MaxDecimalPlaces {
		return 0, newRequestValidationError("invalid_number_format", "Operands can contain at most ten decimal places.")
	}
	if significantDigitCount(value) > calculation.MaxSignificantDigits {
		return 0, newRequestValidationError("invalid_number_format", "Operands can contain at most fifteen significant digits.")
	}

	operand, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, newRequestValidationError("invalid_number_format", "Operands must be valid finite decimal numbers.")
	}
	if err := calculation.ValidateOperandValue(operand); err != nil {
		var calculationError *calculation.CalculationError
		if errors.As(err, &calculationError) {
			return 0, newRequestValidationError(calculationError.Code, calculationError.Message)
		}
		return 0, err
	}
	return operand, nil
}

func significantDigitCount(value string) int {
	digits := strings.NewReplacer("-", "", ".", "").Replace(value)
	digits = strings.TrimLeft(digits, "0")
	return len(digits)
}

func newRequestValidationError(code, message string) *requestValidationError {
	return &requestValidationError{code: code, message: message}
}

func writeCalculationErrorResponse(w http.ResponseWriter, err error) {
	var calculationError *calculation.CalculationError
	if errors.As(err, &calculationError) {
		writeErrorResponse(w, http.StatusBadRequest, calculationError.Code, calculationError.Message)
		return
	}
	writeErrorResponse(w, http.StatusInternalServerError, "internal_error", "An unexpected error occurred.")
}

func writeRequestErrorResponse(w http.ResponseWriter, err error) {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		writeErrorResponse(w, http.StatusRequestEntityTooLarge, "request_too_large", "The request body must not exceed 1 MiB.")
		return
	}

	var validationError *requestValidationError
	if errors.As(err, &validationError) {
		writeErrorResponse(w, http.StatusBadRequest, validationError.code, validationError.message)
		return
	}

	writeErrorResponse(w, http.StatusBadRequest, "invalid_json", "The request body must contain one valid JSON object.")
}

func requestErrorCode(err error) string {
	var maxBytesError *http.MaxBytesError
	if errors.As(err, &maxBytesError) {
		return "request_too_large"
	}
	var validationError *requestValidationError
	if errors.As(err, &validationError) {
		return validationError.code
	}
	return "invalid_json"
}

func logCalculationError(r *http.Request, err error, operation string, operandCount int) {
	var calculationError *calculation.CalculationError
	if errors.As(err, &calculationError) {
		logRequestRejection(r, calculationError.Code, operation, operandCount)
		return
	}
	slog.Error("unexpected calculation failure",
		"request_id", requestIDFromContext(r.Context()),
		"operation", operation,
		"operand_count", operandCount,
		"error", err,
	)
}

func writeErrorResponse(w http.ResponseWriter, status int, code, message string) {
	writeJSONResponse(w, status, errorResponse{Error: apiError{Code: code, Message: message}})
}

func writeJSONResponse(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
