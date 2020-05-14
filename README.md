# go-adp-bridge
A Golang bridge adapter.

This kind of adapter act like a bridge between topics, subscribing in one or N topics, transforming data and publishing it into another topic(s). Optionally, it can call http clients to enrich the original data.

# requirements
    - [x] consume topics
    - [x] publish on topics
    - [x] call http clients
    - [x] expose por 3000 to check health
    - [ ] use in data models:
        - [x] strings
        - [x] slices
        - [x] dates (date time and only date)
        - [x] int
        - [x] bool
        - [x] float
    - [ ] tests with coverage
