'use strict';

var os = require('node:os');
var path = require('node:path');
var fs = require('node:fs');
var fsP = require('node:fs/promises');

const NODE_VERSION_API = "https://nodejs.org/dist/index.json";
let _versions = [];
function getAllVersions() {
  if (_versions.length) return _versions;
  return fetch(NODE_VERSION_API).then((resp) => resp.json()).then((data) => _versions = data);
}
async function getVersion(version) {
  const versions = await getAllVersions();
  if (/^\d/.test(version)) {
    version = `v${version}`;
    return versions.find((v) => v.version.startsWith(version));
  }
  if (version === "lts") {
    return versions.find((v) => v.lts);
  }
  if (version === "latest") {
    return versions[0];
  }
  version = version[0].toUpperCase() + version.slice(1).toLowerCase();
  return versions.find((v) => v.lts === version);
}
async function getShasum(version, os, arch) {
  const file = `node-${version}-${os}-${arch}.tar.gz`;
  const sourceFile = `node-${version}.tar.gz`;
  const data = await fetch(
    `https://nodejs.org/dist/${version}/SHASUMS256.txt`
  ).then((resp) => resp.text());
  const hash = data.split("\n").find((el) => el.includes(file))?.split(/\s+/)[0].trim();
  if (!hash) return null;
  const sourceHash = data.split("\n").find((el) => el.includes(sourceFile)).split(/\s+/)[0].trim();
  return { hash, sourceHash };
}

const osName = os.type().toLowerCase();
const arch = process.arch;
const args = process.argv.slice(2);
const versionsDir = args[0];
async function main() {
  const versionQuery = args[1];
  if (fs.existsSync(path.join(versionsDir, versionQuery))) return;
  const versionInfo = await getVersion(versionQuery);
  if (!versionInfo) return;
  const hash = await getShasum(versionInfo.version, osName, arch);
  if (!hash) return;
  console.log(versionInfo.version.slice(1));
  if (fs.existsSync(path.join(versionsDir, versionInfo.version.slice(1))))
    return;
  await writeVersionFile(versionInfo, hash.hash, hash.sourceHash);
}
function writeVersionFile(versionInfo, hash, sourceHash) {
  const { version } = versionInfo;
  return fsP.writeFile(
    path.join(versionsDir, version.slice(1)),
    `binary ${osName}-${arch} "https://nodejs.org/dist/${version}/node-${version}-${osName}-${arch}.tar.gz#${hash}"

install_package "node-${version}" "https://nodejs.org/dist/${version}/node-${version}.tar.gz#${sourceHash}"`
  );
}
main();
