package controller

import (
	"context"
	v1 "fastdfs_operator/api/v1"

	"github.com/fearlesschenc/operator-utils/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ConfigMap map[string]string

func (r *FastDFSReconciler) ReconcileConfig(ctx context.Context, object metav1.Object) (reconcile.Result, error) {
	cluster, _ := object.(*v1.FastDFS)
	logr.FromContext(ctx).Info("reconcile cluster configmap")

	cm := makeConfigmap(cluster)
	if result, err := controllerutil.CreateOrUpdate(ctx, r.Client, cm, func() error {
		return r.mutateConfigmap(cluster, cm)
	}); err != nil {
		return reconcile.RequeueOnError(err)
	} else {
		switch result {
		case controllerutil.OperationResultCreated:
			logr.FromContext(ctx).Info("created configmap")
			r.Eventf(cluster, corev1.EventTypeNormal, "ConfigCreated", "created fastdfs configuration")
		case controllerutil.OperationResultUpdated:
			logr.FromContext(ctx).Info("updated configmap")
			r.Eventf(cluster, corev1.EventTypeNormal, "ConfigUpdated", "updated fastdfs configuration")
		}
	}
	return reconcile.Continue()
}
