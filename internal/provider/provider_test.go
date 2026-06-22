package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"googleworkspace": providerserver.NewProtocol6WithError(New("test")()),
}

// testProviderConfig is a minimal provider config that passes schema validation.
// The testAPIClient override bypasses all authentication, so these values are unused.
const testProviderConfig = `
provider "googleworkspace" {
  access_token            = "test-token"
  service_account         = "test@test.iam.gserviceaccount.com"
  impersonated_user_email = "admin@test.com"
  customer_id             = "C00000000"
}
`

// setupTestServer creates an httptest.Server with the given handler and registers
// cleanup. Returns the server.
func setupTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

// setupTestClient sets the package-level testAPIClient with an HTTP client pointing
// at the test server. Must be called before running terraform-plugin-testing steps.
func setupTestClient(t *testing.T, server *httptest.Server) {
	t.Helper()
	testAPIClient = &apiClient{
		client:     server.Client(),
		customerID: "C00000000",
		basePath:   server.URL,
	}
	t.Cleanup(func() { testAPIClient = nil })
}

// jsonResponse writes a JSON response with the given status code.
func jsonResponse(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}
