// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"reflect"
)

// Conditions regroup all the precondition methods
type Conditions struct {}

// Procedures regroup all the procedure methods
type Procedures struct {}

func checkCondition(fileName string, data []byte, t Transformation) bool {
	ok := true
	var c Conditions
	for _, pre := range t.Pre {
		m := reflect.ValueOf(&c).MethodByName(pre)
		ok = m.Call([]reflect.Value{reflect.ValueOf(fileName), reflect.ValueOf(data)})[0].Bool()
		if !ok {
			break
		}
	}
	return ok
}

func applyProcs(data []byte, t Transformation) []byte {
	var p Procedures
	for _, proc := range t.Proc {
		vals := []reflect.Value{reflect.ValueOf(data)}
		for _, param := range proc.Params {
			vals = append(vals, reflect.ValueOf(param))
		}
		
		data = reflect.ValueOf(&p).MethodByName(proc.Name).
			Call(vals)[0].Bytes()
	}
	return data
}

// func (c *Conditions) AlwaysTrue(fileName string, data []byte) bool {
// 	return true
// }

// func (p *Procedures) Insert(dat []byte, s string) []byte {
// 	return append(dat, []byte(s)...)
// }
