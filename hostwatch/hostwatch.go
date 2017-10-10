package hostwatch

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
	// ActivatingState specifies the default value of activating state of host in Ranceher Metadata
	ActivatingState = "activating"
	// DeactivatingState specifies the default value of deactivating state of host in Ranceher Metadata
	DeactivatingState = "deactivating"
	// EvacuatingState specifies the default value of evacuating state of host in Ranceher Metadata
	EvacuatingState = "evacuating"
)

type hostSyncer struct {
	kClient            *kubernetesclient.Client
	metadataClient     metadata.Client
	cache              *cache.Cache
	cacheExpiryMinutes time.Duration
}

// StartHostSync ...
func StartHostSync(interval int, kClient *kubernetesclient.Client) error {
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

	// Start Host Label Syncer
	expiringCache := cache.New(cacheExpiryMinutes, 1*time.Minute)
	h := &hostSyncer{
		kClient:        kClient,
		metadataClient: metadataClient,
		cache:          expiringCache,
	}
	metadataClient.OnChange(interval, h.syncHosts)
	return nil
}

func (h *hostSyncer) syncHosts(version string) {
	err := labelSync(h.kClient, h.metadataClient, h.cache)
	if err != nil {
		log.Errorf("Error syncing host labels: [%v]", err)
	}
	statusSync(h.kClient, h.metadataClient)
}
