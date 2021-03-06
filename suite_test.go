/*
Copyright 2020 The Kubernetes Authors.
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

package main

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/cluster-api/controllers/remote"
	crscontrollers "sigs.k8s.io/cluster-api/exp/addons/controllers"
	"sigs.k8s.io/cluster-api/test/helpers"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	testEnv *helpers.TestEnvironment //nolint:gochecknoglobals //upstream test standard
	ctx     = context.Background()   //nolint:gochecknoglobals //upstream test standard
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	By("bootstrapping test environment")
	testEnv = helpers.NewTestEnvironment()
	trckr, err := remote.NewClusterCacheTracker(log.NullLogger{}, testEnv.Manager)
	Expect(err).NotTo(HaveOccurred())

	Expect((&crscontrollers.ClusterResourceSetReconciler{
		Log:     log.Log,
		Client:  testEnv,
		Tracker: trckr,
	}).SetupWithManager(testEnv.Manager, controller.Options{MaxConcurrentReconciles: 1})).To(Succeed())

	Expect((&crscontrollers.ClusterResourceSetBindingReconciler{
		Log:    log.Log,
		Client: testEnv,
	}).SetupWithManager(testEnv.Manager, controller.Options{MaxConcurrentReconciles: 1})).To(Succeed())

	By("starting the manager")
	go func() {
		defer GinkgoRecover()
		Expect(testEnv.StartManager()).To(Succeed())
	}()

	close(done)
}, 60)

var _ = AfterSuite(func() {
	if testEnv != nil {
		By("tearing down the test environment")
		Expect(testEnv.Stop()).To(Succeed())
	}
})
