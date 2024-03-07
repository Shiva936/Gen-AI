package datapipeline

import (
	"fmt"
)

func EmbedVectos(chunks <-chan *DataChunk) {
	for {
		for chunk := range chunks {
			fmt.Println(chunk)
		}
	}
}
