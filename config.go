package kafkaesque

import(
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"fmt"
	"strings"
)

type AutoOffsetResetOption string

const AutoOffsetReset = struct {
	Earliest:	AutoOffsetResetOption("earliest"),
	Latest:		AutoOffsetResetOption("latest"),
	Largest:	AutoOffsetResetOption("largest"),
	Smallest:	AutoOffsetResetOption("smallest"),
}

type SecurityProtocolOption string

const SecurityProtocol = struct {
	Plaintext:		SecurityProtocolOption("plaintext"),
	SSL:			SecurityProtocolOption("ssl"),
	SASL_Plaintext:	SecurityProtocolOption("sasl_plaintext"),
	SASL_SSL:		SecurityProtocolOption("sasl_ssl"),
}

type IsolationLevelOption string

const IsolationLevel = struct {
	ReadUncommitted:	IsolationLevelOption("read_uncommited"),
	ReadCommitted:		IsolationLevelOption("read_commited"),
}

type Config struct {
	Acks						int,
	AutoOffsetReset				AutoOffsetResetOption,
	BootstrapServers			[]string,
	ClientId					string,
	GroupId						string,
	IsolationLevel				IsolationLevelOption,
	MessageMaxBytes				uint32,
	SecurityProtocol			SecurityProtocolOption,
	SslCertificateLocation		string,
	SslCaLocation				string,
	SslKeyLocation				string,
	SslValidate					bool,
	//TODO: complete config
}

func (c *Config) Map() *kafka.ConfigMap {
	 return &kafka.ConfigMap {
		"acks":									fmt.Sprintf("%d",c.Acks),
		"auto.offset.reset": 					string(c.AutoOffsetReset),
		"client.id":							c.ClientId
		"bootstrap.servers": 					strings.Join(c.BootstrapServers, ','),
		"group.id":          					c.GroupId,
		"isolation.level":						string(c.IsolationLevel)
		"message.max.bytes":					fmt.Sprintf("%d",c.MessageMaxBytes),
		"security.protocol":					string(c.SecurityProtocol),
		"ssl.ca.location"						c.SslCaLocation,
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
 