package api

import (
	"git.tor.ph/hiveon/pool/config"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	 kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"context"
	"encoding/json"
)

type PoolHandler struct {
	config      *config.PoolConfig

}

func New(config *config.PoolConfig) *PoolHandler {
	return &PoolHandler{config}
}

func (h *PoolHandler) log() *logrus.Logger {
	return h.config.Log
}

func (h *PoolHandler) Bind(r *gin.Engine) {
	h.MakePoolHandlers(r)
}

func (h *PoolHandler) MakePoolHandlers(r *gin.Engine) {

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	router := mux.NewRouter()

	poolGetIndexHandler := kithttp.NewServer(
		poolGetIndexEndpoint(),
		decodeEmptyRequest,
		encodeResponse, opts...,
	)
	router.Handle("/api/pool/index", poolGetIndexHandler).Methods("GET")

	poolGetIncomeHistoryHandler := kithttp.NewServer(
		poolGetIncomeHistoryEndpoint(),
		decodeEmptyRequest,
		encodeResponse, opts...,
	)
	router.Handle("/api/pool/incomeHistory", poolGetIncomeHistoryHandler).Methods("GET")

	r.Use(gin.WrapH(router))
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return r, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
	//log.Info(response)
	return json.NewEncoder(w).Encode(response)
}

