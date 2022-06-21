/*
 * Copyright (C) 2022 Canonical Ltd
 *
 *  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 *  in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 * SPDX-License-Identifier: Apache-2.0'
 */

package main

import (
	"os"
	"strings"

	hooks "github.com/canonical/edgex-snap-hooks/v2"
	"github.com/canonical/edgex-snap-hooks/v2/env"
	"github.com/canonical/edgex-snap-hooks/v2/log"
	"github.com/canonical/edgex-snap-hooks/v2/options"
	"github.com/canonical/edgex-snap-hooks/v2/snapctl"
)

func main() {
	subCommand := os.Args[1]
	switch subCommand {
	case "install":
		install()
	case "configure":
		configure()
	default:
		panic("Unknown hook sub-command: " + subCommand)
	}
}

// installProfiles copies the profile configuration.toml files from $SNAP to $SNAP_DATA.
func installConfig() error {
	path := "/config/device-virtual/res"

	err := hooks.CopyDir(
		env.Snap+path,
		env.SnapData+path)
	if err != nil {
		return err
	}

	return nil
}

func installDevices() error {
	path := "/config/device-virtual/res/devices"

	err := hooks.CopyDir(
		hooks.Snap+path,
		hooks.SnapData+path)
	if err != nil {
		return err
	}

	return nil
}

func installDevProfiles() error {
	path := "/config/device-virtual/res/profiles"

	err := hooks.CopyDir(
		hooks.Snap+path,
		hooks.SnapData+path)
	if err != nil {
		return err
	}

	return nil
}

func install() {
	log.SetComponentName("install")

	err := installConfig()
	if err != nil {
		log.Fatalf("error installing config file: %s", err)
	}

	err = installDevices()
	if err != nil {
		log.Fatalf("error installing devices config: %s", err)
	}

	err = installDevProfiles()
	if err != nil {
		log.Fatalf("error installing device profiles config: %s", err)
	}
}

//configure
func configure() {
	log.SetComponentName("configure")

	log.Info("Enabling config options")
	err := snapctl.Set("app-options", "true").Run()
	if err != nil {
		log.Fatalf("could not enable config options: %v", err)
	}

	log.Info("Processing options")
	err = options.ProcessAppConfig("device-virtual")
	if err != nil {
		log.Fatalf("could not process options: %v", err)
	}

	// If autostart is not explicitly set, default to "no"
	// as only example service configuration and profiles
	// are provided by default.
	autostart, err := snapctl.Get("autostart").Run()
	if err != nil {
		log.Fatalf("Reading config 'autostart' failed: %v", err)
	}
	if autostart == "" {
		log.Debug("autostart is NOT set, initializing to 'no'")
		autostart = "no"
	}
	autostart = strings.ToLower(autostart)
	log.Debugf("autostart=%s", autostart)

	// services are stopped/disabled by default in the install hook
	switch autostart {
	case "true", "yes":
		err = snapctl.Start("device-virtual").Enable().Run()
		if err != nil {
			log.Fatalf("Can't start service: %s", err)
		}
	case "false", "no":
		// no action necessary
	default:
		log.Fatalf("Invalid value for 'autostart': %s", autostart)
	}
}
