package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shokunin/contrib/ginrus"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/kong/consul-vault-exporter/handlers/healthcheck"
	"github.com/kong/consul-vault-exporter/handlers/metrics"
)

func main() {
	router := gin.New()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	router.Use(ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true, "consul-vault-exporter"))

	// Start routes
	router.GET("/health", healthcheck.HealthCheck)
	router.GET("/", metrics.Redirect)
	router.GET("/metrics", metrics.Metrics)

	// RUN rabit run
	router.Run(":9402") // listen and serve on 0.0.0.0:9402
}
