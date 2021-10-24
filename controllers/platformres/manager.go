// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package platformres

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
	"os"
	"smartmobilelabs.com/evo/evo-operator/controllers/k8sdynamic"
	kubelib2 "smartmobilelabs.com/evo/evo-operator/controllers/kubelib"
	"strings"
	"sync"
	"time"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("helm")

const (
	ResourceRequestPath = "RESREQ_DIR"

	StatusField         = "status"
	ApprovalStatusField = "approvalStatus"
)

func ApplyPlatformResourceRequests(namespace string) ([]k8sdynamic.ResourceDescriptor, error) {
	logger := log.WithName("ApplyPlatformResourceRequests")
	logger.Info("Called")

	dynClient := k8sdynamic.New(kubelib2.GetKubeAPI())

	dir := os.Getenv(ResourceRequestPath)
	if dir == "" {
		logger.Info("directory not set")
		return nil, errors.New(ResourceRequestPath + " is not set")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Info("cannot read dir")
		return nil, errors.Wrapf(err, "failed to read dir: %v", dir)
	}

	var descList []k8sdynamic.ResourceDescriptor
	for _, file := range files {
		if !file.IsDir() {
			logger.Info("read file " + file.Name())
			fileContent, err := ioutil.ReadFile(dir + "/" + file.Name())
			if err != nil {
				return nil, errors.Wrap(err, "failed to read file")
			}

			if strings.TrimSpace(string(fileContent)) == "" {
				logger.Info("File is empty skip it", "path", dir+"/"+file.Name())
				continue
			}
			logger.Info("try to apply request in  " + file.Name())
			resourceDesc, err := dynClient.ApplyYamlResource(string(fileContent), namespace)
			descList = append(descList, resourceDesc)
			if err != nil {
				logger.Info("error for resource with name  " + file.Name())
				return nil, errors.Wrap(err, "failed to apply the request in k8s")
			}
		}
	}

	return descList, nil
}

func WaitUntilResourcesGranted(resourceList []k8sdynamic.ResourceDescriptor, timeout time.Duration) error {
	logger := log.WithName("WaitUntilResourcesGranted")

	var stopperList []chan struct{}
	var waitGroup sync.WaitGroup
	var results []*bool

	for _, resource := range resourceList {
		stopper := make(chan struct{})
		waitGroup.Add(1)
		var result bool
		results = append(results, &result)
		startWatchResourceRequest(
			resource.Name,
			resource.Namespace,
			"",
			resource.Gvr.GetGvr(),
			stopper,
			&waitGroup,
			&result,
		)
		stopperList = append(stopperList, stopper)
	}

	if waitTimeout(&waitGroup, timeout) {
		for _, stopper := range stopperList {
			close(stopper)
		}
		return errors.New("waiting for the approval of the platform resource requests timed out")
	} else {
		for _, result := range results {
			if !*result {
				return errors.New("some of the platform resource request(s) has rejected")
			}
		}
	}
	logger.Info("All of the requested platform resources have been granted")
	return nil
}

const (
	ApprovalStatusApproved = "Approved"
	ApprovalStatusRejected = "Rejected"
)

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		wg.Wait()
	}()
	select {
	case <-finished:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func startWatchResourceRequest(name string, namespace string, resourceVersion string, gvr schema.GroupVersionResource, stopper chan struct{}, waitGroup *sync.WaitGroup, result *bool) {
	logger := log.WithName("StartWatchResourceRequest").WithValues("resource", name)

	logger.Info("Watching resource")

	go k8sdynamic.WatchInformer(
		name, namespace, resourceVersion, gvr,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				close(stopper)
				waitGroup.Done()
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				logger.V(1).Info("Resource request modification detected")

				newValue, _ := getApprovalStatus(newObj)
				oldValue, _ := getApprovalStatus(oldObj)

				if oldValue != newValue {
					switch newValue {
					case ApprovalStatusApproved:
						*result = true
						break
					default:
						*result = false
						break
					}
					waitGroup.Done()
					close(stopper)
				}
			},
			AddFunc: func(obj interface{}) {
				value, _ := getApprovalStatus(obj)
				if value == ApprovalStatusApproved {
					*result = true
					waitGroup.Done()
					close(stopper)
				}
			},
		},
		stopper)
}
func getApprovalStatus(obj interface{}) (string, bool) {
	unstructObj := obj.(*unstructured.Unstructured)
	value, found, _ := unstructured.NestedString(unstructObj.Object, StatusField, ApprovalStatusField)
	return value, found
}
