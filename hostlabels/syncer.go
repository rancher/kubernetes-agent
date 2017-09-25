package hostlabels

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	cache "github.com/patrickmn/go-cache"

	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/client-go/pkg/api/v1"
)

type hostLabelSyncer struct {
	kClient            *kubernetesclient.Client
	metadataClient     metadata.Client
	cache              *cache.Cache
	cacheExpiryMinutes time.Duration
}

func (h *hostLabelSyncer) syncHostLabels(version string) {
	err := sync(h.kClient, h.metadataClient, h.cache)
	if err != nil {
		log.Errorf("Error syncing host labels: [%v]", err)
	}
}

func getKubeNode(kClient *kubernetesclient.Client, hostname string) (*v1.Node, error) {
	node, err := kClient.Node.ByName(hostname)
	if err != nil {
		log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", hostname, err)
		// This node might not have been added to kuberentes cluster yet, so skip it
	}
	return node, err
}

func sync(kClient *kubernetesclient.Client, metadataClient metadata.Client, c *cache.Cache) error {
	hosts, err := metadataClient.GetHosts()
	if err != nil {
		log.Errorf("Error reading host list from metadata service: [%v], retrying", err)
		return err
	}
	for _, host := range hosts {
		nodeInt, ok := c.Get(host.Hostname)
		if !ok {
			tempNode, err := getKubeNode(kClient, host.Hostname)
			if err != nil {
				log.Errorf("Error getting node: [%s] by name from kubernetes, err: [%v]", host.Hostname, err)
				// This node might not have been added to kuberentes cluster yet, so skip it
				continue
			}
			c.Set(host.Hostname, tempNode, 0)
			nodeInt = tempNode
		}
		node := nodeInt.(*v1.Node)
		if node.Annotations == nil {
			node.Annotations = make(map[string]string)
		}
		rancherLabelsMetadataStore := node.Annotations
		changed := false
		//check for new/updated labels
		for k, v1 := range host.Labels {
			if !isValidLabelValue(v1) {
				continue
			}
			if changed {
				break
			}
			v2, ok := node.Labels[k]
			if !ok {
				// This label doesn't exist
				changed = true
			} else {
				v2String := v2
				if !ok {
					changed = true
				} else if v1 != v2String {
					changed = true
				}
			}
		}
		for k := range node.Labels {
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
			if node.Annotations == nil {
				node.Annotations = make(map[string]string)
			}
			rancherLabelsMetadataStore := node.Annotations
			for k, v1 := range host.Labels {
				if !isValidLabelValue(v1) {
					log.Infof("skipping invalid label %s=%s", k, v1)
					continue
				}
				node.Labels[k] = v1
				rancherLabelsMetadataStore[toKMetaLabel(k)] = ""

			}
			for k := range node.Labels {
				if _, ok := rancherLabelsMetadataStore[toKMetaLabel(k)]; !ok {
					// This is not a rancher managed label
					continue
				}
				if _, ok := host.Labels[k]; !ok {
					delete(node.Labels, k)
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
