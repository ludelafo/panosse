/*
Copyright © 2023 Ludovic Delafontaine <@ludelafo>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package utils

import (
	"errors"
	"os/exec"
	"regexp"
)

func Clean(metaflacCommand string, arguments []string, flacFile string) error {
	commandExec := exec.Command(metaflacCommand, append(arguments, flacFile)...)
	err := commandExec.Run()

	return err
}

func Normalize(metaflacCommand string, arguments []string, flacFiles []string) error {
	commandExec := exec.Command(metaflacCommand, append(arguments, flacFiles...)...)
	err := commandExec.Run()

	return err
}

func GetTag(metaflacCommand string, tag string, flacFile string) (string, error) {
	commandExec := exec.Command(metaflacCommand, "--show-tag", tag, flacFile)
	commandOutput, err := commandExec.CombinedOutput()

	if err != nil {
		return "", err
	}

	tagContent := string(commandOutput)

	return tagContent, nil
}

func SetTag(metaflacCommand string, tag string, tagContent string, flacFile string) error {
	commandExec := exec.Command(metaflacCommand, "--set-tag", tag+"="+tagContent, flacFile)
	err := commandExec.Run()

	return err
}

func RemoveAllTags(metaflacCommand string, flacFile string) error {
	commandExec := exec.Command(metaflacCommand, "--remove-all-tags", flacFile)
	err := commandExec.Run()

	return err
}

func RemoveTag(metaflacCommand string, tag string, flacFile string) error {
	commandExec := exec.Command(metaflacCommand, "--remove-tag", tag, flacFile)
	err := commandExec.Run()

	return err
}

const FlacVersionFromFlacFileRegex = "reference libFLAC ([\\d]+.[\\d]+.[\\d]+) [\\d]+"

func GetFlacVersionFromFlacFile(metaflacCommand string, flacFile string) (string, error) {
	commandExec := exec.Command(metaflacCommand, "--show-vendor-tag", flacFile)
	commandOutput, err := commandExec.CombinedOutput()

	if err != nil {
		return "", err
	}

	// Define the regular expression
	re := regexp.MustCompile(FlacVersionFromFlacFileRegex)

	// Find the match in the command output
	matches := re.FindStringSubmatch(string(commandOutput))

	var flacVersion string

	// Check if there is a match
	if len(matches) >= 2 {
		// Extract the version from the second capturing group
		flacVersion = matches[1]
	} else {
		return "", errors.New("unable to extract version")
	}

	return flacVersion, nil
}
