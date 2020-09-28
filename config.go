package kafkaesque

import(
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"fmt"
	"strings"
)

// AutoOffsetResetOption wether to commit the message you brought from Kafka using Zookeeper to persist the last 'offset' which it read.
type AutoOffsetResetOption string

// AutoOffsetReset commit the message you brought from Kafka using Zookeeper to persist the last 'offset' which it read.
var AutoOffsetReset = struct {
	Earliest	AutoOffsetResetOption
	Latest		AutoOffsetResetOption
	Largest		AutoOffsetResetOption
	Smallest	AutoOffsetResetOption
} {
	Earliest:	"earliest",
	Latest:		"latest",
	Largest:	"largest",
	Smallest:	"smallest",
}

// SecurityProtocolOption what security protocol to use.
type SecurityProtocolOption string

// SecurityProtocol the security protocols for kafka.
var SecurityProtocol = struct {
	Plaintext		SecurityProtocolOption
	SSL				SecurityProtocolOption
	SASLPlaintext	SecurityProtocolOption
	SASLSSL			SecurityProtocolOption
} {
	Plaintext:		"plaintext",
	SSL:			"ssl",
	SASLPlaintext:	"sasl_plaintext",
	SASLSSL:		"sasl_ssl",
}

// IsolationLevelOption what isolation level to use.
type IsolationLevelOption string

// IsolationLevel the kafka isolation level for transactions.
var IsolationLevel = struct {
	ReadUncommitted		IsolationLevelOption
	ReadCommitted		IsolationLevelOption
} {
	ReadUncommitted:	"read_uncommitted",
	ReadCommitted:		"read_committed",
}

// Config the configuration for your kafka consumer.
type Config struct {
	Acks						int
	AutoOffsetReset				AutoOffsetResetOption
	BootstrapServers			[]string
	ClientID					string
	GroupID						string
	IsolationLevel				IsolationLevelOption
	MessageMaxBytes				uint32
	SecurityProtocol			SecurityProtocolOption
	SslCertificateLocation		string
	SslCaLocation				string
	SslKeyLocation				string
	SslValidate					bool
	//TODO: complete config
}

// Map the Config-struct to a ConfigMap for confluent.
func (c *Config) Map() *kafka.ConfigMap {
	 return &kafka.ConfigMap {
		"acks":									fmt.Sprintf("%d",c.Acks),
		"auto.offset.reset": 					string(c.AutoOffsetReset),
		"client.id":							c.ClientID,
		"bootstrap.servers": 					strings.Join(c.BootstrapServers, ","),
		"group.id":          					c.GroupID,
		"isolation.level":						string(c.IsolationLevel),
		"message.max.bytes":					fmt.Sprintf("%d",c.MessageMaxBytes),
		"security.protocol":					string(c.SecurityProtocol),
		"ssl.ca.location":						c.SslCaLocation,
		"ssl.certificate.location":				c.SslCertificateLocation,
		"ssl.key.location":						c.SslKeyLocation,
		"enable.ssl.certificate.verification":	func(validate bool) string {
													if(validate) {
														return "true"
													}
													return "false"
												}(c.SslValidate ),
		"enable.partition.eof":					"true",
	 }
}
 