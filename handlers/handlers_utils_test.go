package handlers

import (
	"testing"

	ob "github.com/rahul-asati/orderbook/orderbook"
	"github.com/shopspring/decimal"
)

func TestParseFloatValid(t *testing.T) {
	actual, _ := parseFloat("price", "4.5")
	expected, _ := decimal.NewFromString("4.5")
	a, _ := actual.Float64()
	e, _ := expected.Float64()
	if a != e {
		t.Errorf("Expected  %s is not same as actual  %s", expected, actual)
	}
}

func TestParseFloatInvalid(t *testing.T) {
	_, err := parseFloat("price", "4.5a")
	if err == nil {
		t.Errorf("parseFloat parsing incorrect strings")
	}
	_, err = parseFloat("price", "")
	if err == nil {
		t.Errorf("parseFloat parsing empty strings")
	}
}
func TestParseSideValid(t *testing.T) {
	actual, _ := parseSide("0")
	if actual != ob.Buy {
		t.Errorf("Expected  %s is not same as actual  %s", ob.Buy, actual)
	}
	actual, _ = parseSide("1")
	if actual != ob.Sell {
		t.Errorf("Expected  %s is not same as actual  %s", ob.Sell, actual)
	}
}

func TestParseSideInvalid(t *testing.T) {
	_, err := parseSide("10")
	if err == nil {
		t.Errorf("parseSide parsing invalid values")
	}
}

func TestOrderBookFromId(t *testing.T) {
	book := ob.NewOrderBook()
	OrderBooks["id1"] = book
	actual, _ := getOrderBookFromId("id1")
	if actual != book {
		t.Errorf("orderbook lookup failed")
	}
	if _, err := getOrderBookFromId("id23"); err == nil {
		t.Errorf("orderbook lookup fetching data for wrong id")
	}

}

func TestOrderBookFromOrderId(t *testing.T) {
	book := ob.NewOrderBook()
	OrderBooks["id1"] = book
	OrdersToBook["oid1"] = "id1"
	actual, _ := getOrderBookFromOrderId("oid1")
	if actual != book {
		t.Errorf("orderbook lookup failed by orderid")
	}
	if _, err := getOrderBookFromOrderId("id23"); err == nil {
		t.Errorf("orderbook lookup fetching data for wrong order id")
	}

}
