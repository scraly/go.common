/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package broker

// Configuration is the struct which you can import in your configuration struct and have it working with templateV2
type Configuration struct {
	Use                string `toml:"use" default:"noop" comment:"Broker to use : nats, stan, rabbitmq, kafka, noop"`
	Hosts              string `toml:"hosts" default:"" comment:"Broker cluster hosts"`
	CertificatePath    string `toml:"certificatePath" default:"" comment:"Certificate path"`
	PrivateKeyPath     string `toml:"privateKeyPath" default:"" comment:"Private Key path"`
	CACertificatePath  string `toml:"caCertificatePath" default:"" comment:"CA Certificate Path"`
	InsecureSkipVerify bool   `toml:"insecureSkipVerify" default:"false" comment:"Disable insecure certificates verification. Ex: self signed certificates usages."`
}
