// services/interfaces/http_client.go
package interfaces

// HTTPClient defines an interface for making HTTP requests.
// You can customize the method signatures based on your application's needs.
// For instance, if you want to return the raw response body, or unmarshal directly into a target struct.
type HTTPClient interface {
	// DoRequest executes an HTTP request with the given method, endpoint, and optional body,
	// then unmarshals the result into target.
	DoRequest(method, endpoint string, body interface{}, target interface{}) error
}
