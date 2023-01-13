package producer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"os"
	model "restAPIKafkaProducer/model"
)

func Producer(message model.Data) {
	topic := "none"
	if message.Flag == "PLN" || message.Flag == "pln" {
		topic = "DataPLN"
	} else if message.Flag == "PDAM" || message.Flag == "pdam" {
		topic = "DataPDAM"
	} else if message.Flag == "ISI PULSA" || message.Flag == "isi pulsa" {
		topic = "DataIsiPulsa"
	} else {
		panic("ERROR")
	}
	fmt.Println(topic)

	bootsrapServer := "172.18.46.121:9092"
	registryUrl := "http://172.18.46.121:8081"

	// Create Producer instance
	p, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": bootsrapServer},
	)
	if err != nil {
		panic(err)
	}
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(registryUrl))

	if err != nil {
		fmt.Printf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	ser, err := avro.NewSpecificSerializer(client, serde.ValueSerde, &avro.SerializerConfig{
		SerializerConfig: serde.SerializerConfig{
			AutoRegisterSchemas: false,
			UseSchemaID:         -1,
			UseLatestVersion:    true}})

	if err != nil {
		fmt.Printf("Failed to serialize payload: %s\n", err)
		os.Exit(1)
	}

	messageAvro := Dataschema{
		Tanggal:     message.Tanggal,
		Nama_vendor: message.NamaVendor,
		Kode_barang: message.KodeBarang,
		Cara_bayar:  message.CaraBayar,
		Qty:         int32(message.QTY),
		Harga:       int64(message.Harga),
		Jumlah:      int64(message.Jumlah),
	}
	// Verbose
	//fmt.Println(&messageAvro)

	// Serialize messageAvro
	payload, err := ser.Serialize(topic, &messageAvro)
	if err != nil {
		fmt.Printf("Failed to create serializer: %s\n", err)
		os.Exit(1)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	// not used (verbose)
	//stringMessage := string(message)
	//reqBodyBytes := new(bytes.Buffer)
	//json.NewEncoder(reqBodyBytes).Encode(message)
	//fmt.Printf(string(reqBodyBytes.Bytes()))
	// Produce messages to topic (asynchronously)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
	}, nil)
	if err != nil {
		return
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
