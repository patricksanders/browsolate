package browsolate

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func randomColor() string {
	return fmt.Sprintf("%d,%d,%d", rand.Intn(255), rand.Intn(255), rand.Intn(255))
}

type InstanceOpts struct {
	ChromePath    string
	TempDirBase   string
	TempDirPrefix string
}

func (b *InstanceOpts) fillDefaults() {
	if b.ChromePath == "" {
		b.ChromePath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	}
	if b.TempDirPrefix == "" {
		b.TempDirPrefix = "browsolate."
	}
}

func StartIsolatedChromeInstance(url string, opts *InstanceOpts) error {
	var err error
	opts.fillDefaults()

	profileDir, err := os.MkdirTemp(opts.TempDirBase, opts.TempDirPrefix)
	if err != nil {
		return fmt.Errorf("could not create temp profile dir: %w", err)
	}
	profileColor := randomColor()

	log.Printf("starting isolated browser in %s with color %s", profileDir, profileColor)
	args := []string{
		opts.ChromePath,
		fmt.Sprintf("--user-data-dir=%s", profileDir),
		"--no-first-run",
		fmt.Sprintf("--install-autogenerated-theme=%s", profileColor),
		url,
	}

	err = startDetachedProcess(args, profileDir)
	if err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}
	return nil
}

// getUidGid returns the UID and GID of the current user
func getUidGid() (uid uint32, gid uint32, err error) {
	var currentUser *user.User
	var uid64, gid64 uint64
	currentUser, err = user.Current()
	if err != nil {
		return 0, 0, fmt.Errorf("could not get current user: %w", err)
	}
	uid64, err = strconv.ParseUint(currentUser.Uid, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse current uid: %w", err)
	}
	gid64, err = strconv.ParseUint(currentUser.Gid, 10, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse current gid: %w", err)
	}

	return uint32(uid64), uint32(gid64), nil
}

// startDetachedProcess creates a process and detaches it.
// Ref: github.com/ik5/fork_process
func startDetachedProcess(args []string, workdir string) error {
	uid, gid, err := getUidGid()
	if err != nil {
		return err
	}
	var cred = &syscall.Credential{
		Uid:         uid,
		Gid:         gid,
		NoSetGroups: true,
	}

	var sysproc = &syscall.SysProcAttr{
		Credential: cred,
		Setsid:     true,
	}

	rpipe, wpipe, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("unable to get read and write files: %w", err)
	}
	defer rpipe.Close()
	defer wpipe.Close()

	attr := os.ProcAttr{
		Dir: workdir,
		Env: os.Environ(),
		Files: []*os.File{
			rpipe,
			wpipe,
			wpipe,
		},
		Sys: sysproc,
	}
	//process, err := os.StartProcess(args[0], args, &attr)
	process, err := os.StartProcess(args[0], args, &attr)
	if err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	err = process.Release()
	if err != nil {
		return fmt.Errorf("failed to release process: %w", err)
	}
	return nil
}
