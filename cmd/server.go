/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/BaseMax/real-time-notifications-nats-go/controllers"
	"github.com/BaseMax/real-time-notifications-nats-go/database"
	"github.com/BaseMax/real-time-notifications-nats-go/models"
	"github.com/BaseMax/real-time-notifications-nats-go/notifications"
	"github.com/BaseMax/real-time-notifications-nats-go/rabbitmq"
	"github.com/BaseMax/real-time-notifications-nats-go/routes"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Serve web application",
	Run: func(cmd *cobra.Command, args []string) {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatal("godotenv: ", err)
		}

		c, err := database.ReadConfig()
		if err != nil {
			log.Fatal("readconfig: ", err)
		}

		conn, err := database.OpenPostgres(c)
		if err != nil {
			log.Fatal("open postgres: ", err)
		}
		if err := models.Init(conn); err != nil {
			log.Fatal("models init: ", err)
		}

		if err := notifications.InitNats(); err != nil {
			log.Fatal("nats init: ", err)
		}

		if err := rabbitmq.Connect(); err != nil {
			log.Fatal("rabbitmq init: ", err)
		}

		if err := notifications.InitKafka(); err != nil {
			log.Fatal("kafka init: ", err)
		}

		r := routes.InitRoutes()
		r.Logger.Fatal(r.Start(controllers.GetRunningAddr()))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
