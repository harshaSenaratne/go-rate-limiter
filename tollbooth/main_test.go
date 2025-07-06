package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tollbooth "github.com/didip/tollbooth/v7"
)

func TestSingleAPICall(t *testing.T) {
	// Test scenario: curl -i http://localhost:8080/ping
	limiter := tollbooth.NewLimiter(1, nil)
	limiter.SetMessageContentType("application/json")
	
	message := Message{
		Status: "Request Failed",
		Body:   "The API is at capacity, try again later.",
	}
	jsonMessage, _ := json.Marshal(message)
	limiter.SetMessage(string(jsonMessage))
	
	handler := tollbooth.LimitFuncHandler(limiter, endpointHandler)
	
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
	limiter := tollbooth.NewLimiter(1, nil)
	limiter.SetMessageContentType("application/json")
	
	message := Message{
		Status: "Request Failed",
		Body:   "The API is at capacity, try again later.",
	}
	jsonMessage, _ := json.Marshal(message)
	limiter.SetMessage(string(jsonMessage))
	
	handler := tollbooth.LimitFuncHandler(limiter, endpointHandler)
	
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
	
	// With a rate limit of 1 per second, only 1 request should succeed
	if successCount != 1 {
		t.Errorf("Expected 1 successful request, got %d", successCount)
	}
	
	// The remaining 5 should be rate limited
	if rateLimitedCount != 5 {
		t.Errorf("Expected 5 rate limited requests, got %d", rateLimitedCount)
	}
	
	t.Logf("Results: %d successful, %d rate limited", successCount, rateLimitedCount)
}