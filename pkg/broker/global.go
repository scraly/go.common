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

var (
	globalBroker Broker = noopBroker{}
)

// SetGlobalBroker sets the [singleton] instance of the broker
func SetGlobalBroker(broker Broker) {
	globalBroker = broker
}

// GlobalBroker returns the broker instance
func GlobalBroker() Broker {
	return globalBroker
}
