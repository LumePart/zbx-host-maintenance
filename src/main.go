package main

import (
	"log"
	"os"
	"time"
	"regexp"
	"strconv"
)

type Config struct {
	URL string `json:"URL"`
	APIKey string `json:"APIKey"`
}

func buildConf(URL, APIKey string) Config { // add params to config
	return Config{URL: URL, APIKey: APIKey}
}

func parsePeriod(currentTime time.Time, lengthPeriod string) (time.Time, error) { // Get maintenance period (epoch) from string
	maintenanceLength := currentTime
	re := regexp.MustCompile(`(\d+)([smhdw])`)
    matches := re.FindAllStringSubmatch(lengthPeriod, -1)

    for _, match := range matches {
        period, err := strconv.Atoi(match[1])
		if err != nil {
			return time.Time{}, err
		}
        unit := match[2]

		switch unit {
		case "s":
			maintenanceLength = maintenanceLength.Add(time.Second * time.Duration(period))
		case "m":
			maintenanceLength = maintenanceLength.Add(time.Minute * time.Duration(period))
		case "h":
			maintenanceLength = maintenanceLength.Add(time.Hour * time.Duration(period))
		case "d":
			maintenanceLength = maintenanceLength.Add(time.Hour * 24 * time.Duration(period))
		case "w":
			maintenanceLength = maintenanceLength.Add(time.Hour * 24 * 7 * time.Duration(period))
		}
    }
	return maintenanceLength, nil
}

func main() {
	URL := os.Args[1]
	APIKey := os.Args[2]
	hostname := os.Args[3]
	length := os.Args[4]
	cfg := buildConf(URL, APIKey)
	currentTime := time.Now()
	ID, err := getHostID(cfg, hostname)
	if err != nil {
		log.Println("failed to find hostID: %w", err)
		os.Exit(1)
	}
	maintenanceTime, err := parsePeriod(currentTime, length)
	err = scheduleMaintenance(cfg, currentTime.Unix(), maintenanceTime.Unix(), ID, hostname)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("Maintenance scheduled for %s", hostname)
}