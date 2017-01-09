// Trademark Kian Faroughi 2017

package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"yql"
)

// Sample : stmt, err := db.Query("select * from yahoo.finance.quote where symbol in (\"YHOO\",\"AAPL\",\"GOOG\",\"MSFT\")")

type stock struct {
	averageDailyVolume   string
	change               string
	daysHigh             string
	daysLow              string
	daysRange            string
	currentPrice         string
	marketCapitalization string
	name                 string
	stockExchange        string
	symbol               string
	volume               string
	yearHigh             string
	yearLow              string
	date                 string
	close                string
}

var stockKeys = stock{
	averageDailyVolume:   "AverageDailyVolume",
	change:               "Change",
	daysHigh:             "DaysHigh",
	daysLow:              "DaysLow",
	daysRange:            "DaysRange",
	currentPrice:         "LastTradePriceOnly",
	marketCapitalization: "MarketCapitalization",
	name:                 "Name",
	stockExchange:        "StockExchange",
	symbol:               "Symbol",
	volume:               "Volume",
	yearHigh:             "YearHigh",
	yearLow:              "YearLow",
	date:                 "Date",
	close:                "Close",
}

var maxRetryCount = 3

func main() {
	initStocks()
}

func initStocks() {
	fmt.Println("\n-----------------------------------------------")
	fmt.Println("Kian Faroughi's Stock Program....Powered by YQL")
	fmt.Println("-----------------------------------------------")

	yql, err := initYQL()
	if err != nil {
		return
	}

	for {
		fmt.Println("\n------------------Main Menu--------------------")
		fmt.Println("Modes:")
		fmt.Println("Single Stock Current Price --> 1")
		fmt.Println("Single Stock Historical Data --> 2")
		fmt.Println("Single Stock Historical Analytics --> 3")
		fmt.Println("-----------------------------------------------")

		fmt.Print("\nEnter Mode Number (\"exit\" to quit) --> ")
		reader := bufio.NewReader(os.Stdin)
		mode, _ := reader.ReadString('\n')
		mode = strings.Replace(mode, "\n", "", -1)

		switch mode {
		case "1":
			currentSingleStockMode(yql)

		case "2":
			historicalSingleStockMode(yql)

		case "3":
			historicalAnalytics(yql)

		case "exit":
			fmt.Println("")
			return
		default:
			fmt.Println("\nOnly three supported modes. Restart the program.\n")
		}
	}
}

func initYQL() (*sql.DB, error) {
	var yqlD yql.YQLDriver
	yqlD.Init()
	db, err := sql.Open("yql", "")
	if err != nil {
		fmt.Println("Failure Opening YQL.")
		return nil, err
	}
	return db, nil
}

func currentSingleStockMode(db *sql.DB) {

	currentSingleStockQuery := "select * from yahoo.finance.quote where symbol in (\"KianFaroughi\",\""
	tail := "\")"

	fmt.Println("\n\n----------Current Single Stock Mode.-----------\n")

	reader := bufio.NewReader(os.Stdin)

	var symbol string
	retryCount := maxRetryCount

	for {

		fmt.Print("ENTER STOCK SYMBOL (\"exit\" to return to Main Menu) --> ")
		symbol, _ = reader.ReadString('\n')
		symbol = strings.Replace(symbol, "\n", "", -1)
		fmt.Println("\n")
		if symbol == "exit" {
			break
		}

		request := currentSingleStockQuery + symbol + tail

		for {
			stmt, err := db.Query(request)
			if err != nil {
				if retryCount == 0 {
					fmt.Println("Query Error. Could not retrieve information for ", symbol, ".\n")
					retryCount = maxRetryCount
					break
				} else {
					retryCount--
				}
			} else {

				for stmt.Next() {
					var data map[string]interface{}
					stmt.Scan(&data)
					name := data[stockKeys.name]
					if name != nil {
						fmt.Println("==============================================")
						fmt.Println(name)
						fmt.Println(data[stockKeys.symbol])
						fmt.Println("==============================================")
						fmt.Println("Current Price: ", data[stockKeys.currentPrice])
						fmt.Println("Last Day's Change: ", data[stockKeys.change])
						fmt.Println("Last Day's Range: ", data[stockKeys.daysRange])
						fmt.Println("Year High: ", data[stockKeys.yearHigh])
						fmt.Println("Year Low: ", data[stockKeys.yearLow])
						fmt.Println("==============================================\n")
					}
				}

				retryCount = maxRetryCount
				break
			}
		}
	}
}

