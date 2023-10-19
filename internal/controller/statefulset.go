package controller

import (
	"context"
	v1 "fastdfs_operator/api/v1"
	"fmt"
	"time"

	"github.com/fearlesschenc/operator-utils/pkg/reconcile"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *FastDFSReconciler) ReconcileStatefulSet(ctx context.Context, object metav1.Object) (reconcile.Result, error) {
	cluster, _ := object.(*v1.FastDFS)
	r.Log.WithFields(cluster.Fields()).Info("reconcile cluster statefulset")

	sts := r.makeStatefulSet(cluster)
	if result, err := controllerutil.CreateOrUpdate(ctx, r.Client, sts, func() error {
		err := r.mutateStatefulSet(cluster, sts)
		if err != nil {
			r.Log.WithFields(cluster.Fields()).Error(err, "failed to mutate statefulset")
			return err
		}

		pvcBeingDeleted, err := r.isPVCBeingDeleted(ctx, cluster, *sts.Spec.Replicas)
		if err != nil {
			r.Log.WithFields(cluster.Fields()).Error(err, "failed to check if pvc is being deleted")
			return err
		}

		if pvcBeingDeleted && cluster.Status.CurrentStatefulSetReplicas < *sts.Spec.Replicas {
			return fmt.Errorf("current replicas: %d, target replicas: %d, need to wait for pvc deleted",
				cluster.Status.CurrentStatefulSetReplicas,
				*sts.Spec.Replicas)
		}

		return err
	}); err != nil {
		return reconcile.RequeueOnError(err)
	} else {
		switch result {
		case controllerutil.OperationResultCreated:
			r.Log.WithFields(cluster.Fields()).Info("created statefulset")
			r.Eventf(cluster, corev1.EventTypeNormal, "StatefulSetCreated", "created fastdfs statefulset")
		case controllerutil.OperationResultUpdated:
			r.Log.WithFields(cluster.Fields()).Info("updated statefulset")
			r.Eventf(cluster, corev1.EventTypeNormal, "StatefulSetUpdated", "updated fastdfs statefulset")

			_ = wait.Poll(time.Second, time.Second*30, func() (done bool, err error) {
				if err = r.Get(ctx, client.ObjectKeyFromObject(sts), sts); err != nil {
					return false, nil
				}
				return isUpdating(sts), nil
			})

			if cluster.Status.CurrentStatefulSetReplicas > *sts.Spec.Replicas {
				if err := r.cleanupPVCs(ctx, cluster, *sts.Spec.Replicas); err != nil {
					return reconcile.RequeueOnError(err)
				}
			}
		}
	}
	cluster.Status.Replicas = sts.Status.Replicas
	observerReplicas := cluster.Status.Replicas - *cluster.Spec.ParticipantReplicas
	if observerReplicas < 0 {
		observerReplicas = 0
	}
	cluster.Status.ObserverReplicas = observerReplicas
	cluster.Status.ParticipantReplicas = cluster.Status.Replicas - cluster.Status.ObserverReplicas
	return r.reconcilePods(ctx, cluster)
}

func isUpdating(sts *appsv1.StatefulSet) bool {
	return sts.Status.UpdateRevision != sts.Status.CurrentRevision || sts.Status.ReadyReplicas != *sts.Spec.Replicas
}
