package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/3cb/ssc"
	"github.com/boltdb/bolt"
)

func startPolling(db *bolt.DB, pool *ssc.SocketPool) {
	ticker := time.NewTicker(time.Minute * 5)
	loc, _ := time.LoadLocation("America/Winnipeg")
	date := strings.Split(fmt.Sprint(time.Now().In(loc)), " ")[0]

	readings, err := getReadings()
	if err != nil {
		log.Fatalf("Unable to get initial air quality readings: %v", err)
	}

	msg := serialize(readings, date)
	err = updateDB(db, msg, date)
	if err != nil {
		log.Fatalf("Error saving air quality readings to database: %v", err)
	}
	pool.Pipes.InboundBytes <- ssc.Data{Type: 2, Payload: msg}

	for {
		<-ticker.C
		date := strings.Split(fmt.Sprint(time.Now().In(loc)), " ")[0]

		readings, err := getReadings()
		if err != nil {
			log.Printf("Unable to update air quality readings: %v", err)
		}

		msg := serialize(readings, date)
		err = updateDB(db, msg, date)
		if err != nil {
			log.Fatalf("Error saving air quality readings to database: %v", err)
		}
		pool.Pipes.InboundBytes <- ssc.Data{Type: 2, Payload: msg}
	}
}

func getReadings() ([]Reading, error) {
	air := []Reading{}

	t1, t2 := getDate()
	api := "https://data.winnipeg.ca/resource/f5p2-7r36.json"
	queryString := "?where=observationtime between '" + t1 + "' and '" + t2 + "'"
	resp, err := http.Get(api + queryString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &air)
	if err != nil {
		return nil, err
	}
	return air, nil
}

func getDate() (string, string) {
	loc, _ := time.LoadLocation("America/Winnipeg")
	t1 := strings.Split(fmt.Sprint(time.Now().In(loc)), " ")[0] + "T00:00:00"
	t2 := strings.Split(fmt.Sprint(time.Now().In(loc).Add(time.Hour*24)), " ")[0] + "T00:00:00"
	return t1, t2
}
