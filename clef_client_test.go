package clefclient

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupHTTPTestServer(t *testing.T, expectedMethod string, response interface{}) (*ClefClient, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req rpcRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, expectedMethod, req.Method)

		resp := rpcResponse{
			Jsonrpc: "2.0",
			ID:      1,
		}

		if response != nil {
			resultBytes, err := json.Marshal(response)
			assert.NoError(t, err)
			resp.Result = resultBytes
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))

	client := NewHTTPClient(server.URL)
	return client, server
}

func setupIPCTestServer(t *testing.T, expectedMethod string, response interface{}) (*ClefClient, net.Listener, string) {
	tmpDir, err := os.MkdirTemp("", "clef-test")
	assert.NoError(t, err)

	socketPath := filepath.Join(tmpDir, "clef.ipc")
	listener, err := net.Listen("unix", socketPath)
	assert.NoError(t, err)

	// Handle connections in a goroutine
	go func() {
		conn, err := listener.Accept()
		assert.NoError(t, err)
		defer conn.Close()

		var req rpcRequest
		err = json.NewDecoder(conn).Decode(&req)
		assert.NoError(t, err)
		assert.Equal(t, expectedMethod, req.Method)

		resp := rpcResponse{
			Jsonrpc: "2.0",
			ID:      1,
		}

		if response != nil {
			resultBytes, err := json.Marshal(response)
			assert.NoError(t, err)
			resp.Result = resultBytes
		}

		err = json.NewEncoder(conn).Encode(resp)
		assert.NoError(t, err)
	}()

	client, err := NewIPCClient(socketPath)
	assert.NoError(t, err)

	return client, listener, tmpDir
}

func TestNewAccountHTTP(t *testing.T) {
	expectedAddress := "0x0000000000000000000000000000000000000001"
	client, server := setupHTTPTestServer(t, "account_new", expectedAddress)
	defer server.Close()

	address, err := client.NewAccount()
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, address)
}

func TestNewAccountIPC(t *testing.T) {
	expectedAddress := "0x0000000000000000000000000000000000000001"
	client, listener, tmpDir := setupIPCTestServer(t, "account_new", expectedAddress)
	defer listener.Close()
	defer os.RemoveAll(tmpDir)
	defer client.Close()

	address, err := client.NewAccount()
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, address)
}

func TestListAccountsHTTP(t *testing.T) {
	expectedAccounts := []string{
		"0x0000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000002",
	}
	client, server := setupHTTPTestServer(t, "account_list", expectedAccounts)
	defer server.Close()

	accounts, err := client.ListAccounts()
	assert.NoError(t, err)
	assert.Equal(t, expectedAccounts, accounts)
}

func TestListAccountsIPC(t *testing.T) {
	expectedAccounts := []string{
		"0x0000000000000000000000000000000000000001",
		"0x0000000000000000000000000000000000000002",
	}
	client, listener, tmpDir := setupIPCTestServer(t, "account_list", expectedAccounts)
	defer listener.Close()
	defer os.RemoveAll(tmpDir)
	defer client.Close()

	accounts, err := client.ListAccounts()
	assert.NoError(t, err)
	assert.Equal(t, expectedAccounts, accounts)
}

