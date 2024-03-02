package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jackcoble/puregym-go"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Load environment variables
	godotenv.Load()

	// Create instance of the PureGym client
	pureGym, err := puregym.NewClient(os.Getenv("PUREGYM_EMAIL"), os.Getenv("PUREGYM_PIN"))
	if err != nil {
		log.Fatalf("unable to create puregym client: %s", err.Error())
		return
	}

	// Authenticate and set the Home Gym the client requires for fetching data
	if err := pureGym.Authenticate(); err != nil {
		log.Fatalf("unable to authenticate with puregym: %s", err.Error())
		return
	}
	pureGym.SetHomeGym()

	// Fetch current total amount of people in the gym
	gymInfo, err := pureGym.GetGymAttendance()
	if err != nil {
		log.Fatalf(err.Error())
		return
	}

	log.Println("amount of people in gym:", gymInfo.TotalPeopleInGym)

	// Start a HTTP server for Metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+os.Getenv("PROMETHEUS_METRICS_PORT"), nil)
}
