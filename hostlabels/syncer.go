package hostlabels

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	cache "github.com/patrickmn/go-cache"

	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"github.com/rancher/kubernetes-model/model"
	"k8s.io/kubernetes/pkg/util/validation"
)

type hostLabelSyncer struct {
	kClient            *kubernetesclient.Client
	metadataClient     *metadata.Client
	cache              *cache.Cache
	cacheExpiryMinutes time.Duration
}

func (h *hostLabelSyncer) syncHostLabels(version string) {
	err := sync(h.kClient, h.metadataClient, h.cache)
	if err != nil {
		log.Errorf("Error syncing host labels: [%v]", err)
	}
}

func getKubeNode(kClient *kubernetesclient.Client, hostname string) (*model.Node, error) {
	node, err := kClient.Node.ByName(hostname)
	if err != nil {
		log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", hostname, err)
		// This node might not have been added to kuberentes cluster yet, so skip it
	}
	return node, err
}

func sync(kClient *kubernetesclient.Client, metadataClient *metadata.Client, c *cache.Cache) error {
	hosts, err := metadataClient.GetHosts()
	if err != nil {
		log.Errorf("Error reading host list from metadata service: [%v], retrying", err)
		return err
	}
	for _, host := range hosts {
		nodeInt, ok := c.Get(host.Hostname)
		if !ok {
			temp_node, err := getKubeNode(kClient, host.Hostname)
			if err != nil {
				log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", host.Hostname, err)
				// This node might not have been added to kuberentes cluster yet, so skip it
				continue
			}
			c.Set(host.Hostname, temp_node, 0)
			nodeInt = temp_node
		}
		node := nodeInt.(*model.Node)
		if node.Metadata.Annotations == nil {
			node.Metadata.Annotations = make(map[string]interface{})
		}
		rancherLabelsMetadataStore := node.Metadata.Annotations
		changed := false
		//check for new/updated labels
		for k, v1 := range host.Labels {
			if !isValidLabelValue(v1) {
				continue
			}
			if changed {
				break
			}
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
		}
		for k := range node.Metadata.Labels {
			if changed {
				break
			}
			if _, ok := rancherLabelsMetadataStore[toKMetaLabel(k)]; !ok {
				// This is not a rancher managed label
				continue
			}
			if _, ok := host.Labels[k]; !ok {
				changed = true
			}
		}
		retryCount := 0
		maxRetryCount := 3
		for changed {
			node, err := getKubeNode(kClient, host.Hostname)
			if err != nil {
				log.Errorf("Error getting node: [%s] by name from kubernetes: [%v]", host.Hostname, err)
				continue
			}
			c.Set(host.Hostname, node, 0)
			if node.Metadata.Annotations == nil {
				node.Metadata.Annotations = make(map[string]interface{})
			}
			rancherLabelsMetadataStore := node.Metadata.Annotations
			for k, v1 := range host.Labels {
				if !isValidLabelValue(v1) {
					log.Infof("skipping invalid label %s=%s", k, v1)
					continue
				}
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
				}
			}

			_, err = kClient.Node.ReplaceNode(node)
			if err != nil {
				log.Errorf("Error updating node [%s] with new host labels, err :[%v]", host.Hostname, err)
				if retryCount < maxRetryCount {
					retryCount = retryCount + 1
					continue
				}
			}
			changed = false
		}
	}
	return nil
}

func isValidLabelValue(label string) bool {
	errs := validation.IsValidLabelValue(label)
	if len(errs) > 0 {
		return false
	}
	return true
}

func toKMetaLabel(label string) string {
	return fmt.Sprintf("%s.%s", rancherLabelKey, label)
}
