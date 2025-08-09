package orders_agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/pkg/utils/delay"
	"io"
	"net/http"
	"strconv"
	"time"
)

const acceptedOrderURL = "/api/orders/"

type StatusUpdaterClient struct {
	*http.Client
	url string
}

type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries uint
}

func NewClientAgent(url string) *StatusUpdaterClient {
	return &StatusUpdaterClient{
		Client: &http.Client{
			Transport: &retryRoundTripper{
				maxRetries: 3,
				next:       http.DefaultTransport,
			},
		},
		url: url,
	}
}

func (rr retryRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	var seconds time.Duration
	newDelay := delay.NewDelay()
	for attempts := 0; attempts < int(rr.maxRetries); attempts++ {
		resp, err = rr.next.RoundTrip(r)
		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			retryAfter := resp.Header.Get("Retry-After")
			if retryAfter != "" {
				retryAfterSeconds, err := strconv.Atoi(retryAfter)
				if err != nil {
					fmt.Println("Could not parse Retry-After header.")
					seconds = newDelay()
				}
				seconds = time.Duration(retryAfterSeconds) * time.Second
			}
		default:
			seconds = newDelay()
		}

		select {
		case <-r.Context().Done():

			return resp, r.Context().Err()
		case <-time.After(seconds):
		}
	}
	return resp, err
}

func (s *StatusUpdaterClient) getOrderStatus(orderNumber string) models.OrderResult {
	var result models.OrderResult
	var order *models.Order
	url := fmt.Sprintf("%s%s%s", s.url, acceptedOrderURL, orderNumber)
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		result.Err = err
		return result
	}

	response, err := s.Do(req)
	if err != nil {
		result.Err = err
		return result
	}

	buf, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		result.Err = err
		return result
	}

	if err = json.Unmarshal(buf, &order); err != nil {
		result.Err = err
		return result
	}
	result.Order = order
	return result
}
