package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

/*
* Cli flags :
*   -d : destination folder for files
*		-i : update interval
 */

var wg sync.WaitGroup

func fetchWrapper(pair CurrencyPair, destination string) {
	if queryResult, err := fetchJson(pair); err != nil {
		log.Print(err)
	} else {
		writeCsv(destination, pair, queryResult)
	}
	wg.Done()
}

func forever(destination string, pairs []CurrencyPair, interval time.Duration) {
//implement time.Ticker here
	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
				case <-ticker.C:
					for _, pair := range pairs {
						wg.Add(1)
						go fetchWrapper(pair, destination)
					}
				case <-quit:
					ticker.Stop()
					return
			}
		}
	}()
}

func fetchJson(pair CurrencyPair) (QueryResult, error) {
	url := "https://api.cryptonator.com/api/ticker/" + pair.Base + "-" + pair.Target + "/"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var queryResult QueryResult
	if err := json.Unmarshal(body, &queryResult); err != nil {
		log.Fatal(err)
	}
	if !queryResult.Success {
		return queryResult, PairNotFoundError{C: pair, Message: queryResult.Error}
	}
	return queryResult, nil
}

func writeCsv(destination string, pair CurrencyPair, query QueryResult) {
	filename := time.Now().Format("20060102") + ".csv"
	path := filepath.Join(destination, pair.Base+"_"+pair.Target)
	if err := os.MkdirAll(path, 0744); err != nil {
		log.Fatal(err)
	}
	csv_file, err := os.OpenFile(filepath.Join(path, filename), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer csv_file.Close()

	reader := csv.NewReader(csv_file)
	writer := csv.NewWriter(csv_file)

	// If file is empty, write the csv headers
	record, err := reader.Read()
	if record == nil && err == io.EOF {
		if err2 := writer.Write([]string{"timestamp", "price", "volume", "change"}); err2 != nil {
			log.Fatal("writeCsv, unable to write CSV headers", err2)
		}
		if err := writer.Write(query.ToArray()); err != nil {
			log.Fatal(err)
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	} else {
		records, err := reader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		timestamp, err := strconv.Atoi(records[len(records)-1][0])
		if err != nil {
			log.Fatal(err)
		}
		if timestamp < query.Timestamp {
			if err := writer.Write(query.ToArray()); err != nil {
				log.Fatal(err)
			}
			writer.Flush()
			if err := writer.Error(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {

	default_destination := filepath.Join(os.Getenv("HOME"), "cryptotrend")
	var destination string
	usage_destination := "The folder in which to save the output files."
	flag.StringVar(&destination, "d", default_destination, usage_destination)

	default_interval := "30m"
	var intervalStr string
	usage_interval := "The interval of time before fetching data again. Formats such as 1h45m or 30s are valid. Accepted units are \"s\", \"m\" and \"h\""
	flag.StringVar(&intervalStr, "i", default_interval, usage_interval)
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	var pairs []CurrencyPair
	for _, arg := range flag.Args() {
		if newPair, err := NewCurrencyPair(arg); err != nil {
			log.Print(err)
		} else {
			pairs = append(pairs, newPair)
		}
	}
	if len(pairs) == 0 {
		log.Fatal("No valid currency pairs were provided. Exiting program")
	}

	wg.Add(1)
	go fetchWrapper(pairs[0], destination)
	go forever(destination, pairs, interval)
	wg.Wait()
}
