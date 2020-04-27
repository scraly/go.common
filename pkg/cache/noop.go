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

import (
	"time"
)

type noopStore struct {
}

// -----------------------------------------------------------------------------
func (s *noopStore) Name() string {
	return "noop"
}

func (s *noopStore) Connect() error {
	// Return no error
	return nil
}

func (s *noopStore) Get(key string, value interface{}) error {
	// Return no error
	return nil
}

func (s *noopStore) Set(key string, value interface{}, expires time.Duration) error {
	// Return no error
	return nil
}

func (s *noopStore) Add(key string, value interface{}, expires time.Duration) error {
	// Return no error
	return nil
}

func (s *noopStore) Replace(key string, value interface{}, expires time.Duration) error {
	// Return no error
	return nil
}

func (s *noopStore) Delete(key string) error {
	// Return no error
	return nil
}

func (s *noopStore) Flush() error {
	// Return no error
	return nil
}
