package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiUrls(t *testing.T) {
	// No environment variables set
	os.Unsetenv("BEACON_API_URLS")
	os.Unsetenv("ETH_API_URLS")

	apiUrls := NewApiUrls()

	require.Len(t, apiUrls.beaconApiUrls, 3, "Expected 3 default beacon API URLs")
	require.Len(t, apiUrls.ethApiUrls, 5, "Expected 5 default ETH API URLs")

	// Environment variables set
	os.Setenv("BEACON_API_URLS", "http://example.com")
	os.Setenv("ETH_API_URLS", "http://example-eth.com,http://example-eth2.com")

	apiUrls = NewApiUrls()

	require.Len(t, apiUrls.beaconApiUrls, 1, "Expected 1 beacon API URL")
	require.Len(t, apiUrls.ethApiUrls, 2, "Expected 2 ETH API URLs")
}
