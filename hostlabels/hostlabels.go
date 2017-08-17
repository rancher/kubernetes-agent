package hostlabels

import (
	"fmt"
	"os"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/kubernetes-agent/kubernetesclient"

	log "github.com/Sirupsen/logrus"
)

const (
	metadataURLTemplate = "http://%v/2015-12-19"
	rancherLabelKey     = "io.rancher.labels"
	cacheExpiryMinutes  = 5 * time.Minute

	// DefaultMetadataAddress specifies the default value to use if nothing is specified
	DefaultMetadataAddress = "169.254.169.250"
)

// StartHostLabelSync ...
func StartHostLabelSync(interval int, kClient *kubernetesclient.Client) error {
	metadataAddress := os.Getenv("RANCHER_METADATA_ADDRESS")
	if metadataAddress == "" {
		metadataAddress = DefaultMetadataAddress
	}
	metadataURL := fmt.Sprintf(metadataURLTemplate, metadataAddress)

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
