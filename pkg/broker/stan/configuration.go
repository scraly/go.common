/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package stan

// FullConfName is the key to use in the context to provide full specific configuration
var FullConfName = "StanFullConfiguration"

// SubsConfName is the key to use in the context to provide subscribe specific configuration
var SubsConfName = "StanSubsConfiguration"

// Configuration is the struct which you can import in your configuration struct and have it working with templateV2
// This is for the specificity of stan broker. For the common Configuration struct, refer to broker.Configuration
type Configuration struct {
	ClientID  string        `toml:"clientid" comment:"ClientID used to subscribe to NATS Streaming"`
	ClusterID string        `toml:"clusterid" comment:"ClusterID used to subscribe to NATS Streaming"`
	Subject   string        `toml:"subject" comment:"subject to subscribe in NATS Streaming consumer"`
	User      string        `toml:"user" comment:"User used to authenticate with NATS Streaming cluster"`
	Password  string        `toml:"password" comment:"password used to authenticate with NATS Streaming cluster"`
	Subscribe SubscribeOpts `toml:"subscribe" comment:"###############################\n Subscribe options \n##############################"`
}

// SubscribeOpts contains the
type SubscribeOpts struct {
	DurableName         string `toml:"durablename" comment:"durable name to set when subscribe in NATS Streaming consumer (could be generate with go get -u github.com/docker/docker/pkg/namesgenerator/cmd/names-generator)"`
	QueueGroup          string `toml:"queuegroup" default:"" comment:"Queue group name shen subscribe in NATS Streaming consumer"`
	DeliverAllAvailable bool   `toml:"deliverAllAvailable" default:"false" comment:"deliver all messages in subject when subscribe in NATS Streaming consumer"`
	ManualAcks          bool   `toml:"manualAcks" default:"false" comment:"Manually ack on new message received"`
	StartSequence       uint64 `toml:"startSequence,omitempty" default:"0" comment:"This should never be managed by static configuration. Optional start sequence number. "`
}
