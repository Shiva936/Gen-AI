package documents

import (
	"api-server/config"
	"api-server/daos"
	"api-server/daos/models"
	"api-server/pubsub"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type IDocument interface {
	SaveAndTriggerDataPipeline(documents map[string]string) error
}

type Documnet struct {
	documnets daos.IDocument
}

func NewDocument() IDocument {
	return &Documnet{
		documnets: daos.NewDocument(),
	}
}

func (t *Documnet) SaveAndTriggerDataPipeline(documents map[string]string) error {
	ctx := context.Background()

	producer, err := pubsub.NewKafkaProducer([]string{config.Get().MessageBrokerURL})
	if err != nil {
		log.Fatal("Error creating Kafka producer:", err)
	}
	defer producer.Close()

	for name, url := range documents {
		// Start Pipeline with GO

		document := &models.Document{
			Id:   uuid.New(),
			Name: name,
			URL:  url,
		}
		err := t.documnets.Save(&ctx, document)
		if err != nil {
			log.Println("[documents] Error while processing ", fmt.Sprintf(":%+v :", document), err.Error())
			continue
		}

		message := &pubsub.Message{
			Key:   name,
			Value: url,
		}
		err = producer.Produce(config.Get().MessageTopic, message)
		if err != nil {
			log.Println("[kafka-producer] uanble to trigger pipeline ", fmt.Sprintf(": %v : %+v :", config.Get().MessageTopic, message), err.Error())
		}
	}

	return nil
}
