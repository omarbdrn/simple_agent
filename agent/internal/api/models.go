package api

type HTTPRequest struct{
	Endpoint    string
	Method		string
	IsJson		bool
	Body		interface{}
	Headers		map[string]string
}