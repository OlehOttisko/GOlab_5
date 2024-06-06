package main

import (
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type BalancerSuite struct{}

var _ = Suite(&BalancerSuite{})

func (s *BalancerSuite) TestBalancer(c *C) {
	healthChecker := &HealthChecker{}
	healthChecker.healthyServers = []string{"server1:8080", "server2:8080", "server3:8080"}

	balancer := &Balancer{}
	balancer.healthChecker = healthChecker

	index1 := balancer.getServerIndexWithLowestLoad(map[string]int64{
		"server1:8080": 100,
		"server2:8080": 200,
		"server3:8080": 150,
	}, []string{"server1:8080", "server2:8080", "server3:8080"})

	index2 := balancer.getServerIndexWithLowestLoad(map[string]int64{
		"server1:8080": 300,
		"server2:8080": 200,
		"server3:8080": 250,
	}, []string{"server1:8080", "server2:8080", "server3:8080"})

	index3 := balancer.getServerIndexWithLowestLoad(map[string]int64{
		"server1:8080": 200,
		"server2:8080": 150,
		"server3:8080": 100,
	}, []string{"server1:8080", "server2:8080", "server3:8080"})

	c.Assert(index1, Equals, 0)
	c.Assert(index2, Equals, 1)
	c.Assert(index3, Equals, 2)
}

func (s *BalancerSuite) TestHealthChecker(c *C) {
	healthChecker := &HealthChecker{}
	healthChecker.health = func(s string) bool {
		if s == "1" {
			return false
		} else {
			return true
		}
	}

	healthChecker.serversPool = []string{"1", "2", "3"}
	healthChecker.healthyServers = []string{"4", "5", "6"}
	healthChecker.checkInterval = 1 * time.Second

	healthChecker.StartHealthCheck()

	time.Sleep(2 * time.Second)

	c.Assert(healthChecker.GetHealthyServers()[0], Equals, "2")
	c.Assert(healthChecker.GetHealthyServers()[1], Equals, "3")
	c.Assert(len(healthChecker.GetHealthyServers()), Equals, 2)
}
