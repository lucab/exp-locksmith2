package server

import (
	"encoding/json"
	"errors"
	"net/http"
)

// HTTPParams contains all parameters for a remote lock
// request.
type HTTPParams struct {
	ClientParams Params `json:"client_params"`
}

// Params contains client parameters for a remote lock
// request.
type Params struct {
	CurrentVersion string `json:"current_version,omitempty"`
	NodeUUID       string `json:"node_uuid"`
	Group          string `json:"group,omitempty"`
}

// NodeIdentity contains validated client identity from
// request parameters.
type NodeIdentity struct {
	UUID  string
	Group string
}

func validateIdentity(req *http.Request) (*NodeIdentity, error) {
	var group string
	var nodeID string

	decoder := json.NewDecoder(req.Body)
	var input HTTPParams
	if err := decoder.Decode(&input); err != nil {
		return nil, err
	}

	if input.ClientParams.Group == "" {
		return nil, errors.New("empty group")
	}
	group = input.ClientParams.Group

	if input.ClientParams.NodeUUID == "" {
		return nil, errors.New("empty node ID")
	}
	nodeID = input.ClientParams.NodeUUID

	identity := NodeIdentity{
		Group: group,
		UUID:  nodeID,
	}

	return &identity, nil
}
