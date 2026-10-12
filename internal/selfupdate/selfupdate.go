// Package selfupdate replaces the binary of the running agent on unix systems.
package selfupdate

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var rename = os.Rename

// Temporary files untouched for this long are leftovers of an interrupted Replace (e.g. the
// process was killed mid-copy). Fresher ones may belong to a Replace still in flight.
const staleTempAge = time.Hour

// Executable returns the path of the running binary with symlinks resolved, so the
// replacement lands on the file that is actually executed rather than on a link to it.
func Executable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(exe)
}

// Replace atomically swaps the file at target (symlinks resolved) for the contents of src,
// keeping the permission bits of target. On failure target is left untouched.
func Replace(target string, src io.Reader) (err error) {
	// Renaming over a symlink would replace the link itself and leave the real binary stale.
	target, err = filepath.EvalSymlinks(target)
	if err != nil {
		return err
	}

	info, err := os.Stat(target)
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", target)
	}

	dir := filepath.Dir(target)
	prefix := "." + filepath.Base(target) + ".new-"

	removeStaleTemps(dir, prefix)

	// The temporary file must live next to target: rename(2) is only atomic within a filesystem.
	tmp, err := os.CreateTemp(dir, prefix+"*")
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmp.Name())
		}
	}()

	n, err := io.Copy(tmp, src)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("refusing to replace the executable with an empty file")
	}

	if err = tmp.Chmod(info.Mode().Perm()); err != nil {
		return err
	}

	if err = tmp.Sync(); err != nil {
		return err
	}

	if err = tmp.Close(); err != nil {
		return err
	}

	if err = rename(tmp.Name(), target); err != nil {
		return err
	}

	// Best effort: the swap has already happened, a failed directory sync only
	// weakens durability of the rename across a power loss.
	_ = syncDir(dir)

	return nil
}

// removeStaleTemps deletes, best effort, the temporary files that interrupted replacements
// left next to the binary: their names are random, so nothing else would ever reclaim them.
func removeStaleTemps(dir, prefix string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if !e.Type().IsRegular() || !strings.HasPrefix(e.Name(), prefix) {
			continue
		}

		info, err := e.Info()
		if err != nil || time.Since(info.ModTime()) < staleTempAge {
			continue
		}

		_ = os.Remove(filepath.Join(dir, e.Name()))
	}
}

func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer d.Close()

	return d.Sync()
}
