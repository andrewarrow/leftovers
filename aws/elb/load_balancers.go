package elb

import (
	"fmt"
	"strings"

	awselb "github.com/aws/aws-sdk-go/service/elb"
	"github.com/genevieve/leftovers/aws/common"
)

type loadBalancersClient interface {
	DescribeLoadBalancers(*awselb.DescribeLoadBalancersInput) (*awselb.DescribeLoadBalancersOutput, error)
	DeleteLoadBalancer(*awselb.DeleteLoadBalancerInput) (*awselb.DeleteLoadBalancerOutput, error)
}

type LoadBalancers struct {
	client loadBalancersClient
	logger logger
}

func NewLoadBalancers(client loadBalancersClient, logger logger) LoadBalancers {
	return LoadBalancers{
		client: client,
		logger: logger,
	}
}

func (l LoadBalancers) ListOnly(filter string) ([]common.Deletable, error) {
	return l.get(filter)
}

func (l LoadBalancers) List(filter string) ([]common.Deletable, error) {
	resources, err := l.get(filter)
	if err != nil {
		return nil, err
	}

	var delete []common.Deletable
	for _, r := range resources {
		proceed := l.logger.PromptWithDetails(r.Type(), r.Name())
		if !proceed {
			continue
		}

		delete = append(delete, r)
	}

	return delete, nil
}

func (l LoadBalancers) get(filter string) ([]common.Deletable, error) {
	loadBalancers, err := l.client.DescribeLoadBalancers(&awselb.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, fmt.Errorf("Describe Elastic Load Balancers: %s", err)
	}

	var resources []common.Deletable
	for _, lb := range loadBalancers.LoadBalancerDescriptions {
		resource := NewLoadBalancer(l.client, lb.LoadBalancerName)

		if !strings.Contains(resource.Name(), filter) {
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
