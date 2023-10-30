package utils

import (
	"context"
	"strings"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetReferencesFromStringSlice generate LocalObjectReference slice from string consist of secrets,
// a slice of length 0 will not make statefulset rolling update.
func GetReferencesFromStringSlice(in string) []corev1.LocalObjectReference {
	slice := strings.Split(in, ";")
	refs := make([]corev1.LocalObjectReference, 0, len(slice))
	for i := range slice {
		if name := strings.TrimSpace(slice[i]); name != "" {
			refs = append(refs, corev1.LocalObjectReference{Name: name})
		}
	}
	return refs
}

func GetJobFinishedStatus(job *batchv1.Job) (bool, batchv1.JobConditionType) {
	for _, cond := range job.Status.Conditions {
		if (cond.Type == batchv1.JobComplete || cond.Type == batchv1.JobFailed) && cond.Status == corev1.ConditionTrue {
			return true, cond.Type
		}
	}
	return false, ""
}

func FinalizeJobPod(ctx context.Context, cli client.Client, job *batchv1.Job) error {
	pods := &corev1.PodList{}
	if err := cli.List(ctx, pods, client.MatchingLabels(job.Labels), client.InNamespace(job.Namespace)); err != nil {
		return err
	}

	for _, pod := range pods.Items {
		if err := cli.Delete(ctx, &pod); err != nil && !apierrors.IsNotFound(err) {
			return err
		}
	}

	return nil
}

func GetDefaultStorageClass(ctx context.Context, client client.Client) (*storagev1.StorageClass, error) {
	scs := &storagev1.StorageClassList{}
	if err := client.List(ctx, scs); err != nil {
		return nil, err
	}
	for _, sc := range scs.Items {
		if sc.Annotations["storageclass.kubernetes.io/is-default-class"] == "true" {
			return &sc, nil
		}
	}
	return nil, nil
}
