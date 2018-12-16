package example

import (
	"bytes"
	crand "crypto/rand"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/go-amino"

	client "github.com/kooksee/krpc/client"
	server "github.com/kooksee/krpc/server"
	types "github.com/kooksee/krpc/types"
	"math/rand"
)

// Client and Server should work over tcp or unix sockets
const (
	tcpAddr = "tcp://0.0.0.0:47768"

	unixSocket = "/tmp/rpc_test.sock"
	unixAddr   = "unix://" + unixSocket
)

type ResultEcho struct {
	Value string `json:"value"`
}

type ResultEchoInt struct {
	Value int `json:"value"`
}

type ResultEchoBytes struct {
	Value []byte `json:"value"`
}

type ResultEchoDataBytes struct {
	Value []byte `json:"value"`
}

// Define some routes
var Routes = map[string]*server.RPCFunc{
	"echo":            server.NewRPCFunc(EchoResult, "arg"),
	"echo_bytes":      server.NewRPCFunc(EchoBytesResult, "arg"),
	"echo_data_bytes": server.NewRPCFunc(EchoDataBytesResult, "arg"),
	"echo_int":        server.NewRPCFunc(EchoIntResult, "arg"),
}

// Amino codec required to encode/decode everything above.
var RoutesCdc = amino.NewCodec()

func EchoResult(v string) (*ResultEcho, error) {
	return &ResultEcho{v}, nil
}

func EchoWSResult(wsCtx types.WSRPCContext, v string) (*ResultEcho, error) {
	return &ResultEcho{v}, nil
}

func EchoIntResult(v int) (*ResultEchoInt, error) {
	return &ResultEchoInt{v}, nil
}

func EchoBytesResult(v []byte) (*ResultEchoBytes, error) {
	return &ResultEchoBytes{v}, nil
}

func EchoDataBytesResult(v []byte) (*ResultEchoDataBytes, error) {
	return &ResultEchoDataBytes{v}, nil
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

// launch unix and tcp servers
func setup() {

	cmd := exec.Command("rm", "-f", unixSocket)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	if err = cmd.Wait(); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	server.RegisterRPCFuncs(mux, Routes, RoutesCdc)
	go func() {
		server.StartHTTPServer(tcpAddr, mux, server.Config{})
	}()

	mux2 := http.NewServeMux()
	server.RegisterRPCFuncs(mux2, Routes, RoutesCdc)
	go func() {
		server.StartHTTPServer(unixAddr, mux2, server.Config{})
	}()

	// wait for servers to start
	time.Sleep(time.Second * 2)
}

func echoViaHTTP(cl client.HTTPClient, val string) (string, error) {
	params := map[string]interface{}{
		"arg": val,
	}
	result := new(ResultEcho)
	if _, err := cl.Call("echo", params, result); err != nil {
		return "", err
	}
	return result.Value, nil
}

func echoIntViaHTTP(cl client.HTTPClient, val int) (int, error) {
	params := map[string]interface{}{
		"arg": val,
	}
	result := new(ResultEchoInt)
	if _, err := cl.Call("echo_int", params, result); err != nil {
		return 0, err
	}
	return result.Value, nil
}

func echoBytesViaHTTP(cl client.HTTPClient, bytes []byte) ([]byte, error) {
	params := map[string]interface{}{
		"arg": bytes,
	}
	result := new(ResultEchoBytes)
	if _, err := cl.Call("echo_bytes", params, result); err != nil {
		return []byte{}, err
	}
	return result.Value, nil
}

func echoDataBytesViaHTTP(cl client.HTTPClient, bytes []byte) ([]byte, error) {
	params := map[string]interface{}{
		"arg": bytes,
	}
	result := new(ResultEchoDataBytes)
	if _, err := cl.Call("echo_data_bytes", params, result); err != nil {
		return []byte{}, err
	}
	return result.Value, nil
}

func testWithHTTPClient(t *testing.T, cl client.HTTPClient) {
	val := "acbd"
	got, err := echoViaHTTP(cl, val)
	require.Nil(t, err)
	assert.Equal(t, got, val)

	val2 := randBytes(t)
	got2, err := echoBytesViaHTTP(cl, val2)
	require.Nil(t, err)
	assert.Equal(t, got2, val2)

	val3 := []byte(randBytes(t))
	got3, err := echoDataBytesViaHTTP(cl, val3)
	require.Nil(t, err)
	assert.Equal(t, got3, val3)

	rand.Seed(time.Now().UnixNano())
	val4 := rand.Int31n(10000)
	got4, err := echoIntViaHTTP(cl, int(val4))
	require.Nil(t, err)
	assert.Equal(t, got4, val4)
}

//-------------

func randBytes(t *testing.T) []byte {
	rand.Seed(time.Now().UnixNano())

	n := rand.Int31n(10) + 2
	buf := make([]byte, n)
	_, err := crand.Read(buf)
	require.Nil(t, err)
	return bytes.Replace(buf, []byte("="), []byte{100}, -1)
}
