package proto_test

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	"github.com/pact-foundation/pact-go/types"

	"github.com/owncloud/ocis-hello/pkg/proto/v0"
	svc "github.com/owncloud/ocis-hello/pkg/service/v0"

	"github.com/go-chi/chi"
)

const PORT = "9150"

func startServer() {
	r := chi.NewRouter()
	proto.RegisterHelloWeb(r, svc.NewService())
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}

func TestProviderContract(t *testing.T) {
	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		LogLevel: "DEBUG",
		Provider: "OCIS-Hello_API",
	}

	// Start provider API in the background
	go startServer()

	pactDir := "../../pact/pacts"

	// Verify the Provider using the locally saved Pact Files
	_, err := pact.VerifyProvider(t, types.VerifyRequest{
		ProviderBaseURL: "http://localhost:" + PORT,
		PactURLs:        []string{filepath.ToSlash(fmt.Sprintf("%s/ocis-hello_ui-ocis-hello_api.json", pactDir))},
		FailIfNoPactsFound: true,
	})

	if err != nil {
		log.Fatalf("%v", err)
	}
}
