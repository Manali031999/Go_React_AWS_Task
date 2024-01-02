package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"new_be/db"
	"new_be/models"
	"time"

	"github.com/gorilla/mux"
)

func GetAllInstances(w http.ResponseWriter, r *http.Request) {
	var instances []models.Instance
	db, err := db.GetDBInstance()
	if err != nil {
		log.Fatalf("Error getting database instance: %v", err)
	}
	result := db.Find(&instances)
	if result.Error != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch instances")
		return
	}
	RespondWithJSON(w, http.StatusOK, instances)
}

func GetMetricData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		RespondWithError(w, http.StatusBadRequest, "Instance ID not provided in the URL")
		return
	}

	db, err := db.GetDBInstance()
	if err != nil {
		log.Fatalf("Error getting database instance: %v", err)
	}
	var metricData models.MetricData
	result := db.First(&metricData, "instance_id = ?", id)
	if result.Error != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to fetch metric data")
		return
	}

	var loadedGraphData map[time.Time]float64
	err = json.Unmarshal(metricData.GraphData, &loadedGraphData)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to Unmarshal Graph data")

	}
	metricDataRes := map[string]interface{}{
		"InstanceID": metricData.InstanceID,
		"CPU":        metricData.CPU,
		"GraphData":  loadedGraphData,
	}

	RespondWithJSON(w, http.StatusOK, metricDataRes)

}
