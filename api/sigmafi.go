package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/binance-exchange/go-binance"
	"github.com/go-kit/log"
	"github.com/julienschmidt/httprouter"
)

var (
	sigmaRouter *httprouter.Router
	logger      log.Logger
	b           binance.Binance
)

func init() {
	sigmaRouter = httprouter.New()
	registerSigmaHttpRouter(sigmaRouter)
	// init binance service
	// var logger log.Logger
	apiKey := os.Getenv("BINANCE_API_KEY")
	if apiKey == "" {
		panic("api key is empty")
	}
	apiSecret := os.Getenv("BINANCE_API_SECRET")
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(apiSecret),
	}
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	ctx := context.Background()
	// use second return value for cancelling request when shutting down the app

	binanceService := binance.NewAPIService(
		"https://testnet.binancefuture.com",
		apiKey,
		hmacSigner,
		logger,
		ctx,
	)
	b = binance.NewBinance(binanceService)
}

func registerSigmaHttpRouter(r *httprouter.Router) {
	// r.NotFound = http.HandlerFunc(NotFoundHandler)
	r.GET("/sigmafi/order/:orderId", getOrderById)
}

func SigmafiHandler(w http.ResponseWriter, r *http.Request) {
	// preprocess
	sigmaRouter.ServeHTTP(w, r)
	// postprocess
}

func getOrderById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	orderIdStr := p.ByName("orderId")
	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("orderId is invalid"))
		return
	}
	po, err := b.QueryOrder(binance.QueryOrderRequest{
		Symbol:    "BTCUSDT",
		OrderID:   orderId,
		Timestamp: time.Now(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s", err)))
		logger.Log(err)
		return
	}
	bs, _ := json.Marshal(po)
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
