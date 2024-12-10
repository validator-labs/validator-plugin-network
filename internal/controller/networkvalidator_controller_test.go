package controller

import (
	"context"
	"os"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/validator-labs/validator-plugin-network/api/v1alpha1"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	vres "github.com/validator-labs/validator/pkg/validationresult"
	// +kubebuilder:scaffold:imports
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
					RuleName: "DNS Rule - Pass",
					Host:     "google.com",
				},
				{
					RuleName: "DNS Rule - Fail",
					Host:     "_invalid_",
				},
			},
			ICMPRules: []v1alpha1.ICMPRule{
				{
					RuleName: "ICMP Rule - Pass",
					Host:     "google.com",
				},
				{
					RuleName: "ICMP Rule - Fail",
					Host:     "_invalid_",
				},
			},
			IPRangeRules: []v1alpha1.IPRangeRule{
				{
					RuleName: "IP Range Rule - Pass",
					StartIP:  "127.0.0.1",
					Length:   1,
				},
				{
					RuleName: "IP Range Rule - Fail",
					StartIP:  "10.0.0.0",
					Length:   4,
				},
			},
			MTURules: []v1alpha1.MTURule{
				{
					RuleName: "MTU Rule - Pass",
					Host:     "google.com",
					MTU:      1200,
				},
				{
					RuleName: "MTU Rule - Fail",
					Host:     "google.com",
					MTU:      4000,
				},
			},
			TCPConnRules: []v1alpha1.TCPConnRule{
				{
					RuleName: "TCP Connection Rule - Pass",
					Host:     "google.com",
					Ports:    []int{443},
				},
				{
					RuleName: "TCP Connection Rule - Fail",
					Host:     "google.com",
					Ports:    []int{9999},
				},
			},
		},
	}

	vr := &vapi.ValidationResult{}
	vrKey := types.NamespacedName{Name: vres.Name(val), Namespace: validatorNamespace}

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
			for _, c := range vr.Status.ValidationConditions {
				if c.Status == corev1.ConditionTrue {
					passConditions++
				}
				if c.Status == corev1.ConditionFalse {
					failConditions++
				}
			}
			return expectedTestResult(failConditions, passConditions)
		}, timeout, interval).Should(BeTrue(), "failed to create a ValidationResult")
	})
})

// expectedTestResult returns whether the overall test expectation is a pass, which varies due to:
//   - ping's lack of MTU discovery support on Darwin
//   - ICMP blocked on GitHub-hosted runners: https://github.com/orgs/community/discussions/26184
func expectedTestResult(failConditions, passConditions int) bool {
	var env string
	var expectedFail, expectedPass int

	gha := os.Getenv("IS_GITHUB_ACTION") == "true"
	selfHosted := os.Getenv("IS_SELF_HOSTED") == "true"

	if runtime.GOOS == "darwin" && !gha {
		env = "darwin / local"
		expectedFail = 6 // all checks related to MTU will fail
		expectedPass = 4
	} else if gha && !selfHosted {
		env = "Github-hosted runner"
		expectedFail = 8 // all checks related to ICMP will fail
		expectedPass = 2
	} else {
		env = "self-hosted runner"
		expectedFail = 5
		expectedPass = 5
	}

	logf.Log.Info(
		"expectedTestResult", "environment", env,
		"expectedFailed", expectedFail, "failed", failConditions,
		"expectedPassed", expectedPass, "passed", passConditions,
	)
	return failConditions == expectedFail && passConditions == expectedPass
}
