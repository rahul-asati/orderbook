package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rahul-asati/orderbook/handlers"
)

func main() {
	var port string
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else {
		log.Fatal("Need only 1 argument: port number")
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatal("port should be in correct range")
	}

	// register handler functions and endpoints.
	http.HandleFunc("/orderbook/create", handlers.HandleCreateOrderBook) // handles orderbook creation
	http.HandleFunc("/order/limit", handlers.HandleLimitOrder)           // handles limit order request
	http.HandleFunc("/order/market", handlers.HandleMarketOrder)         // handles market order request
	http.HandleFunc("/order/cancel", handlers.HandleCancelOrder)         // handles cancel order request
	http.HandleFunc("/orderbook/", handlers.HandleOrderBookDetails)      // handles orderbook details request
	http.HandleFunc("/order/", handlers.HandleOrderDetails)              // handles order details request
	http.HandleFunc("/orderbook/marketview/", handlers.HandleMarketView) // handles orderbook marketview request

	log.Println("Server is running on :" + port)
	http.ListenAndServe(":"+port, nil)
}
