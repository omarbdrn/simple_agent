package reporter

import (
	"github.com/omarbdrn/simple_agent/internal/api"
	"github.com/omarbdrn/simple_agent/pkg/constants"
	"github.com/omarbdrn/simple_agent/pkg/models"
)

func Report(hostService models.HostServices) {
	if len(hostService.Ports) > 0 {
		request := api.HTTPRequest{
			Endpoint: constants.SubmitHostEndpoint,
			Method:   "POST",
			IsJson:   true,
			Body:     hostService,
		}

		_, _ = api.PerformRequest(request)
	}
}
