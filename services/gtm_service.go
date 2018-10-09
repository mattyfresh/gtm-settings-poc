package services

import (
	"context"
	"net/http"

	"golang.org/x/oauth2/google"

	gtm "google.golang.org/api/tagmanager/v2"
)

// returns an OAuth2 client
func createHTTPClient(scope ...string) (*http.Client, error) {
	// Use Application Default Credentials to authenticate.
	// In our case it should default to the environment variable GOOGLE_APPLICATION_CREDENTIALS
	// For more info, see https://developers.google.com/accounts/docs/application-default-credentials
	client, err := google.DefaultClient(context.Background(), scope...)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// ConfigureGTMService returns a reference to the GTM service
func ConfigureGTMService() (*gtm.Service, error) {
	// create HTTP client and add the correct scope/permissions
	client, createClientErr := createHTTPClient(
		gtm.TagmanagerEditContainersScope,
		gtm.TagmanagerDeleteContainersScope,
		gtm.TagmanagerEditContainerversionsScope,
		gtm.TagmanagerPublishScope,
	)
	if createClientErr != nil {
		return nil, createClientErr
	}

	// our service that will allow us to use the GTM apis
	service, err := gtm.New(client)
	if err != nil {
		return nil, err
	}

	return service, nil
}
