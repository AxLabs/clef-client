package clefclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// rpcClient represents a client to interact with the clef JSON-RPC interface.
type rpcClient struct {
	url string
}

// newRPCClient creates a new clef JSON-RPC client.
func newRPCClient(url string) *rpcClient {
	return &rpcClient{url: url}
}

// rpcRequest represents a JSON-RPC request.
type rpcRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

// rpcResponse represents a JSON-RPC response.
type rpcResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error"`
	ID      int             `json:"id"`
}

// rpcError represents a JSON-RPC error.
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// call sends a JSON-RPC request and returns the response.
func (c *rpcClient) call(method string, params interface{}) (*rpcResponse, error) {
	reqBody, err := json.Marshal(rpcRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.url, "application/json", bytes.NewBuffer(reqBody))
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

// ClefClient represents a higher-level client to interact with clef.
type ClefClient struct {
	transport transport
}

// NewHTTPClient creates a new ClefClient using HTTP transport
func NewHTTPClient(url string) *ClefClient {
	return &ClefClient{transport: newHTTPTransport(url)}
}

// NewIPCClient creates a new ClefClient using IPC transport
func NewIPCClient(socketPath string) (*ClefClient, error) {
	transport, err := newIPCTransport(socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPC transport: %w", err)
	}
	return &ClefClient{transport: transport}, nil
}

// Close closes the underlying transport
func (cc *ClefClient) Close() error {
	return cc.transport.close()
}

// NewAccount creates a new account
func (cc *ClefClient) NewAccount() (string, error) {
	resp, err := cc.transport.call("account_new", nil)
	if err != nil {
		return "", err
	}

	var address string
	if err := json.Unmarshal(resp.Result, &address); err != nil {
		return "", err
	}
	return address, nil
}

// ListAccounts returns the list of available accounts
func (cc *ClefClient) ListAccounts() ([]string, error) {
	resp, err := cc.transport.call("account_list", nil)
	if err != nil {
		return nil, err
	}

	var accounts []string
	if err := json.Unmarshal(resp.Result, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// SignTransaction signs the given transaction
func (cc *ClefClient) SignTransaction(tx *Transaction) (*SignTxResponse, error) {
	resp, err := cc.transport.call("account_signTransaction", tx)
	if err != nil {
		return nil, err
	}

	var result SignTxResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SignData signs the given data
func (cc *ClefClient) SignData(req *SignDataRequest) (*SignDataResponse, error) {
	resp, err := cc.transport.call("account_signData", req)
	if err != nil {
		return nil, err
	}

	var result SignDataResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SignTypedData signs the given typed data
func (cc *ClefClient) SignTypedData(req *TypedDataRequest) (*SignDataResponse, error) {
	resp, err := cc.transport.call("account_signTypedData", req)
	if err != nil {
		return nil, err
	}

	var result SignDataResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EcRecover recovers the address from the given signature
func (cc *ClefClient) EcRecover(req *EcRecoverRequest) (*EcRecoverResponse, error) {
	resp, err := cc.transport.call("account_ecRecover", req)
	if err != nil {
		return nil, err
	}

	var result EcRecoverResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Version returns the version of the clef service
func (cc *ClefClient) Version() (*VersionResponse, error) {
	resp, err := cc.transport.call("account_version", nil)
	if err != nil {
		return nil, err
	}

	var result VersionResponse
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
