package controller

import (
	v1 "fastdfs_operator/api/v1"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *FastDFSReconciler) makeStatefulSet(cluster *v1.FastDFS) *appsv1.StatefulSet {
	nn := cluster.GetStatefulSetNamespacedName()
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nn.Name,
			Namespace: nn.Namespace,
		},
	}
}

func (r *FastDFSReconciler) makePVCStorageSize(cluster *v1.FastDFS) resource.Quantity {
	return resource.MustParse(fmt.Sprintf(v1.StorageValueUnit, cluster.Spec.Storage.DiskSize, cluster.Spec.Storage.Unit))
}

func (r *FastDFSReconciler) mutateStatefulSet(cluster *v1.FastDFS, sts *appsv1.StatefulSet) error {
	if sts.ObjectMeta.CreationTimestamp.IsZero() {
		// sts resource was not created yet, or happened any error
		sts.ObjectMeta.Labels = cluster.ResourceLabels()
		sts.Spec.Selector = &metav1.LabelSelector{MatchLabels: cluster.ResourceMatchingLabels()}
		sts.Spec.ServiceName = cluster.GetHeadlessServiceName()
		sts.Spec.PodManagementPolicy = appsv1.ParallelPodManagement

		if sts.Spec.VolumeClaimTemplates == nil || len(sts.Spec.VolumeClaimTemplates) == 0 {
			sts.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{{}}
		}
		pvc := &sts.Spec.VolumeClaimTemplates[0]
		pvc.Name = v1.PvcName
		pvc.Namespace = cluster.GetNamespace()
		pvc.Labels = cluster.ResourceLabels()
		pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}
		pvc.Spec.Resources = corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: r.makePVCStorageSize(cluster),
			},
		}
		pvc.Spec.StorageClassName = cluster.Spec.Storage.StorageClass
	}
	return nil
}
