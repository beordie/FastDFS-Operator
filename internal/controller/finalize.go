package controller

import (
	"context"
	v1 "fastdfs_operator/api/v1"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func (r *FastDFSReconciler) isPVCBeingDeleted(ctx context.Context, cluster *v1.FastDFS, replicas int32) (deleted bool, err error) {
	var pvcList corev1.PersistentVolumeClaimList
	pvcList, err = r.getPVCList(ctx, cluster)
	if err != nil {
		r.Log.WithFields(cluster.Fields()).Error("Failed to get PVC list")
		return false, err
	}
	for _, pvcItem := range pvcList.Items {
		if isPVCOrphan(pvcItem.Name, replicas) {
			continue
		}

		if !pvcItem.GetDeletionTimestamp().IsZero() {
			deleted = true
			break
		}
	}
	return
}

// To determine whether a Kubernetes Persistent Volume Claim (PVC) is isolated (that is, not associated with any Kubernetes resource),
// the specific method of determination is to compare the ordinal value in the name of the PVC with the expected number of replicas.
func isPVCOrphan(pvcName string, replicas int32) bool {
	index := strings.LastIndexAny(pvcName, "-")
	if index == -1 {
		return false
	}

	ordinal, err := strconv.Atoi(pvcName[index+1:])
	if err != nil {
		return false
	}

	return int32(ordinal) >= replicas
}
