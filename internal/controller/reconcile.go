package controller

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	ValuePaused = "true"
)

func ShouldReconcile(object metav1.Object) bool {
	annotations := object.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	return annotations[AnnotationPauseReconcile] != ValuePaused
}

func RequeueRequestAfter(duration time.Duration, err error) (reconcile.Result, error) {
	return reconcile.Result{RequeueAfter: duration}, err
}
