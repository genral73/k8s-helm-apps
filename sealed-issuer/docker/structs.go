package main

// SealedRequest struct
type SealedRequest struct {
	Name      string    `json:"name,omitempty"`
	Namespace string    `json:"namespace,omitempty"`
	KVPairs   []KVPairs `json:"kvPairs,omitempty"`
	YAML      string    `json:"yaml,omitempty"`
}

// KVPairs struct is part of SealedRequest
type KVPairs struct {
	K string `json:"k,omitempty"`
	V string `json:"v,omitempty"`
}

// SealedResponse struct
type SealedResponse struct {
	FileContent  string `json:"fileContent,omitempty"`
	FileName     string `json:"fileName,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}
