package provider

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

func newRetryableClient(retryOn []int) *http.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 5
	client.CheckRetry = retryPolicy(retryOn)
	client.Logger = nil
	return client.StandardClient()
}

func retryPolicy(retryOn []int) retryablehttp.CheckRetry {
	retrySet := make(map[int]bool, len(retryOn))
	for _, code := range retryOn {
		retrySet[code] = true
	}

	return func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		if err != nil {
			return true, nil
		}
		if resp.StatusCode == 429 {
			return true, nil
		}
		if resp.StatusCode == 403 {
			return isQuotaError(resp), nil
		}
		if resp.StatusCode >= 500 && resp.StatusCode != 501 {
			return true, nil
		}
		if retrySet[resp.StatusCode] {
			return true, nil
		}
		return false, nil
	}
}

func isQuotaError(resp *http.Response) bool {
	if resp.Body == nil {
		return false
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return false
	}
	s := strings.ToLower(string(body))
	return strings.Contains(s, "quotaexceeded") ||
		strings.Contains(s, "ratelimitexceeded") ||
		strings.Contains(s, "userratelimitexceeded") ||
		strings.Contains(s, "rate limit")
}
