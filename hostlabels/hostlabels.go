package hostlabels

import (
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"

	log "github.com/Sirupsen/logrus"
)

const (
	metadataURL     = "http://rancher-metadata/2015-12-19"
	rancherLabelKey = "io.rancher.labels"
)

func StartHostLabelSync(interval int, kClient *kubernetesclient.Client) error {
	metadataClient, err := metadata.NewClientAndWait(metadataURL)
	if err != nil {
		log.Errorf("Error initializing metadata client: [%v]", err)
		return err
	}
	h := &hostLabelSyncer{
		kClient:        kClient,
		metadataClient: metadataClient,
	}
	metadataClient.OnChange(interval, h.syncHostLabels)
	return nil
}
