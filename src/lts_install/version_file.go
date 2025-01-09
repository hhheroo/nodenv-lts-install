package lts_install

import (
	"bytes"
	"html/template"
)

var versionFileTemplate = template.Must(
	template.New("versionFile").Parse(`binary {{.Os}}-{{.Arch}} "https://nodejs.org/dist/{{.Version}}/node-{{.Version}}-{{.Os}}-{{.Arch}}.tar.gz#{{.Hash}}"

install_package "node-{{.Version}}" "https://nodejs.org/dist/{{.Version}}/node-{{.Version}}.tar.gz#{{.SourceHash}}"`))

func MakeVersionFileContent(nodeVersion NodeVersion, hash VersionHash, os string, arch string) []byte {
	version := nodeVersion.Version
	var builder bytes.Buffer

	data := struct {
		Os         string
		Arch       string
		Version    string
		Hash       string
		SourceHash string
	}{
		os,
		arch,
		version,
		hash.Hash,
		hash.SourceHash,
	}

	versionFileTemplate.Execute(&builder, data)

	return builder.Bytes()
}
