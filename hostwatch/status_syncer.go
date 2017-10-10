package hostwatch

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubectld/cli"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
)

var (
	maxRetryCount int = 3
)

const (
	DrainLabelName   = "io.rancher.host.state"
	DrainLabelValue  = "evacuated"
	DrainSyncTimeout = "120s"
)

func statusSync(kClient *kubernetesclient.Client, metadataClient metadata.Client) {
	hosts, err := metadataClient.GetHosts()
	if err != nil {
		log.Errorf("Error reading host list from metadata service: [%v]", err)
		return
	}
	for _, host := range hosts {
		mNodeStatus := host.State
		switch mNodeStatus {
		case ActivatingState:
			cordonUncordon(host, kClient, metadataClient, false)
		case DeactivatingState:
			cordonUncordon(host, kClient, metadataClient, true)
		case EvacuatingState:
			drain(host, kClient)
		}
	}
}

func cordonUncordon(host metadata.Host, kClient *kubernetesclient.Client, metadataClient metadata.Client, unschedulable bool) {
	changed := false
	for retryCount := 0; retryCount <= maxRetryCount; retryCount++ {
		node, err := kClient.Node.ByName(host.Hostname)
		if err != nil {
			log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", host.Hostname, err)
			continue
		}
		if node.Spec.Unschedulable == unschedulable {
			changed = true
			break
		}
		node.Spec.Unschedulable = unschedulable
		delete(node.ObjectMeta.Labels, DrainLabelName)
		_, err = kClient.Node.ReplaceNode(node)
		if err != nil {
			log.Errorf("Error updating node [%s] with new schedulable state, err :[%v]", host.Hostname, err)
			continue
		}
		changed = true
		break
	}
	if !changed {
		log.Errorf("Failed to cordon/uncordon node: [%s]", host.Hostname)
	}
}

func drain(host metadata.Host, kClient *kubernetesclient.Client) {
	node, err := kClient.Node.ByName(host.Hostname)
	if err != nil {
		log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", host.Hostname, err)
		return
	}
	if state, ok := node.ObjectMeta.Labels[DrainLabelName]; !ok || (state != DrainLabelValue) {
		err = execDrainCmd(host.Hostname)
		if err != nil {
			log.Errorf("Error draining node: [%s], err: [%v]", host.Hostname, err)
			return
		}
		node, err = kClient.Node.ByName(host.Hostname)
		if err != nil {
			log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", host.Hostname, err)
			return
		}
		node.ObjectMeta.Labels[DrainLabelName] = DrainLabelValue
		_, err = kClient.Node.ReplaceNode(node)
		if err != nil {
			log.Errorf("Error updating node [%s] with new label, err :[%v]", host.Hostname, err)
			return
		}
	}
}

func execDrainCmd(hostname string) error {
	cmd := "kubectl"
	args := []string{
		"drain",
		hostname,
		"--ignore-daemonsets",
		"--force",
		"--delete-local-data",
		"--timeout",
		DrainSyncTimeout}
	output := cli.Execute(cmd, args...)
	if output.ExitCode > 0 {
		return fmt.Errorf("%s", output.StdErr)
	}
	if output.Err != nil {
		return fmt.Errorf("%s", output.Err)
	}
	fmt.Printf("%s", output.StdOut)
	return nil
}
