package reconcile

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Func func(ctx context.Context, object metav1.Object) (Result, error)

type Funcs []Func

func (fns Funcs) GetReconcileSteps() []Func {
	return fns
}

var _ Reconciler = Funcs{}

type Reconciler interface {
	GetReconcileSteps() []Func
}

type StatusUpdater interface {
	UpdateStatus(ctx context.Context, object metav1.Object) error
}

type Validator interface {
	Validate(ctx context.Context, object metav1.Object) (Result, error)
}

type Initializer interface {
	Initialize(ctx context.Context, object metav1.Object) (Result, error)
}

type Finalizer interface {
	Finalize(ctx context.Context, object metav1.Object) (Result, error)
}
