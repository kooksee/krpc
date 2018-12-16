package example

import "github.com/kooksee/krpc/client"

type _client struct {
	cl krpcc.HTTPClient
}

func (t *_client) Echo(val string) (string, error) {
	params := map[string]interface{}{
		"arg": val,
	}
	result := new(ResultEcho)
	if _, err := t.cl.Call("echo", params, result); err != nil {
		return "", err
	}
	return result.Value, nil
}

func echoIntViaHTTP(cl krpcc.HTTPClient, val int) (int, error) {
	params := map[string]interface{}{
		"arg": val,
	}
	result := new(ResultEchoInt)
	if _, err := cl.Call("echo_int", params, result); err != nil {
		return 0, err
	}
	return result.Value, nil
}

func echoBytesViaHTTP(cl krpcc.HTTPClient, bytes []byte) ([]byte, error) {
	params := map[string]interface{}{
		"arg": bytes,
	}
	result := new(ResultEchoBytes)
	if _, err := cl.Call("echo_bytes", params, result); err != nil {
		return []byte{}, err
	}
	return result.Value, nil
}

func echoDataBytesViaHTTP(cl krpcc.HTTPClient, bytes []byte) ([]byte, error) {
	params := map[string]interface{}{
		"arg": bytes,
	}
	result := new(ResultEchoDataBytes)
	if _, err := cl.Call("echo_data_bytes", params, result); err != nil {
		return []byte{}, err
	}
	return result.Value, nil
}
