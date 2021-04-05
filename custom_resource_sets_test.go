package main

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	addonsv1 "sigs.k8s.io/cluster-api/exp/addons/api/v1alpha3"
)

var _ = Describe("CNI Cluster Resource Sets", func() {
	It("Should be able to install", func() {
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
		Expect(testEnv.Create(context.TODO(), clusterResourceSetInstance)).To(Succeed())
	})
})
