/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package prepullinitialized

import (
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
)

// PrepullInitialized is a plugin that checks if node has finished prepull.
type PrepullInitialized struct{}

var _ framework.FilterPlugin = &PrepullInitialized{}

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "PrepullInitialized"
	// ErrReason returned when node name doesn't match.
	ErrReason = "node(s) didn't match prepull initialization filter"
)

// Name returns name of the plugin. It is used in logs, etc.
func (pl *PrepullInitialized) Name() string {
	return Name
}

// Filter invoked at the filter extension point.
func (pl *PrepullInitialized) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
	if nodeInfo.Node() == nil {
		return framework.NewStatus(framework.Error, "node not found")
	}
	if !Fits(pod, nodeInfo) {
		return framework.NewStatus(framework.UnschedulableAndUnresolvable, ErrReason)
	}
	return nil
}

// Fits actually checks if the pod fits the node.
func Fits(pod *v1.Pod, nodeInfo *framework.NodeInfo) bool {
	fmt.Println("=============")

	fmt.Println("Node:")
	fmt.Println(nodeInfo.Node().Name)

	for _, pi := range nodeInfo.Pods {
		fmt.Println(pi.Pod.Name)
		fmt.Println(pi.Pod.Status.Phase)

		if strings.Contains(pi.Pod.Name, "placeholder") && pi.Pod.Status.Phase == "Running" {
			return true
		}
	}
	return false
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, _ framework.FrameworkHandle) (framework.Plugin, error) {
	return &PrepullInitialized{}, nil
}
