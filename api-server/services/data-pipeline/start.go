package datapipeline

import (
	"api-server/config"
	"api-server/pubsub"
	"log"
	"sync"

	"github.com/pgvector/pgvector-go"
	"github.com/tmc/langchaingo/schema"
)

type IDataPipeline interface {
	Start()
}

type DataPipeline struct {
}

func NewDataPipeline() IDataPipeline {
	return &DataPipeline{}
}

var BasePath = "/documents/"

type DataFile struct {
	FileNameWithPath string
	sync.Mutex
}

var Files = make(chan *DataFile, 10)

type DataChunk struct {
	Chunk schema.Document
	DataFile
	sync.Mutex
}

var Chunks = make(chan *DataChunk, 10)

type DataEmbed struct {
	Embed pgvector.Vector
	DataChunk
	sync.Mutex
}

var Embeds = make(chan *DataEmbed, 10)

func (t *DataPipeline) Start() {
	// Create a Kafka consumer
	consumer, err := pubsub.NewKafkaConsumer([]string{config.Get().MessageBrokerURL}, []string{config.Get().MessageTopic})
	if err != nil {
		log.Fatal("Error creating Kafka consumer:", err)
	}
	defer consumer.Close()

	go SaveEmbeddings(Embeds)
	go EmbedVectos(Chunks)
	go GenerateChunks(Files)
	go ConsumeMessages(consumer)
}
