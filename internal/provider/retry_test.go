package provider

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestRetryPolicy_429_Retries(t *testing.T) {
	policy := retryPolicy(nil)
	resp := &http.Response{StatusCode: 429}
	retry, err := policy(context.Background(), resp, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !retry {
		t.Error("expected retry on 429")
	}
}

func TestRetryPolicy_500_Retries(t *testing.T) {
	policy := retryPolicy(nil)
	resp := &http.Response{StatusCode: 500}
	retry, err := policy(context.Background(), resp, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !retry {
		t.Error("expected retry on 500")
	}
}

func TestRetryPolicy_502_Retries(t *testing.T) {
	policy := retryPolicy(nil)
	resp := &http.Response{StatusCode: 502}
	retry, err := policy(context.Background(), resp, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !retry {
		t.Error("expected retry on 502")
	}
}

func TestRetryPolicy_501_DoesNotRetry(t *testing.T) {
	policy := retryPolicy(nil)
	resp := &http.Response{StatusCode: 501}
	retry, err := policy(context.Background(), resp, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if retry {
		t.Error("expected no retry on 501")
	}
}

func TestRetryPolicy_200_DoesNotRetry(t *testing.T) {
	policy := retryPolicy(nil)
	resp := &http.Response{StatusCode: 200}
	retry, err := policy(context.Background(), resp, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if retry {
		t.Error("expected no retry on 200")
	}
}

func TestRetryPolicy_403_QuotaError_Retries(t *testing.T) {
	policy := retryPolicy(nil)
	body := `{"error":{"errors":[{"domain":"usageLimits","reason":"rateLimitExceeded"}]}}`
	resp := &http.Response{
		StatusCode: 403,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	retry, err := policy(context.Background(), resp, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !retry {
		t.Error("expected retry on 403 with quota error")
	}
}

func TestRetryPolicy_403_NonQuota_DoesNotRetry(t *testing.T) {
	policy := retryPolicy(nil)
	body := `{"error":{"errors":[{"domain":"global","reason":"forbidden"}]}}`
	resp := &http.Response{
		StatusCode: 403,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
	retry, err := policy(context.Background(), resp, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if retry {
		t.Error("expected no retry on 403 without quota error")
	}
}

func TestRetryPolicy_CustomCodes_Retries(t *testing.T) {
	policy := retryPolicy([]int{404, 502})
	resp := &http.Response{StatusCode: 404}
	retry, err := policy(context.Background(), resp, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !retry {
		t.Error("expected retry on custom 404 code")
	}
}

func TestRetryPolicy_ConnectionError_Retries(t *testing.T) {
	policy := retryPolicy(nil)
	retry, err := policy(context.Background(), nil, io.ErrUnexpectedEOF)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !retry {
		t.Error("expected retry on connection error")
	}
}

func TestRetryPolicy_CancelledContext_DoesNotRetry(t *testing.T) {
	policy := retryPolicy(nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	retry, err := policy(ctx, nil, nil)
	if err == nil {
		t.Fatal("expected context error")
	}
	if retry {
		t.Error("expected no retry on cancelled context")
	}
}
