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

	"github.com/spectrocloud-labs/valid8or-plugin-network/api/v1alpha1"
	"github.com/spectrocloud-labs/valid8or-plugin-network/internal/constants"
	"github.com/spectrocloud-labs/valid8or-plugin-network/internal/validators"
	v8or "github.com/spectrocloud-labs/valid8or/api/v1alpha1"
	"github.com/spectrocloud-labs/valid8or/pkg/types"
	v8ores "github.com/spectrocloud-labs/valid8or/pkg/validationresult"
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
	vr := &v8or.ValidationResult{}
	nn := ktypes.NamespacedName{
		Name:      fmt.Sprintf("valid8or-plugin-aws-%s", validator.Name),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		res, err := v8ores.HandleExistingValidationResult(nn, vr, r.Log)
		if res != nil {
			return *res, err
		}
	} else {
		if !apierrs.IsNotFound(err) {
			r.Log.V(0).Error(err, "unexpected error getting ValidationResult", "name", nn.Name, "namespace", nn.Namespace)
		}
		res, err := v8ores.HandleNewValidationResult(r.Client, constants.PluginCode, nn, vr, r.Log)
		if res != nil {
			return *res, err
		}
	}

	networkService := validators.NewNetworkService(r.Log)
	failed := &types.MonotonicBool{}

	// DNS rules
	for _, rule := range validator.Spec.DNSRules {
		validationResult, err := networkService.ReconcileDNSRule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile DNS rule")
		}
		v8ores.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
	}

	// ICMP rules
	for _, rule := range validator.Spec.ICMPRules {
		validationResult, err := networkService.ReconcileICMPRule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile ICMP rule")
		}
		v8ores.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
	}

	// IP range rules
	for _, rule := range validator.Spec.IPRangeRules {
		validationResult, err := networkService.ReconcileIPRangeRule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile IP range rule")
		}
		v8ores.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
	}

	// MTU rules
	for _, rule := range validator.Spec.MTURules {
		validationResult, err := networkService.ReconcileMTURule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile MTU rule")
		}
		v8ores.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
	}

	// Netcat rules
	for _, rule := range validator.Spec.NetcatRules {
		validationResult, err := networkService.ReconcileNetcatRule(nn, rule)
		if err != nil {
			r.Log.V(0).Error(err, "failed to reconcile netcat rule")
		}
		v8ores.SafeUpdateValidationResult(r.Client, nn, validationResult, failed, err, r.Log)
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
