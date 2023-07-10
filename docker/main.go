package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"lyric/common/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gokrazy/gokrazy"
	"github.com/tidwall/gjson"
)

func podman(args ...string) error {
	podman := exec.Command("/user/podman", args...)
	podman.Env = expandPath(os.Environ())
	podman.Env = append(podman.Env, "TMPDIR=/tmp")
	podman.Stdin = os.Stdin
	podman.Stdout = os.Stdout
	podman.Stderr = os.Stderr
	if err := podman.Run(); err != nil {
		logger.Error(err)
		return fmt.Errorf("%v: %v", podman.Args, err)
	}
	return nil
}

func start() error {
	// Ensure we have an up-to-date clock, which in turn also means that
	// networking is up. This is relevant because podman takes what’s in
	// /etc/resolv.conf (nothing at boot) and holds on to it, meaning your
	// container will never have working networking if it starts too early.
	gokrazy.WaitForClock()

	if err := mountVar(); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

/*
/user/podman run -itd --name nginx -p 8080:80 docker.io/library/nginx
*/
func main() {
	if err := start(); err != nil {
		logger.Error(err)
		return
	}
}

func getContainerId(name string) string {
	data, err := ioutil.ReadFile("/var/lib/containers/storage/overlay-containers/containers.json")
	if err != nil {
		logger.Error(err)
		return ""
	}
	contains := gjson.Parse(string(data)).Array()
	for _, contain := range contains {
		n := contain.Get("metadata.name").String()
		if n == name {
			return contain.Get("id").String()
		}
	}
	return ""
}
func reStartService(name string) error {
	podman("kill", name)

	path := "/var/lib/cni/networks/podman"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(err)
		return err
	}

	id := getContainerId(name)
	for _, f := range files {
		filePath := filepath.Join(path, f.Name())
		data, _ := ioutil.ReadFile(filePath)
		if strings.Contains(string(data), id) {
			os.Remove(filePath)
			break
		}
	}

	if err := podman("start", name); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// mountVar bind-mounts /perm/container-storage to /var if needed.
// This could be handled by an fstab(5) feature in gokrazy in the future.
func mountVar() error {
	os.Mkdir("/perm/container-storage", 0777)
	os.Mkdir("/var/tmp", 0777)

	b, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		logger.Error(err)
		return err
	}
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		parts := strings.Fields(line)
		if len(parts) < 5 {
			continue
		}
		mountpoint := parts[4]
		log.Printf("Found mountpoint %q", parts[4])
		if mountpoint == "/var" {
			logger.Error("/var file system already mounted, nothing to do")
			return nil
		}
	}

	if err := syscall.Mount("/perm/container-storage", "/var", "", syscall.MS_BIND, ""); err != nil {
		logger.Error(err)
		return fmt.Errorf("mounting /perm/container-storage to /var: %v", err)
	}

	return nil
}

// expandPath returns env, but with PATH= modified or added
// such that both /user and /usr/local/bin are included, which podman needs.
func expandPath(env []string) []string {
	extra := "/user:/usr/local/bin"
	found := false
	for idx, val := range env {
		parts := strings.Split(val, "=")
		if len(parts) < 2 {
			continue // malformed entry
		}
		key := parts[0]
		if key != "PATH" {
			continue
		}
		val := strings.Join(parts[1:], "=")
		env[idx] = fmt.Sprintf("%s=%s:%s", key, extra, val)
		found = true
	}
	if !found {
		const busyboxDefaultPATH = "/usr/local/sbin:/sbin:/usr/sbin:/usr/local/bin:/bin:/usr/bin:/user"
		env = append(env, fmt.Sprintf("PATH=%s:%s", extra, busyboxDefaultPATH))
	}
	return env
}

/*
	if err := podman("run",
		"-itd",
		"--name", "nginx",
		"-p", "8080:80",
		"docker.io/library/nginx"); err != nil {
		logger.Error(err)
		return err
	}

*/
