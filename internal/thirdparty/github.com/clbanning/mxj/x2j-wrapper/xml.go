// Copyright 2012-2018 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package x2j

var castNanInf bool

// Cast "Nan", "Inf", "-Inf" XML values to 'float64'.
// By default, these values will be decoded as 'string'.
func CastNanInf(b bool) {
	castNanInf = b
}
