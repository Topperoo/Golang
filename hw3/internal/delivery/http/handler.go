package http

import (
	"encoding/json"
	"homework3/internal/domain"
	"homework3/internal/usecase"
	"net/http"
)

type BalanceHandler struct {
	useCase *usecase.BalanceUseCase
}

func NewBalanceHandler(useCase *usecase.BalanceUseCase) *BalanceHandler {
	return &BalanceHandler{useCase: useCase}
}

type CreditRequest struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type TransferRequest struct {
	FromUserID string  `json:"from_user_id"`
	ToUserID   string  `json:"to_user_id"`
	Amount     float64 `json:"amount"`
}

type BalanceResponse struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *BalanceHandler) Credit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		h.sendError(w, "user_id is required", http.StatusBadRequest)
		return
	}

	if err := h.useCase.CreditBalance(req.UserID, req.Amount); err != nil {
		h.handleUseCaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *BalanceHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.FromUserID == "" || req.ToUserID == "" {
		h.sendError(w, "from_user_id and to_user_id are required", http.StatusBadRequest)
		return
	}

	if err := h.useCase.TransferBalance(req.FromUserID, req.ToUserID, req.Amount); err != nil {
		h.handleUseCaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *BalanceHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		h.sendError(w, "user_id query parameter is required", http.StatusBadRequest)
		return
	}

	balance, err := h.useCase.GetBalance(userID)
	if err != nil {
		h.handleUseCaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(BalanceResponse{
		UserID:  userID,
		Balance: balance,
	})
}

func (h *BalanceHandler) handleUseCaseError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrInsufficientBalance:
		h.sendError(w, err.Error(), http.StatusBadRequest)
	case domain.ErrNegativeAmount:
		h.sendError(w, err.Error(), http.StatusBadRequest)
	case domain.ErrSelfTransfer:
		h.sendError(w, err.Error(), http.StatusBadRequest)
	case domain.ErrUserNotFound:
		h.sendError(w, err.Error(), http.StatusNotFound)
	default:
		h.sendError(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *BalanceHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
