// Copyright (c) OpenFaaS Project 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/openfaas/faas-cli/version"
	"github.com/spf13/cobra"
)

const (
	githubApiRelease = "https://api.github.com/repos/openfaas/faas-cli/releases"
)

var selfupdate struct {
	beta bool
}

var (
	ErrorReleaseNotFound = errors.New("could not determine the release to download")
)

func init() {
	versionCmd.AddCommand(newSelfUpdateCmd())
}

func newSelfUpdateCmd() *cobra.Command {
	selfUpdateCmd := &cobra.Command{
		Use:   `update`,
		Short: "Self-update faas-cli",
		Long:  "Self-update faas-cli",
		RunE:  selfUpdate,
	}

	selfUpdateCmd.Flags().BoolVar(&selfupdate.beta, "beta", false, "Include beta releases")

	return selfUpdateCmd
}

func selfUpdate(cmd *cobra.Command, args []string) error {
	file, err := filepath.Abs(os.Args[0])
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Get(githubApiRelease)
	defer resp.Body.Close()

	type asset struct {
		Name        string `json:"name"`
		DownloadUrl string `json:"browser_download_url"`
	}
	type githubApiReleaseResponse struct {
		TagName string `json:"tag_name"`
		Assets  []asset
	}

	var githubApiReleaseResp []githubApiReleaseResponse

	if jsonBytes, err := ioutil.ReadAll(resp.Body); err != nil {
		return err
	} else {
		if err := json.Unmarshal(jsonBytes, &githubApiReleaseResp); err != nil {
			return err
		}
	}

	var newVersionResp githubApiReleaseResponse
	for _, r := range githubApiReleaseResp {
		if !selfupdate.beta && strings.HasSuffix(r.TagName, "-beta") {
			continue
		}

		if isNewerVersion(version.Version, r.TagName) {
			newVersionResp = r
			break
		}
	}

	if len(newVersionResp.TagName) == 0 {
		return ErrorReleaseNotFound
	}

	fmt.Printf("Found newer version %s. Updating...\n", newVersionResp.TagName)
	var a asset
	for _, e := range newVersionResp.Assets {
		if strings.Contains(e.DownloadUrl, runtime.GOOS) || strings.Contains(e.DownloadUrl, runtime.GOARCH) ||
			runtime.GOOS == "windows" && strings.Contains(e.DownloadUrl, ".exe") {
			a = e
			break
		}
	}

	if len(a.DownloadUrl) == 0 {
		return ErrorReleaseNotFound
	}

	tmpFile := a.Name + "-" + newVersionResp.TagName + ".tmp"
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		if err := DownloadFile(tmpFile, a.DownloadUrl); err != nil {
			return err
		}
		if runtime.GOOS != "windows" {
			if err := os.Chmod(tmpFile, 0750); err != nil {
				return err
			}
		}
	}

	fileBkp := "." + file + ".old"
	err = os.Rename(file, fileBkp)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFile, file)
	if err != nil {
		return err
	}

	fmt.Printf("Update to version: %s\n", newVersionResp.TagName)

	return nil
}

// DownloadFile downloads an url to filepath
func DownloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// isNewerVersion compares old and new versions
func isNewerVersion(old, new string) bool {
	// remove -beta
	new = strings.Replace(new, "-beta", "", 1)

	// remove .
	newParts := strings.Split(new, ".")
	oldParts := strings.Split(old, ".")

	// SemVer : 3 parts
	if newParts[0] > oldParts[0] {
		return true
	} else if newParts[0] == oldParts[0] {
		if newParts[1] > oldParts[1] {
			return true
		} else if newParts[1] == oldParts[1] {
			return newParts[2] > oldParts[2]
		} else {
			return false
		}
	} else {
		return false
	}
}
