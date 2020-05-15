# go-adp-bridge
A Golang bridge adapter.

This kind of adapter act like a bridge between topics, subscribing in one or N topics, transforming data and publishing it into another topic(s). Optionally, it can call http clients to enrich the original data.

# run whole app on docker
1) if you don't have yet, install the docker and docker-compose (links bellow)
2) on the project's root folder type "docker-compose -f docker/compose-app.yml up --build" (optionally you can use -d to unlock the terminal)
3) access http://localhost:15672/#/queues
4) click on the "q-merchants" or "q-sellers" to send data (that can be found into "tests/support/mock/")
5) go to "Publish message", fill the "Payload" field with correspondent data and click in "Publish message" button
6) go to http://localhost:15672/#/queues/ and click on "q-accounts"
7) go to "Get messages", increase the number of desired messages to get on "Messages" field and click in "Get Message(s)" button
8) the data sent to topics should appear transformed into an account

# run only infrastructure on the docker and app locally
1) if you don't have yet, install the docker and docker-compose (links bellow)
2) install Golang (links bellow)
3) on the project's root folder type "docker-compose -f docker/compose-infra.yml up --build" (optionally you can use -d to unlock the terminal)
4) on the project's root folder type "go mod download" then "go run *.go"

# application health check routes
Once running, you can call http://localhos:3000/live (also /health and /ready) to check the app status

# requirements
    - [x] consume topics
    - [x] publish on topics
    - [x] call http clients
    - [x] expose por 3000 to check health
    - [x] use in data models:
        - [x] strings
        - [x] slices
        - [x] dates (date time and only date)
        - [x] int
        - [x] bool
        - [x] float
    - [ ] tests with coverage

# links
- Golang: https://golang.org
- Docker installation: https://docs.docker.com/install
- Docker Compose installation: https://docs.docker.com/compose/install