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

func Encode(flacCommand string, encodeSettings []string, flacFile string) (string, error) {
	commandExec := exec.Command(flacCommand, append(encodeSettings, flacFile)...)
	commandOutput, err := commandExec.CombinedOutput()

	return string(commandOutput), err
}

const FlacVersionFromFlacCommandRegex = "flac ([\\d]+.[\\d]+.[\\d]+)"

func GetFlacVersionFromFlacCommand(flacCommand string) (string, error) {
	commandExec := exec.Command(flacCommand, "--version")
	commandOutput, err := commandExec.CombinedOutput()

	if err != nil {
		return string(commandOutput), err
	}

	// Define the regular expression
	re := regexp.MustCompile(FlacVersionFromFlacCommandRegex)

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

func Verify(flacCommand string, verifyCommandArguments []string, flacFile string) (string, error) {
	commandExec := exec.Command(flacCommand, append(verifyCommandArguments, flacFile)...)
	commandOutput, err := commandExec.CombinedOutput()

	return string(commandOutput), err
}
