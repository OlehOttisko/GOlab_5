package integration

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

func sendRequest(baseAddress string, responseSize int, client *http.Client) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/some-data?key=team", baseAddress), nil)
	if err != nil {
		log.Printf("error creating request: %s", err)
		return nil, err
	}
	req.Header.Set("Response-Size", strconv.Itoa(responseSize))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error: %s", err)
		return nil, err
	}

	log.Printf("response %d", resp.StatusCode)
	return resp, nil
}

func TestIntegration(t *testing.T) { TestingT(t) }

type IntegrationSuite struct{}

var _ = Suite(&IntegrationSuite{})

func (s *IntegrationSuite) TestBalancer(c *C) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		c.Skip("Integration test is not enabled")
	}

	responseSize := 0
	serverNum := [6]int{1, 2, 3, 1, 3, 2}
	for i := 0; i < 6; i++ {
		if i%2 == 0 {
			responseSize = 1000
		} else {
			responseSize = 2000
		}
		server, _ := sendRequest(baseAddress, responseSize, &client)
		c.Assert(server.Header.Get("lb-from"), Equals, fmt.Sprintf("server%d:8080", serverNum[i]))

	}
}

func (s *IntegrationSuite) BenchmarkBalancer(c *C) {
	if _, exists := os.LookupEnv("INTEGRATION_TEST"); !exists {
		c.Skip("Integration test is not enabled")
	}

	for i := 0; i < c.N; i++ {
		resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data?key=team", baseAddress))
		if err != nil {
			log.Printf("error: %s", err)
			c.FailNow()
		}
		c.Assert(resp.StatusCode, Equals, http.StatusOK)
	}
}
