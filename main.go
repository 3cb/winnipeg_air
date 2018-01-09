package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/3cb/ssc"

	"github.com/boltdb/bolt"
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = createBucket(db, "Readings")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// start websocket pool
	config := ssc.PoolConfig{
		IsReadable: true,
		IsWritable: true,
		IsJSON:     false,
	}
	pool, err := ssc.NewSocketPool(config)
	if err != nil {
		log.Fatal(err)
	}

	go startPolling(db, pool)

	// routes
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/dist").Handler(http.FileServer(http.Dir("./static/")))

	// start server
	log.Fatal(http.ListenAndServe(":4040", r))
}
