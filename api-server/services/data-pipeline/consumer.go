package datapipeline

import (
	"api-server/config"
	"api-server/pubsub"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ConsumeMessages(consumer *pubsub.KafkaConsumer) {
	for {
		consumer.Consume([]string{config.Get().MessageTopic})
		for msg := range consumer.Messages() {
			fileName := string(msg.Key)
			fileUrl := string(msg.Value)
			fmt.Printf("Received message: Key - %s, Value - %s\n", fileName, fileUrl)

			fileNameWithPath := BasePath + string(fileName)
			resp, err := http.Get(fileUrl)
			if err != nil {
				panic(err)
			}

			file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(file, resp.Body)
			if err != nil {
				panic(err)
			}

			resp.Body.Close()
			file.Close()

			Files <- &DataFile{
				FileNameWithPath: fileNameWithPath,
			}
		}
	}
}
