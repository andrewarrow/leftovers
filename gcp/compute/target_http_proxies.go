package compute

import (
	"fmt"
	"strings"

	gcpcompute "google.golang.org/api/compute/v1"
)

type targetHttpProxiesClient interface {
	ListTargetHttpProxies() (*gcpcompute.TargetHttpProxyList, error)
	DeleteTargetHttpProxy(targetHttpProxy string) error
}

type TargetHttpProxies struct {
	client targetHttpProxiesClient
	logger logger
}

func NewTargetHttpProxies(client targetHttpProxiesClient, logger logger) TargetHttpProxies {
	return TargetHttpProxies{
		client: client,
		logger: logger,
	}
}

func (t TargetHttpProxies) List(filter string) (map[string]string, error) {
	targetHttpProxies, err := t.client.ListTargetHttpProxies()
	if err != nil {
		return nil, fmt.Errorf("Listing target http proxies: %s", err)
	}

	delete := map[string]string{}
	for _, targetHttpProxy := range targetHttpProxies.Items {
		if !strings.Contains(targetHttpProxy.Name, filter) {
			continue
		}

		proceed := t.logger.Prompt(fmt.Sprintf("Are you sure you want to delete target http proxy %s?", targetHttpProxy.Name))
		if !proceed {
			continue
		}

		delete[targetHttpProxy.Name] = ""
	}

	return delete, nil
}

func (t TargetHttpProxies) Delete(targetHttpProxies map[string]string) {
	var resources []TargetHttpProxy
	for name, _ := range targetHttpProxies {
		resources = append(resources, NewTargetHttpProxy(t.client, name))
	}

	t.delete(resources)
}

func (t TargetHttpProxies) delete(resources []TargetHttpProxy) {
	for _, resource := range resources {
		err := resource.Delete()

		if err != nil {
			t.logger.Printf("%s\n", err)
		} else {
			t.logger.Printf("SUCCESS deleting target http proxy %s\n", resource.name)
		}
	}
}
