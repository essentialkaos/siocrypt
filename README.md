<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/y/siocrypt"><img src="https://kaos.sh/y/31ef70b4136e4b48aab5d4b934e11eac.svg" alt="Codacy badge" /></a>
  <br/>
  <a href="https://kaos.sh/w/siocrypt/ci"><img src="https://kaos.sh/w/siocrypt/ci-push.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/siocrypt/codeql"><img src="https://kaos.sh/w/siocrypt/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`siocrypt` is a tool for encrypting/decrypting arbitrary data streams using [Data At Rest Encryption](https://github.com/essentialkaos/sio/blob/master/DARE.md) (_DARE_).

### Installation

#### From source

To build the `siocrypt` from scratch, make sure you have a working Go 1.23+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/siocrypt@latest
```

#### Container Image

The latest version of `siocrypt` also available as container image on [GitHub Container Registry](https://kaos.sh/p/siocrypt) and [Docker Hub](https://kaos.sh/d/siocrypt):

```bash
podman run --rm -it ghcr.io/essentialkaos/siocrypt:latest
# or
docker run --rm -it ghcr.io/essentialkaos/siocrypt:latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/siocrypt/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) siocrypt
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo siocrypt --completion=bash 1> /etc/bash_completion.d/siocrypt
```

ZSH:
```bash
sudo siocrypt --completion=zsh 1> /usr/share/zsh/site-functions/siocrypt
```

Fish:
```bash
sudo siocrypt --completion=fish 1> /usr/share/fish/vendor_completions.d/siocrypt.fish
```

### Man documentation

You can generate man page using next command:

```bash
siocrypt --generate-man | sudo gzip > /usr/share/man/man1/siocrypt.1.gz
```

### Usage

<img src=".github/images/usage.svg"/>

### CI Status

| Branch | Status |
|--------|----------|
| `master` | [![CI](https://kaos.sh/w/siocrypt/ci-push.svg?branch=master)](https://kaos.sh/w/siocrypt/ci-push?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/siocrypt/ci-push.svg?branch=develop)](https://kaos.sh/w/siocrypt/ci-push?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/.github/blob/master/CONTRIBUTING.md).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://kaos.dev"><img src="https://raw.githubusercontent.com/essentialkaos/.github/refs/heads/master/images/ekgh.svg"/></a></p>
