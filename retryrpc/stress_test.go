package retryrpc

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/swiftstack/ProxyFS/retryrpc/rpctest"
)

func TestStress(t *testing.T) {

	testLoop(t)
}

func testLoop(t *testing.T) {
	var (
		agentCount = 100
		sendCount  = 100
	)
	assert := assert.New(t)
	zero := 0
	assert.Equal(0, zero)

	// Create new rpctest server - needed for calling
	// RPCs
	myJrpcfs := rpctest.NewServer()

	rrSvr, ipAddr, port := getNewServer()
	assert.NotNil(rrSvr)

	// Register the Server - sets up the methods supported by the
	// server
	err := rrSvr.Register(myJrpcfs)
	assert.Nil(err)

	// Start listening for requests on the ipaddr/port
	startErr := rrSvr.Start()
	assert.Nil(startErr, "startErr is not nil")

	// Tell server to start accepting and processing requests
	rrSvr.Run()

	// Start up the agents
	parallelAgentSenders(t, ipAddr, port, agentCount, sendCount, rrSvr.Creds.RootCAx509CertificatePEM)
}

func sendIt(t *testing.T, client *Client, i int, sendWg *sync.WaitGroup) {
	assert := assert.New(t)
	defer sendWg.Done()

	// Send a ping RPC and print the results
	msg := fmt.Sprintf("Ping Me - %v", i)
	pingRequest := &rpctest.PingReq{Message: msg}
	pingReply := &rpctest.PingReply{}
	expectedReply := fmt.Sprintf("pong %d bytes", len(msg))
	err := client.Send("RpcPing", pingRequest, pingReply)
	assert.Nil(err, "client.Send() returned an error")
	if expectedReply != pingReply.Message {
		fmt.Printf("ERR ==== client: '%+v'\n", client)
		fmt.Printf("         client.Send(RpcPing) reply '%+v'\n", pingReply)
		fmt.Printf("         client.Send(RpcPing) expected '%s' but received '%s'\n", expectedReply, pingReply.Message)
		fmt.Printf("         client.Send(RpcPing) SENT: msg '%v' but received '%s'\n", msg, pingReply.Message)
		fmt.Printf("         client.Send(RpcPing) len(pingRequest.Message): '%d' i: %v\n", len(pingRequest.Message), i)
	}
	assert.Equal(expectedReply, pingReply.Message, "Received different output then expected")
}

// Represents a pfsagent - sepearate client
func pfsagent(t *testing.T, ipAddr string, port int, agentID uint64, agentWg *sync.WaitGroup,
	sendCnt int, rootCAx509CertificatePEM []byte) {
	defer agentWg.Done()

	clientID := fmt.Sprintf("client - %v", agentID)
	client, err := NewClient(clientID, ipAddr, port, rootCAx509CertificatePEM)
	if err != nil {
		fmt.Printf("Dial() failed with err: %v\n", err)
		return
	}
	defer client.Close()

	var sendWg sync.WaitGroup

	var z int
	for i := 0; i < sendCnt; i++ {
		z = (z + i) * 10

		sendWg.Add(1)
		go sendIt(t, client, z, &sendWg)
	}
	sendWg.Wait()
}

// Start a bunch of "pfsagents" in parallel
func parallelAgentSenders(t *testing.T, ipAddr string, port int, agentCnt int,
	sendCnt int, rootCAx509CertificatePEM []byte) {

	var agentWg sync.WaitGroup

	// Figure out random seed for runs
	r := rand.New(rand.NewSource(99))
	clientSeed := r.Uint64()

	// Start parallel pfsagents - each agent doing sendCnt parallel sends
	var agentID uint64
	for i := 0; i < agentCnt; i++ {
		agentID = clientSeed + uint64(i)

		agentWg.Add(1)
		go pfsagent(t, ipAddr, port, agentID, &agentWg, sendCnt, rootCAx509CertificatePEM)
	}
	agentWg.Wait()
}
