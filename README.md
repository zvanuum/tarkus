# tarkus

A simple implementation of a blockchain primarily implemented by following [this](https://hackernoon.com/learn-blockchains-by-building-one-117428612f46) Medium article.

TODO:
- API documentation
- Persist chain to file or database somehow (Redis list structure maybe?)
- Use GraphQL for API
- Dockerize
- Logrus for better logging

### Building and Running


- To build the binary:
`$ go build -o tarkus ./main.go -port 8080`

- To run the program:
`$ go run ./main.go -port 8080`
