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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	fastdfsv1 "fastdfs_operator/api/v1"
	v1 "fastdfs_operator/api/v1"

	"github.com/fearlesschenc/operator-utils/pkg/reconcile"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// FastDFSReconciler reconciles a FastDFS object
type FastDFSReconciler struct {
	client.Client
	Log      logr.Logger
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
}

func (r *FastDFSReconciler) GetReconcileSteps() []reconcile.Func {
	return reconcile.Funcs{
		r.ReconcileConfig,
		r.ReconcileStatefulSet,
	}
}

//+kubebuilder:rbac:groups=fastdfs.beordie.cn,resources=fastdfs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fastdfs.beordie.cn,resources=fastdfs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fastdfs.beordie.cn,resources=fastdfs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the FastDFS object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *FastDFSReconciler) Reconcile(ctx context.Context, request ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	r.Log.Info("Reconciling FastDFS started")
	cluster := &v1.FastDFS{}
	if err := r.Get(context.Background(), request.NamespacedName, cluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if ShouldReconcile(cluster) {
		result, err := reconcile.Reconcile(ctx, cluster).WithReconciler(r)
		if result.RequeueRequest {
			return reconcile.RequeueRequestAfter(result.RequeueDelay, err)
		}
		if result.CancelReconciliation {
			return reconcile.DoNotRequeueRequest(err)
		}
	}

	return RequeueRequestAfter(time.Second*10, nil)
}

func (r *FastDFSReconciler) Eventf(cluster *v1.FastDFS, eventType, reason, message string) {
	r.Recorder.AnnotatedEventf(cluster, map[string]string{
		"cloud.netease.com/app":          "fastDFS",
		"cloud.netease.com/cluster-name": cluster.Name,
	}, eventType, reason, message)
}

// SetupWithManager sets up the controller with the Manager.
func (r *FastDFSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{MaxConcurrentReconciles: 5}).
		For(&fastdfsv1.FastDFS{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&corev1.Pod{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
