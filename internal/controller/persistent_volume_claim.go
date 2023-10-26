package controller

import (
	"context"
	v1 "fastdfs_operator/api/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
