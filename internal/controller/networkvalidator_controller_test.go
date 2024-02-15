package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/spectrocloud-labs/validator-plugin-network/api/v1alpha1"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

const awsValidatorName = "network-validator"

var _ = Describe("NetworkValidator controller", Ordered, func() {

	BeforeEach(func() {
		// toggle true/false to enable/disable the NetworkValidator controller specs
		if false {
			Skip("skipping")
		}
	})

	val := &v1alpha1.NetworkValidator{
		ObjectMeta: metav1.ObjectMeta{
			Name:      awsValidatorName,
			Namespace: validatorNamespace,
		},
		Spec: v1alpha1.NetworkValidatorSpec{
			DNSRules: []v1alpha1.DNSRule{
				{
					RuleName: "DNS Rule",
					Host:     "google.com",
				},
			},
			ICMPRules: []v1alpha1.ICMPRule{
				{
					RuleName: "ICMP Rule",
					Host:     "google.com",
				},
			},
			IPRangeRules: []v1alpha1.IPRangeRule{
				{
					RuleName: "IP Range Rule",
					StartIP:  "10.0.0.0",
					Length:   4,
				},
			},
			MTURules: []v1alpha1.MTURule{
				{
					RuleName: "MTU Rule",
					Host:     "google.com",
					MTU:      1200,
				},
			},
			TCPConnRules: []v1alpha1.TCPConnRule{
				{
					RuleName: "TCP Connection Rule",
					Host:     "google.com",
					Ports:    []int{443},
				},
			},
		},
	}

	vr := &vapi.ValidationResult{}
	vrKey := types.NamespacedName{Name: validationResultName(val), Namespace: validatorNamespace}

	It("Should create a ValidationResult and update its Status", func() {
		By("By creating a new NetworkValidator")
		ctx := context.Background()

		Expect(k8sClient.Create(ctx, val)).Should(Succeed())

		// Wait for the ValidationResult's Status to be updated
		Eventually(func() bool {
			if err := k8sClient.Get(ctx, vrKey, vr); err != nil {
				return false
			}
			passConditions, failConditions := 0, 0
			for _, c := range vr.Status.Conditions {
				if c.Status == corev1.ConditionTrue {
					passConditions++
				}
				if c.Status == corev1.ConditionFalse {
					failConditions++
				}
			}
			stateFailed := vr.Status.State == vapi.ValidationFailed
			// OR required here due to:
			// - ping's lack of MTU discovery support on Darwin
			// - ICMP blocked on GHA runners: https://github.com/orgs/community/discussions/26184
			darwinOk := stateFailed && failConditions == 1 && passConditions == 4
			ghOnlineOk := stateFailed && failConditions == 2 && passConditions == 3
			ghSelfHostedOk := !stateFailed && failConditions == 0 && passConditions == 5
			return darwinOk || ghOnlineOk || ghSelfHostedOk
		}, timeout, interval).Should(BeTrue(), "failed to create a ValidationResult")
	})
})
