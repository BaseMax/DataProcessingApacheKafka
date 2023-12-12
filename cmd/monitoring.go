/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"

	"github.com/BaseMax/real-time-notifications-nats-go/models"
	"github.com/BaseMax/real-time-notifications-nats-go/notifications"
)

var gauges = make(map[string]prometheus.Gauge, 0)

func prepareMonitoring() {
	prom := prometheus.NewRegistry()

	for _, action := range models.ACTIONS {
		gauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "consumer",
			Name:      action,
		})
		gauges[action] = gauge
		prom.MustRegister(gauge)
	}

	go func() {
		for {
			for _, gauge := range gauges {
				time.Sleep(time.Minute)
				gauge.Set(0)
			}
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))
		http.ListenAndServe(":8080", nil)
	}()
}

func StreamAndMonitoring(cmd *cobra.Command) {
	conf, err := notifications.GetKafkaConfig()
	if err != nil {
		log.Fatal(err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{conf.Server},
		Topic:     conf.Topic,
		Partition: conf.Partition,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	r.SetOffset(0)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}

		var activity map[string]any
		err = json.Unmarshal(m.Value, &activity)
		if err != nil {
			log.Println(err)
			continue
		}

		for action := range gauges {
			if action == activity["action"] {
				gauges[action].Inc()
			}
		}
	}

}

// monitoringCmd represents the monitoring command
var monitoringCmd = &cobra.Command{
	Use:   "monitoring",
	Short: "Recieve data from kafka and process them",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("godotenv: ", err)
		}

		prepareMonitoring()
		StreamAndMonitoring(cmd)
	},
}

func init() {
	rootCmd.AddCommand(monitoringCmd)
}
