package shares

import (
	"encoding/json"
	"fmt"

	"github.com/omarbdrn/simple_agent/internal/api"
	"github.com/omarbdrn/simple_agent/internal/configuration"
	"github.com/omarbdrn/simple_agent/pkg/constants"
)

type Share struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	CIDRs []string `json:"cidrs"`
}

func GetShare() (*Share, error) {
	systemCapabilities := configuration.GetSystemCapabilities()

	request := api.HTTPRequest{
		Endpoint: constants.GetShareEndpoint,
		Method:   "GET",
		IsJson:   true,
		Body:     "",
		Headers: map[string]string{
			"Max-CIDRs":        fmt.Sprintf("%d", systemCapabilities.MaxParallelScans),
			"Available-Memory": fmt.Sprintf("%d", systemCapabilities.AvailableMemory),
		},
	}

	response, err := api.PerformRequest(request)
	if err != nil {
		return &Share{}, err
	}

	defer response.Body.Close()

	var result Share
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return &Share{}, err
	}

	return &result, nil
}
