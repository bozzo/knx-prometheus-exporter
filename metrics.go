/*
 *
 *    Copyright 2020 Boris Barnier <bozzo@users.noreply.github.com>
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var metrics = map[string]*prometheus.GaugeVec{}

func getMetric(dptName string, dptID string) *prometheus.GaugeVec {
	if metric := metrics[dptName]; metric != nil {
		return metric
	}
	metrics[dptName] = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "knx_" + strcase.ToSnake(dptName),
		Help: "The " + dptName + " represented by DPT " + dptID,
	}, []string{"source", "destination"})
	return metrics[dptName]
}

func runMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PROMETHEUS_PORT")
	if port == "" {
		port = "13671"
	}

	logrus.Infof("Metrics server listening on port %s", port)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
