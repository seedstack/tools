// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
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
//
// proc:
//  -
//    name: Insert
//    params: "endOfFile"
func (p *Procedures) Insert(dat []byte, s string) []byte {
	return append(dat, []byte(s)...)
}

// Replace the old string by the new one. You can use it as follows in your transformation file.
//
// proc:
//  -
//    name: Replace
//    params:
//      - "myStringToModify"
//      - "myModifiedString"
//      # After you can add other pairs
//      - "x"
//      - "y"
//      ...
func (p *Procedures) Replace(dat []byte, pairs ...string) []byte {
	new := dat
	for i := 0; i < len(pairs); i += 2 {
		new = []byte(strings.Replace(string(new), pairs[i], pairs[i+1], -1))
		if vverbose && bytes.Compare(new, dat) != 0 {
			fmt.Printf("\t%s -> %s\n", pairs[i], pairs[i+1])
		}
	}

	return new
}

// ReplaceMavenDependency replaces a maven dependency by a new one.
// The dependency to update are passed as pairs. For instance you want to update the following dependency:
//
// <dependency>
//   <groupId>org.mycompany</groupId>
//   <artifactId>myApp1</artifactId>
// </dependency>
//
// Call the ReplaceMavenDependency with a pair of the old dependency and the new one.
//
// proc:
//  -
//    name: ReplaceMavenDependency
//    params:
//      - "org.mycompany:myApp1"
//      - "com.mycompany:myApp2"
//      # After you can add other pairs
//      - "x:x"
//      - "y:y"
//      ...
func (p *Procedures) ReplaceMavenDependency(data []byte, pairs ...string) []byte {
	for i := 0; i < len(pairs); i += 2 {
		data = []byte(matchDependency(string(data), pairs[i], pairs[i+1]))
	}
	return data
}

// ReplaceMavenDependencyWithVersion replaces a maven dependency by a new one including the version tag.
// The dependency to update are passed as pairs. For instance you want to update the following dependency:
//
// <dependency>
//   <groupId>org.mycompany</groupId>
//   <artifactId>myApp1</artifactId>
//   <version>1.0.0</version>
// </dependency>
//
// Call the ReplaceMavenDependencyWithVersion with a pair of the old dependency and the new one.
//
// proc:
//  -
//    name: ReplaceMavenDependencyWithVersion
//    params:
//      - "org.mycompany:myApp1:1.0.0"
//      - "com.mycompany:myApp2:2.0.0"
//      # After you can add other pairs
//      - "x:x:x"
//      - "y:y:y"
//      ...
func (p *Procedures) ReplaceMavenDependencyWithVersion(data []byte, pairs ...string) []byte {
	for i := 0; i < len(pairs); i += 2 {
		data = []byte(matchDependencyWithVersion(string(data), pairs[i], pairs[i+1]))
	}
	return data
}

func matchDependencyWithVersion(pom, old, new string) string {
	currentDep := strings.Split(old, ":")
	newDep := strings.Split(new, ":")

	if len(currentDep) != 3 && len(newDep) != 3 {
		log.Fatalf(`ReplaceMavenDependencyWithVersion takes dependencies with`+
			` the following format "groupId:artifactId:vesion". But "%s" and "%s" where found.`, old, new)
	}

	regex := "(<groupId>)" + currentDep[0] + "(<\\/groupId>.*?\\n.*?" +
		"<artifactId>)" + currentDep[1] + "(<\\/artifactId>.*?\\n.*?" +
		"<version>)" + currentDep[2] + "(<\\/version>)"

	depRegex := regexp.MustCompile(regex)

	return depRegex.ReplaceAllString(pom, "${1}"+newDep[0]+"${2}"+newDep[1]+"${3}"+newDep[2]+"${4}")
}

