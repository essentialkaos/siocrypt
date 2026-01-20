<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/y/siocrypt"><img src="https://kaos.sh/y/31ef70b4136e4b48aab5d4b934e11eac.svg" alt="Codacy badge" /></a>
  <a href="https://kaos.sh/w/siocrypt/ci"><img src="https://kaos.sh/w/siocrypt/ci-push.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/siocrypt/codeql"><img src="https://kaos.sh/w/siocrypt/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`siocrypt` is a tool for encrypting/decrypting arbitrary data streams using [Data At Rest Encryption](https://github.com/essentialkaos/sio/blob/master/DARE.md) (_DARE_).

### Installation

#### From source

To build the `siocrypt` from scratch, make sure you have a working Go 1.24+ workspace (_[instructions](https://go.dev/doc/install)_), then:

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

### Upgrading

You can update prebuilt `siocrypt` binary to the latest release using [self-update feature](https://github.com/essentialkaos/.github/blob/master/APPS-UPDATE.md):

```bash
siocrypt --update
```

This command will runs a self-update in interactive mode. If you want to run a quiet update (_no output_), use the following command:

```bash
siocrypt --update=quiet
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
siocrypt --completion=bash | sudo tee /etc/bash_completion.d/siocrypt > /dev/null
```

ZSH:
```bash
siocrypt --completion=zsh | sudo tee /usr/share/zsh/site-functions/siocrypt > /dev/null
```

Fish:
```bash
siocrypt --completion=fish | sudo tee /usr/share/fish/vendor_completions.d/siocrypt.fish > /dev/null
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
