package lts_install

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

const NODE_VERSION_API = "https://nodejs.org/dist/index.json"

func NODE_SHASUM_API(version string) string {
	return "https://nodejs.org/dist/" + version + "/SHASUMS256.txt"
}

type NodeVersion struct {
	Version string `json:"version"`
	Date    string `json:"date"`
	Lts     string
	PLts    interface{} `json:"lts"`
}

var _versions []NodeVersion

func unmarshalVersions(data []byte, v *[]NodeVersion) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	for i, version := range *v {
		switch lts := version.PLts.(type) {
		case string:
			(*v)[i].Lts = lts
		default:
			version.Lts = ""
		}
	}

	return nil
}

func GetVersions() ([]NodeVersion, error) {
	if _versions != nil {
		return _versions, nil
	}

	resp, err := http.Get(NODE_VERSION_API)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = unmarshalVersions(body, &_versions)

	if err != nil {
		return nil, err
	}

	return _versions, nil
}

var NumberVersionPattern = regexp.MustCompile(`^\d`)

func GetVersion(version string) (*NodeVersion, error) {
	versions, err := GetVersions()

	if err != nil {
		return nil, err
	}

	if NumberVersionPattern.MatchString(version) {
		version = NormalizeVersion(version)

		el, ok := Find(versions, func(el NodeVersion) bool {
			return strings.HasPrefix(el.Version, version)
		})

		if !ok {
			return nil, nil
		}

		return &el, nil
	}

	if version == "latest" {
		return &versions[0], nil
	}

	if version == "lts" {
		el, ok := Find(versions, func(el NodeVersion) bool {
			return el.Lts != ""
		})

		if !ok {
			return nil, nil
		}

		return &el, nil
	}

	version = strings.ToUpper(version[0:1]) + strings.ToLower(version[1:])

	el, ok := Find(versions, func(el NodeVersion) bool {
		return el.Lts == version
	})

	if !ok {
		return nil, nil
	}

	return &el, nil
}

type VersionHash struct {
	Hash       string
	SourceHash string
}

func GetHash(version NodeVersion, os string, arch string) (VersionHash, error) {
	binaryTarName := NodeBinaryTarName(version.Version, os, arch)
	sourceTarName := NodeSourceTarName(version.Version)
	versionHash := VersionHash{}

	resp, err := http.Get(NODE_SHASUM_API(version.Version))

	if err != nil {
		return VersionHash{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return VersionHash{}, err
	}

	tarHashes := strings.Split(string(body), "\n")

	hash, ok := Find(tarHashes, func(s string) bool {
		return strings.Contains(s, binaryTarName)
	})

	if ok {
		versionHash.Hash = strings.TrimSpace(strings.Split(hash, " ")[0])
	}

	sourceHash, ok := Find(tarHashes, func(s string) bool {
		return strings.Contains(s, sourceTarName)
	})

	if ok {
		versionHash.SourceHash = strings.TrimSpace(strings.Split(sourceHash, " ")[0])
	}

	return versionHash, nil
}

func NormalizeVersion(version string) string {
	if version[0] == 'v' {
		return version
	}

	return "v" + version
}

func NodeBinaryTarName(version string, os string, arch string) string {
	return "node-" + version + "-" + os + "-" + arch + ".tar.gz"
}

func NodeSourceTarName(version string) string {
	return "node-" + version + ".tar.gz"
}
