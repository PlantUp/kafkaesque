
![kafkaesque logo](assets/logo.png)
# Reactive consumers for Apache Kafka
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
    go run .
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
2020/09/28 12:12:08 Showing events:
2020/09/28 12:12:11 Moby Dick
```