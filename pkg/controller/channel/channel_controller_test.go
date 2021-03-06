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

package channel

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	chv1 "github.com/open-cluster-management/multicloud-operators-channel/pkg/apis/apps/v1"
)

var c client.Client

var targetNamespace = "default"
var tragetChannelName = "foo"
var targetChannelType = chv1.ChannelType("namespace")

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: tragetChannelName, Namespace: targetNamespace}}

const timeout = time.Second * 5

func TestChannelControllerReconcile(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	channelInstance := &chv1.Channel{
		ObjectMeta: metav1.ObjectMeta{Name: tragetChannelName, Namespace: targetNamespace},
		Spec: chv1.ChannelSpec{
			Type:     targetChannelType,
			Pathname: targetNamespace,
		},
	}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: "0"})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	c = mgr.GetClient()

	//create events handler on hub cluster. All the deployable events will be written to the root deploable on hub cluster.
	hubClientSet, _ := kubernetes.NewForConfig(cfg)

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: hubClientSet.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: "channel"})

	recFn, requests := SetupTestReconcile(newReconciler(mgr, recorder))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	// Create the Channel object and expect the Reconcile
	err = c.Create(context.TODO(), channelInstance)

	g.Expect(err).NotTo(gomega.HaveOccurred())

	defer c.Delete(context.TODO(), channelInstance)

	time.Sleep(time.Second * 1)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))
}
