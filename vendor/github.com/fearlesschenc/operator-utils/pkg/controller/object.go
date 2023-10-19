package controller

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func IsObjectBeingDeleted(obj metav1.Object) bool {
	return !obj.GetDeletionTimestamp().IsZero()
}

// IsObjectHaveFinalizer returns whether this object has the passed finalizer
func IsObjectHaveFinalizer(obj metav1.Object, finalizer string) bool {
	for _, fin := range obj.GetFinalizers() {
		if fin == finalizer {
			return true
		}
	}
	return false
}
