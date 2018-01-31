/*
Copyright (C) 2017 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package parser

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/minishift/minishift/pkg/minishift/addon"
	"github.com/minishift/minishift/pkg/minishift/addon/command"
	"github.com/minishift/minishift/pkg/util/filehelper"
)

const (
	commentChar = "#"

	noAddOnDefinitionFoundError         = "There needs to be an addon file per addon directory. Found none in '%s'"
	multipleAddOnDefinitionsError       = "There can only be one addon file per addon directory. Found '%s'"
	multipleAddOnRemoveDefinitionsError = "There can only be one addon.remove file per addon directory. Found '%s'"
	regexToGetMetaTagInfo               = `^# ?([a-zA-Z-]*):(.*)`
)

// AddOnParser is responsible for loading an addon from file and converting it into an AddOn
type AddOnParser struct {
	handler CommandHandler
}

func NewAddOnParser() *AddOnParser {
	parser := AddOnParser{}

	// build the responsibility chain
	ocHandler := &OcCommandHandler{&defaultCommandHandler{}}

	openshiftHandler := &OpenShiftCommandHandler{&defaultCommandHandler{}}
	ocHandler.next = openshiftHandler

	dockerHandler := &DockerCommandHandler{&defaultCommandHandler{}}
	openshiftHandler.SetNext(dockerHandler)

	sleepHandler := &SleepCommandHandler{&defaultCommandHandler{}}
	dockerHandler.SetNext(sleepHandler)

	sshHandler := &SSHCommandHandler{&defaultCommandHandler{}}
	sleepHandler.SetNext(sshHandler)

	echoHandler := &EchoCommandHandler{&defaultCommandHandler{}}
	sshHandler.SetNext(echoHandler)

	parser.handler = ocHandler

	return &parser
}

// Parse takes as parameter a reader containing an addon definition and returns an AddOn instance.
// If an error occurs, the error is returned.
func (parser *AddOnParser) Parse(addOnDir string) (addon.AddOn, error) {
	addonInstallReader, err := parser.getAddOnContentReader(addOnDir, ".addon")
	if err != nil {
		return nil, err
	}
	addonRemoveReader, err := parser.getAddOnContentReader(addOnDir, ".addon.remove")
	if err != nil {
		return nil, err
	}
	meta, commands, err := parser.parseAddOnContent(addonInstallReader)
	if err != nil {
		name := ""
		if meta != nil {
			name = meta.Name()
		}
		return nil, NewParseError(err.Error(), name, addOnDir)
	}

	var removeCommands []command.Command
	if addonRemoveReader != nil {
		_, removeCommands, err = parser.parseAddOnContent(addonRemoveReader)
		if err != nil {
			name := ""
			if meta != nil {
				name = meta.Name()
			}
			return nil, NewParseError(err.Error(), name, addOnDir)
		}
	}

	addOn := addon.NewAddOn(meta, commands, removeCommands, addOnDir)

	return addOn, nil
}

func (parser *AddOnParser) getAddOnContentReader(addOnDir string, fileSuffix string) (io.Reader, error) {
	if !filehelper.Exists(addOnDir) {
		return nil, NewParseError("Addon directory does not exist", addOnDir, "")
	}

	files, err := ioutil.ReadDir(addOnDir)
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("Unexpected error reading addon content in '%s'", addOnDir), addOnDir, "")
	}
	var addOnFiles []string
	for _, fileInfo := range files {
		if strings.HasSuffix(fileInfo.Name(), fileSuffix) {
			addOnFiles = append(addOnFiles, fileInfo.Name())
		}
	}

	if fileSuffix == ".addon.remove" {
		if len(addOnFiles) == 0 {
			return nil, nil
		}
	}

	if len(addOnFiles) == 0 {
		return nil, NewParseError(fmt.Sprintf(noAddOnDefinitionFoundError, addOnDir), addOnDir, "")
	}

	if len(addOnFiles) > 1 {
		if fileSuffix == ".addon" {
			return nil, NewParseError(fmt.Sprintf(multipleAddOnDefinitionsError, strings.Join(addOnFiles, ", ")), addOnDir, "")
		}
		return nil, NewParseError(fmt.Sprintf(multipleAddOnRemoveDefinitionsError, strings.Join(addOnFiles, ", ")), addOnDir, "")
	}

	file, err := ioutil.ReadFile(filepath.Join(addOnDir, addOnFiles[0]))
	if err != nil {
		return nil, NewParseError(fmt.Sprintf("Unable to open addon definition '%s'", addOnFiles[0]), addOnDir, "")
	}
	reader := strings.NewReader(string(file))
	return bufio.NewReader(reader), nil
}

func (parser *AddOnParser) parseAddOnContent(reader io.Reader) (addon.AddOnMeta, []command.Command, error) {
	scanner := bufio.NewScanner(reader)
	meta, err := parser.parseHeader(scanner)
	if err != nil {
		return nil, nil, err
	}

	commands, err := parser.parseCommands(scanner)
	if err != nil {
		return meta, nil, err
	}

	return meta, commands, nil
}

func (parser *AddOnParser) parseHeader(scanner *bufio.Scanner) (addon.AddOnMeta, error) {
	var header []string
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if !strings.HasPrefix(line, commentChar) {
			break
		}
		header = append(header, line)
	}

	headerMeta, err := createMetaData(header)
	if err != nil {
		return nil, err
	}

	return headerMeta, nil
}

func (parser *AddOnParser) parseCommands(scanner *bufio.Scanner) ([]command.Command, error) {
	var commands []command.Command
	for scanner.Scan() {
		line := scanner.Text()

		// skip blank and comment lines
		line = strings.Trim(line, " ")
		if len(line) == 0 || strings.HasPrefix(line, commentChar) {
			continue
		}

		newCommand, err := parser.handler.Handle(parser.handler, line)
		if err != nil {
			return nil, err
		}

		commands = append(commands, newCommand)
	}

	return commands, nil
}

func createMetaData(header []string) (addon.AddOnMeta, error) {
	regex, _ := regexp.Compile(regexToGetMetaTagInfo)
	metaMap := make(map[string]interface{})
	var key string
	var value []string
	for _, line := range header {
		matches := regex.FindAllStringSubmatch(line, -1)
		if len(matches) == 0 && key == addon.DescriptionMetaTagName {
			line = strings.TrimPrefix(line, commentChar)
			metaMap[key] = append(metaMap[key].([]string), strings.TrimSpace(line))
			continue
		}
		for _, match := range matches {
			key = strings.Trim(match[1], " ")
			if key == addon.DescriptionMetaTagName {
				value = append(value, strings.Trim(match[2], " "))
				metaMap[key] = value
				continue
			}
			metaMap[key] = strings.Trim(match[2], " ")
		}
	}

	metaData, err := addon.NewAddOnMeta(metaMap)
	if err != nil {
		return nil, err
	}

	return metaData, nil
}
