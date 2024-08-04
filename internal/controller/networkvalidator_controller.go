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
	"time"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ktypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/validator-labs/validator-plugin-network/api/v1alpha1"
	"github.com/validator-labs/validator-plugin-network/internal/secrets"
	"github.com/validator-labs/validator-plugin-network/pkg/validate"
	vapi "github.com/validator-labs/validator/api/v1alpha1"
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
		Name:      vres.Name(validator),
		Namespace: req.Namespace,
	}
	if err := r.Get(ctx, nn, vr); err == nil {
		vres.HandleExisting(vr, r.Log)
	} else {
		if !apierrs.IsNotFound(err) {
			l.Error(err, "unexpected error getting ValidationResult")
		}
		if err := vres.HandleNew(ctx, r.Client, p, vres.Build(validator), r.Log); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: time.Millisecond}, nil
	}

	// Always update the expected result count in case the validator's rules have changed
	vr.Spec.ExpectedResults = validator.Spec.ResultCount()

	// Fetch additional CAs to augment the system cert pool
	caPems := validator.Spec.CACerts.RawCerts()

	for _, secretRef := range validator.Spec.CACerts.SecretRefs {
		caPem, err := secrets.ReadKeys(secretRef.Name, req.Namespace, []string{secretRef.Key}, r.Client)
		if err != nil {
			r.Log.Error(err, "failed to read CA certificate secret")
			return ctrl.Result{}, err
		}
		caPems = append(caPems, caPem[0])
	}

	// Fetch HTTP basic auth credentials
	auths := make([][][]byte, len(validator.Spec.HTTPFileRules))
	for _, rule := range validator.Spec.HTTPFileRules {
		var auth [][]byte
		if rule.AuthSecretRef != nil {
			auth, err = secrets.ReadKeys(rule.AuthSecretRef.Name, req.Namespace, rule.AuthSecretRef.Keys(), r.Client)
			if err != nil {
				r.Log.Error(err, "failed to parse HTTP basic auth", "rule", rule.RuleName)
				return ctrl.Result{}, err
			}
		}
		auths = append(auths, auth)
	}

	// Validate the rules
	resp := validate.Validate(validator.Spec, caPems, auths, r.Log)

	// Patch the ValidationResult with the latest ValidationRuleResults
	if err := vres.SafeUpdate(ctx, p, vr, resp, r.Log); err != nil {
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
