package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/renanqts/xdpdropper/pkg/config"
	"github.com/renanqts/xdpdropper/pkg/logger"
	"github.com/renanqts/xdpdropper/pkg/xdp"
	"go.uber.org/zap"
)

type API interface {
	Start() error
	Close()
}

type api struct {
	wg         sync.WaitGroup
	httpServer *http.Server
	router     *mux.Router
	xdp        xdp.XDP
}

type drop struct {
	IP string `json:"ip"`
}

func New(config config.Config) (API, error) {
	xdp, err := xdp.New(config.Iface)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter().StrictSlash(true)
	server := http.Server{
		Addr:              config.Address,
		Handler:           router,
		ReadHeaderTimeout: time.Second * 1,
	}

	return api{
		httpServer: &server,
		router:     router,
		xdp:        xdp,
	}, nil
}

func reqUnmarshal(w http.ResponseWriter, r *http.Request) (d drop, err error) {
	logger.Log.Debug(
		"Request received",
		zap.String("method", r.Method),
		zap.String("host", r.Host),
		zap.String("uri", r.RequestURI),
	)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error("Error while unmarshaling", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return d, err
	}

	err = json.Unmarshal(reqBody, &d)
	if err != nil {
		logger.Log.Error("Error while unmarshaling", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return d, err
	}

	logger.Log.Debug("Request body", zap.Any("payload", d))

	return d, err
}

func (a api) Close() {
	if a.httpServer != nil {
		if err := a.httpServer.Close(); err != nil {
			logger.Log.Error("Failed to stop the HTTP server", zap.Error(err))
		}
	}
	a.xdp.Close()
	a.wg.Wait()
}

func (a api) Start() error {
	a.wg.Add(1)

	a.router.HandleFunc("/health", a.health).Methods("GET")
	a.router.HandleFunc("/add", a.add).Methods("POST")
	a.router.HandleFunc("/remove", a.remove).Methods("POST")

	logger.Log.Info("Starting http server", zap.String("address", a.httpServer.Addr))
	go func() {
		defer a.wg.Done()
		err := a.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Failed to start the HTTP server", zap.Error(err))
			return
		}
	}()

	return nil
}

func (a api) health(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	}
}

func (a api) add(w http.ResponseWriter, r *http.Request) {
	req, err := reqUnmarshal(w, r)
	if err == nil {
		logger.Log.Debug("api drop invocation", zap.String("ip", req.IP))
		err = a.xdp.AddToDrop(req.IP)
		if err == nil {
			w.WriteHeader(http.StatusCreated)
			logger.Log.Debug("add request", zap.Int("statusCode", http.StatusCreated))
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	logger.Log.Error(
		"add request",
		zap.Int("statusCode", http.StatusInternalServerError),
		zap.Error(err),
	)
}

func (a api) remove(w http.ResponseWriter, r *http.Request) {
	req, err := reqUnmarshal(w, r)
	if err == nil {
		logger.Log.Debug("remove drop invocation", zap.String("ip", req.IP))
		err = a.xdp.RemoveFromDrop(req.IP)
		if err == nil {
			w.WriteHeader(http.StatusNoContent)
			logger.Log.Debug("add request", zap.Int("statusCode", http.StatusNoContent))
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	logger.Log.Error(
		"remove request",
		zap.Int("statusCode", http.StatusInternalServerError),
		zap.Error(err),
	)
}
