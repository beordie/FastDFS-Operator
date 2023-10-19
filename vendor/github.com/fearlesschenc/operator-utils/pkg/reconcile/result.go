package reconcile

import (
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type Result struct {
	RequeueDelay         time.Duration
	RequeueRequest       bool
	CancelReconciliation bool
}

// Continue to continue the execution of the whole reconciliation
func Continue() (result Result, err error) {
	result = Result{
		RequeueDelay:         0,
		RequeueRequest:       false,
		CancelReconciliation: false,
	}
	return
}

// StopOnError stop the whole reconciliation when error happen
func StopOnError(errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:         0,
		RequeueRequest:       false,
		CancelReconciliation: true,
	}
	err = errIn
	return
}

// Stop the whole reconciliation
func Stop() (result Result, err error) {
	return StopOnError(nil)
}

func requeue(delay time.Duration, errIn error) (result Result, err error) {
	result = Result{
		RequeueDelay:         delay,
		RequeueRequest:       true,
		CancelReconciliation: false,
	}
	err = errIn
	return
}

func RequeueOnError(errIn error) (result Result, err error) {
	return requeue(0, errIn)
}

func Requeue() (result Result, err error) {
	return RequeueOnError(nil)
}

// RequeueAfter will requeue request after delay.
func RequeueAfter(delay time.Duration, errIn error) (result Result, err error) {
	return requeue(delay, errIn)
}

// controller-runtime results

func DoNotRequeueRequest(err error) (reconcile.Result, error) {
	return reconcile.Result{}, err
}

func RequeueRequestOnErr(err error) (reconcile.Result, error) {
	// note: reconcile will auto requeue failed request
	return reconcile.Result{}, err
}

func RequeueRequestAfter(duration time.Duration, err error) (reconcile.Result, error) {
	return reconcile.Result{RequeueAfter: duration}, err
}