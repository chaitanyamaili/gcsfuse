// Copyright 2023 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Provides integration tests for read large files sequentially and randomly.
package read_large_files

import (
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/util/mounting/static_mounting"
	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/util/operations"
	"github.com/googlecloudplatform/gcsfuse/tools/integration_tests/util/setup"
)

const OneMB = 1024 * 1024
const FiveHundredMB = 500 * OneMB
const FiveHundredMBFile = "fiveHundredMBFile.txt"
const ChunkSize = 200 * OneMB
const NumberOfRandomReadCalls = 200
const MinReadableByteFromFile = 0
const MaxReadableByteFromFile = 500 * OneMB

func TestMain(m *testing.M) {
	setup.ParseSetUpFlags()

	flags := [][]string{{"--implicit-dirs"}}

	setup.ExitWithFailureIfBothTestBucketAndMountedDirectoryFlagsAreNotSet()

	if setup.TestBucket() != "" && setup.MountedDirectory() != "" {
		log.Print("Both --testbucket and --mountedDirectory can't be specified at the same time.")
		os.Exit(1)
	}

	// Run tests for mountedDirectory only if --mountedDirectory flag is set.
	setup.RunTestsForMountedDirectoryFlag(m)

	// Run tests for testBucket
	setup.SetUpTestDirForTestBucketFlag()

	successCode := static_mounting.RunTests(flags, m)

	setup.RemoveBinFileCopiedForTesting()

	os.Exit(successCode)
}

func createFileOnDiskAndCopyToMntDir(fileInLocalDisk string, fileInMntDir string, fileSize int, t *testing.T) {
	setup.RunScriptForTestData("testdata/write_content_of_fix_size_in_file.sh", fileInLocalDisk, strconv.Itoa(fileSize))
	err := operations.CopyFile(fileInLocalDisk, fileInMntDir)
	if err != nil {
		t.Errorf("Error in copying file:%v", err)
	}
}