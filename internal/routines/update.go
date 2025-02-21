package routines

import (
	"bytes"
	"context"
	"crypto"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/inconshreveable/go-update"
	"github.com/sirupsen/logrus"
	"golang.org/x/mod/semver"

	"pkg.redcarbon.ai/internal/build"
	"pkg.redcarbon.ai/internal/utils"
)

func (r RoutineConfig) UpdateRoutine(ctx context.Context) {
	rel, _, err := r.gh.Repositories.GetLatestRelease(ctx, "redcarbon-dev", "redcarbon-agent")
	if err != nil {
		return
	}

	if build.Version == "DEV" || strings.Contains(build.Version, "SNAPSHOT") {
		logrus.Info("Skipping update as the agent is running in development mode")
		return
	}

	if semver.Compare(fmt.Sprintf("v%s", build.Version), *rel.TagName) >= 0 {
		logrus.Info("Skipping update as the agent is already at the last version")
		return
	}

	for _, asset := range rel.Assets {
		if asset.Name == nil {
			continue
		}

		if strings.Contains(*asset.Name, fmt.Sprintf("%s.tar.gz", build.Architecture)) {
			checksum, err := r.retrieveChecksumForAsset(*asset.Name, rel.Assets)
			if err != nil {
				logrus.Errorf("Error while retrieving the checksum for the latest version for error %v", err)
				return
			}

			err = r.doUpdate(*asset.BrowserDownloadURL, *asset.Name, checksum)
			if err != nil {
				logrus.Errorf("Unexpected error while updating the binary %v", err)
				return
			}

			logrus.Info("Update executed successfully! Shutting down the agent...")

			r.done <- true

			return
		}
	}
}

func (r RoutineConfig) doUpdate(url string, name string, hexChecksum string) error {
	executable, err := r.downloadAsset(url, name, hexChecksum)
	if err != nil {
		return err
	}

	logrus.Info("Performing update...")

	err = update.Apply(executable, update.Options{})
	if err != nil {
		logrus.Errorf("Rollbacking for error found during update %v", err)

		if rErr := update.RollbackError(err); rErr != nil {
			logrus.Fatalf("Failed to rollback from bad update: %v", rErr)
		}

		logrus.Errorf("Rollback executed successfully")
	}

	return nil
}

func (r RoutineConfig) retrieveChecksumForAsset(assetName string, assets []*github.ReleaseAsset) (string, error) {
	for _, asset := range assets {
		if asset.Name == nil {
			continue
		}

		if strings.Contains(*asset.Name, "checksums.txt") {
			checksums, err := r.retrieveChecksumsList(asset)
			if err != nil {
				return "", err
			}

			for _, checksum := range checksums {
				if strings.Contains(checksum, assetName) {
					return strings.Split(checksum, " ")[0], nil
				}
			}
		}
	}

	return "", fmt.Errorf("checksum not found")
}

func (r RoutineConfig) retrieveChecksumsList(asset *github.ReleaseAsset) ([]string, error) {
	resp, err := http.Get(*asset.BrowserDownloadURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil
}

func (r RoutineConfig) downloadAsset(url string, name string, hexChecksum string) (io.Reader, error) {
	logrus.Info("Downloading new version...")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	archive, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = r.verifyChecksum(bytes.NewReader(archive), hexChecksum)
	if err != nil {
		return nil, err
	}

	logrus.Info("Extracting the binary...")

	outDir, err := os.MkdirTemp("", "*_redcarbon_agent")
	if err != nil {
		return nil, err
	}

	if err := utils.Untar(bytes.NewReader(archive), outDir); err != nil {
		return nil, err
	}

	target := filepath.Join(outDir, strings.Replace(name, ".tar.gz", "", -1), "bin", "redcarbon")

	file, err := os.Open(target)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (r RoutineConfig) verifyChecksum(file io.Reader, hexChecksum string) error {
	checksum, err := hex.DecodeString(hexChecksum)
	if err != nil {
		return err
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	h := crypto.SHA256.New()
	h.Write(content)
	fileChecksum := h.Sum([]byte{})

	if !bytes.Equal(fileChecksum, checksum) {
		return fmt.Errorf("updated file has wrong checksum. Expected: %x, got: %x", checksum, fileChecksum)
	}

	return nil
}
