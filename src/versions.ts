const NODE_VERSION_API = 'https://nodejs.org/dist/index.json';

export type NodeVersion = {
  version: string;
  date: string;
  lts: boolean | string;
};

let _versions: NodeVersion[] = [];

export function getAllVersions() {
  if (_versions.length) return _versions;

  return fetch(NODE_VERSION_API)
    .then(resp => resp.json() as Promise<NodeVersion[]>)
    .then(data => (_versions = data));
}

export async function getVersion(version: string) {
  const versions = await getAllVersions();

  // nodenv install 20
  // nodenv install 20.10
  // nodenv install 20.10.0
  if (/^\d/.test(version)) {
    version = `v${version}`;
    return versions.find(v => v.version.startsWith(version));
  }

  if (version === 'lts') {
    return versions.find(v => v.lts);
  }

  if (version === 'latest') {
    return versions[0];
  }

  // other lts version
  version = version[0].toUpperCase() + version.slice(1).toLowerCase();
  return versions.find(v => v.lts === version);
}

export async function getShasum(
  version: string,
  os: string,
  arch: string
): Promise<{ hash: string; sourceHash: string } | null> {
  const file = `node-${version}-${os}-${arch}.tar.gz`;
  const sourceFile = `node-${version}.tar.gz`;
  const data = await fetch(
    `https://nodejs.org/dist/${version}/SHASUMS256.txt`
  ).then(resp => resp.text());

  const hash = data
    .split('\n')
    .find(el => el.includes(file))
    ?.split(/\s+/)[0]
    .trim();

  if (!hash) return null;

  const sourceHash = data
    .split('\n')
    .find(el => el.includes(sourceFile))!
    .split(/\s+/)[0]
    .trim();

  return { hash, sourceHash };
}
