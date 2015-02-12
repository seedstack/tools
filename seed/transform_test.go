// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"testing"
)

func TestPrecondition(t *testing.T) {
	tt := Transformation{Pre: []string{"AlwaysTrue"}}
	tf := Transformation{Pre: []string{"AlwaysFalse"}}	
	
	if !checkCondition("", []byte{}, tt) {
		t.Error("Precondition should be always true")
	}
	if checkCondition("", []byte{}, tf) {
		t.Error("Precondition should be always false")
	}

}

func (c *Conditions) AlwaysFalse(fileName string, data []byte) bool {
	return false
}

func TestProcedures(t *testing.T) {
	tn := Transformation{Proc: []Procedure{Procedure{Name: "DoNothing"}}}
	ti := Transformation{Proc: []Procedure{Procedure{Name: "Insert", Params: []string{"bar"}}}}

	res := applyProcs([]byte("foo"), tn)
	if string(res) != "foo" {
		t.Errorf("Procedure should do nothing, %s was expected but found %s", "foo", res)
	}

	res = applyProcs([]byte("foo"), ti)
	if string(res) != "foobar" {
		t.Errorf("Procedure should insert bar, %s was expected but found %s", "foobar", res)
	}

}

func (p *Procedures) DoNothing(dat []byte) []byte {
	return dat
}
