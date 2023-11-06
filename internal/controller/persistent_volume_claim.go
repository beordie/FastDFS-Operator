package controller

import (
	"context"
	v1 "fastdfs_operator/api/v1"
	"fastdfs_operator/pkg/utils"

	util "github.com/fearlesschenc/operator-utils/pkg/controller"
	"github.com/fearlesschenc/operator-utils/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *FastDFSReconciler) getStorageClass(ctx context.Context, cluster *v1.FastDFS) (*storagev1.StorageClass, error) {
	var sc *storagev1.StorageClass

	if cluster.Spec.Storage.StorageClass != nil {
		sc = &storagev1.StorageClass{}
		if err := r.Get(ctx, client.ObjectKey{Name: *cluster.Spec.Storage.StorageClass}, sc); err != nil {
			return nil, err
		}
	} else {
		var err error
		if sc, err = utils.GetDefaultStorageClass(ctx, r.Client); err != nil {
			return nil, err
		}
	}

	return sc, nil
}

func (r *FastDFSReconciler) ReconcilePersistentVolumeClaim(ctx context.Context, object metav1.Object) (reconcile.Result, error) {
	cluster, _ := object.(*v1.FastDFS)
	logr.FromContext(ctx).Info("reconcile cluster pvc")

	if sc, err := r.getStorageClass(ctx, cluster); err != nil {
		return reconcile.RequeueOnError(err)
	} else if sc == nil || sc.AllowVolumeExpansion == nil || !*sc.AllowVolumeExpansion {
		return reconcile.Continue()
	}

	for ord := 0; ord < int(cluster.Status.Replicas); ord++ {
		pvc := &corev1.PersistentVolumeClaim{}
		if err := r.Get(ctx, cluster.GetPersistentVolumeClaimNamespacedName(ord), pvc); err != nil && !apierrors.IsNotFound(err) {
			return reconcile.RequeueOnError(err)
		} else if err != nil || util.IsObjectBeingDeleted(pvc) {
			continue
		}

		currentSize := pvc.Spec.Resources.Requests[corev1.ResourceStorage]
		expectSize := r.makePVCStorageSize(cluster)
		if (&currentSize).Cmp(expectSize) < 0 {
			pvc.Spec.Resources.Requests[corev1.ResourceStorage] = expectSize

			if err := r.Update(ctx, pvc); err != nil {
				return reconcile.RequeueOnError(err)
			}
			logr.FromContext(ctx).Info("updated persistent volume claim")
			r.Eventf(cluster, corev1.EventTypeNormal, "PersistentVolumeClaimUpdated",
				"updated persistent volume claim")
		}
	}

	return reconcile.Continue()
}

func (r *FastDFSReconciler) getPVCList(ctx context.Context, cluster *v1.FastDFS) (pvList corev1.PersistentVolumeClaimList, err error) {
	pvcList := &corev1.PersistentVolumeClaimList{}
	err = r.List(ctx, pvcList,
		client.InNamespace(cluster.Namespace), client.MatchingLabels(cluster.ResourceMatchingLabels()))
	return *pvcList, err
}

func (r *FastDFSReconciler) cleanupPVCs(ctx context.Context, cluster *v1.FastDFS, replicas int32) error {
	pvcList, err := r.getPVCList(ctx, cluster)
	if err != nil {
		return err
	}

	for _, pvcItem := range pvcList.Items {
		// delete only Orphan PVCs
		if isPVCOrphan(pvcItem.Name, replicas) {
			r.Log.Info("removing orphan pvc", "pvc", pvcItem)
			if err := r.deletePVC(pvcItem); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *FastDFSReconciler) deletePVC(pvcItem corev1.PersistentVolumeClaim) error {
	pvcDelete := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pvcItem.Name,
			Namespace: pvcItem.Namespace,
		},
	}
	return r.Delete(context.TODO(), pvcDelete)
}