func TestSignTransactionHTTP(t *testing.T) {
	tx := &Transaction{
		From:     "0x0000000000000000000000000000000000000001",
		To:       "0x0000000000000000000000000000000000000002",
		Gas:      "0x5208",
		GasPrice: "0x4a817c800",
		Value:    "0xde0b6b3a7640000",
		Nonce:    "0x0",
		Data:     "0x",
	}

	expected := &SignTxResponse{
		Raw: "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675",
		Tx: struct {
			Nonce    string `json:"nonce"`
			GasPrice string `json:"gasPrice"`
			Gas      string `json:"gas"`
			To       string `json:"to"`
			Value    string `json:"value"`
			Input    string `json:"input"`
			V        string `json:"v"`
			R        string `json:"r"`
			S        string `json:"s"`
			Hash     string `json:"hash"`
		}{
			Nonce:    "0x0",
			GasPrice: "0x4a817c800",
			Gas:      "0x5208",
			To:       "0x0000000000000000000000000000000000000002",
			Value:    "0xde0b6b3a7640000",
			Input:    "0x",
			V:        "0x25",
			R:        "0x4f355c7f6c7f7a4c9a0874ab8a8b98b2c97d43e7a208b8474b7b0d11f857c003",
			S:        "0x6f7e456609e6e797d1b4e9d5b4482e9c778b3d3ca7e8a8b4d2d3e7a8c8d2e4f5",
			Hash:     "0x123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0",
		},
	}

	client, server := setupHTTPTestServer(t, "account_signTransaction", expected)
	defer server.Close()

	result, err := client.SignTransaction(tx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSignTransactionIPC(t *testing.T) {
	tx := &Transaction{
		From:     "0x0000000000000000000000000000000000000001",
		To:       "0x0000000000000000000000000000000000000002",
		Gas:      "0x5208",
		GasPrice: "0x4a817c800",
		Value:    "0xde0b6b3a7640000",
		Nonce:    "0x0",
		Data:     "0x",
	}

	expected := &SignTxResponse{
		Raw: "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675",
		Tx: struct {
			Nonce    string `json:"nonce"`
			GasPrice string `json:"gasPrice"`
			Gas      string `json:"gas"`
			To       string `json:"to"`
			Value    string `json:"value"`
			Input    string `json:"input"`
			V        string `json:"v"`
			R        string `json:"r"`
			S        string `json:"s"`
			Hash     string `json:"hash"`
		}{
			Nonce:    "0x0",
			GasPrice: "0x4a817c800",
			Gas:      "0x5208",
			To:       "0x0000000000000000000000000000000000000002",
			Value:    "0xde0b6b3a7640000",
			Input:    "0x",
			V:        "0x25",
			R:        "0x4f355c7f6c7f7a4c9a0874ab8a8b98b2c97d43e7a208b8474b7b0d11f857c003",
			S:        "0x6f7e456609e6e797d1b4e9d5b4482e9c778b3d3ca7e8a8b4d2d3e7a8c8d2e4f5",
			Hash:     "0x123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0",
		},
	}

	client, listener, tmpDir := setupIPCTestServer(t, "account_signTransaction", expected)
	defer listener.Close()
	defer os.RemoveAll(tmpDir)
	defer client.Close()

	result, err := client.SignTransaction(tx)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSignDataHTTP(t *testing.T) {
	req := &SignDataRequest{
		Address: "0x0000000000000000000000000000000000000001",
		Data:    "0x48656c6c6f20576f726c64", // "Hello World" in hex
	}

	expected := &SignDataResponse{
		Signature: "0x4f355c7f6c7f7a4c9a0874ab8a8b98b2c97d43e7a208b8474b7b0d11f857c0036f7e456609e6e797d1b4e9d5b4482e9c778b3d3ca7e8a8b4d2d3e7a8c8d2e4f500",
	}

	client, server := setupHTTPTestServer(t, "account_signData", expected)
	defer server.Close()

	result, err := client.SignData(req)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSignDataIPC(t *testing.T) {
	req := &SignDataRequest{
		Address: "0x0000000000000000000000000000000000000001",
		Data:    "0x48656c6c6f20576f726c64", // "Hello World" in hex
	}

	expected := &SignDataResponse{
		Signature: "0x4f355c7f6c7f7a4c9a0874ab8a8b98b2c97d43e7a208b8474b7b0d11f857c0036f7e456609e6e797d1b4e9d5b4482e9c778b3d3ca7e8a8b4d2d3e7a8c8d2e4f500",
	}

	client, listener, tmpDir := setupIPCTestServer(t, "account_signData", expected)
	defer listener.Close()
	defer os.RemoveAll(tmpDir)
	defer client.Close()

	result, err := client.SignData(req)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSignTypedDataHTTP(t *testing.T) {
	typedData := []byte(`{
		"types": {
			"EIP712Domain": [
				{"name": "name", "type": "string"},
				{"name": "version", "type": "string"},
				{"name": "chainId", "type": "uint256"},
				{"name": "verifyingContract", "type": "address"}
			],
			"Person": [
				{"name": "name", "type": "string"},
				{"name": "wallet", "type": "address"}
			]
		},
		"primaryType": "Person",
		"domain": {
			"name": "Test",
			"version": "1",
			"chainId": 1,
			"verifyingContract": "0x0000000000000000000000000000000000000000"
		},
		"message": {
			"name": "John Doe",
			"wallet": "0x0000000000000000000000000000000000000001"
		}
	}`)

	req := &TypedDataRequest{
		Address:    "0x0000000000000000000000000000000000000001",
		TypedData:  typedData,
		RawVersion: "V4",
	}

	expected := &SignDataResponse{
		Signature: "0x4f355c7f6c7f7a4c9a0874ab8a8b98b2c97d43e7a208b8474b7b0d11f857c0036f7e456609e6e797d1b4e9d5b4482e9c778b3d3ca7e8a8b4d2d3e7a8c8d2e4f500",
	}

	client, server := setupHTTPTestServer(t, "account_signTypedData", expected)
	defer server.Close()

	result, err := client.SignTypedData(req)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestSignTypedDataIPC(t *testing.T) {
	typedData := []byte(`{
		"types": {
			"EIP712Domain": [
				{"name": "name", "type": "string"},
				{"name": "version", "type": "string"},
				{"name": "chainId", "type": "uint256"},
				{"name": "verifyingContract", "type": "address"}
			],
			"Person": [
				{"name": "name", "type": "string"},
				{"name": "wallet", "type": "address"}
			]
		},
		"primaryType": "Person",
		"domain": {
			"name": "Test",
			"version": "1",
			"chainId": 1,
			"verifyingContract": "0x0000000000000000000000000000000000000000"
		},
		"message": {
			"name": "John Doe",
			"wallet": "0x0000000000000000000000000000000000000001"
		}
	}`)

	req := &TypedDataRequest{
		Address:    "0x0000000000000000000000000000000000000001",
		TypedData:  typedData,
		RawVersion: "V4",
	}

	expected := &SignDataResponse{
		Signature: "0x4f355c7f6c7f7a4c9a0874ab8a8b98b2c97d43e7a208b8474b7b0d11f857c0036f7e456609e6e797d1b4e9d5b4482e9c778b3d3ca7e8a8b4d2d3e7a8c8d2e4f500",
	}

	client, listener, tmpDir := setupIPCTestServer(t, "account_signTypedData", expected)
	defer listener.Close()
	defer os.RemoveAll(tmpDir)
	defer client.Close()

	result, err := client.SignTypedData(req)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestEcRecoverHTTP(t *testing.T) {
	req := &EcRecoverRequest{
		Data:      "0x48656c6c6f20576f726c64", // "Hello World" in hex
		Signature: "0x4f355c7f6c7f7a4c9a0874ab8a8b98b2c97d43e7a208b8474b7b0d11f857c0036f7e456609e6e797d1b4e9d5b4482e9c778b3d3ca7e8a8b4d2d3e7a8c8d2e4f500",
	}

	expected := &EcRecoverResponse{
		Address: "0x0000000000000000000000000000000000000001",
	}

	client, server := setupHTTPTestServer(t, "account_ecRecover", expected)
	defer server.Close()

	result, err := client.EcRecover(req)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestEcRecoverIPC(t *testing.T) {
	req := &EcRecoverRequest{
		Data:      "0x48656c6c6f20576f726c64", // "Hello World" in hex
		Signature: "0x4f355c7f6c7f7a4c9a0874ab8a8b98b2c97d43e7a208b8474b7b0d11f857c0036f7e456609e6e797d1b4e9d5b4482e9c778b3d3ca7e8a8b4d2d3e7a8c8d2e4f500",
	}

	expected := &EcRecoverResponse{
		Address: "0x0000000000000000000000000000000000000001",
	}

	client, listener, tmpDir := setupIPCTestServer(t, "account_ecRecover", expected)
	defer listener.Close()
	defer os.RemoveAll(tmpDir)
	defer client.Close()

	result, err := client.EcRecover(req)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestVersionHTTP(t *testing.T) {
	expected := &VersionResponse{
		Version: "6.1.0",
	}

	client, server := setupHTTPTestServer(t, "account_version", expected)
	defer server.Close()

	result, err := client.Version()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestVersionIPC(t *testing.T) {
	expected := &VersionResponse{
		Version: "6.1.0",
	}

	client, listener, tmpDir := setupIPCTestServer(t, "account_version", expected)
	defer listener.Close()
	defer os.RemoveAll(tmpDir)
	defer client.Close()

	result, err := client.Version()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
