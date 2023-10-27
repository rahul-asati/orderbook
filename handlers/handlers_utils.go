/*
Package handlers contains business logic for each rest endpoint
*/

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	ob "github.com/rahul-asati/orderbook/orderbook"
	"github.com/shopspring/decimal"
)

// Generate unique id.
func generateId() string {
	id := uuid.New()
	return id.String()
}

// Try to parse val as float.
func parseFloat(field, val string) (decimal.Decimal, error) {
	if val == "" {
		return decimal.New(0, 0), errors.New(field + "  can not be empty")
	}
	v, err := decimal.NewFromString(val)
	if err != nil {
		return decimal.New(0, 0), errors.New("Invalid " + field)
	}
	return v, nil
}

// Try to parse input as ob.Side.
func parseSide(side string) (ob.Side, error) {
	if side == "0" {
		return ob.Buy, nil
	} else if side == "1" {
		return ob.Sell, nil
	} else {
		return -1, errors.New("Side should be either 0 or 1")
	}
}

// Checks if the http request is Post request or not.
func checkHttpPost(method string, w http.ResponseWriter) bool {
	if method != http.MethodPost {
		http.Error(w, "Invalid request method. Post Method is required", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// Get order book from cache for the given orderbook id.
func getOrderBookFromId(orderBookId string) (*ob.OrderBook, error) {
	if orderBookId == "" {
		return nil, errors.New("Missing orderbook_id")
	}
	o, ok := OrderBooks[orderBookId]
	if !ok {
		return nil, errors.New("orderbook_id is invalid")
	}
	return o, nil
}

// Get orderbook from the cache for the given orderid.
func getOrderBookFromOrderId(orderId string) (*ob.OrderBook, error) {
	if orderId == "" {
		return nil, errors.New("Missing order_id")
	}
	orderBookId, ok := OrdersToBook[orderId]
	if !ok {
		return nil, errors.New("order_id is invalid")
	}
	return getOrderBookFromId(orderBookId)
}

// Remove order entry from cache.
func removeOrderEntry(orderId string) {
	delete(OrdersToBook, orderId)
}

// Convert limit and market order output to common structure and marshall that to return to user.
func marshalProcessedOrder(done []*ob.Order, partial *ob.Order, partialQuantityProcessed, quantityLeft decimal.Decimal) ([]byte, error) {
	s := OrderStatus{
		PartialQuantityProcessed: partialQuantityProcessed,
		QuantityLeft:             quantityLeft,
	}
	for _, o := range done {
		s.Done = append(s.Done, *o)
	}
	if partial != nil {
		s.Partial = *partial
	}
	return json.Marshal(&s)
}

// Update orderid to orderbookid cache for given orderbookId,orderid and orders.
func updateOrderToOrderBookMap(orderBookId string, orderID string, orders []*ob.Order) {
	if orderID != "" {
		OrdersToBook[orderID] = orderBookId
	}
	for _, o := range orders {
		if o != nil {
			OrdersToBook[o.ID()] = orderBookId
		}
	}
}
