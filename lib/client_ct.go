package lib

// CTClient -
type CTClient struct {
	*RequestClient
}

// NewCTClient -
func NewCTClient(c *RequestClient) *CTClient {
	return &CTClient{
		c,
	}
}
