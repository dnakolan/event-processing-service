# event-processing-service
A Go service that processes user activity events (clicks, views, purchases)

# Features
* Batch Event Ingestion - POST /events (accept arrays of events)
* Analytics API - GET /analytics/summary?window=1h|24h|7d
* Real-time WebSocket - Stream live events to connected clients
* Event Validation - Schema validation, deduplication
* In-memory Aggregation - Count events by type, user, time windows
* Docker containerization

# Tech Stack
* Language: Go 1.24+
* Router: chi
* UUIDs: github.com/google/uuid
* Testing: Go standard library

# Project Structure
```
trail-data-service/
├── models/
├──── cmd/
│     └── main.go     // application entry point
├── handlers/
│   └── events.go     // http handlers for /events
│   └── health.go     // http handlers for /health
│   └── analytics.go  // http handlers for /analytics
├── models/
│   └── events.go     // data models and validation for /events
├── services/
│   └── events.go     // supporting service for /trails endpoints
│   └── analytics.go  // supporting service for /trails endpoints
├── storage/
│   └── memory.go     // in memory store of event data
└── build.sh          // builds docker images for the application
└── Dockerfile        // core application dependencies
└── Dockerfile.deps   // isolated base image to speed up docker build
```

# Event Schema
```
{
  "event_id": "uuid",
  "user_id": "string",
  "event_type": "page_view|click|purchase|signup",
  "timestamp": "RFC3339",
  "properties": {
    "page": "/home",
    "amount": 29.99,
    "product_id": "xyz"
  }
}
```

# Running the Service
First run the included build.sh script to build the container images
```
./build.sh
```

Then start the application in docker with the following command
```
docker run -d -p 8080:8080 --name event-service event-processing-service
```


# Example Usage (cURL)
## POST /events - create events
```
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
        {
          "event_id": "e58ed763-928c-4155-bee9-fdbaaadc15f3",
          "user_id": "123",
          "event_type": "page_view",
          "timestamp": "2006-01-02T15:04:05Z07:00",
          "properties": {
            "page": "/home",
            "amount": 29.99,
            "product_id": "xyz"
          }
        }
}'
```

GET /analytics/summary?window=1h|24h|7d
```
curl http:///analytics/summary?window=24h

Returns
{
  "time_window": "24h",
  "total_events": 1250,
  "events_by_type": {
    "page_view": 800,
    "click": 350,
    "purchase": 100
  },
  "unique_users": 45,
  "events_per_hour": [
    {"hour": "2025-05-26T14:00:00Z", "count": 120},
    {"hour": "2025-05-26T15:00:00Z", "count": 95}
  ]
}
```

# Design Considerations
* Dependency Injection is used for loose coupling between components.
* Interface-Driven Architecture enables testability and future extensibility (e.g., database-backed repo).
* Validation is handled at the request model level to separate concerns cleanly.
* The service layer enforces any domain-specific business rules.
* Websocket interface for processing batches of events

# Tests
`go test ./...`
Tests cover handler logic, service behavior, and in-memory repo operations.

# Future Improvements / Next Steps
TBD

# Time Spent
TBD

# Author
David Nakolan - david.nakolan@gmail.com
