package controller

import (
	"context"
	v1 "fastdfs_operator/api/v1"

	"github.com/fearlesschenc/operator-utils/pkg/reconcile"
)

func (r *FastDFSReconciler) reconcilePods(ctx context.Context, cluster *v1.FastDFS) (reconcile.Result, error) {
	return reconcile.Result{}, nil
}
