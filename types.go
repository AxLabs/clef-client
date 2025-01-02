package clefclient

import (
	"encoding/json"
)

// Transaction represents an Ethereum transaction
type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
	Data     string `json:"data,omitempty"`
}

// SignDataRequest represents the parameters for signing data
type SignDataRequest struct {
	Address string `json:"address"`
	Data    string `json:"data"`
}

// TypedDataRequest represents the parameters for signing typed data
type TypedDataRequest struct {
	Address    string          `json:"address"`
	TypedData  json.RawMessage `json:"data"`
	RawVersion string          `json:"raw_version,omitempty"`
}

// SignTxResponse represents the response from signing a transaction
type SignTxResponse struct {
	Raw string `json:"raw"`
	Tx  struct {
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
	} `json:"tx"`
}

// SignDataResponse represents the response from signing data
type SignDataResponse struct {
	Signature string `json:"signature"`
}

// VersionResponse represents the response from version query
type VersionResponse struct {
	Version string `json:"version"`
}

// EcRecoverRequest represents the parameters for ecRecover
type EcRecoverRequest struct {
	Data      string `json:"data"`
	Signature string `json:"sig"`
}

// EcRecoverResponse represents the response from ecRecover
type EcRecoverResponse struct {
	Address string `json:"address"`
}
