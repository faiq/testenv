package test2

import (
	"context"
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api/controllers/remote"
	addonsv1 "sigs.k8s.io/cluster-api/exp/addons/api/v1alpha3"

	crscontrollers "sigs.k8s.io/cluster-api/exp/addons/controllers"
	"sigs.k8s.io/cluster-api/test/helpers"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func TestMain(t *testing.T) {
	testEnv := helpers.NewTestEnvironment()
	fmt.Printf("%v \n", testEnv)

	trckr, err := remote.NewClusterCacheTracker(log.NullLogger{}, testEnv.Manager)
	if err != nil {
		t.Error(err)
	}

	err = (&crscontrollers.ClusterResourceSetReconciler{
		Log:     log.Log,
		Client:  testEnv,
		Tracker: trckr,
	}).SetupWithManager(testEnv.Manager, controller.Options{MaxConcurrentReconciles: 1})

	if err != nil {
		t.Error(err)
	}

	err = (&crscontrollers.ClusterResourceSetBindingReconciler{
		Log:    log.Log,
		Client: testEnv,
	}).SetupWithManager(testEnv.Manager, controller.Options{MaxConcurrentReconciles: 1})
	if err != nil {
		t.Error(err)
	}

	go func() {
		err := testEnv.StartManager()
		if err != nil {
			t.Error(err)
		}
	}()
	clusterResourceSetInstance := &addonsv1.ClusterResourceSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "crs",
			Namespace: "default",
		},
		Spec: addonsv1.ClusterResourceSetSpec{
			ClusterSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"test": "foo",
				},
			},
			Resources: []addonsv1.ResourceRef{{Name: "crscm", Kind: "ConfigMap"}, {Name: "crssecret", Kind: "Secret"}},
		},
	}
	// Create the ClusterResourceSet.
	err = testEnv.Create(context.TODO(), clusterResourceSetInstance)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(clusterResourceSetInstance)
}
