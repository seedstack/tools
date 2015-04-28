// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// Conditions regroup all the precondition methods
type Conditions struct{}

// Procedures regroup all the procedure methods
type Procedures struct{}

func checkFileName(fileName string, tr Transformation) bool {
	matched := false
	// Include files
	for _, patt := range strings.Split(tr.Filter, "|") {
		res, err := filepath.Match(patt, filepath.Base(fileName))
		matched = res || matched
		if err != nil {
			log.Fatalf("Failed to parse pattern: %s\n%v", tr.Filter, err)
		}
	}
	return matched
}

func checkCondition(fileName string, data []byte, t Transformation) bool {
	ok := true
	var c Conditions
	for _, pre := range t.Pre {
		m := reflect.ValueOf(&c).MethodByName(pre)
		if m.IsValid() {
			ok = m.Call([]reflect.Value{reflect.ValueOf(fileName), reflect.ValueOf(data)})[0].Bool()
		} else {
			log.Fatalf(`Cannot find the precondition method "%s"`, pre)
		}
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
		m := reflect.ValueOf(&p).MethodByName(proc.Name)
		if !m.IsValid() {
			log.Fatalf("Cannot find method to proc name: %s\n", proc.Name)
		}
		data = m.Call(vals)[0].Bytes()
	}
	return data
}

// -----------------

// AlwaysTrue is a precondition which will be true for all the files.
func (c *Conditions) AlwaysTrue(fileName string, data []byte) bool {
	return true
}

// -----------------

// Insert the string s at the end of the given data.
func (p *Procedures) Insert(dat []byte, s string) []byte {
	return append(dat, []byte(s)...)
}

// Replace the old string by the new one.
func (p *Procedures) Replace(dat []byte, pairs ...string) []byte {
	for i := 0; i < len(pairs); i += 2 {
		dat = []byte(strings.Replace(string(dat), pairs[i], pairs[i+1], -1))
	}
	return dat
}

// ReplaceMavenDependency replaces a maven dependency by a new one.
func (p *Procedures) ReplaceMavenDependency(data []byte, pairs ...string) []byte {
	for i := 0; i < len(pairs); i += 2 {
		data = []byte(matchDependency(string(data), pairs[i], pairs[i+1]))
	}
	return data
}

// ReplaceMavenDependencyWithVersion replaces a maven dependency by a new one including the version tag.
func (p *Procedures) ReplaceMavenDependencyWithVersion(data []byte, pairs ...string) []byte {
	for i := 0; i < len(pairs); i += 2 {
		data = []byte(matchDependencyWithVersion(string(data), pairs[i], pairs[i+1]))
	}
	return data
}

func matchDependencyWithVersion(pom, old, new string) string {
	currentDep := strings.Split(old, ":")
	newDep := strings.Split(new, ":")

	regex := "(<groupId>)" + currentDep[0] + "(<\\/groupId>.*?\\n.*?" +
		"<artifactId>)" + currentDep[1] + "(<\\/artifactId>.*?\\n.*?" +
		"<version>)" + currentDep[2] + "(<\\/version>)"

	depRegex := regexp.MustCompile(regex)

	return depRegex.ReplaceAllString(pom, "${1}"+newDep[0]+"${2}"+newDep[1]+"${3}"+newDep[2]+"${4}")
}

func matchDependency(pom, old, new string) string {
	currentDep := strings.Split(old, ":")
	newDep := strings.Split(new, ":")

	regex := "(<groupId>)" + currentDep[0] + "(<\\/groupId>.*?\\n.*?" +
		"<artifactId>)" + currentDep[1] + "(<\\/artifactId>)"

	depRegex := regexp.MustCompile(regex)

	return depRegex.ReplaceAllString(pom, "${1}"+newDep[0]+"${2}"+newDep[1]+"${3}")
}
