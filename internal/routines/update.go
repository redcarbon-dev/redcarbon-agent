package routines

import (
	"context"
	"crypto"
	"encoding/hex"
	"fmt"
	"github.com/google/go-github/v50/github"
	"golang.org/x/mod/semver"
	"io"
	"net/http"
	"pkg.redcarbon.ai/internal/build"
	"strings"

	"github.com/inconshreveable/go-update"
	"github.com/sirupsen/logrus"
)

func (r routineConfig) UpdateRoutine() {
	rel, _, err := r.gh.Repositories.GetLatestRelease(context.Background(), "redcarbon-dev", "redcarbon-agent")
	if err != nil {
		return
	}

	if build.Version == "DEV" {
		logrus.Info("Skipping update as the agent is running in a development status")
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
				return
			}

			err = r.doUpdate(*asset.BrowserDownloadURL, checksum)
			if err != nil {
				return
			}
		}
	}
}

func (r routineConfig) doUpdate(url string, hexChecksum string) error {
	logrus.Info("Downloading new version...")

	checksum, err := hex.DecodeString(hexChecksum)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	logrus.Info("Performing update...")

	err = update.Apply(resp.Body, update.Options{
		Hash:     crypto.SHA256,
		Checksum: checksum,
	})
	if err != nil {
		if rErr := update.RollbackError(err); rErr != nil {
			logrus.Fatalf("Failed to rollback from bad update: %v", rErr)
		}
	}

	return nil
}

func (r routineConfig) retrieveChecksumForAsset(assetName string, assets []*github.ReleaseAsset) (string, error) {
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

func (r routineConfig) retrieveChecksumsList(asset *github.ReleaseAsset) ([]string, error) {
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
