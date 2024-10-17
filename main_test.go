
package main
import (
	"errors"
	"fmt"
	"log"

	"strconv"
	"strings"
	"testing"
	"time"

	as "github.com/aerospike/aerospike-client-go/v6"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
)

type testSuite struct {
	pool      *dockertest.Pool
	aerospike *aerospike
}

type aerospike struct {
	Host           string
	Port           string
	UserName       string
	UserPassword   string
	dockerResource *dockertest.Resource
}

func newTestSuite() *testSuite {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	pool.MaxWait = time.Minute * 2

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not create docker test pool : %s", err)
	}
	return &testSuite{
		pool: pool,
	}
}

func (suite *testSuite) AerospikeUp() (*aerospike, error) {
	if suite.aerospike != nil {
		return suite.aerospike, nil
	}
	// auths, err := docker.NewAuthConfigurationsFromDockerCfg()
	// if err != nil {
	// 	return nil, fmt.Errorf("Could not find config from docker cfg %s", err)
	// }

	resource, err := suite.pool.RunWithOptions(&dockertest.RunOptions{
		Hostname: "aerospike",
		Repository: "aerospike/aerospike-server",
		Tag:        "6.1.0.1",
		Env:        []string{
			"AS_AUTH_USER=admin",
			"AS_AUTH_MODE=internal",
			"AS_AUTH_PASSWORD=admin",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, fmt.Errorf("Could not start resource: %s", err)
	}

	
	host, port := getHostPort(resource, "3000/tcp")
	log.Println("host:", host, "port:", port)
	log.Println(resource.GetBoundIP("3000/tcp"))
	log.Println("ip addr",resource.Container.NetworkSettings)
	if err = suite.pool.Retry(func() error {
		portInt, _ := strconv.Atoi(port)
		return pingToAerospike("admin", "admin", host, portInt)
	}); err != nil {
		suite.pool.Purge(resource)
		return nil, fmt.Errorf("Could not connect to docker: %s", err)
	}

	suite.aerospike = &aerospike{
		Host:           "localhost",
		Port:           port,
		UserName:       "admin",
		UserPassword:   "secret",
		dockerResource: resource,
	}
	return suite.aerospike, nil
}

func (suite *testSuite) AerospikeDown() error {
	if suite.aerospike == nil {
		return nil
	}
	return suite.pool.Purge(suite.aerospike.dockerResource)
}

func getHostPort(resource *dockertest.Resource, id string) (string, string) {
	return strings.Split(resource.GetHostPort(id), ":")[0], strings.Split(resource.GetHostPort(id), ":")[1]
}

func pingToAerospike(user string, password string, host string, port int) error {
	cp := as.NewClientPolicy()
	cp.User = user
	cp.Password = password
	cp.Timeout = 5 * time.Minute

	client , err := as.NewClient(host,port)
	// asHost := as.NewHost(host, port)
	// client, err := as.CreateClientWithPolicyAndHost(as.CTNative, cp, asHost)
	if err != nil {
		return err
	}
	if client.IsConnected() {
		return nil
	}
	return errors.New("not connected")
}

func TestIntegrationServing(t *testing.T) {
	fmt.Println(" integration test of serving")
	testSuite := newTestSuite()
	_, err := testSuite.AerospikeUp()
	log.Println("Aerospike was up")
	testSuite.AerospikeDown()
	assert.Nil(t, err)
}
