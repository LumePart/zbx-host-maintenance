package main

import (
	"net/http"
	"bytes"
	"fmt"
	"io"
	"encoding/json"
)

type Hosts struct {
	Results []Result `json:"result"`
	Error   Error  `json:"error"`
}

type Result struct {
	Hostid string `json:"hostid"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func getHostID(cfg Config, hostname string) (string, error) {
	payload := fmt.Appendf(nil, `{
		"jsonrpc": "2.0",
		"method": "host.get",
		"params": {
		"filter": {
			"host": ["%s"]
		},
		"output": ["hostid"]
	 },
		"id": 1 }`, hostname)
	body, err := makeZbxRequest(cfg, payload)
	if err != nil {
		return "", err
	}
	var host Hosts
	if err = json.Unmarshal(body, &host); err != nil {
		return "", err
	}
	
	return host.Results[0].Hostid, nil
}

func scheduleMaintenance(cfg Config, currentTime, maintenanceTime int64, hostID, hostname string) error {

	period := maintenanceTime - currentTime
	payload := fmt.Appendf(nil, `{     
		"jsonrpc": "2.0",
		"method": "maintenance.create",
		"params": {
			"name": "Maintenance for %s",
			"active_since": %v,
			"active_till": %v,
			"description": "scheduled from manual action",
			"hosts": [{"hostid": "%s"}],
			"timeperiods": [
					{
					"period": %v,
					"start_date": %v,
					"timeperiod_type": 0
					}
				]
				},
				"id": 1
			}`, hostname, currentTime, maintenanceTime, hostID, period, currentTime)

		body, err := makeZbxRequest(cfg, payload)
		if err != nil {
			return fmt.Errorf("failed to create maintenance: %s", err.Error())
		}
		var res Hosts
		if err = json.Unmarshal(body, &res); err != nil {
			return err
		}
		
		if res.Error.Data != "" {
			return fmt.Errorf("failed to create maintenance: %s", res.Error.Data)
		}
		return nil

}

func makeZbxRequest(cfg Config, payload []byte) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", cfg.URL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+cfg.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}