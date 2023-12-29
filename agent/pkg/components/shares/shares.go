package shares

import (
	"encoding/json"
	"fmt"

	"github.com/omarbdrn/simple_agent/internal/api"
	"github.com/omarbdrn/simple_agent/internal/configuration"
	"github.com/omarbdrn/simple_agent/pkg/constants"
	"github.com/omarbdrn/simple_agent/pkg/global"
	"github.com/omarbdrn/simple_agent/pkg/models"
)

func GetShare() (*models.Share, error) {
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
		return &models.Share{}, err
	}

	defer response.Body.Close()

	var result models.Share
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return &models.Share{}, err
	}

	global.SetCurrentShare(result)

	return &result, nil
}
