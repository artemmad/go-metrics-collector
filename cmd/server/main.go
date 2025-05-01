package main

import (
	"net/http"
	"strconv"
	"strings"
)

var (
	gauge   = map[string]float64{}
	counter = map[string]int64{}
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /update/", metricCalc)
	mux.HandleFunc("GET /", metricList)
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "404 page not found", http.StatusNotFound)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func metricList(writer http.ResponseWriter, request *http.Request) {
	response := ""
	response += strings.ToUpper(GAUGE) + ":\n"
	for key, value := range gauge {
		response += "\t" + key + ": " + strconv.FormatFloat(value, 'f', -1, 64) + "\n"
	}
	response += strings.ToUpper(COUNTER) + ":\n"
	for key, value := range counter {
		response += "\t" + key + ": " + strconv.FormatInt(value, 10) + "\n"
	}
	response += "\n"
	writer.Write([]byte(response))
}

func metricCalc(writer http.ResponseWriter, request *http.Request) {

	//http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>, Content-Type: text/plain.
	parts := strings.Split(strings.Trim(request.URL.Path, "/"), "/")

	// валидировать что как минимум имеем 4 секции с / и оно начинается с update
	if len(parts) < 4 || parts[0] != "update" {
		http.Error(writer, "metric calc requires at least 4 parts", http.StatusNotFound)
		return
	}

	// разобраться реквест чтобы понять ТИП_МЕТРИКИ, если не существует такого типа метрики - 400
	metricType := strings.ToLower(strings.TrimSpace(parts[1]))
	if metricType != GAUGE && metricType != COUNTER {
		http.Error(writer, "metric type is unknown", http.StatusBadRequest)
		return
	}

	// разобрать реквест и вытащить ИМЯ_МЕТРИКИ, если нет - вернуть 404
	metricName := strings.TrimSpace(parts[2])
	if metricName == "" {
		http.Error(writer, "metric name is required", http.StatusNotFound)
		return
	}

	// разобрать реквест и вытащить ЗНАЧЕНИЕ_МЕТРИКИ
	metricValueStr := strings.TrimSpace(parts[3])
	switch metricType {
	case GAUGE:
		gaugeVal, err := strconv.ParseFloat(metricValueStr, 64)
		if err != nil {
			http.Error(writer, "metric value is invalid for type GAUGE", http.StatusBadRequest)
			return
		}
		gauge[metricName] = gaugeVal
	case COUNTER:
		counterVal, err := strconv.ParseInt(metricValueStr, 10, 64)
		if err != nil {
			http.Error(writer, "metric value is invalid for type COUNTER", http.StatusBadRequest)
			return
		}
		counter[metricName] = counterVal
	default:
		http.Error(writer, "unsupported metric type", http.StatusBadRequest)
		return
	}

	// вернуть 200ОК
	writer.WriteHeader(http.StatusOK)
}
