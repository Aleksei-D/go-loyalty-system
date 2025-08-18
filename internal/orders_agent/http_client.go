package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/delay"
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
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusTooManyRequests:
				retryAfter := resp.Header.Get("Retry-After")
				if retryAfter != "" {
					retryAfterSeconds, err := strconv.Atoi(retryAfter)
					if err != nil {
						logger.Log.Warn("Could not parse Retry-After header.")
						seconds = newDelay()
					} else {
						seconds = time.Duration(retryAfterSeconds) * time.Second
					}

				}
			default:
				seconds = newDelay()
			}
		}

		select {
		case <-r.Context().Done():

			return resp, r.Context().Err()
		case <-time.After(seconds):
		}
	}
	return resp, err
}

func (s *StatusUpdaterClient) getOrderStatus(orderNumber string) (*models.OrderStatusResponse, error) {
	var order models.OrderStatusResponse
	url := fmt.Sprintf("%s%s%s", s.url, acceptedOrderURL, orderNumber)
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		return &order, err
	}

	response, err := s.Do(req)
	if err != nil {
		return &order, err
	}

	buf, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return &order, err
	}

	if err = json.Unmarshal(buf, &order); err != nil {
		return &order, err
	}
	return &order, nil
}
