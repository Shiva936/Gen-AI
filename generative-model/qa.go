package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Source struct {
	SourceName  string `json:"source_name"`
	SourceUrl   string `json:"source_url"`
	SourceTitle string `json:"source_title"`
	SourceType  string `json:"source_type"`
	PageNo      int32  `json:"page_no"`
}

type Message struct {
	Artifacts     interface{}         `json:"artifacts"`
	ChatCreatedAt string              `json:"chat_created_at"`
	ChatID        string              `json:"chat_id"`
	ChatResponse  string              `json:"chat_response"`
	FollowUps     []string            `json:"follow_ups"`
	Sources       []Source            `json:"sources"`
	ThreadID      string              `json:"thread_id"`
	UserQuestion  string              `json:"user_question"`
	Keywords      map[string][]string `json:"keywords"`
}

func createSummaryPayload(ticker string, theme string, chunks []string) []byte {
	user_prompt_template := "Raw Data from %s Annual Report:\n%s\n----\nSummary:\n---"
	message_template := map[string]interface{}{
		"messages": []interface{}{
			map[string]interface{}{
				"role": "system",
				"content": fmt.Sprintf(`
				As a professional summarizer, create a concise and comprehensive summary of the provided text, be it an article, post, conversation, or passage, while adhering to these guidelines:
				0. Craft a summary that is detailed, thorough, in-depth, and complex, while maintaining clarity and conciseness.
				1. Incorporate main ideas and essential information, eliminating extraneous language and focusing on critical aspects.
				2. Rely strictly on the provided text, without including external information.
				3. Format the summary in paragraph form for easy understanding UNDER 100 WORDS.
				4. The raw data below deals with the %s stock as it relates to %s, keep this context in mind.
				`, ticker, theme),
			},
			map[string]interface{}{
				"role":    "user",
				"content": fmt.Sprintf(user_prompt_template, ticker, strings.Join(chunks[:int(math.Min(float64(4), float64(len(chunks))))], "\n")),
			},
		},
	}

	payload, err := json.Marshal(message_template)
	if err != nil {
		log.Fatalf("Error marshaling request data: %v", err)
	}

	return []byte(payload)
}

func (db *DBOps) RealChat(c *gin.Context) {
	var wg sync.WaitGroup

	theme, tickersRaw := QuerySplitter(c.Query("search_str"))
	engine := "heavyLM"

	if tickersRaw == "" {
		c.JSON(203, gin.H{
			"message":  "Data Missing. No tickers have been defined for this question",
			"keywords": gin.H{},
		})
	}
	theme = strings.TrimSpace(theme)
	fmt.Println(theme)
	tickers := StringListCapitalize(StripWhitespaceFromSlice(strings.Split(tickersRaw, ",")))
	valid_existing_tickers := db.ValidateTickers(tickers)

	copilot := make([]interface{}, len(valid_existing_tickers))
	kw_map := make(map[string][]string, len(valid_existing_tickers))

	for idx, ticker := range valid_existing_tickers {
		wg.Add(1)
		go func(idx int, finalKw []interface{}, wg *sync.WaitGroup, research_topic string, ticker []string) {
			defer wg.Done()
			lm_output := ChatLLMCall(CreatePayload(ticker[1], research_topic), os.Getenv(engine))
			parts := strings.Split(lm_output, "---")
			result := StripWhitespaceFromSlice(strings.Split(parts[len(parts)-1], ","))
			copilot[idx] = map[string]interface{}{
				"ticker":     ticker[0],
				"chunks":     db.ReadDocuments(result, ticker[0]),
				"core_topic": theme,
			}
			kw_map[ticker[0]] = result
		}(idx, copilot, &wg, theme, ticker)
	}
	wg.Wait()

	var sources []Source

	summaries := make([]string, len(copilot))

	// Cross Encoder to pick top 3 chunks per ticker
	for idx, item := range copilot {
		chunks := item.(map[string]interface{})["chunks"]
		ctx := make([]string, len(chunks.([]bson.M)))
		for chunk_no, chunk := range chunks.([]bson.M) {
			sources = append(
				sources,
				Source{
					SourceName:  fmt.Sprintf("%s*%d", chunk["doc_name"].(string), chunk["page_no"].(int32)),
					SourceUrl:   fmt.Sprintf("https://nsedocs.blob.core.windows.net/agm-reports-2223/%s", chunk["doc_name"].(string)),
					SourceTitle: fmt.Sprintf("%s AGM Report", item.(map[string]interface{})["ticker"].(string)),
					SourceType:  "pdf",
					PageNo:      chunk["page_no"].(int32),
				},
			)
			ctx[chunk_no] = chunk["answer"].(string)
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, idx int, summaries []string, ctx []string, ticker string, theme string) {
			defer wg.Done()
			summaries[idx] = fmt.Sprintf(
				"Executive Summary (%s):\n%s\n",
				ticker,
				ChatLLMCall(createSummaryPayload(ticker, theme, ctx), os.Getenv("heavyLM")),
			)
		}(&wg, idx, summaries, ctx, item.(map[string]interface{})["ticker"].(string), theme)
		// Summarize at Ticker Level using GPT4
	}
	wg.Wait()

	chatMessage := map[string]interface{}{
		"message": Message{
			UserQuestion:  c.Query("search_str"),
			Artifacts:     nil,
			FollowUps:     []string{},
			ChatID:        "c40ad2ef-d818-45b9-a006-60409f9a0b08",
			ChatCreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			Sources:       sources,
			ChatResponse:  strings.Join(summaries, "\n"),
			Keywords:      kw_map,
		},
	}

	c.JSON(http.StatusOK, chatMessage)

}
