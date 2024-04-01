package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Metrics struct {
	Name      string `json:"name"`
	Threshold string `json:"threshold"`
	Value     string `json:"value"`
}

type Result struct {
	ResultID int    `json:"result_id"`
	SiteName string `json:"site_name"`
	Metrics  any    `json:"metrics"`
}
type ResultMap struct {
	ResultID int                `json:"result_id"`
	SiteName string             `json:"site_name"`
	Metrics  map[string]Metrics `json:"metrics"`
}

func main() {
	http.HandleFunc("POST /", processRequest)
	http.ListenAndServe(":8080", nil)
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading request body"))
		return
	}

	formattedOutput := parseAndFormat(string(body))

	sendNotification(formattedOutput)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification sent successfully"))
}

func parseAndFormat(jsonData string) string {
	var result Result
	err := json.Unmarshal([]byte(jsonData), &result)
	if err != nil {
		return "Error parsing JSON: " + err.Error()
	}

	output := fmt.Sprintf("ℹ️ Result ID: %d\n🌐 Site Name: %s\n📊 Metrics:\n", result.ResultID, result.SiteName)

	var metricsParsed []Metrics
	// Access metrics based on type
	if metrics, ok := result.Metrics.([]interface{}); ok {
		metricsParsed = make([]Metrics, 0, len(metrics))
		for _, metric := range metrics {
			metricMap := metric.(map[string]interface{})
			metricsParsed = append(
				metricsParsed, Metrics{
					Name:      metricMap["name"].(string),
					Threshold: metricMap["threshold"].(string),
					Value:     metricMap["value"].(string),
				},
			)
		}
	} else {
		resultMap := ResultMap{}
		err := json.Unmarshal([]byte(jsonData), &resultMap)
		if err != nil {
			return "Error parsing JSON: " + err.Error()
		}
		metricsParsed = make([]Metrics, 0, len(resultMap.Metrics))
		for _, metric := range resultMap.Metrics {
			metricsParsed = append(metricsParsed, metric)
		}
	}

	for key, metric := range metricsParsed {
		output += fmt.Sprintf("   - %d:\n", key)
		output += fmt.Sprintf("     - Name: %s\n", metric.Name)
		output += fmt.Sprintf("       - Threshold: %s\n", metric.Threshold)
		output += fmt.Sprintf("       - Value: %s\n", metric.Value)
	}

	return output
}

func sendNotification(message string) {
	req, _ := http.NewRequest(
		"POST", os.Getenv("NOTIFICATION_URL"),
		strings.NewReader(message),
	)
	req.Header.Set("Title", "Speedtest Results Notification")
	req.Header.Set("Tags", "warning,skull")
	do, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending notification:", err)
		return
	}
	fmt.Println("Notification sent successfully:", do.Status)
}
