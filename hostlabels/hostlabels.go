package hostlabels

import (
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"

	log "github.com/Sirupsen/logrus"
)

const (
	metadataURL        = "http://rancher-metadata/2015-12-19"
	rancherLabelKey    = "io.rancher.labels"
	cacheExpiryMinutes = 5 * time.Minute
)

func StartHostLabelSync(interval int, kClient *kubernetesclient.Client) error {
	metadataClient, err := metadata.NewClientAndWait(metadataURL)
	if err != nil {
		log.Errorf("Error initializing metadata client: [%v]", err)
		return err
	}
	expiringCache := cache.New(cacheExpiryMinutes, 1*time.Minute)
	h := &hostLabelSyncer{
		kClient:        kClient,
		metadataClient: metadataClient,
		cache:          expiringCache,
	}
	metadataClient.OnChange(interval, h.syncHostLabels)
	return nil
}