func historicalSingleStockMode(db *sql.DB) {

	historicalSingleStockQuery := "select * from yahoo.finance.historicaldata where symbol = \""
	tail := "\" and startDate = \""
	tail2 := "\" and endDate = \""
	tail3 := "\""

	fmt.Println("\n\n---------Historical Single Stock Mode.---------\n")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("ENTER STOCK SYMBOL (\"exit\" to return to Main Menu) --> ")
		symbol, _ := reader.ReadString('\n')
		symbol = strings.Replace(symbol, "\n", "", -1)
		fmt.Println("\n")
		if symbol == "exit" {
			break
		}

		fmt.Print("ENTER START DATE (YYYY-MM-DD) (\"exit\" to return to Main Menu) --> ")
		start, _ := reader.ReadString('\n')
		start = strings.Replace(start, "\n", "", -1)
		fmt.Println("\n")
		if symbol == "exit" {
			break
		}

		fmt.Print("ENTER END DATE (YYYY-MM-DD) (\"exit\" to return to Main Menu) --> ")
		end, _ := reader.ReadString('\n')
		end = strings.Replace(end, "\n", "", -1)
		fmt.Println("\n")
		if symbol == "exit" {
			break
		}

		request := historicalSingleStockQuery + symbol + tail + start + tail2 + end + tail3

		retryCount := maxRetryCount

		for {
			stmt, err := db.Query(request)
			if err != nil {
				if retryCount == 0 {
					fmt.Println("Query Error. Could not retrieve historical information for ", symbol, ".\n")
					retryCount = maxRetryCount

					break
				} else {
					retryCount--
				}
			} else {

				index := 0
				var finalData []map[string]interface{}

				for stmt.Next() {
					var data map[string]interface{}
					stmt.Scan(&data)

					finalData = append(finalData, data)
					index++
				}

				first := true
				index--

				fmt.Println("==============================================")

				for i := index; i > -1; i-- {
					if first {
						fmt.Println(finalData[i][stockKeys.symbol], "\n")
						first = false
					}
					fmt.Println("Date: ", finalData[i][stockKeys.date])
					fmt.Println("Price: ", finalData[i][stockKeys.close])
					fmt.Println("==============================================\n")
				}

				retryCount = maxRetryCount
				break
			}
		}
	}
}

func historicalAnalytics(db *sql.DB) {

	historicalSingleStockQuery := "select * from yahoo.finance.historicaldata where symbol = \""
	tail := "\" and startDate = \""
	tail2 := "\" and endDate = \""
	tail3 := "\""

	fmt.Println("\n\nFaroughi Analytics.\n")

	lines, err := readLines("/Users/kianfaroughi/Documents/goWorkspace/stocks/symbols/nasdaqlisted.txt")
	if err != nil {
		fmt.Println("readLines: %s", err)
	}

	var symbols []string
	for i := 0; i < len(lines); i++ {
		symbol := strings.Split(lines[i], "|")
		symbols = append(symbols, symbol[0])
	}

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print("ENTER START DATE (YYYY-MM-DD) (\"exit\" to return to Main Menu) --> ")
		start, _ := reader.ReadString('\n')
		start = strings.Replace(start, "\n", "", -1)
		fmt.Println("\n")
		if start == "exit" {
			break
		}

		fmt.Print("ENTER END DATE (YYYY-MM-DD) (\"exit\" to return to Main Menu) --> ")
		end, _ := reader.ReadString('\n')
		end = strings.Replace(end, "\n", "", -1)
		fmt.Println("\n")
		if end == "exit" {
			break
		}

		retryCount := maxRetryCount

		failures := 0

		for i := 1; i < 50; i++ {

			request := historicalSingleStockQuery + symbols[i] + tail + start + tail2 + end + tail3

			for {
				stmt, err := db.Query(request)
				if err != nil {
					if retryCount == 0 {
						fmt.Println("Query Error. Could not retrieve information for ", symbols[i], ".\n")
						retryCount = maxRetryCount
						failures++
						break
					} else {
						retryCount--
					}
				} else {

					index := 0
					var finalData []map[string]interface{}

					for stmt.Next() {
						var data map[string]interface{}
						stmt.Scan(&data)

						finalData = append(finalData, data)
						index++
					}
					index--

					startPrice, _ := strconv.ParseFloat(finalData[index][stockKeys.close].(string), 64)
					endPrice, _ := strconv.ParseFloat(finalData[0][stockKeys.close].(string), 64)
					difference := endPrice - startPrice

					fmt.Print("\n", symbols[i])
					fmt.Print("\nPrice on ", finalData[index][stockKeys.date].(string), ": ")
					fmt.Print(startPrice)
					fmt.Print("\nPrice on ", finalData[0][stockKeys.date].(string), ": ")
					fmt.Print(endPrice)
					fmt.Print("\nDifference: ")
					fmt.Print(difference)
					fmt.Println("\n")

					retryCount = maxRetryCount
					break
				}
			}
		}

		fmt.Println("\nNumber of failure: ", failures, "\n")

	}
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
