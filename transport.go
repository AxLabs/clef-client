package clefclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
)

// transport defines the interface for different transport mechanisms
type transport interface {
	call(method string, params interface{}) (*rpcResponse, error)
	close() error
}

// httpTransport implements transport interface for HTTP JSON-RPC
type httpTransport struct {
	url string
}

func newHTTPTransport(url string) *httpTransport {
	return &httpTransport{url: url}
}

func (t *httpTransport) call(method string, params interface{}) (*rpcResponse, error) {
	reqBody, err := json.Marshal(rpcRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(t.url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rpcResp rpcResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return nil, err
	}

	if rpcResp.Error != nil {
		return nil, errors.New(rpcResp.Error.Message)
	}

	return &rpcResp, nil
}

func (t *httpTransport) close() error {
	return nil // HTTP transport doesn't need explicit cleanup
}

// ipcTransport implements transport interface for IPC
type ipcTransport struct {
	conn net.Conn
}

func newIPCTransport(socketPath string) (*ipcTransport, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}
	return &ipcTransport{conn: conn}, nil
}

func (t *ipcTransport) call(method string, params interface{}) (*rpcResponse, error) {
	reqBody, err := json.Marshal(rpcRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})
	if err != nil {
		return nil, err
	}

	_, err = t.conn.Write(append(reqBody, '\n'))
	if err != nil {
		return nil, err
	}

	var rpcResp rpcResponse
	if err := json.NewDecoder(t.conn).Decode(&rpcResp); err != nil {
		return nil, err
	}

	if rpcResp.Error != nil {
		return nil, errors.New(rpcResp.Error.Message)
	}

	return &rpcResp, nil
}

func (t *ipcTransport) close() error {
	return t.conn.Close()
}
