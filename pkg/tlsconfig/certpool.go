/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package tlsconfig

import (
	"crypto/x509"
)

// SystemCertPool returns an new empty cert pool,
// accessing system cert pool
func SystemCertPool() (*x509.CertPool, error) {
	return x509.NewCertPool(), nil
}
