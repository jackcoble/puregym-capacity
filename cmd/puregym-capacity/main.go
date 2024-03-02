package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackcoble/puregym-go"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus metrics
var (
	capacityCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "capacity_count",
		Help: "The total number of people inside the gym",
	})
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

	// Fetch the capacity within the gym every minute
	go func() {
		for {
			gymInfo, err := pureGym.GetGymAttendance()
			if err != nil {
				log.Fatalf(err.Error())
				return
			}

			capacityCount.Set(float64(gymInfo.TotalPeopleInGym))

			// Do it again in 30 seconds
			time.Sleep(30 * time.Second)
		}
	}()

	// Start a HTTP server for Metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":"+os.Getenv("PROMETHEUS_METRICS_PORT"), nil)
}
