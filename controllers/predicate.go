// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

/*the purpose of this module to filter every event which doesn't have any meaning for the applications.
eg: status updates, restart cases, etc*/
package controllers

import (
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	app "smartmobilelabs.com/evo/evo-operator/api/v1alpha1"
	"smartmobilelabs.com/evo/evo-operator/controllers/finalizer"
)

type SmlCustomPredicate struct{}

// predicate filter added to extend for future extensions

func smlEventPredicate() predicate.Predicate {
	return SmlCustomPredicate{}
}

func (SmlCustomPredicate) Create(event event.CreateEvent) bool {
	logger := log.WithName("predicate").WithName("create_event")
	instance, ok := event.Object.(*app.SmlEvo)

	logger.V(1).Info("Event received", "object", instance)

	if ok {
		logger.Info("Event can be reconciled")
		return true
	}

	logger.V(1).Info("Event skipped")
	return false
}

func (SmlCustomPredicate) Delete(event event.DeleteEvent) bool {
	logger := log.WithName("predicate").WithName("delete_event")
	logger.V(1).Info("Event received")
	instance, ok := event.Object.(*app.SmlEvo)

	if ok && finalizer.HasFinalizers(instance) {
		logger.Info("Event can be reconciled")
		return true
	}
	logger.V(1).Info("Event skipped")
	return false
}

func (SmlCustomPredicate) Update(event event.UpdateEvent) bool {
	logger := log.WithName("predicate").WithName("update_event")
	logger.V(1).Info("Event received")

	oldInstance, okOld := event.ObjectOld.(*app.SmlEvo)
	newInstance, okNew := event.ObjectNew.(*app.SmlEvo)

	if okOld && okNew {
		logger.V(1).Info("New object content", "object", newInstance)
		logger.V(1).Info("Old object content", "object", oldInstance)
		if isChangeInSpec(oldInstance, newInstance) {
			logger.Info("Event can be reconciled")
			return true
		} else if isFinalizerAddition(oldInstance, newInstance) {
			logger.Info("Finalizer added allow the reconciliation")
			return true
		} else if isDeleteEvent(oldInstance, newInstance) {
			logger.Info("DeleteTimestamp changed")
			return true
		}
	}
	logger.V(1).Info("Update event skipped")
	return false
}

func isChangeInSpec(oldInstance *app.SmlEvo, newInstance *app.SmlEvo) bool {
	return !reflect.DeepEqual(oldInstance.Spec, newInstance.Spec)
}

func isFinalizerAddition(oldInstance *app.SmlEvo, newInstance *app.SmlEvo) bool {
	return !finalizer.HasFinalizers(oldInstance) && finalizer.HasFinalizers(newInstance)
}

func isDeleteEvent(oldInstance *app.SmlEvo, newInstance *app.SmlEvo) bool {
	return oldInstance.DeletionTimestamp == nil && newInstance.DeletionTimestamp != nil
}

func (SmlCustomPredicate) Generic(event.GenericEvent) bool {
	logger := log.WithName("predicate").WithName("generic event")
	logger.V(1).Info("Generic event received, skip it.")
	return false
}
