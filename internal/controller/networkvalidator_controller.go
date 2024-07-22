/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package controller defines a controller for reconciling NetworkValidator objects.
package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/validator-labs/validator-plugin-network/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-network/internal/constants"
	"github.com/validator-labs/validator-plugin-network/internal/http"
	"github.com/validator-labs/validator-plugin-network/internal/secrets"
	"github.com/validator-labs/validator-plugin-network/internal/validators"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
	"github.com/validator-labs/validator/pkg/types"
	"github.com/validator-labs/validator/pkg/util"
	vres "github.com/validator-labs/validator/pkg/validationresult"
)

// NetworkValidatorReconciler reconciles a NetworkValidator object
type NetworkValidatorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=networkvalidators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=networkvalidators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=validation.spectrocloud.labs,resources=networkvalidators/finalizers,verbs=update

// Reconcile reconciles each rule found in each NetworkValidator in the cluster and creates ValidationResults accordingly
func (r *NetworkValidatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := r.Log.V(0).WithValues("name", req.Name, "namespace", req.Namespace)
	l.Info("Reconciling NetworkValidator")

	validator := &v1alpha1.NetworkValidator{}
	if err := r.Get(ctx, req.NamespacedName, validator); err != nil {
		if !apierrs.IsNotFound(err) {
			l.Error(err, "failed to fetch NetworkValidator")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the active validator's validation result
	vr := &vapi.ValidationResult{}
	p, err := patch.NewHelper(vr, r.Client)
	if err != nil {
		l.Error(err, "failed to create patch helper")
		return ctrl.Result{}, err
	}
	nn := ktypes.NamespacedName{
		Name:      validationResultName(validator),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExistingValidationResult(vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			l.Error(err, "unexpected error getting ValidationResult")
		}
		if err := vres.HandleNewValidationResult(ctx, r.Client, p, buildValidationResult(validator), r.Log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Millisecond}, nil
	}

	// Always update the expected result count in case the validator's rules have changed
	vr.Spec.ExpectedResults = validator.Spec.ResultCount()

	resp := types.ValidationResponse{
		ValidationRuleResults: make([]*types.ValidationRuleResult, 0, vr.Spec.ExpectedResults),
		ValidationRuleErrors:  make([]error, 0, vr.Spec.ExpectedResults),
	}

	networkService := validators.NewNetworkService(r.Log)

	// If CACert config provided, use the inline certs and secret refs.
	caPems := make([][]byte, 0)
	for _, cert := range validator.Spec.CACerts.Certs {
		caPems = append(caPems, []byte(cert))
	}
	for _, secretRef := range validator.Spec.CACerts.SecretRefs {
		caPem, err := secrets.Read(secretRef.Name, req.Namespace, secretRef.Key, r.Client)
		if err != nil {
			r.Log.Error(err, "failed to read CA certificate secret")
			return ctrl.Result{}, err
		}
		caPems = append(caPems, caPem)
	}

	// DNS rules
	for _, rule := range validator.Spec.DNSRules {
		vrr := networkService.ReconcileDNSRule(rule)
		resp.AddResult(vrr, err)
	}

	// ICMP rules
	for _, rule := range validator.Spec.ICMPRules {
		vrr := networkService.ReconcileICMPRule(rule)
		resp.AddResult(vrr, err)
	}

	// IP range rules
	for _, rule := range validator.Spec.IPRangeRules {
		vrr := networkService.ReconcileIPRangeRule(rule)
		resp.AddResult(vrr, err)
	}

	// MTU rules
	for _, rule := range validator.Spec.MTURules {
		vrr := networkService.ReconcileMTURule(rule)
		resp.AddResult(vrr, err)
	}

	// TCP connection rules
	for _, rule := range validator.Spec.TCPConnRules {
		tlsConfig, err := http.TLSConfig(caPems, rule.InsecureSkipTLSVerify, r.Log)
		if err != nil {
			vrr := validators.BuildValidationResult(rule, constants.ValidationTypeTCPConn)
			resp.AddResult(vrr, fmt.Errorf("failed to create TLS config: %w", err))
			continue
		}
		ruleSvc := validators.NewNetworkService(r.Log, validators.WithTLSConfig(tlsConfig))
		vrr := ruleSvc.ReconcileTCPConnRule(rule)
		resp.AddResult(vrr, err)
	}

	// HTTP file rules
	for _, rule := range validator.Spec.HTTPFileRules {
		transport, err := http.Transport(caPems, rule.InsecureSkipTLSVerify, r.Log)
		if err != nil {
			vrr := validators.BuildValidationResult(rule, constants.ValidationTypeHTTPFile)
			resp.AddResult(vrr, fmt.Errorf("failed to create HTTP transport: %w", err))
			continue
		}
		ruleSvc := validators.NewNetworkService(r.Log, validators.WithTransport(transport))
		vrr := ruleSvc.ReconcileHTTPFileRule(rule)
		resp.AddResult(vrr, err)
	}

	// Patch the ValidationResult with the latest ValidationRuleResults
	if err := vres.SafeUpdateValidationResult(ctx, p, vr, resp, r.Log); err != nil {
		return ctrl.Result{}, err
	}

	l.Info("Requeuing for re-validation in two minutes.")
	return ctrl.Result{RequeueAfter: time.Second * 120}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetworkValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NetworkValidator{}).
		Complete(r)
}

func buildValidationResult(validator *v1alpha1.NetworkValidator) *vapi.ValidationResult {
	return &vapi.ValidationResult{
		ObjectMeta: metav1.ObjectMeta{
			Name:      validationResultName(validator),
			Namespace: validator.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: validator.APIVersion,
					Kind:       validator.Kind,
					Name:       validator.Name,
					UID:        validator.UID,
					Controller: util.Ptr(true),
				},
			},
		},
		Spec: vapi.ValidationResultSpec{
			Plugin:          constants.PluginCode,
			ExpectedResults: validator.Spec.ResultCount(),
		},
	}
}

func validationResultName(validator *v1alpha1.NetworkValidator) string {
	return fmt.Sprintf("validator-plugin-network-%s", validator.Name)
}
