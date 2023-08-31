// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package fetchers

import (
	"context"
	"fmt"
	"io/fs"
	"strconv"
	"testing"
	"testing/fstest"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/json"

	"github.com/elastic/cloudbeat/resources/fetching"
	"github.com/elastic/cloudbeat/resources/utils/testhelper"
)

const (
	statContent             = `1167 (containerd-shim) S 1 1167 198 0 -1 1077952768 223005 9831 39 0 665 1329 8 10 20 0 12 0 76222 730476544 2268 18446744073709551615 1 1 0 0 0 0 1006249984 0 2143420159 0 0 0 17 2 0 0 0 0 0 0 0 0 0 0 0 0 0`
	VanillaCmdLineDelimiter = "="
	EksCmdLineDelimiter     = " "
)

var Status = `Name:   %s`
var CmdLine = `/usr/bin/%s --kubeconfig=/etc/kubernetes/kubelet.conf --%s%s%s`

type TextProcessContext struct {
	Pid               string
	Name              string
	ConfigFileFlagKey string
	ConfigFilePath    string
}

type ProcessConfigTestStruct struct {
	A string
	B int
}

type ProcessFetcherTestSuite struct {
	suite.Suite

	resourceCh chan fetching.ResourceInfo
}

func TestProcessFetcherTestSuite(t *testing.T) {
	s := new(ProcessFetcherTestSuite)

	suite.Run(t, s)
}

func (s *ProcessFetcherTestSuite) SetupTest() {
	s.resourceCh = make(chan fetching.ResourceInfo, 50)
}

func (s *ProcessFetcherTestSuite) TearDownTest() {
	close(s.resourceCh)
}

func (s *ProcessFetcherTestSuite) TestFetchWhenFlagExistsButNoFile() {
	testProcess := TextProcessContext{
		Pid:               "3",
		Name:              "kubelet",
		ConfigFileFlagKey: "fetcherConfig",
		ConfigFilePath:    "test/path",
	}
	sysfs := createProcess(testProcess, VanillaCmdLineDelimiter)
	var procCfg = ProcessesConfigMap{
		testProcess.Name: {ConfigFileArguments: []string{"fetcherConfig"}},
	}
	processesFetcher := &ProcessesFetcher{log: testhelper.NewLogger(s.T()), Fs: sysfs, resourceCh: s.resourceCh, processes: procCfg}

	err := processesFetcher.Fetch(context.TODO(), fetching.CycleMetadata{})
	results := testhelper.CollectResources(s.resourceCh)

	s.Equal(1, len(results))
	s.Nil(err)

	processResource := results[0].Resource
	evalRes := processResource.GetData().(EvalProcResource)

	s.Equal(testProcess.Pid, evalRes.PID)
	s.Equal("kubelet", evalRes.Stat.Name)
	s.Contains(evalRes.Cmd, "/usr/bin/kubelet")
}

func (s *ProcessFetcherTestSuite) TestFetchWhenProcessDoesNotExist() {
	testProcess := TextProcessContext{
		Pid:               "3",
		Name:              "kubelet",
		ConfigFileFlagKey: "fetcherConfig",
		ConfigFilePath:    "test/path",
	}

	fsys := createProcess(testProcess, VanillaCmdLineDelimiter)
	var procCfg = ProcessesConfigMap{
		"someProcess": {ConfigFileArguments: []string{"fetcherConfig"}},
	}
	processesFetcher := &ProcessesFetcher{
		log:        testhelper.NewLogger(s.T()),
		Fs:         fsys,
		resourceCh: s.resourceCh,
		processes:  procCfg,
	}

	err := processesFetcher.Fetch(context.TODO(), fetching.CycleMetadata{})
	results := testhelper.CollectResources(s.resourceCh)

	s.Equal(0, len(results))
	s.Nil(err)
}

