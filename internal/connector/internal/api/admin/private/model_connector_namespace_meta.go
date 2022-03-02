/*
 * Connector Service Fleet Manager Admin APIs
 *
 * Connector Service Fleet Manager Admin is a Rest API to manage connector clusters.
 *
 * API version: 0.0.3
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package private

import (
	"time"
)

// ConnectorNamespaceMeta struct for ConnectorNamespaceMeta
type ConnectorNamespaceMeta struct {
	Owner       string                                     `json:"owner,omitempty"`
	CreatedAt   time.Time                                  `json:"created_at,omitempty"`
	ModifiedAt  time.Time                                  `json:"modified_at,omitempty"`
	Name        string                                     `json:"name"`
	Annotations []ConnectorNamespaceRequestMetaAnnotations `json:"annotations,omitempty"`
}