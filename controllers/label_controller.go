/*
Copyright 2022.

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

package controllers

import (
	"context"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var matchLabelName = "test"
var matchLabelValue = "test"
var addLabelName = "test_add"
var addLabelValue = "test_add"

func init() {
	getEnvDefault := func(env, defaultVal string) string {
		if v := os.ExpandEnv("${" + env + "}"); v != "" {
			return v
		}
		return defaultVal
	}
	matchLabelName = getEnvDefault("MATCH_LABEL_NAME", matchLabelName)
	matchLabelValue = getEnvDefault("MATCH_LABEL_VALUE", matchLabelValue)
	addLabelName = getEnvDefault("ADD_LABEL_NAME", addLabelName)
	addLabelValue = getEnvDefault("ADD_LABEL_VALUE", addLabelValue)
}

// LabelReconciler reconciles a Label object
type LabelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Label object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *LabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// your logic here
	var node corev1.Node
	if err := r.Get(ctx, req.NamespacedName, &node); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	log.Info("got node", "name", node.Name)
	if node.GetLabels()[matchLabelName] == matchLabelValue {
		log.Info("match node label", "name", node.Name, "label_name", matchLabelName, "label_value", matchLabelValue)
		if v, ok := node.GetLabels()[addLabelName]; !ok || v != addLabelValue {
			log.Info("add node label", "name", node.Name, "label_name", addLabelName, "label_value", addLabelName)
			node.Labels[addLabelName] = addLabelValue
			if err := r.Update(ctx, &node); err != nil {
				log.Error(err, "update node err")
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Node{}, builder.WithPredicates(predicate.Or(predicate.LabelChangedPredicate{}))).
		Complete(r)
}
