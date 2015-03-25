// Copyright 2015 The Vanadium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package util

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"time"

	"v.io/x/devtools/internal/tool"
)

type BuildCopRotation struct {
	Shifts []struct {
		Primary string `xml:"primary"`
		Date    string `xml:"startDate"`
	} `xml:"shift"`
}

// LoadBuildCopRotation parses the default build cop schedule file.
func LoadBuildCopRotation(ctx *tool.Context) (*BuildCopRotation, error) {
	buildCopRotationsFile, err := BuildCopRotationPath(ctx)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadFile(buildCopRotationsFile)
	if err != nil {
		return nil, fmt.Errorf("ReadFile(%q) failed: %v", buildCopRotationsFile, err)
	}
	rotation := BuildCopRotation{}
	if err := xml.Unmarshal(content, &rotation); err != nil {
		return nil, fmt.Errorf("Unmarshal(%q) failed: %v", string(content), err)
	}
	return &rotation, nil
}

// BuildCop finds the build cop at the given time from the buildcop
// configuration file by comparing timestamps.
func BuildCop(ctx *tool.Context, targetTime time.Time) (string, error) {
	// Parse buildcop.xml file.
	rotation, err := LoadBuildCopRotation(ctx)
	if err != nil {
		return "", err
	}

	// Find the build cop at targetTime.
	layout := "Jan 2, 2006 3:04:05 PM"
	for i := len(rotation.Shifts) - 1; i >= 0; i-- {
		shift := rotation.Shifts[i]
		t, err := time.Parse(layout, shift.Date)
		if err != nil {
			return "", fmt.Errorf("Parse(%q, %v) failed: %v", layout, shift.Date, err)
		}
		if targetTime.Unix() >= t.Unix() {
			return shift.Primary, nil
		}
	}
	return "", nil
}
