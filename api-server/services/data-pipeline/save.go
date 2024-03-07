package datapipeline

import "fmt"

func SaveEmbeddings(embeds <-chan *DataEmbed) {
	for {
		for embed := range embeds {
			fmt.Println(embed)
		}
	}
}
