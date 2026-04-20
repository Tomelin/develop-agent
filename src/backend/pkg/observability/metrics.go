package observability

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type key struct {
	Method string
	Path   string
	Status string
}

type latencyKey struct {
	Method string
	Path   string
}

type llmKey struct {
	Provider string
	Model    string
	Status   string
}

type llmPair struct {
	Provider string
	Model    string
}

var (
	inFlight int64

	httpCountsMu sync.RWMutex
	httpCounts   = map[key]int64{}

	httpLatencyMu sync.RWMutex
	httpLatency   = map[latencyKey]float64{}

	llmCountMu sync.RWMutex
	llmCounts  = map[llmKey]int64{}

	llmPromptMu sync.RWMutex
	llmPrompt   = map[llmPair]int64{}

	llmCompletionMu sync.RWMutex
	llmCompletion   = map[llmPair]int64{}

	llmDurationMu sync.RWMutex
	llmDuration   = map[llmPair]float64{}
)

func HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		atomic.AddInt64(&inFlight, 1)
		start := time.Now()
		defer func() {
			atomic.AddInt64(&inFlight, -1)
			path := c.FullPath()
			if path == "" {
				path = c.Request.URL.Path
			}
			status := strconv.Itoa(c.Writer.Status())

			httpCountsMu.Lock()
			httpCounts[key{Method: c.Request.Method, Path: path, Status: status}]++
			httpCountsMu.Unlock()

			httpLatencyMu.Lock()
			httpLatency[latencyKey{Method: c.Request.Method, Path: path}] += time.Since(start).Seconds()
			httpLatencyMu.Unlock()
		}()
		c.Next()
	}
}

func ObserveLLMCall(provider, model string, duration time.Duration, promptTokens, completionTokens int, callErr error) {
	status := "ok"
	if callErr != nil {
		status = "error"
	}

	llmCountMu.Lock()
	llmCounts[llmKey{Provider: provider, Model: model, Status: status}]++
	llmCountMu.Unlock()

	pair := llmPair{Provider: provider, Model: model}
	llmPromptMu.Lock()
	llmPrompt[pair] += int64(promptTokens)
	llmPromptMu.Unlock()

	llmCompletionMu.Lock()
	llmCompletion[pair] += int64(completionTokens)
	llmCompletionMu.Unlock()

	llmDurationMu.Lock()
	llmDuration[pair] += duration.Seconds()
	llmDurationMu.Unlock()
}

func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; version=0.0.4", []byte(snapshot()))
	}
}

func snapshot() string {
	var b strings.Builder

	fmt.Fprintf(&b, "# HELP http_requests_in_flight Current number of in-flight HTTP requests.\n")
	fmt.Fprintf(&b, "# TYPE http_requests_in_flight gauge\n")
	fmt.Fprintf(&b, "http_requests_in_flight %d\n", atomic.LoadInt64(&inFlight))

	fmt.Fprintf(&b, "# HELP http_requests_total Total HTTP requests by method, path and status.\n")
	fmt.Fprintf(&b, "# TYPE http_requests_total counter\n")
	httpCountsMu.RLock()
	keys := make([]key, 0, len(httpCounts))
	for k := range httpCounts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j]) })
	for _, k := range keys {
		fmt.Fprintf(&b, "http_requests_total{method=%q,path=%q,status=%q} %d\n", k.Method, k.Path, k.Status, httpCounts[k])
	}
	httpCountsMu.RUnlock()

	fmt.Fprintf(&b, "# HELP http_request_duration_seconds_accum Accumulated HTTP request duration in seconds by method and path.\n")
	fmt.Fprintf(&b, "# TYPE http_request_duration_seconds_accum counter\n")
	httpLatencyMu.RLock()
	latKeys := make([]latencyKey, 0, len(httpLatency))
	for k := range httpLatency {
		latKeys = append(latKeys, k)
	}
	sort.Slice(latKeys, func(i, j int) bool { return fmt.Sprint(latKeys[i]) < fmt.Sprint(latKeys[j]) })
	for _, k := range latKeys {
		fmt.Fprintf(&b, "http_request_duration_seconds_accum{method=%q,path=%q} %.6f\n", k.Method, k.Path, httpLatency[k])
	}
	httpLatencyMu.RUnlock()

	fmt.Fprintf(&b, "# HELP llm_requests_total Total LLM requests by provider, model and status.\n")
	fmt.Fprintf(&b, "# TYPE llm_requests_total counter\n")
	llmCountMu.RLock()
	lk := make([]llmKey, 0, len(llmCounts))
	for k := range llmCounts {
		lk = append(lk, k)
	}
	sort.Slice(lk, func(i, j int) bool { return fmt.Sprint(lk[i]) < fmt.Sprint(lk[j]) })
	for _, k := range lk {
		fmt.Fprintf(&b, "llm_requests_total{provider=%q,model=%q,status=%q} %d\n", k.Provider, k.Model, k.Status, llmCounts[k])
	}
	llmCountMu.RUnlock()

	writeLLMSeries(&b, "llm_tokens_prompt_total", &llmPromptMu, llmPrompt)
	writeLLMSeries(&b, "llm_tokens_completion_total", &llmCompletionMu, llmCompletion)
	writeLLMDuration(&b)

	return b.String()
}

func writeLLMSeries(b *strings.Builder, metric string, mu *sync.RWMutex, values map[llmPair]int64) {
	fmt.Fprintf(b, "# HELP %s Total accumulated value by provider and model.\n", metric)
	fmt.Fprintf(b, "# TYPE %s counter\n", metric)
	mu.RLock()
	keys := make([]llmPair, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j]) })
	for _, k := range keys {
		fmt.Fprintf(b, "%s{provider=%q,model=%q} %d\n", metric, k.Provider, k.Model, values[k])
	}
	mu.RUnlock()
}

func writeLLMDuration(b *strings.Builder) {
	fmt.Fprintf(b, "# HELP llm_request_duration_seconds_accum Accumulated LLM request duration in seconds by provider and model.\n")
	fmt.Fprintf(b, "# TYPE llm_request_duration_seconds_accum counter\n")
	llmDurationMu.RLock()
	keys := make([]llmPair, 0, len(llmDuration))
	for k := range llmDuration {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return fmt.Sprint(keys[i]) < fmt.Sprint(keys[j]) })
	for _, k := range keys {
		fmt.Fprintf(b, "llm_request_duration_seconds_accum{provider=%q,model=%q} %.6f\n", k.Provider, k.Model, llmDuration[k])
	}
	llmDurationMu.RUnlock()
}
