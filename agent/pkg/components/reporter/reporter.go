package reporter

import (
	"fmt"

	"github.com/omarbdrn/simple_agent/pkg/models"
)

func Report(hostService models.HostServices) {
	fmt.Println(hostService.IP, hostService.Ports)

	// if len(hostService.Ports) > 0 {
	// 	request := api.HTTPRequest{
	// 		Endpoint: constants.SubmitHostEndpoint,
	// 		Method:   "POST",
	// 		IsJson:   true,
	// 		Body:     hostService,
	// 	}

	// 	_, _ = api.PerformRequest(request)
	// }
}