func (s *ProcessFetcherTestSuite) TestFetchWhenNoFlagRequired() {
	testProcess := TextProcessContext{
		Pid:               "3",
		Name:              "kubelet",
		ConfigFileFlagKey: "fetcherConfig",
		ConfigFilePath:    "test/path",
	}
	fsys := createProcess(testProcess, VanillaCmdLineDelimiter)
	var procCfg = ProcessesConfigMap{
		"kubelet": {ConfigFileArguments: []string{}}}
	processesFetcher := &ProcessesFetcher{log: testhelper.NewLogger(s.T()), Fs: fsys, resourceCh: s.resourceCh, processes: procCfg}
	err := processesFetcher.Fetch(context.TODO(), fetching.CycleMetadata{})

	results := testhelper.CollectResources(s.resourceCh)
	s.Equal(1, len(results))
	s.Nil(err)

	processResource := results[0].Resource
	evalRes := processResource.GetData().(EvalProcResource)

	s.Equal(testProcess.Pid, evalRes.PID)
	s.Equal("kubelet", evalRes.Stat.Name)
	s.Contains(evalRes.Cmd, "/usr/bin/kubelet")
}

func (s *ProcessFetcherTestSuite) TestFetchWhenFlagExistsWithConfigFile() {
	testCases := []struct {
		configFileName string
		marshal        func(in interface{}) (out []byte, err error)
		configType     string
		delimiter      string
	}{
		{"kubeletConfig.yaml", yaml.Marshal, "yaml", EksCmdLineDelimiter},
		{"kubeletConfig.yaml", yaml.Marshal, "yaml", VanillaCmdLineDelimiter},
		{"kubeletConfig.json", json.Marshal, "json", EksCmdLineDelimiter},
		{"kubeletConfig.json", json.Marshal, "json", VanillaCmdLineDelimiter},
	}

	for _, test := range testCases {
		configFlagKey := "fetcherConfig"
		// Creating a yaml file for the process fetcherConfig
		processConfig := ProcessConfigTestStruct{
			A: "A",
			B: 2,
		}
		configData, err := test.marshal(&processConfig)
		s.Nil(err)

		testProcess := TextProcessContext{
			Pid:               "3",
			Name:              "kubelet",
			ConfigFileFlagKey: configFlagKey,
			ConfigFilePath:    test.configFileName,
		}

		sysfs := createProcess(testProcess, test.delimiter).(fstest.MapFS)
		sysfs[test.configFileName] = &fstest.MapFile{
			Data: configData,
		}
		procCfg := ProcessesConfigMap{
			testProcess.Name: {ConfigFileArguments: []string{"fetcherConfig"}},
		}
		processesFetcher := &ProcessesFetcher{log: testhelper.NewLogger(s.T()), Fs: sysfs, resourceCh: s.resourceCh, processes: procCfg}
		err = processesFetcher.Fetch(context.TODO(), fetching.CycleMetadata{})
		results := testhelper.CollectResources(s.resourceCh)

		s.Equal(1, len(results))
		s.Nil(err)

		processResource := results[0].Resource
		evalRes := processResource.GetData().(EvalProcResource)
		procCD, err := processResource.GetElasticCommonData()
		s.NoError(err)

		s.Equal(testProcess.Pid, evalRes.PID)
		s.Equal("kubelet", evalRes.Stat.Name)
		s.Contains(evalRes.Cmd, "/usr/bin/kubelet")

		s.Equal(testProcess.Pid, strconv.FormatInt((procCD["process.pid"].(int64)), 10))
		s.True(procCD["process.args_count"].(int64) > 0)
		s.Contains(procCD["process.command_line"], "/usr/bin/kubelet")
		s.Equal("kubelet", procCD["process.name"])

		configResource := evalRes.ExternalData[configFlagKey]
		var result ProcessConfigTestStruct
		decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &result})
		s.Nil(err, "Could not decode process fetcherConfig result from %s type", test.configType)
		err = decoder.Decode(configResource)
		s.Nil(err, "Could not decode process fetcherConfig result from file %s", test.configFileName)

		s.Equal(processConfig.A, result.A)
		s.Equal(processConfig.B, result.B)
	}
}

func createProcess(process TextProcessContext, cmdDelimiter string) fs.FS {
	return fstest.MapFS{
		fmt.Sprintf("proc/%s/stat", process.Pid): {
			Data: []byte(statContent),
		},
		fmt.Sprintf("proc/%s/status", process.Pid): {
			Data: []byte(fmt.Sprintf(Status, process.Name)),
		},
		fmt.Sprintf("proc/%s/cmdline", process.Pid): {
			Data: []byte(fmt.Sprintf(CmdLine, process.Name, process.ConfigFileFlagKey, cmdDelimiter, process.ConfigFilePath)),
		},
	}
}
