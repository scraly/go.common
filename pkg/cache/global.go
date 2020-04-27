/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package cache

var (
	globalCache Store = &noopStore{}
)

// SetGlobalBroker sets the [singleton] instance of the broker
func SetGlobalBroker(store Store) {
	globalCache = store
}

// GlobalStore returns the cache manager instance
func GlobalStore() Store {
	return globalCache
}
