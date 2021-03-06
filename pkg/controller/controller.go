// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controller

import (
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	gitsync "github.com/open-cluster-management/multicloud-operators-channel/pkg/synchronizer/githubsynchronizer"
	helmsync "github.com/open-cluster-management/multicloud-operators-channel/pkg/synchronizer/helmreposynchronizer"
	"github.com/open-cluster-management/multicloud-operators-channel/pkg/utils"
)

// AddToManagerFuncs is a list of functions to add all Controllers to the Manager
var AddToManagerFuncs []func(manager.Manager, record.EventRecorder, *utils.ChannelDescriptor, *helmsync.ChannelSynchronizer, *gitsync.ChannelSynchronizer) error

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager, recorder record.EventRecorder,
	chdesc *utils.ChannelDescriptor, sync *helmsync.ChannelSynchronizer,
	gsync *gitsync.ChannelSynchronizer) error {
	for _, f := range AddToManagerFuncs {
		if err := f(m, recorder, chdesc, sync, gsync); err != nil {
			return err
		}
	}

	return nil
}
