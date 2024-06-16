# Usage

panosse can be used as a [standalone binary](#standalone-binary) or with
[Docker](#docker).

## Standalone binary

To use panosse as a standalone binary, download the latest release from the
GitHub Releases page: <https://github.com/ludelafo/panosse/releases>.

> [!IMPORTANT]
>
> flac and metaflac must be installed on your computer in order to use panosse.
> You can install them using your package manager or download them from the
> [Xiph website](https://xiph.org/flac/download.html).
>
> panosse was tested with flac version 1.4.2 and metaflac version 1.4.2.

```sh
# Run panosse as a standalone binary
./panosse --help
```

## Docker

To use panosse with Docker, pull the Docker image from GitHub Container Registry
and run it as a container: <https://ghcr.io/ludelafo/panosse>.

> [!IMPORTANT]
>
> flac and metaflac are already installed in the Docker image. No need to
> install them on your computer.

```sh
# Run panosse as a Docker container
docker run --rm ghcr.io/ludelafo/panosse --help

# Run panosse as a Docker container with a volume
docker run --rm \
  --volume "$(pwd)/custom-config.yaml:/config/custom-config.yaml" \
  --volume "$(pwd):/files" \
  ghcr.io/ludelafo/panosse --config-file /config/custom-config.yaml  verify /files/file.flac
```
