package compute

import (
	"fmt"
	"strings"

	"github.com/genevieve/leftovers/gcp/common"
	gcpcompute "google.golang.org/api/compute/v1"
)

type globalForwardingRulesClient interface {
	ListGlobalForwardingRules() (*gcpcompute.ForwardingRuleList, error)
	DeleteGlobalForwardingRule(rule string) error
}

type GlobalForwardingRules struct {
	client globalForwardingRulesClient
	logger logger
}

func NewGlobalForwardingRules(client globalForwardingRulesClient, logger logger) GlobalForwardingRules {
	return GlobalForwardingRules{
		client: client,
		logger: logger,
	}
}

func (g GlobalForwardingRules) List(filter string) ([]common.Deletable, error) {
	rules, err := g.client.ListGlobalForwardingRules()
	if err != nil {
		return nil, fmt.Errorf("Listing global forwarding rules: %s", err)
	}

	var resources []common.Deletable
	for _, rule := range rules.Items {
		resource := NewGlobalForwardingRule(g.client, rule.Name)

		if !strings.Contains(rule.Name, filter) {
			continue
		}

		proceed := g.logger.Prompt(fmt.Sprintf("Are you sure you want to delete global forwarding rule %s?", rule.Name))
		if !proceed {
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
