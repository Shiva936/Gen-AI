package datapipeline

import (
	"context"
	"log"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/textsplitter"
)

func GenerateChunks(files <-chan *DataFile) {
	for {
		for filedata := range files {
			stat, err := os.Stat(filedata.FileNameWithPath)
			if err != nil {
				panic(err)
			}

			file, err := os.Open(filedata.FileNameWithPath)
			if err != nil {
				panic(err)
			}

			// Can be made chunk for each page as if deviding chunks by seperators can cause
			// tables data to be distorted.
			// Need to think of a proper solution to sustain tables related data inside documents.
			doc := documentloaders.NewPDF(file, stat.Size())
			pageChunks, err := doc.LoadAndSplit(context.Background(), textsplitter.NewRecursiveCharacter(
				textsplitter.WithChunkSize(1024),
				textsplitter.WithChunkOverlap(0),
				textsplitter.WithSeparators([]string{"\n\n", "\n", " ", "", "."}),
			))
			if err != nil {
				log.Println("[data-pipeline] error while chunking")
				return
			}

			for _, chunk := range pageChunks {
				Chunks <- &DataChunk{
					Chunk: chunk,
				}
			}
		}
	}
}
