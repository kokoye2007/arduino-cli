// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package inventory

import (
	"os"
	"path/filepath"

	"github.com/arduino/arduino-cli/cli/feedback"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
)

// Store is the Read Only config storage
var Store = viper.New()

var (
	// Type is the inventory file type
	Type = "yaml"
	// Name is the inventory file Name with Type as extension
	Name = "inventory" + "." + Type
)

// Init configures the Read Only config storage
func Init(configPath string) {
	configFilePath := filepath.Join(configPath, Name)
	Store.SetConfigName(Name)
	Store.SetConfigType(Type)
	Store.AddConfigPath(configPath)
	// Attempt to read config file
	if err := Store.ReadInConfig(); err != nil {
		// ConfigFileNotFoundError is acceptable, anything else
		// should be reported to the user
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			generateInstallationData()
			writeStore(configFilePath)
		} else {
			feedback.Errorf("Error reading inventory file: %v", err)
		}
	}
}

func generateInstallationData() {
	installationID, err := uuid.NewV4()
	if err != nil {
		feedback.Errorf("Error generating installation.id: %v", err)
	}
	Store.Set("installation.id", installationID.String())

	installationSecret, err := uuid.NewV4()
	if err != nil {
		feedback.Errorf("Error generating installation.secret: %v", err)
	}
	Store.Set("installation.secret", installationSecret.String())
}

func writeStore(configFilePath string) {
	configPath := filepath.Dir(configFilePath)

	// Create config dir if not present,
	// MkdirAll will retrun no error if the path already exists
	if err := os.MkdirAll(configPath, os.FileMode(0755)); err != nil {
		feedback.Errorf("Error creating inventory dir: %v", err)
	}

	// Create file if not present
	err := Store.WriteConfigAs(configFilePath)
	if err != nil {
		feedback.Errorf("Error writing inventory file: %v", err)
	}
}
