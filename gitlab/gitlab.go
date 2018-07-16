package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client holds the information needed for Gitlab API authentication.
type Client struct {
	token string
	url   string
}

// New creates a new AfterShip API client.
func New(url, token string) *Client {
	return &Client{
		url:   url,
		token: token,
	}
}

func (c *Client) doRequest(method, endpoint string, data interface{}) (*Data, error) {
	client := http.DefaultClient

	// Encode data if we are passed an object.
	b := bytes.NewBuffer(nil)
	if data != nil {
		// Create the encoder.
		enc := json.NewEncoder(b)
		if err := enc.Encode(data); err != nil {
			return nil, fmt.Errorf("json encoding data for doRequest failed: %v", err)
		}
	}

	// Create the request.
	uri := fmt.Sprintf("%s/%s", c.url, strings.Trim(endpoint, "/"))
	req, err := http.NewRequest(method, uri, b)
	if err != nil {
		return nil, fmt.Errorf("creating %s request to %s failed: %v", method, uri, err)
	}

	// Set the correct headers.
	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Do the request.
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("performing %s request to %s failed: %v", method, uri, err)
	}
	defer resp.Body.Close()

	// Check that the response status code was OK.
	if resp.StatusCode > 400 {
		// Read the body of the request, ignore the error since we are already in the error state.
		body, _ := ioutil.ReadAll(resp.Body)

		// Create a friendly error message based off the status code returned.
		// These come from: https://docs.aftership.com/api/4/errors
		var message string
		switch resp.StatusCode {
		case http.StatusUnauthorized: // 401
			message = "Invalid Token."
		case http.StatusForbidden: // 403
			message = "The request is understood, but it has been refused or access is not allowed."
		case http.StatusNotFound: // 404
			message = "The URI requested is invalid or the resource requested does not exist."
		case http.StatusTooManyRequests: // 429
			message = "You have exceeded the API call rate limit."
		case http.StatusInternalServerError: // 500
			message = "Something went wrong on Gitlab's end. (500: Internal error)"
		case http.StatusNotImplemented: // 501
			message = "Something went wrong on Gitlab's end. (501: Not implemented)"
		case http.StatusBadGateway: // 502
			message = "Something went wrong on Gitlab's end. (502: Bad gateway)"
		case http.StatusServiceUnavailable: // 503
			message = "Something went wrong on Gitlab's end. (503: Service unavailable)"
		}

		return nil, fmt.Errorf("%s request to %s returned status code %d: message -> %s\nbody -> %s", method, uri, resp.StatusCode, message, string(body))
	}

	// Decode the response into a AfterShip response object.
	var r Response
	if err := decodeResponse(resp, &r); err != nil {
		// Read the body of the request, ignore the error since we are already in the error state.
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("decoding response from %s request to %s failed: body -> %s\nerr -> %v", method, uri, string(body), err)
	}

	// Return errors on the API errors.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return &r.Data, fmt.Errorf("API error %s [%d]: %s", r.Meta.Type, r.Meta.Code, r.Meta.Message)
	}

	return &r.Data, nil
}
