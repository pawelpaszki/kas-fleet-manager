package types

import "github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/api"

// OpenIDIdentityProviderInfo information for configuring an OpenID identity provider
type OpenIDIdentityProviderInfo struct {
	ID           string
	Name         string
	ClientID     string
	ClientSecret string
	Issuer       string
}

// IdentityProviderInfo contains information about identity providers
type IdentityProviderInfo struct {
	OpenID *OpenIDIdentityProviderInfo
}

// ResourceSet A set of OpenShift/k8s resources
type ResourceSet struct {
	Name      string
	Resources []interface{}
}

// ClusterRequest information about the cluster creation request
type ClusterRequest struct {
	// cloud provider requirement
	CloudProvider string
	// region of the cluster
	Region string
	// if the cluster is multi-az
	MultiAZ bool
	// AdditionalSpec Additional information that can be used by the cloud provider to use when creating a new OpenShift/k8s cluster
	AdditionalSpec api.JSON
}

// ClusterSpec Information about the openshift/k8s cluster
type ClusterSpec struct {
	// internal id of the cluster. Used when making requests to the provider
	InternalID string `json:"internal_id"`
	// external id of the cluster. Some providers provide an additional id for external usage. If not provided, the InternalID will be used.
	ExternalID string `json:"external_id"`
	// the status of the cluster
	Status api.ClusterStatus `json:"status"`
	// additional information related to the cluster, can vary depending on the provider
	AdditionalInfo api.JSON `json:"additional_info"`
}

type CloudProviderInfo struct {
	ID          string
	Name        string
	DisplayName string
}

type CloudProviderInfoList struct {
	Items []CloudProviderInfo
}

type CloudProviderRegionInfo struct {
	ID              string
	CloudProviderID string
	Name            string
	DisplayName     string
	SupportsMultiAZ bool
}

type CloudProviderRegionInfoList struct {
	Items []CloudProviderRegionInfo
}

type ComputeNodesInfo struct {
	Actual  int
	Desired int
}