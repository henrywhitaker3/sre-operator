package handlers

import (
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	fn string = "sre.henrywhitaker.com/finalizer"
)

type Handler interface {
	// Check the object exists
	Get() error
	// Code to execute when the resource is created/updated
	CreateOrUpdate() (error, bool)
	// Code to execute when the resource is deleted
	Delete() error
	// Is the objects deletion timestamp zero?
	DeletionTimestampIsZero() bool
	// Get the objects finalizers
	AddFinalizer(string) error
	// Remove a given finalizer from the object
	RemoveFinalizer(string) error
	// Update the resource with a successful status
	SuccessStatus() error
	// Update the resource with a failed status
	ErrorStatus(error) error
}

func RunHandler(l logr.Logger, h Handler) (reconcile.Result, error) {
	if err := h.Get(); err != nil {
		if client.IgnoreNotFound(err) != nil {
			l.Error(err, "could not find webgook configuration")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	l.Info("processing resource")

	if !h.DeletionTimestampIsZero() {
		l.Info("deleting resource")
		if err := h.Delete(); err != nil {
			h.ErrorStatus(err)
			return ctrl.Result{}, err
		}

		err := h.RemoveFinalizer(fn)
		return ctrl.Result{}, err
	}

	if err := h.AddFinalizer(fn); err != nil {
		return ctrl.Result{}, err
	}

	if err, prop := h.CreateOrUpdate(); err != nil {
		h.ErrorStatus(err)
		if !prop {
			err = nil
		}
		return ctrl.Result{}, err
	}

	err := h.SuccessStatus()
	return ctrl.Result{}, err
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
