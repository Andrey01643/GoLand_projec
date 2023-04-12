package handlers

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"web_service/models"
)

func GetItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idsParam := vars["id"]
	ids, err := strconv.Atoi(idsParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := os.Open("ueba.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := make([]models.Record, 0)

	for _, record := range records[1:] {
		id, err := strconv.Atoi(record[1])
		if err != nil {
			log.Printf("failed to parse id: %v", err)
			continue
		}

		if id == ids {
			result = append(result, models.Record{
				Id:         id,
				Uid:        record[2],
				Domain:     record[3],
				Cn:         record[4],
				Department: record[5],
				Title:      record[6],
				Who:        record[7],
			})
		}
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}
