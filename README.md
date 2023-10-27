# Order Booking
This is an improved matching engine web server framework. It uses following matching engine module: https://github.com/i25959341/orderbook

# Setup
1. Building from repository
    - Clone the repository.
    - Go inside repository.
    - Run - `go mod download`
    - Run - `go build `
    - run - `./orderbooking <port>`
This will start server at given port.

2. For Windows: 
Just run the build/orderbook.exe with port number as an argument on windows system.

# Endpoints
All endpoints returns json structure
1. Create orderbook
- Target: /orderbook/create
- Type: Post
- Example: `curl http://localhost:8080/orderbook/create`
- Sample Return:
`
{
    "orderbook_id": "de8bdf60-e4ea-4079-beb2-33583af248ba"
}`

2. Get Orderbook details
- Target: /orderbook/<orderbookid>
- Type: Get
- Example: `curl http://localhost:8080/orderbook/de8bdf60-e4ea-4079-beb2-33583af248ba`
- Sample Return:

```json
{
    "asks": {
        "numOrders": 3,
        "depth": 1,
        "prices": {
            "4.5": {
                "volume": "300",
                "price": "4.5",
                "orders": [
                    {
                        "side": "sell",
                        "id": "7e22530d-2203-483b-8008-9b79e66d682a",
                        "timestamp": "2023-10-27T07:15:23.8461309Z",
                        "quantity": "100",
                        "price": "4.5"
                    },
                   
                    {
                        "side": "sell",
                        "id": "eba8e1b3-e4b4-44bd-bf05-3c5418f63887",
                        "timestamp": "2023-10-27T07:15:53.9445147Z",
                        "quantity": "100",
                        "price": "4.5"
                    }
                ]
            }
        }
    },
    "bids": {
        "numOrders": 0,
        "depth": 0,
        "prices": {}
    }
}
```

3. Get Order details
- Target: /order/<orderid>
- Type: Get
- Example: `curl http://localhost:8080/order/f3cd7ee1-58cb-42e9-b7a4-179d9c3a8c5b`
- Sample Return:
```json
{
    "side": "sell",
    "id": "f3cd7ee1-58cb-42e9-b7a4-179d9c3a8c5b",
    "timestamp": "2023-10-27T11:00:36.3562441Z",
    "quantity": "100",
    "price": "4.5"
}
```

4. Post limit order
- Target: /order/limit
- Type: Post
- Form Body:
    - side - 0 for buy, 1 for sell
    - quantity
    - price
    - orderbook_id
- Example: `curl -X POST http://localhost:8080/order/limit -d "side=0,quantity=2000,price=100,orderbook_id=de8bdf60-e4ea-4079-beb2-33583af248ba"`
- Sample Return:
```json
{
    "done": null,
    "partial": {
        "side": "sell",
        "id": "",
        "timestamp": "0001-01-01T00:00:00Z",
        "quantity": "0",
        "price": "0"
    },
    "partialQuantityProcessed": "0",
    "quantityLeft": "0"
}
```

5. Post market order
- Target: /order/market
- Type: Post
- Form Body:
    - side - 0 for buy, 1 for sell
    - quantity
    - orderbook_id
- Example: `curl -X POST http://localhost:8080/order/market -d "side=0,quantity=2000,orderbook_id=de8bdf60-e4ea-4079-beb2-33583af248ba"`
- Sample Return:
```json
{
    "done": null,
    "partial": {
        "side": "sell",
        "id": "",
        "timestamp": "0001-01-01T00:00:00Z",
        "quantity": "0",
        "price": "0"
    },
    "partialQuantityProcessed": "0",
    "quantityLeft": "10980"
}
```

6. Get Orderbook marketview
- Target: /orderbook/<orderbookid>
- Type: Get
- Example: `curl localhost:8080/orderbook/marketview/de8bdf60-e4ea-4079-beb2-33583af248ba`
- Sample Return:
```json
{
    "asks": {
        "4.5": "100"
    },
    "bids": {}
}
```