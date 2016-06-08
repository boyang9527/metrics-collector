package server

import (
	"crypto/tls"
	"encoding/json"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/metrics"
	"github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/security"
	. "github.com/cloudfoundry-incubator/app-autoscaler/metrics-collector/util"
	"github.com/cloudfoundry/noaa/consumer"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	noaa *consumer.Consumer
}

func NewHandler(doppler string) *Handler {
	Logger.Info("create-noaa-client", map[string]interface{}{"doppler": doppler})

	noaaClient := consumer.New(doppler, &tls.Config{InsecureSkipVerify: true}, nil)
	return &Handler{
		noaa: noaaClient,
	}
}

func (h *Handler) GetMemoryMetric(w http.ResponseWriter, r *http.Request) {
	Logger.Debug("request-to-get-memory-metric", map[string]interface{}{"Request": DumpRequest(r)})

	appId := mux.Vars(r)["appid"]
	containerMetrics, err := h.noaa.ContainerMetrics(appId, "bearer "+security.GetOAuthToken())

	Logger.Debug("get-container-metrics-from-doppler", map[string]interface{}{"container metrics": containerMetrics})

	if err == nil {
		metric := metrics.GetMemoryMetricFromContainerMetrics(appId, containerMetrics)
		body, err := json.Marshal(metric)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(body)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write(CreateJsonErrorResponse("Error-Get-Metrics-From-Doppler", err.Error()))

	Logger.Error("failed-to-get-memory-metric", err)
}
