# nodenv-lts-install

This nodenv plugin can install node version with "lts", "latest", lts-codenames and any node version prefix.

It is highly inspired by [n](https://github.com/tj/n).

It requires the `node-build` plugin to be installed.

## Installation

### Prepare

Make sure you have at least one node version installed, because it needs node to run the plugin.

```shell
nodenv install 20.18.1
```

### Installing as a nodenv plugin

Make sure you have the latest nodenv and node-build versions, then run:

```shell
git clone https://github.com/hhheroo/nodenv-lts-install.git $(nodenv root)/plugins/nodenv-lts-install
```


## Installing Node.js Versions

### latest

Download latest node active version with `latest`.

```shell
# install current latest node version
nodenv install latest
```

### lts

Download latest lts node active version with `lts`.

```shell
# install current latest lts node version
nodenv install lts
```

### lts-codenames

Download a specify lts node version with its [codenames](https://github.com/nodejs/Release/blob/main/CODENAMES.md).

```shell
# install latest node22
nodenv install jod
# install latest node20
nodenv install iron
```


### numeric version numbers

Download a specify node version with numeric version numbers.

```shell
# install latest node20.x.x
nodenv install 20
# install latest node20.1.x
nodenv install 20.1
# install node20.1.0
nodenv install 20.1.0
```
