package bitbucket

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider will create the necessary terraform provider to talk to the Bitbucket APIs you should
// specify a USERNAME and PASSWORD
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Required:    true,
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BITBUCKET_PASSWORD", nil),
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"bitbucket_hook":                resourceHook(),
			"bitbucket_default_reviewers":   resourceDefaultReviewers(),
			"bitbucket_repository":          resourceRepository(),
			"bitbucket_repository_variable": resourceRepositoryVariable(),
			"bitbucket_deployment":          resourceDeployment(),
			"bitbucket_deployment_variable": resourceDeploymentVariable(),
			"bitbucket_project":             resourceProject(),
			"bitbucket_branch_restriction":  resourceBranchRestriction(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bitbucket_user":       dataUser(),
			"bitbucket_repository": dataRepository(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client := &Client{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		HTTPClient: &http.Client{},
	}

	return client, nil
}
