/*
 * Kafka Management API
 *
 * Kafka Management API is a REST API to manage Kafka instances
 *
 * API version: 1.16.0
 * Contact: rhosak-support@redhat.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package public

// VersionMetadata struct for VersionMetadata
type VersionMetadata struct {
	Id            string            `json:"id"`
	Kind          string            `json:"kind"`
	Href          string            `json:"href"`
	ServerVersion string            `json:"server_version,omitempty"`
	Collections   []ObjectReference `json:"collections,omitempty"`
}
