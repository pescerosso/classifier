/*
Copyright 2022. projectsveltos.io. All rights reserved.

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

package fv_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	libsveltosv1alpha1 "github.com/projectsveltos/libsveltos/api/v1alpha1"
)

var _ = Describe("Classifier: deployment", func() {
	const (
		namePrefix = "deploy"
	)

	It("Deploy Classifier instance in CAPI clusters", Label("FV"), func() {
		clusterLabels := map[string]string{randomString(): randomString(), randomString(): randomString()}
		classifier := getClassifier(namePrefix, clusterLabels)

		Byf("Creating classifier instance %s in the management cluster", classifier.Name)
		Expect(k8sClient.Create(context.TODO(), classifier)).To(Succeed())

		Byf("Getting client to access the workload cluster")
		workloadClient, err := getKindWorkloadClusterKubeconfig()
		Expect(err).To(BeNil())
		Expect(workloadClient).ToNot(BeNil())

		Byf("Verifying Classifier CRD is installed in the workload cluster")
		Eventually(func() error {
			classifierCRD := &apiextensionsv1.CustomResourceDefinition{}
			return workloadClient.Get(context.TODO(),
				types.NamespacedName{Name: "classifiers.lib.projectsveltos.io"}, classifierCRD)
		}, timeout, pollingInterval).Should(BeNil())

		Byf("Verifying ClassifierReport CRD is installed in the workload cluster")
		Eventually(func() error {
			classifierReportCRD := &apiextensionsv1.CustomResourceDefinition{}
			return workloadClient.Get(context.TODO(),
				types.NamespacedName{Name: "classifierreports.lib.projectsveltos.io"}, classifierReportCRD)
		}, timeout, pollingInterval).Should(BeNil())

		Byf("Verifying Classifier instance is deployed in the workload cluster")
		Eventually(func() error {
			currentClassifier := &libsveltosv1alpha1.Classifier{}
			return workloadClient.Get(context.TODO(),
				types.NamespacedName{Name: classifier.Name}, currentClassifier)
		}, timeout, pollingInterval).Should(BeNil())

		Byf("Deleting classifier instance %s in the management cluster", classifier.Name)
		currentClassifier := &libsveltosv1alpha1.Classifier{}
		Expect(k8sClient.Get(context.TODO(),
			types.NamespacedName{Name: classifier.Name}, currentClassifier)).To(Succeed())
		Expect(k8sClient.Delete(context.TODO(), currentClassifier)).To(Succeed())

		Byf("Verifying Classifier instance is removed from the workload cluster")
		Eventually(func() bool {
			currentClassifier := &libsveltosv1alpha1.Classifier{}
			err = workloadClient.Get(context.TODO(),
				types.NamespacedName{Name: classifier.Name}, currentClassifier)
			return err != nil && apierrors.IsNotFound(err)
		}, timeout, pollingInterval).Should(BeTrue())

		Byf("Verifying Classifier instance is removed from the management cluster")
		Eventually(func() bool {
			currentClassifier := &libsveltosv1alpha1.Classifier{}
			err = k8sClient.Get(context.TODO(),
				types.NamespacedName{Name: classifier.Name}, currentClassifier)
			return err != nil && apierrors.IsNotFound(err)
		}, timeout, pollingInterval).Should(BeTrue())
	})

})
