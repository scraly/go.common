/*
 * Copyright (c) Continental Corporation - All Rights Reserved
 *
 * This file is a part of Entry project.
 * ITS France - Entry squad members
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package ern

import (
	"errors"
	"strings"
)

const (
	ernDelimiter = ":"
	ernSections  = 4
	ernPrefix    = "ern:"

	// zero-indexed
	sectionTenant   = 1
	sectionType     = 2
	sectionResource = 3

	// errors
	invalidPrefix   = "ern: invalid prefix"
	invalidSections = "ern: not enough sections"
)

// ERN represents Entry Resource Name components
// Example ERNs:
// ern:12:account:123456789
// ern:12:group:entry/administrators
// ern:12:application:123456789
type ERN struct {
	Tenant   string
	Type     string
	Resource string
}

// Parse parses an ERN into its constituent parts.
func Parse(ern string) (ERN, error) {
	if !strings.HasPrefix(ern, ernPrefix) {
		return ERN{}, errors.New(invalidPrefix)
	}
	sections := strings.SplitN(ern, ernDelimiter, ernSections)
	if len(sections) != ernSections {
		return ERN{}, errors.New(invalidSections)
	}
	return ERN{
		Tenant:   sections[sectionTenant],
		Type:     sections[sectionType],
		Resource: sections[sectionResource],
	}, nil
}

// String returns the canonical representation of the ERN
func (ern ERN) String() string {
	return ernPrefix +
		ern.Tenant + ernDelimiter +
		ern.Type + ernDelimiter +
		ern.Resource
}
