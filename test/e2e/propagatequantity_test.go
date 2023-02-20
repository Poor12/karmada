package e2e

import (
	"context"
	"fmt"
	"time"

	"github.com/karmada-io/karmada/test/e2e/framework"
	testhelper "github.com/karmada-io/karmada/test/helper"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BasicPropagation focus on basic propagation functionality testing.
var _ = ginkgo.Describe("namespace scope resource propagation quantity testing", func() {
	var deploymentSlice []*appsv1.Deployment
	var quantityNum int
	var targetClusters []string
	var startTime, endTime time.Time

	ginkgo.BeforeEach(func() {
		quantityNum = 1000
		testNamespace = "default"
		targetClusters = []string{"member1"}

		deploymentSlice = make([]*appsv1.Deployment, quantityNum)

		for index := 0; index < quantityNum; index++ {
			deploymentName := fmt.Sprintf("%s-%d", deploymentNamePrefix, index)
			deployment := testhelper.NewDeployment(testNamespace, deploymentName)
			deploymentSlice[index] = deployment
		}
	})

	ginkgo.BeforeEach(func() {
		for index := 0; index < quantityNum; index++ {
			go func(i int) {
				deployment := deploymentSlice[index]
				framework.CreateDeployment(kubeClient, deployment)
			}(index)
		}
		startTime = time.Now()
	})

	ginkgo.AfterEach(func() {
		endTime = time.Now()
		fmt.Printf("Running time: %f\n", endTime.Sub(startTime).Minutes())
	})

	ginkgo.It("deployment propagation testing", func() {
		for _, cluster := range targetClusters {
			clusterClient := framework.GetClusterClient(cluster)
			gomega.Expect(clusterClient).ShouldNot(gomega.BeNil())

			gomega.Eventually(func() bool {
				deploys, err := clusterClient.AppsV1().Deployments(testNamespace).List(context.TODO(), metav1.ListOptions{})
				if err != nil {
					return false
				}

				return len(deploys.Items) == quantityNum
			}, pollTimeout, pollInterval).Should(gomega.Equal(true))
		}
	})
})
