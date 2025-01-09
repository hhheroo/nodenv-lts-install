import os from 'os';
import { getShasum, getVersion, NodeVersion } from './versions';
import path from 'path';
import fs from 'fs';
import fsP from 'fs/promises';

const osName = os.type().toLowerCase();
const arch = process.arch;
const args = process.argv.slice(2);

const versionsDir = args[0]

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

function writeVersionFile(
  versionInfo: NodeVersion,
  hash: string,
  sourceHash: string
) {
  const { version } = versionInfo;

  return fsP.writeFile(
    path.join(versionsDir, version.slice(1)),
    `binary ${osName}-${arch} "https://nodejs.org/dist/${version}/node-${version}-${osName}-${arch}.tar.gz#${hash}"

install_package "node-${version}" "https://nodejs.org/dist/${version}/node-${version}.tar.gz#${sourceHash}"`
  );
}

main();