func matchDependencyWithVersionAndProps(pom, old, new string) string {
	currentDep := strings.Split(old, ":")
	newDep := strings.Split(new, ":")

	if len(currentDep) != 3 && len(newDep) != 3 {
		log.Fatalf(`ReplaceMavenDependencyWithVersion takes dependencies with`+
			` the following format "groupId:artifactId:vesion". But "%s" and "%s" where found.`, old, new)
	}

	regex := "(<groupId>)" + regexp.QuoteMeta(currentDep[0]) + "(<\\/groupId>.*?\\n.*?" +
		"<artifactId>)" + regexp.QuoteMeta(currentDep[1]) + "(<\\/artifactId>.*?\\n.*?" +
		"<version>)(.*?)(<\\/version>)"

	depRegex := regexp.MustCompile(regex)

	propsDefinition := regexp.MustCompile("\\$\\{(.*?)\\}")
	// Find the version of the dependency to replace
	version := depRegex.FindStringSubmatch(pom)[4]
	// Check if it is a property
	match := propsDefinition.FindStringSubmatch(version)
	fmt.Printf("match %v", match)
	var props string
	if match != nil {
		props = match[1]
	}

	if props != "" {
		propsToReplace := regexp.MustCompile("(<" + regexp.QuoteMeta(props) + ">).*?(</" + regexp.QuoteMeta(props) + ">)")
		pom = propsToReplace.ReplaceAllString(pom, "${1}"+newDep[2]+"${2}")
		return depRegex.ReplaceAllString(pom, "${1}"+newDep[0]+"${2}"+newDep[1]+"${3}"+"${4}"+"${5}")
		// "(<seed-bom\.version>).*?(<\/seed-bom\.version>)"
	} else {
		return depRegex.ReplaceAllString(pom, "${1}"+newDep[0]+"${2}"+newDep[1]+"${3}"+newDep[2]+"${5}")
	}
}

func matchDependency(pom, old, new string) string {
	currentDep := strings.Split(old, ":")
	newDep := strings.Split(new, ":")

	if len(currentDep) != 2 && len(newDep) != 2 {
		log.Fatalf(`ReplaceMavenDependencyWithVersion takes dependencies with`+
			` the following format "groupId:artifactId". But "%s" and "%s" where found.`, old, new)
	}

	regex := "(<groupId>)" + currentDep[0] + "(<\\/groupId>.*?\\n.*?" +
		"<artifactId>)" + currentDep[1] + "(<\\/artifactId>)"

	depRegex := regexp.MustCompile(regex)

	return depRegex.ReplaceAllString(pom, "${1}"+newDep[0]+"${2}"+newDep[1]+"${3}")
}

func matchDependencyAndRemoveVersion(pom, old, new string) string {
	currentDep := strings.Split(old, ":")
	newDep := strings.Split(new, ":")

	if len(currentDep) != 2 && len(newDep) != 2 {
		log.Fatalf(`ReplaceMavenDependencyWithVersion takes dependencies with`+
			` the following format "groupId:artifactId". But "%s" and "%s" where found.`, old, new)
	}

	regex := "(<groupId>)" + currentDep[0] + "(<\\/groupId>.*?\\n.*?" +
		"<artifactId>)" + currentDep[1] + "(<\\/artifactId>.*?)\\n.*?" +
		"<version>.*?<\\/version>"

	depRegexWithVersion := regexp.MustCompile(regex)
	if depRegexWithVersion.FindString(pom) != "" {
		return depRegexWithVersion.ReplaceAllString(pom, "${1}"+newDep[0]+"${2}"+newDep[1]+"${3}")
	} else {
		regex := "(<groupId>)" + currentDep[0] + "(<\\/groupId>.*?\\n.*?" +
			"<artifactId>)" + currentDep[1] + "(<\\/artifactId>)"

		depRegex := regexp.MustCompile(regex)

		return depRegex.ReplaceAllString(pom, "${1}"+newDep[0]+"${2}"+newDep[1]+"${3}")
	}
}
