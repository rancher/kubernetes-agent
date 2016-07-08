package hostlabels

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
)

type hostLabelSyncer struct {
	kClient        *kubernetesclient.Client
	metadataClient *metadata.Client
}

func (h *hostLabelSyncer) syncHostLabels(version string) {
	err := sync(h.kClient, h.metadataClient)
	if err != nil {
		log.Errorf("Error syncing host labels: [%v]", err)
	}
}

func sync(kClient *kubernetesclient.Client, metadataClient *metadata.Client) error {
	hosts, err := metadataClient.GetHosts()
	if err != nil {
		log.Errorf("Error reading host list from metadata service: [%v], retrying", err)
		return err
	}
	for _, host := range hosts {
		node, err := kClient.Node.ByName(host.Hostname)
		if err != nil {
			log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", host.Hostname, err)
			// This node might not have been added to kuberentes cluster yet, so skip it
			continue
		}
		if node.Metadata.Annotations == nil {
			node.Metadata.Annotations = make(map[string]interface{})
		}
		rancherLabelsMetadataStore := node.Metadata.Annotations
		changed := false
		//check for new/updated labels
		for k, v1 := range host.Labels {
			v2, ok := node.Metadata.Labels[k]
			if !ok {
				// This label doesn't exist
				changed = true
			} else {
				v2String, ok := v2.(string)
				if !ok {
					changed = true
				} else if v1 != v2String {
					changed = true
				}
			}
			if changed {
				node.Metadata.Labels[k] = v1
				rancherLabelsMetadataStore[toKMetaLabel(k)] = ""
			}
			for k := range node.Metadata.Labels {
				if _, ok := rancherLabelsMetadataStore[toKMetaLabel(k)]; !ok {
					// This is not a rancher managed label
					continue
				}
				if _, ok := host.Labels[k]; !ok {
					delete(node.Metadata.Labels, k)
					delete(rancherLabelsMetadataStore, toKMetaLabel(k))
					changed = true
				}
			}
			if changed {
				_, err := kClient.Node.ReplaceNode(node)
				if err != nil {
					log.Errorf("Error updating node [%s] with new host labels, err :[%v]", host.Hostname, err)
					return err
				}
				changed = false
			}
		}
	}
	return nil
}

func toKMetaLabel(label string) string {
	return fmt.Sprintf("%s.%s", rancherLabelKey, label)
}
