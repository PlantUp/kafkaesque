
![kafkaesque logo](assets/logo.png)
# Reactive consumers for Apache Kafka
[![Go Report Card](https://goreportcard.com/badge/github.com/PlantUp/kafkaesque)](https://goreportcard.com/report/github.com/PlantUp/kafkaesque)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/PlantUp/kafkaesque)](https://pkg.go.dev/mod/github.com/PlantUp/kafkaesque)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)  
This is a simple package for using Kafka in Golang with the Observer-pattern.

## Example
We have two Kafka-topics in JSON and want to serialize the events in a pipeline.
### Create new type  
in order to serialize our events
```go
type BookRequest struct {
	Book   string `json:"book"`
	Person string `json:"person"`
}
```
### Create a new config
```go
const config = &kafkaesque.Config{
		AutoOffsetReset:  kafkaesque.AutoOffsetReset.Latest,
		BootstrapServers: []string{"localhost:9092"},
		IsolationLevel:   kafkaesque.IsolationLevel.ReadCommitted,
		MessageMaxBytes:  1048588,
		ClientID:         "test_1",
		GroupID:          "test_group",
		SecurityProtocol: kafkaesque.SecurityProtocol.Plaintext,
    }
```
### Create the consumer
```go
testConsumer, err := kafkaesque.NewConsumer(config, "book-requests", "external-book-requests")
if err != nil {
	    panic(err)
    }
```
### Serialize BookRequest
```go
bookRequests := testConsumer.Events.
		Map(func(_ context.Context, item interface{}) (interface{}, error) {
			var request BookRequest
			err := json.Unmarshal(kafkaesque.ValueToBytes(item), &request)
			return request, err
		})
```
### Show Types
```go
for event := range bookRequests.Observe() {
		log.Println(fmt.Sprintf("%T", event.V))
	}
```
### Test
[Test Program](example/main.go)
1. Start your Confluent Plattform:
    ```bash
    cd example
    docker-compose up
    ```
2. Start the program:
    ```bash
    go run . -name test_consumer -topics test
    ```
3. In a new terminal run:
    ```bash
    docker exec -it broker bin/kafka-console-producer --topic test --bootstrap-server localhost:9092
    ```
4. Insert a BookRequest in JSON:
    ```json
    {"book": "Moby Dick", "person": "Max Musterknabe"}
    ```
    The output should be:
    ```
    2020/10/03 15:32:08 Showing events for [test test1 test2] on [localhost:9092]:
    2020/10/03 15:32:16 Book:	Moby Dick	Person:	Max Musterknabe
    ```
   **Or** start the example as producer:
     ```bash
	 go run . -mode producer -name test_producer -topics test
	 ```

    The output should be:
    ```
    2020/10/03 15:32:08 Showing events for [test test1 test2] on [localhost:9092]:
    2020/10/03 15:32:16 Book:	10 Awesome Ways to Photograph Blue Bottles	Person:	Pamela Osborne
    2020/10/03 15:33:02 Book:	10 Awesome Ways to Photograph Blue Bottles	Person:	Pamela Osborne
    2020/10/03 15:33:02 Book:	21 Myths About Blue bottles Debunked	Person:	Inayah Brown
    2020/10/03 15:33:02 Book:	How to Make Your Own Vast Hat for less than £5	Person:	Ikrah Blair
    2020/10/03 15:33:03 Book:	An analysis of handsome jugs	Person:	Aurora Stafford
                                    ...
    ```