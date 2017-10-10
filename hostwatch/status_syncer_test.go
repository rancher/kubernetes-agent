package hostwatch

import (
	"testing"

	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/api/v1"
)

func TestActiveInactive(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := kubernetesclient.NewClient(kubeURL)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test1",
			Hostname: "test1",
			State:    DeactivatingState,
		},
	}

	kubeHandler.nodes["test1"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test1",
		},
		Spec: v1.NodeSpec{
			Unschedulable: false,
		},
		Status: v1.NodeStatus{
			Conditions: []v1.NodeCondition{
				v1.NodeCondition{
					Type:   v1.NodeReady,
					Status: v1.ConditionTrue,
				},
			},
		},
	}

	statusSync(kubeClient, metadataClient)

	if kubeHandler.nodes["test1"].Spec.Unschedulable == false {
		t.Error("Node test1 was not cordoned correctly")
	}
}

func TestInactiveActive(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := kubernetesclient.NewClient(kubeURL)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test2",
			Hostname: "test2",
			State:    ActivatingState,
		},
	}

	kubeHandler.nodes["test2"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test2",
		},
		Spec: v1.NodeSpec{
			Unschedulable: true,
		},
		Status: v1.NodeStatus{
			Conditions: []v1.NodeCondition{
				v1.NodeCondition{
					Type:   v1.NodeReady,
					Status: v1.ConditionTrue,
				},
			},
		},
	}

	statusSync(kubeClient, metadataClient)

	if kubeHandler.nodes["test2"].Spec.Unschedulable == true {
		t.Error("Node test2 was not uncordoned correctly")
	}
}

func TestEvacuate(t *testing.T) {
	metadataClient := metadata.NewClient(fakeMetadataURL)
	kubeClient := kubernetesclient.NewClient(kubeURL)

	metadataHandler.hosts = []metadata.Host{
		{
			Name:     "test3",
			Hostname: "test3",
			State:    EvacuatingState,
		},
	}

	kubeHandler.nodes["test3"] = &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test3",
			Labels: map[string]string{
				"node": "test3",
			},
		},
		Spec: v1.NodeSpec{
			Unschedulable: false,
		},
		Status: v1.NodeStatus{
			Conditions: []v1.NodeCondition{
				v1.NodeCondition{
					Type:   v1.NodeReady,
					Status: v1.ConditionTrue,
				},
			},
		},
	}

	statusSync(kubeClient, metadataClient)

	if kubeHandler.nodes["test3"].Labels[DrainLabelName] != DrainLabelValue {
		t.Error("Node test3 was not evacuated correctly")
	}
}
