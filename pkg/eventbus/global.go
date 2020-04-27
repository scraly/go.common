/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package eventbus

var (
	globalBus = NewLocal()
)

// SetGlobalBus sets the [singleton] instance of the event bus
func SetGlobalBus(bus EventBus) {
	globalBus = bus
}

// GlobalBus returns the event bus instance
func GlobalBus() EventBus {
	return globalBus
}
