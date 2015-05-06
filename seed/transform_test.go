// Copyright (c) 2013-2015 by The SeedStack authors. All rights reserved.

// This file is part of SeedStack, An enterprise-oriented full development stack.

// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	verbose = true
	vverbose = true
	res := m.Run()
	os.Exit(res)
}

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

func TestFile(t *testing.T) {
	tg := Transformation{Filter: "*.go"}
	tgy := Transformation{Filter: "*.go|*.yml"}
	matched := checkFileName("test\\src\\bla\\bla\\cmd.go", tg) &&
		!checkFileName("./test/bloat.java.", tg)
	if !matched {
		t.Errorf("The file ./test/cmd.go should match the pattern '*.go' but not 'bloat.java'")
	}

	matched = checkFileName("./test/cmd.go", tgy) &&
		checkFileName("./test/conf.yml", tgy) &&
		!checkFileName("./test/bloat.java.", tgy)
	if !matched {
		t.Errorf("The file 'cmd.go' and 'conf.yml' should match the pattern '*.go|*.yml' but not 'bloat.java'")
	}
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

func TestReplace(t *testing.T) {
	var p *Procedures
	news := string(p.Replace([]byte("foo"), "foo", "bar", "bar", "toto"))
	if news != "toto" {
		t.Errorf("Procedure should replace 'foo' with 'toto' but %s was found\n", news)
	}
}

func TestInsertAndRemove(t *testing.T) {
	var p *Procedures
	ori := []byte("foo")
	
	inc := p.Insert(ori, "bar")
	if  string(inc) != "foobar" {
		t.Errorf("removeAtEnd: %s was expected but found %s", "foobar", inc)
	}

	clean := string(p.RemoveAtEnd(inc, "bar"))
	if clean != "foo" {
		t.Errorf("removeAtEnd: %s was expected but found %s", ori, clean)
	}
}

func TestReplaceMavenDependency(t *testing.T) {
	var p *Procedures
	news := string(p.ReplaceMavenDependency([]byte(pom), "com.inetpsa.fnd:seed-bom", "org.seedstack:bom", "org.seedstack:bom", "org.seedstack:seedstack-bom"))
	if news != expectedPom {
		t.Errorf("Procedure should replace 'com.inetpsa.fnd:seed-bom' with 'org.seedstack:seedstack-bom' but found:\n %s", news)
	}
}

var pom = `

    <dependencyManagement>
        <dependencies>
            <!-- SEED Distribution BOM-->
            <dependency>
                <groupId>com.inetpsa.fnd</groupId>  <!-- test ->
               <artifactId>seed-bom</artifactId>
                <version>14.11</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>
`

var expectedPom = `

    <dependencyManagement>
        <dependencies>
            <!-- SEED Distribution BOM-->
            <dependency>
                <groupId>org.seedstack</groupId>  <!-- test ->
               <artifactId>seedstack-bom</artifactId>
                <version>14.11</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>
`

func TestMatchDependency(t *testing.T) {
	old := "com.inetpsa.fnd:seed-bom"
	new := "org.seedstack:seedstack-bom"

	result := matchDependency(pom, old, new)
	if result != expectedPom {
		fmt.Println("found:\n" + result)
		t.Error("Fail to replace maven dependency")
	}
}

var expectedPomWithVersion = `

    <dependencyManagement>
        <dependencies>
            <!-- SEED Distribution BOM-->
            <dependency>
                <groupId>org.seedstack</groupId>  <!-- test ->
               <artifactId>seedstack-bom</artifactId>
                <version>15.4-M2-SNAPSHOT</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>
`

func TestMatchDependencyWithVersion(t *testing.T) {
	old := "com.inetpsa.fnd:seed-bom:14.11"
	new := "org.seedstack:seedstack-bom:15.4-M2-SNAPSHOT"

	result := matchDependency(pom, old, new)
	if result != expectedPomWithVersion {
		fmt.Println("found:\n" + result)
		t.Error("Fail to replace maven dependency with version")
	}
}

func TestReplaceMavenDependencyWithVersion(t *testing.T) {
	var p *Procedures
	news := string(p.ReplaceMavenDependency([]byte(pom), "com.inetpsa.fnd:seed-bom:14.11", "org.seedstack:bom:15.4", "org.seedstack:bom:15.4", "org.seedstack:seedstack-bom:15.4-M2-SNAPSHOT"))
	if news != expectedPomWithVersion {
		t.Errorf("Procedure should replace 'com.inetpsa.fnd:seed-bom:14.11' with 'org.seedstack:seedstack-bom:15.4-M2-SNAPSHOT' but found:\n %s", news)
	}
}

var pomWithProperty = `
    <properties>
        <!-- blabla -->
		<seed-bom.version>14.11.2</seed-bom.version>
	</properties>

    <dependencyManagement>
        <dependencies>
            <!-- SEED Distribution BOM-->
            <dependency>
                <groupId>com.inetpsa.fnd</groupId>  <!-- test ->
               <artifactId>seed-bom</artifactId>
                <version>${seed-bom.version}</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>
`

var expectedPomWithProperty = `
    <properties>
        <!-- blabla -->
		<seed-bom.version>15.4-M2-SNAPSHOT</seed-bom.version>
	</properties>

    <dependencyManagement>
        <dependencies>
            <!-- SEED Distribution BOM-->
            <dependency>
                <groupId>org.seedstack</groupId>  <!-- test ->
               <artifactId>seedstack-bom</artifactId>
                <version>${seed-bom.version}</version>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>
`

func TestMatchDependencyWithVersionAndProps(t *testing.T) {
	old := "com.inetpsa.fnd:seed-bom:14.11"
	new := "org.seedstack:seedstack-bom:15.4-M2-SNAPSHOT"

	result := matchDependency(pomWithProperty, old, new)
	if result != expectedPomWithProperty {
		fmt.Println("found:\n" + result)
		t.Error("Fail to replace maven dependency with version")
	}
}

var expectedPomWithRemovedVersion = `

    <dependencyManagement>
        <dependencies>
            <!-- SEED Distribution BOM-->
            <dependency>
                <groupId>org.seedstack</groupId>  <!-- test ->
               <artifactId>seedstack-bom</artifactId>
                <type>pom</type>
                <scope>import</scope>
            </dependency>
        </dependencies>
    </dependencyManagement>
`

func TestMatchDependencyWithRemovedVersion(t *testing.T) {
	old := "com.inetpsa.fnd:seed-bom:*"
	new := "org.seedstack:seedstack-bom"

	result := matchDependency(pom, old, new)
	if result != expectedPomWithRemovedVersion {
		fmt.Println("found:\n" + result)
		t.Error("Fail to replace maven dependency and removing its version")
	}
}
