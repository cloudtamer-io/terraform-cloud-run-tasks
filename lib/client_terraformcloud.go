package lib

// TerraformCloudClient -
type TerraformCloudClient struct {
	*RequestClient
}

// NewTerraformCloudClient -
func NewTerraformCloudClient(c *RequestClient) *TerraformCloudClient {
	return &TerraformCloudClient{
		c,
	}
}
