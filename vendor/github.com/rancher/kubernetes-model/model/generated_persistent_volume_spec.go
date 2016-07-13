package model

const (
	PERSISTENT_VOLUME_SPEC_TYPE = "v1.PersistentVolumeSpec"
)

type PersistentVolumeSpec struct {
	AccessModes []PersistentVolumeAccessMode `json:"accessModes,omitempty" yaml:"access_modes,omitempty"`

	AwsElasticBlockStore *AWSElasticBlockStoreVolumeSource `json:"awsElasticBlockStore,omitempty" yaml:"aws_elastic_block_store,omitempty"`

	AzureFile *AzureFileVolumeSource `json:"azureFile,omitempty" yaml:"azure_file,omitempty"`

	Capacity map[string]interface{} `json:"capacity,omitempty" yaml:"capacity,omitempty"`

	Cephfs *CephFSVolumeSource `json:"cephfs,omitempty" yaml:"cephfs,omitempty"`

	Cinder *CinderVolumeSource `json:"cinder,omitempty" yaml:"cinder,omitempty"`

	ClaimRef *ObjectReference `json:"claimRef,omitempty" yaml:"claim_ref,omitempty"`

	Fc *FCVolumeSource `json:"fc,omitempty" yaml:"fc,omitempty"`

	FlexVolume *FlexVolumeSource `json:"flexVolume,omitempty" yaml:"flex_volume,omitempty"`

	Flocker *FlockerVolumeSource `json:"flocker,omitempty" yaml:"flocker,omitempty"`

	GcePersistentDisk *GCEPersistentDiskVolumeSource `json:"gcePersistentDisk,omitempty" yaml:"gce_persistent_disk,omitempty"`

	Glusterfs *GlusterfsVolumeSource `json:"glusterfs,omitempty" yaml:"glusterfs,omitempty"`

	HostPath *HostPathVolumeSource `json:"hostPath,omitempty" yaml:"host_path,omitempty"`

	Iscsi *ISCSIVolumeSource `json:"iscsi,omitempty" yaml:"iscsi,omitempty"`

	Nfs *NFSVolumeSource `json:"nfs,omitempty" yaml:"nfs,omitempty"`

	PersistentVolumeReclaimPolicy string `json:"persistentVolumeReclaimPolicy,omitempty" yaml:"persistent_volume_reclaim_policy,omitempty"`

	Rbd *RBDVolumeSource `json:"rbd,omitempty" yaml:"rbd,omitempty"`
}
