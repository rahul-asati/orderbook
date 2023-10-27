/*
Package handlers contains business logic for each rest endpoint
*/

package handlers

import (
	"encoding/json"
	"net/http"

	ob "github.com/rahul-asati/orderbook/orderbook"
	"github.com/shopspring/decimal"
)

// Stores order details to send back to user
type OrderStatus struct {
	Done                     []ob.Order      `json:"done"`
	Partial                  ob.Order        `json:"partial"`
	PartialQuantityProcessed decimal.Decimal `json:"partialQuantityProcessed"`
	QuantityLeft             decimal.Decimal `json:"quantityLeft"`
}

// In memory cache for fast searching
var (
	OrderBooks   map[string]*ob.OrderBook // stores order book id to order book object
	OrdersToBook map[string]string        // stores order id to orderbook id
)

// Initialize cache
func init() {
	OrderBooks = map[string]*ob.OrderBook{}
	OrdersToBook = map[string]string{}
}

// Handler function for '/orderbook/create'. Create new orderbook and returns its Id
func HandleCreateOrderBook(w http.ResponseWriter, r *http.Request) {
	if !checkHttpPost(r.Method, w) {
		return
	}
	orderBook := ob.NewOrderBook()
	id := generateId() // generate unique id
	OrderBooks[id] = orderBook
	// generate response
	response, err := json.Marshal(map[string]string{"orderbook_id": id})
	if err != nil {
		http.Error(w, "Inernal server error.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// Handler function for '/order/cancel' to cancel the order.
func HandleCancelOrder(w http.ResponseWriter, r *http.Request) {
	if !checkHttpPost(r.Method, w) {
		return
	}
	orderId := r.FormValue("order_id")
	orderBook, err := getOrderBookFromOrderId(orderId) // get orderbook for given order id from cache
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result := orderBook.CancelOrder(orderId); result != nil { // try to cancel order
		removeOrderEntry(orderId) // update cache
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Order Cancel request is successful."))
	} else {
		http.Error(w, "Failed to cancel order", http.StatusInternalServerError)
	}
	return
}

// Handler function for '/orderbook/<id>'. This returns full view of orderbook.
func HandleOrderBookDetails(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/orderbook/"):]
	orderBook, err := getOrderBookFromId(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := orderBook.MarshalJSON()
	if err != nil {
		http.Error(w, "Inernal server erro.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

// Handler function for '/order/<id>'. This returns order details.
func HandleOrderDetails(w http.ResponseWriter, r *http.Request) {
	orderId := r.URL.Path[len("/order/"):]
	orderBook, err := getOrderBookFromOrderId(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response, err := orderBook.Order(orderId).MarshalJSON()
	if err != nil {
		http.Error(w, "Inernal server error.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

// Handler function for '/order/market'. This places the market order and returns the order details back to user.
func HandleMarketOrder(w http.ResponseWriter, r *http.Request) {
	if !checkHttpPost(r.Method, w) {
		return
	}
	side, err := parseSide(r.FormValue("side"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	quantity, err := parseFloat("quantity", r.FormValue("quantity"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderBookId := r.FormValue("orderbook_id")
	if orderBookId == "" {
		http.Error(w, "Missing orderbook_id", http.StatusBadRequest)
		return
	}
	orderBook, err := getOrderBookFromId(orderBookId) // get orderbook for given orderbook id from cache
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	done, partial, partialQuantityProcessed, quantityLeft, err := orderBook.ProcessMarketOrder(side, quantity)
	if err != nil {
		http.Error(w, "Encountered error in placing market order:"+err.Error(), http.StatusInternalServerError)
		return
	}
	updateOrderToOrderBookMap(orderBookId, "", append(done, partial)) // update cache
	response, err := marshalProcessedOrder(done, partial, partialQuantityProcessed, quantityLeft)
	if err != nil {
		http.Error(w, "Order is place. Internal server error in returning the details.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

// Handler function for '/order/limit'. This places the limit order and returns the order details back to user.
func HandleLimitOrder(w http.ResponseWriter, r *http.Request) {
	if !checkHttpPost(r.Method, w) {
		return
	}
	side, err := parseSide(r.FormValue("side"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	quantity, err := parseFloat("quantity", r.FormValue("quantity"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	price, err := parseFloat("price", r.FormValue("price"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderBookId := r.FormValue("orderbook_id")
	if orderBookId == "" {
		http.Error(w, "Missing orderbook_id", http.StatusBadRequest)
		return
	}
	orderBook, err := getOrderBookFromId(orderBookId) // get orderbook for given orderbook id from cache
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderId := generateId()
	done, partial, partialQuantityProcessed, err := orderBook.ProcessLimitOrder(side, orderId, quantity, price)
	if err != nil {
		http.Error(w, "Encountered error in placing limit order:"+err.Error(), http.StatusInternalServerError)
		return
	}
	updateOrderToOrderBookMap(orderBookId, orderId, append(done, partial)) // update cache
	response, err := marshalProcessedOrder(done, partial, partialQuantityProcessed, decimal.New(0, 0))
	if err != nil {
		http.Error(w, "Order is placed. Internal server error in returning the details.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

// Handles '/orderbook/marletview/' request and returns market view to user.
func HandleMarketView(w http.ResponseWriter, r *http.Request) {
	orderBookId := r.URL.Path[len("/orderbook/marketview/"):]
	if orderBookId == "" {
		http.Error(w, "Missing orderbook_id", http.StatusBadRequest)
		return
	}
	orderBook, err := getOrderBookFromId(orderBookId) // get orderbook for given orderbook id from cache
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	mktview := orderBook.MarketOverview()
	response, err := json.Marshal(mktview)
	if err != nil {
		http.Error(w, "Order is placed. Internal server error in returning the details.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}
