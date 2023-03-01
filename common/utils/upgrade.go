package utils

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"runtime"

	"go-lyric/common/logger"

	"github.com/kardianos/osext"
	"github.com/kr/binarydist"
	"github.com/sanbornm/go-selfupdate/selfupdate"
	"gopkg.in/inconshreveable/go-update.v0"
)

var plat = runtime.GOOS + "-" + runtime.GOARCH
var up = update.New()
var ErrHashMismatch = errors.New("new file hash mismatch after patch")
var defaultHTTPRequester = selfupdate.HTTPRequester{}

func DoUpgrade(upgradeUrl, app, version string) error {
	var updater = &Updater{
		CurrentVersion: version,
		ApiURL:         upgradeUrl,
		BinURL:         upgradeUrl,
		DiffURL:        upgradeUrl,
		CmdName:        app,
	}

	if err := updater.BackgroundUpdate(); err != nil {
		logger.Info("升级失败")
		return errors.New("升级失败")
	} else {
		logger.Info("升级成功")
		return nil
	}
}

type Updater struct {
	CurrentVersion string               // Currently running version.
	ApiURL         string               // Base URL for API requests (json files).
	CmdName        string               // Command name is appended to the ApiURL like http://apiurl/CmdName/. This represents one binary.
	BinURL         string               // Base URL for full binary downloads.
	DiffURL        string               // Base URL for diff downloads.
	Dir            string               // Directory to store selfupdate state.
	Requester      selfupdate.Requester //Optional parameter to override existing http request handler
	Info           struct {
		Version string
		Sha256  []byte
	}
}

func (u *Updater) BackgroundUpdate() error {
	path, err := osext.Executable()
	if err != nil {
		logger.Error(err)
		return err
	}
	old, err := os.Open(path)
	if err != nil {
		logger.Error(err)

		return err
	}
	defer old.Close()

	err = u.fetchInfo()
	if err != nil {
		logger.Error(err)

		return err
	}
	if u.Info.Version == u.CurrentVersion {
		logger.Error("errrrrrrrrrrrrr!")

		return nil
	}
	bin, err := u.fetchAndVerifyPatch(old)
	if err != nil {
		if err == ErrHashMismatch {
			log.Println("update: hash mismatch from patched binary")
		} else {
			if u.DiffURL != "" {
				log.Println("update: patching binary,", err)
			}
		}

		bin, err = u.fetchAndVerifyFullBin()
		if err != nil {
			if err == ErrHashMismatch {
				log.Println("update: hash mismatch from full binary")
			} else {
				log.Println("update: fetching full binary,", err)
			}
			logger.Error(err)

			return err
		}
	}

	// close the old binary before installing because on windows
	// it can't be renamed if a handle to the file is still open
	old.Close()

	err, errRecover := up.FromStream(bytes.NewBuffer(bin))
	if errRecover != nil {
		logger.Error("recover error!!!")
		return err
	}
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
func (u *Updater) fetchInfo() error {
	logger.Info(u.ApiURL + url.QueryEscape(plat) + ".json")
	//r, err := u.fetch(u.ApiURL + url.QueryEscape(u.CmdName) + "/" + url.QueryEscape(plat) + ".json")

	r, err := u.fetch(u.ApiURL + url.QueryEscape(plat) + ".json")
	if err != nil {
		logger.Error(err)

		return err
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&u.Info)
	if err != nil {
		return err
	}
	if len(u.Info.Sha256) != sha256.Size {
		return errors.New("bad cmd hash in info")
	}
	return nil
}

func (u *Updater) fetchAndVerifyPatch(old io.Reader) ([]byte, error) {
	bin, err := u.fetchAndApplyPatch(old)
	if err != nil {
		return nil, err
	}
	if !verifySha(bin, u.Info.Sha256) {
		return nil, ErrHashMismatch
	}
	return bin, nil
}

func verifySha(bin []byte, sha []byte) bool {
	h := sha256.New()
	h.Write(bin)
	return bytes.Equal(h.Sum(nil), sha)
}

func (u *Updater) fetchAndApplyPatch(old io.Reader) ([]byte, error) {
	//r, err := u.fetch(u.DiffURL + url.QueryEscape(u.CmdName) + "/" + url.QueryEscape(u.CurrentVersion) + "/" + url.QueryEscape(u.Info.Version) + "/" + url.QueryEscape(plat))

	r, err := u.fetch(u.DiffURL + url.QueryEscape(u.CurrentVersion) + "/" + url.QueryEscape(u.Info.Version) + "/" + url.QueryEscape(plat))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var buf bytes.Buffer
	err = binarydist.Patch(old, &buf, r)
	return buf.Bytes(), err
}

func (u *Updater) fetchAndVerifyFullBin() ([]byte, error) {
	bin, err := u.fetchBin()
	if err != nil {
		return nil, err
	}
	verified := verifySha(bin, u.Info.Sha256)
	if !verified {
		return nil, ErrHashMismatch
	}
	return bin, nil
}

func (u *Updater) fetchBin() ([]byte, error) {
	//r, err := u.fetch(u.BinURL + url.QueryEscape(u.CmdName) + "/" + url.QueryEscape(u.Info.Version) + "/" + url.QueryEscape(plat) + ".gz")

	r, err := u.fetch(u.BinURL + url.QueryEscape(u.Info.Version) + "/" + url.QueryEscape(plat) + ".gz")
	if err != nil {
		return nil, err
	}
	defer r.Close()
	buf := new(bytes.Buffer)
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(buf, gz); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (u *Updater) fetch(url string) (io.ReadCloser, error) {
	if u.Requester == nil {
		return defaultHTTPRequester.Fetch(url)
	}

	readCloser, err := u.Requester.Fetch(url)
	if err != nil {
		return nil, err
	}

	if readCloser == nil {
		return nil, fmt.Errorf("Fetch was expected to return non-nil ReadCloser")
	}

	return readCloser, nil
}
