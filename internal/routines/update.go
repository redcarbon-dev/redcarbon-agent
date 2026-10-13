package routines

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/semver"

	"pkg.redcarbon.ai/internal/build"
	"pkg.redcarbon.ai/internal/github"
	"pkg.redcarbon.ai/internal/selfupdate"
)

var (
	downloadClient = &http.Client{Timeout: 10 * time.Minute}
	executable     = selfupdate.Executable

	// Release archives are a few tens of MiB and are buffered in memory to verify their checksum.
	maxDownloadSize int64 = 256 << 20

	// Every profile schedules its own update routine; only one may swap the binary at a time.
	updateMu sync.Mutex
)

func (r RoutineConfig) UpdateRoutine(ctx context.Context) {
	if !updateMu.TryLock() {
		r.Logger().Info("Skipping update as another profile is already updating the agent")
		return
	}
	defer updateMu.Unlock()

	r.Logger().Info("Starting update routine")

	rel, err := r.gh.LatestRelease(ctx, "redcarbon-dev", "redcarbon-agent")
	if err != nil {
		r.Logger().Error("Error while retrieving the latest release", "error", err)
		return
	}

	if build.Version == "DEV" || strings.Contains(build.Version, "SNAPSHOT") {
		r.Logger().Info("Skipping update as the agent is running in development mode")
		return
	}

	if semver.Compare(fmt.Sprintf("v%s", build.Version), rel.TagName) >= 0 {
		r.Logger().Info("Skipping update as the agent is already at the last version")
		return
	}

	for _, asset := range rel.Assets {
		if !strings.HasSuffix(asset.Name, fmt.Sprintf("_%s.tar.gz", build.Architecture)) {
			continue
		}

		checksum, err := retrieveChecksumForAsset(ctx, asset.Name, rel.Assets)
		if err != nil {
			r.Logger().Error("Error while retrieving the checksum for the latest version", "error", err)
			return
		}

		if err := r.doUpdate(ctx, asset, checksum); err != nil {
			r.Logger().Error("Unexpected error while updating the binary", "error", err)
			return
		}

		r.Logger().Info("Update executed successfully! Shutting down the agent...")

		// Only one receiver ever reads done: once shutdown has started, ctx is cancelled instead.
		select {
		case r.done <- true:
		case <-ctx.Done():
		}

		return
	}

	r.Logger().Error("No asset found for the agent architecture", "arch", build.Architecture, "release", rel.TagName)
}

func (r RoutineConfig) doUpdate(ctx context.Context, asset github.Asset, hexChecksum string) error {
	target, err := executable()
	if err != nil {
		return fmt.Errorf("locating the running executable: %w", err)
	}

	r.Logger().Info("Downloading new version...")

	archive, err := download(ctx, asset.BrowserDownloadURL)
	if err != nil {
		return err
	}

	if err := verifyChecksum(archive, hexChecksum); err != nil {
		return err
	}

	binary, err := openInArchive(archive, path.Join(strings.TrimSuffix(asset.Name, ".tar.gz"), "bin", "redcarbon"))
	if err != nil {
		return err
	}
	defer binary.Close()

	r.Logger().Info("Replacing the executable...", "path", target)

	return selfupdate.Replace(target, binary)
}

func retrieveChecksumForAsset(ctx context.Context, assetName string, assets []github.Asset) (string, error) {
	for _, asset := range assets {
		if !strings.HasSuffix(asset.Name, "checksums.txt") {
			continue
		}

		checksums, err := download(ctx, asset.BrowserDownloadURL)
		if err != nil {
			return "", err
		}

		return findChecksum(checksums, assetName)
	}

	return "", errors.New("checksums file not found in the release")
}

func findChecksum(checksums []byte, assetName string) (string, error) {
	for line := range strings.Lines(string(checksums)) {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == assetName {
			return fields[0], nil
		}
	}

	return "", fmt.Errorf("checksum for %s not found", assetName)
}

func download(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", build.UserAgent())

	resp, err := downloadClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("downloading %s: unexpected status %s", url, resp.Status)
	}

	b, err := io.ReadAll(io.LimitReader(resp.Body, maxDownloadSize+1))
	if err != nil {
		return nil, err
	}

	if int64(len(b)) > maxDownloadSize {
		return nil, fmt.Errorf("downloading %s: larger than %d bytes", url, maxDownloadSize)
	}

	return b, nil
}

func verifyChecksum(content []byte, hexChecksum string) error {
	expected, err := hex.DecodeString(hexChecksum)
	if err != nil {
		return err
	}

	actual := sha256.Sum256(content)

	if !bytes.Equal(actual[:], expected) {
		return fmt.Errorf("updated file has wrong checksum. Expected: %x, got: %x", expected, actual)
	}

	return nil
}

func openInArchive(archive []byte, name string) (_ io.ReadCloser, err error) {
	gzr, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = gzr.Close()
		}
	}()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("%s not found in the archive", name)
		}

		if err != nil {
			return nil, err
		}

		if header.Typeflag == tar.TypeReg && path.Clean(header.Name) == name {
			return struct {
				io.Reader
				io.Closer
			}{tr, gzr}, nil
		}
	}
}
