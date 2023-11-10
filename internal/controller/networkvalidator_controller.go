/*
Copyright 2023.

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

package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ktypes "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/spectrocloud-labs/validator-plugin-network/api/v1alpha1"
	"github.com/spectrocloud-labs/validator-plugin-network/internal/constants"
	"github.com/spectrocloud-labs/validator-plugin-network/internal/validators"
	vapi "github.com/spectrocloud-labs/validator/api/v1alpha1"
	vres "github.com/spectrocloud-labs/validator/pkg/validationresult"
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
	r.Log.V(0).Info("Reconciling NetworkValidator", "name", req.Name, "namespace", req.Namespace)

	validator := &v1alpha1.NetworkValidator{}
	if err := r.Get(ctx, req.NamespacedName, validator); err != nil {
		// Ignore not-found errors, since they can't be fixed by an immediate requeue
		if apierrs.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "failed to fetch NetworkValidator")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the active validator's validation result
	vr := &vapi.ValidationResult{}
	nn := ktypes.NamespacedName{
		Name:      fmt.Sprintf("validator-plugin-aws-%s", validator.Name),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExistingValidationResult(nn, vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			r.Log.V(0).Error(err, "unexpected error getting ValidationResult", "name", nn.Name, "namespace", nn.Namespace)
		}
		if err := vres.HandleNewValidationResult(
			r.Client, constants.PluginCode, validator.Spec.ResultCount(), nn, r.Log,
		); err != nil {
			return ctrl.Result{}, err
		}
	}

	networkService := validators.NewNetworkService(r.Log)

	// DNS rules
	for _, rule := range validator.Spec.DNSRules {
		validationResult, err := networkService.ReconcileDNSRule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile DNS rule")
		}
		vres.SafeUpdateValidationResult(r.Client, nn, validationResult, err, r.Log)
	}

	// ICMP rules
	for _, rule := range validator.Spec.ICMPRules {
		validationResult, err := networkService.ReconcileICMPRule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile ICMP rule")
		}
		vres.SafeUpdateValidationResult(r.Client, nn, validationResult, err, r.Log)
	}

	// IP range rules
	for _, rule := range validator.Spec.IPRangeRules {
		validationResult, err := networkService.ReconcileIPRangeRule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile IP range rule")
		}
		vres.SafeUpdateValidationResult(r.Client, nn, validationResult, err, r.Log)
	}

	// MTU rules
	for _, rule := range validator.Spec.MTURules {
		validationResult, err := networkService.ReconcileMTURule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile MTU rule")
		}
		vres.SafeUpdateValidationResult(r.Client, nn, validationResult, err, r.Log)
	}

	// TCP connection rules
	for _, rule := range validator.Spec.TCPConnRules {
		validationResult, err := networkService.ReconcileTCPConnRule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile netcat rule")
		}
		vres.SafeUpdateValidationResult(r.Client, nn, validationResult, err, r.Log)
	}

	r.Log.V(0).Info("Requeuing for re-validation in two minutes.", "name", req.Name, "namespace", req.Namespace)
	return ctrl.Result{RequeueAfter: time.Second * 120}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetworkValidatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.NetworkValidator{}).
		Complete(r)
}
