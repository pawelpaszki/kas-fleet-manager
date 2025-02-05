/*
 * Connector Service Fleet Manager
 *
 * Connector Service Fleet Manager is a Rest API to manage connectors.
 *
 * API version: 0.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package public

// ConnectorConfiguration struct for ConnectorConfiguration
type ConnectorConfiguration struct {
	Kafka          KafkaConnectionSettings          `json:"kafka"`
	ServiceAccount ServiceAccount                   `json:"service_account"`
	SchemaRegistry SchemaRegistryConnectionSettings `json:"schema_registry,omitempty"`
	Connector      map[string]interface{}           `json:"connector"`
}
