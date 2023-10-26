package controller

import (
	"context"
	v1 "fastdfs_operator/api/v1"

	"github.com/fearlesschenc/operator-utils/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *FastDFSReconciler) reconcilePods(ctx context.Context, cluster *v1.FastDFS) (reconcile.Result, error) {
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.MatchingLabels(cluster.ResourceMatchingLabels())); err != nil {
		return reconcile.RequeueOnError(err)
	}

	gvk := cluster.GroupVersionKind()
	isController := false
	blockOwnerDeletion := false
	clusterRef := metav1.OwnerReference{
		APIVersion:         gvk.GroupVersion().String(),
		Kind:               gvk.Kind,
		Name:               cluster.Name,
		UID:                cluster.GetUID(),
		Controller:         &isController,
		BlockOwnerDeletion: &blockOwnerDeletion,
	}
	for _, pod := range podList.Items {
		// add OwnerReference
		haveOwnerReference := false
		for _, ref := range pod.GetOwnerReferences() {
			haveOwnerReference = ref.APIVersion == clusterRef.APIVersion &&
				ref.Kind == clusterRef.Kind &&
				ref.Name == clusterRef.Name
			if haveOwnerReference {
				break
			}
		}
		if !haveOwnerReference {
			pod.OwnerReferences = append(pod.OwnerReferences, clusterRef)
			if err := r.Update(ctx, &pod); err != nil {
				return reconcile.RequeueOnError(err)
			}
		}
	}
	return reconcile.Continue()
}
