# Clef Client

[![Go Reference](https://pkg.go.dev/badge/github.com/AxLabs/clef-client.svg)](https://pkg.go.dev/github.com/AxLabs/clef-client)
[![Go Version](https://img.shields.io/github/go-mod/go-version/AxLabs/clef-client)](https://go.dev/)
[![CI Status](https://github.com/AxLabs/clef-client/workflows/Tests/badge.svg)](https://github.com/AxLabs/clef-client/actions)

A Go client library for interacting with the [Ethereum Clef external signer](https://geth.ethereum.org/docs/tools/clef/). This library provides a clean, type-safe interface to communicate with Clef through both HTTP JSON-RPC and Unix Domain Socket (IPC) protocols. It supports all major [Clef](https://github.com/ethereum/go-ethereum/tree/master/cmd/clef) operations including account management, transaction signing, and data signing following various Ethereum standards.

## Installation

```bash
go get github.com/AxLabs/clef-client
```

## Usage

### Creating a Client

You can create a client using either HTTP or IPC:

```go
// Using HTTP
httpClient := clefclient.NewHTTPClient("http://localhost:8550")
defer httpClient.Close()

// Using IPC
ipcClient, err := clefclient.NewIPCClient("/path/to/clef.ipc")
if err != nil {
    log.Fatal(err)
}
defer ipcClient.Close()
```

### Account Management

```go
// List accounts
accounts, err := client.ListAccounts()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Available accounts: %v\n", accounts)

// Create new account
address, err := client.NewAccount()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("New account created: %s\n", address)
```

### Signing Transactions

```go
tx := &clefclient.Transaction{
    From:     "0x0000000000000000000000000000000000000001",
    To:       "0x0000000000000000000000000000000000000002",
    Gas:      "0x5208",
    GasPrice: "0x4a817c800",
    Value:    "0xde0b6b3a7640000", // 1 ETH in wei
    Nonce:    "0x0",
    Data:     "0x",
}

response, err := client.SignTransaction(tx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Signed transaction: %s\n", response.Raw)
```

### Signing Data

```go
// Sign plain data
dataReq := &clefclient.SignDataRequest{
    Address: "0x0000000000000000000000000000000000000001",
    Data:    "0x48656c6c6f20576f726c64", // "Hello World" in hex
}

signature, err := client.SignData(dataReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Signature: %s\n", signature.Signature)

// Sign typed data (EIP-712)
typedData := []byte(`{
    "types": {
        "EIP712Domain": [
            {"name": "name", "type": "string"},
            {"name": "version", "type": "string"}
        ],
        "Person": [
            {"name": "name", "type": "string"},
            {"name": "wallet", "type": "address"}
        ]
    },
    "primaryType": "Person",
    "domain": {
        "name": "My App",
        "version": "1.0"
    },
    "message": {
        "name": "Alice",
        "wallet": "0x0000000000000000000000000000000000000001"
    }
}`)

typedReq := &clefclient.TypedDataRequest{
    Address:    "0x0000000000000000000000000000000000000001",
    TypedData:  typedData,
    RawVersion: "V4",
}

signature, err = client.SignTypedData(typedReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("EIP-712 Signature: %s\n", signature.Signature)
```

### EC Recover

```go
// Recover address from signature
recoverReq := &clefclient.EcRecoverRequest{
    Data:      "0x48656c6c6f20576f726c64",
    Signature: "0x4f355c7f6c7f7a4c...",
}

recovered, err := client.EcRecover(recoverReq)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Recovered address: %s\n", recovered.Address)
```

### Version

```go
// Get Clef version
version, err := client.Version()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Clef version: %s\n", version.Version)
```

## License

Apache License 2.0 - See [LICENSE](LICENSE) file for details.
