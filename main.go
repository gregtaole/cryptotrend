package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

/*
* Cli flags :
*   -d : destination folder for files
 */



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
	fmt.Printf("%s\n", body)
	var queryResult QueryResult
	if err := json.Unmarshal(body, &queryResult); err != nil {
		log.Fatal(err)
	}
	if !queryResult.Success {
		return queryResult, PairNotFoundError{C: pair}
	}
	return queryResult, nil
}

func writeCsv(pair CurrencyPair, query QueryResult) {
	csv_file, err := os.OpenFile(pair.Base+pair.Target+".csv", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
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
	} else if err != nil {
		log.Fatal(err)
	}

	if err := writer.Write(query.ToArray()); err != nil {
		log.Fatal(err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	default_destination := filepath.Join(os.Getenv("HOME"), "cryptotrend")
	var destination string
	usage_destination := "The folder in which to save the output files."
	flag.StringVar(&destination, "d", default_destination, usage_destination)

	flag.Parse()

	var pairs []CurrencyPair
	for _, arg := range flag.Args() {
		if newPair, err := NewCurrencyPair(arg); err != nil {
			log.Print(err)
		} else {
			pairs = append(pairs, newPair)
		}
	}

	for _, pair := range pairs {
		if queryResult, err := fetchJson(pair); err != nil {
			log.Print(err)
		} else {
			fmt.Println(queryResult)
			writeCsv(pair, queryResult)
		}
	}
}
