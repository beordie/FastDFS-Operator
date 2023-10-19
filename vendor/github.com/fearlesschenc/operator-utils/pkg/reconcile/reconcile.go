package reconcile

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/fearlesschenc/operator-utils/pkg/controller"
)

type job struct {
	ctx    context.Context
	object metav1.Object
}

func Reconcile(ctx context.Context, obj metav1.Object) *job {
	return &job{ctx: ctx, object: obj}
}

func (j *job) WithReconciler(reconciler Reconciler) (result Result, err error) {
	steps := []Func{}
	if initializer, ok := reconciler.(Initializer); ok {
		steps = append(steps, initializer.Initialize)
	}
	if finalizer, ok := reconciler.(Finalizer); ok {
		steps = append(steps, finalizer.Finalize)
	}
	if validator, ok := reconciler.(Validator); ok {
		steps = append(steps, validator.Validate)
	}
	if genericSteps := reconciler.GetReconcileSteps(); genericSteps != nil && len(genericSteps) > 0 {
		steps = append(steps, genericSteps...)
	}

	func() {
		if statusUpdater, ok := reconciler.(StatusUpdater); ok {
			defer func() {
				if controller.IsObjectBeingDeleted(j.object) {
					return
				}

				if updateErr := statusUpdater.UpdateStatus(j.ctx, j.object); updateErr != nil {
					err = updateErr
				}
			}()
		}

		for _, f := range steps {
			result, err = f(j.ctx, j.object)
			if err != nil || result.RequeueRequest || result.CancelReconciliation {
				return
			}
		}
	}()

	return
}
