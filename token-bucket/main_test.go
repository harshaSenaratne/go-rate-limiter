package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSingleAPICall(t *testing.T) {
	// Test scenario: curl -i http://localhost:8080/ping
	handler := rateLimiter(endpointHandler)
	
	req := httptest.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()
	
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK {
		t.Errorf("Single API call should succeed: got %d, want %d", rr.Code, http.StatusOK)
	}
	
	var response Message
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}
	
	if response.Status != "Successful" {
		t.Errorf("Expected 'Successful', got '%s'", response.Status)
	}
}

func TestMultipleAPICalls(t *testing.T) {
	// Test scenario: for i in {1..6}; do curl http://localhost:8080/ping; done
	handler := rateLimiter(endpointHandler)
	
	successCount := 0
	rateLimitedCount := 0
	
	// Make 6 rapid requests
	for i := 1; i <= 6; i++ {
		req := httptest.NewRequest("GET", "/ping", nil)
		rr := httptest.NewRecorder()
		
		handler.ServeHTTP(rr, req)
		
		if rr.Code == http.StatusOK {
			successCount++
		} else if rr.Code == http.StatusTooManyRequests {
			rateLimitedCount++
			
			// Verify rate limit response
			var response Message
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Errorf("Could not unmarshal rate limit response: %v", err)
			}
			
			if response.Status != "Request Failed" {
				t.Errorf("Rate limited response should have 'Request Failed' status, got '%s'", response.Status)
			}
			
			expectedBody := "The API is at capacity, try again later."
			if response.Body != expectedBody {
				t.Errorf("Rate limited response body mismatch: got '%s', want '%s'", response.Body, expectedBody)
			}
		}
	}
	
	// With rate limit of 2 per second and burst of 4, first 4 should succeed
	if successCount != 4 {
		t.Errorf("Expected 4 successful requests, got %d", successCount)
	}
	
	// The remaining 2 should be rate limited
	if rateLimitedCount != 2 {
		t.Errorf("Expected 2 rate limited requests, got %d", rateLimitedCount)
	}
	
	t.Logf("Results: %d successful, %d rate limited", successCount, rateLimitedCount)
}

func TestTokenBucketRecovery(t *testing.T) {
	handler := rateLimiter(endpointHandler)
	
	// Exhaust the burst limit (4 tokens)
	for i := 1; i <= 4; i++ {
		req := httptest.NewRequest("GET", "/ping", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusOK {
			t.Errorf("Request %d should succeed: got %d", i, rr.Code)
		}
	}
	
	// Next request should be rate limited
	req := httptest.NewRequest("GET", "/ping", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Request should be rate limited: got %d", rr.Code)
	}
	
	// Wait for tokens to be added back (2 per second rate)
	// Wait 600ms to get at least 1 token back
	time.Sleep(600 * time.Millisecond)
	
	req2 := httptest.NewRequest("GET", "/ping", nil)
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	
	if rr2.Code != http.StatusOK {
		t.Errorf("Request after recovery should succeed: got %d", rr2.Code)
	}
}