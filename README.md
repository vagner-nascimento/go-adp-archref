# go-enriching-adp
A Golang enriching adapter.

This kind of adapter receives data from topics or queues (one or many), transform the data and publishing it into another topics or queues (one or many). Usually it call other sources of data (like other http clients) to enrich the original data.

# run whole app on docker
1) if you don't have yet, install the docker and docker-compose (links bellow)
2) on the project's root folder type "docker-compose -f docker/compose-app.yml up --build" (optionally you can use -d to unlock the terminal)
3) access http://localhost:15672/#/queues
4) click on the "q-merchants" or "q-sellers" to send data (that can be found into "tests/support/mocks/")
5) go to "Publish message", fill the "Payload" field with correspondent data and click in "Publish message" button
6) go to http://localhost:15672/#/queues/ and click on "q-accounts"
7) go to "Get messages", increase the number of desired messages to get on "Messages" field and click in "Get Message(s)" button
8) the data sent to topics should appear transformed into an account

# run only infrastructure on the docker and app locally
1) if you don't have yet, install the docker and docker-compose (links bellow)
2) install Golang (links bellow)
3) on the project's root folder type "docker-compose -f docker/compose-infra.yml up --build" (optionally you can use -d to unlock the terminal)
4) on the project's root folder type "go mod download" then "go run *.go"

# run stress tests
1) go to "tests" folder
2) open the file "compose-stress.yml" and, into "go-stress-test" service, set "QTD_SELL" (quantity of sellers event to send), "QTD_MERCH" (quantity of merchants event to send) and "MINUTES_TIMEOUT" (timeout to complete test) 
3) save and run this command to start tests: docker-compose -f compose-stress.yml up --build -d
4) when the terminal was free, run this command to watch tests results: docker logs go-stress-test --follow
5) you can also access http://localhost:15672/#/queues/ to check the rabbit mq queues

# application health check routes
Once running, you can call http://localhos:3000/live (also /health and /ready) to check the app status

# stress test result running on docker
    - Tests with 50k merchants and 50k sellers, total of 100k messages
    - PC configs: Intel i7 9th gen and 8GB of memory
    - Running producer and consumer togheter:
        - 1st: completed in 9:43 minutes
        - 2nd (after improvments): completed in 6:59 minuts
    - Running only consumer with data already into topics: completed in 6:00 minutes

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
