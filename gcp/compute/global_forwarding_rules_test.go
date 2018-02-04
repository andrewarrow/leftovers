package compute_test

import (
	"errors"

	"github.com/genevieve/leftovers/gcp/compute"
	"github.com/genevieve/leftovers/gcp/compute/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gcpcompute "google.golang.org/api/compute/v1"
)

var _ = Describe("GlobalForwardingRules", func() {
	var (
		client *fakes.GlobalForwardingRulesClient
		logger *fakes.Logger

		globalForwardingRules compute.GlobalForwardingRules
	)

	BeforeEach(func() {
		client = &fakes.GlobalForwardingRulesClient{}
		logger = &fakes.Logger{}

		globalForwardingRules = compute.NewGlobalForwardingRules(client, logger)
	})

	Describe("List", func() {
		var filter string

		BeforeEach(func() {
			logger.PromptCall.Returns.Proceed = true
			client.ListGlobalForwardingRulesCall.Returns.Output = &gcpcompute.ForwardingRuleList{
				Items: []*gcpcompute.ForwardingRule{{
					Name: "banana-rule",
				}},
			}
			filter = "banana"
		})

		It("lists, filters, and prompts for global forwarding rules to delete", func() {
			list, err := globalForwardingRules.List(filter)
			Expect(err).NotTo(HaveOccurred())

			Expect(client.ListGlobalForwardingRulesCall.CallCount).To(Equal(1))

			Expect(logger.PromptCall.Receives.Message).To(Equal("Are you sure you want to delete global forwarding rule banana-rule?"))

			Expect(list).To(HaveLen(1))
		})

		Context("when the client fails to list global forwarding rules", func() {
			BeforeEach(func() {
				client.ListGlobalForwardingRulesCall.Returns.Error = errors.New("some error")
			})

			It("returns the error", func() {
				_, err := globalForwardingRules.List(filter)
				Expect(err).To(MatchError("Listing global forwarding rules: some error"))
			})
		})

		Context("when the global forwarding rule name does not contain the filter", func() {
			It("does not add it to the list", func() {
				list, err := globalForwardingRules.List("grape")
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PromptCall.CallCount).To(Equal(0))
				Expect(list).To(HaveLen(0))
			})
		})

		Context("when the user says no to the prompt", func() {
			BeforeEach(func() {
				logger.PromptCall.Returns.Proceed = false
			})

			It("does not add it to the list", func() {
				list, err := globalForwardingRules.List(filter)
				Expect(err).NotTo(HaveOccurred())

				Expect(list).To(HaveLen(0))
			})
		})
	})
})
