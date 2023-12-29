package rest

type IPRangeDTO struct {
	CIDR string `json:"cidr"`
}

type Share struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	CIDRs []string `json:"cidrs"`
}
