package processor

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/progotman/multiplexer/requester"
)

type UrlsProcessor struct {
	Requester                   requester.Interface
	MaxNumberOfUrls             int
	MaxNumberOfOutgoingRequests int
	mu                          sync.Mutex
}

func (s *UrlsProcessor) Process(ctx context.Context, r ProcessUrlsRequest) *ProcessUrlsResult {
	if len(r.Urls) > s.MaxNumberOfUrls && s.MaxNumberOfUrls > 0 {
		tooManyURLsError := TooManyURLsError{MaxAvailableUrls: s.MaxNumberOfUrls, GivenUrls: len(r.Urls)}
		return &ProcessUrlsResult{Error: tooManyURLsError.Error()}
	}

	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	guard := make(chan struct{}, s.MaxNumberOfOutgoingRequests)

	result := &ProcessUrlsResult{
		Urls: make([]UrlResult, 0, len(r.Urls)),
	}
	for _, url := range r.Urls {
		guard <- struct{}{}
		wg.Add(1)
		go func(url string, result *ProcessUrlsResult) {
			defer func() {
				wg.Done()
				<-guard
			}()
			content, err := s.Requester.GetContent(ctx, url)

			select {
			case <-ctx.Done():
				return
			default:
				s.mu.Lock()
				defer s.mu.Unlock()
				if err != nil {
					result.Error = UrlProcessingError{Url: url}.Error()
					log.Println("Error processing url " + url + ". Error: " + err.Error())
					cancel()
					return
				}
				result.AddUrlResult(UrlResult{Url: url, Content: content})
			}
		}(url, result)
	}
	wg.Wait()

	if result.Error != "" {
		result.Urls = nil
	}

	cancel()
	return result
}

type ProcessUrlsRequest struct {
	Urls []string
}

type ProcessUrlsResult struct {
	Urls  []UrlResult `json:"urls,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (c *ProcessUrlsResult) AddUrlResult(result UrlResult) {
	c.Urls = append(c.Urls, result)
}

type UrlResult struct {
	Url     string `json:"url"`
	Content string `json:"content"`
}

type TooManyURLsError struct {
	MaxAvailableUrls int
	GivenUrls        int
}

func (e TooManyURLsError) Error() string {
	return fmt.Sprintf("Too many urls. Maximum available %d urls. Given %d", e.MaxAvailableUrls, e.GivenUrls)
}

type UrlProcessingError struct {
	Url string
}

func (e UrlProcessingError) Error() string {
	return "Error processing url: " + e.Url
}
