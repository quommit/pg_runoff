## pg_runoff CLI development guidelines

### Overview

The `pg_runoff` extension provides a Command Line Interface (CLI) as a Go binary named `runoff`. This binary consists of a series of subcommands to assist users in managing and computing a threshold runoff model for a given drainage basin. The sections below give some guidelines on how to configure the build environment and compile a single static binary file in development.

### Build environment configuration

The `pg_runoff` CLI is packed along with the `pg_runoff` PostgreSQL extension itself, which is in turn distributed as a Docker image. Several `pg_runoff` images of the same version may exist, the difference among them being the SIOSE data set which is included along with the software. Every `pg_runoff` image is based on a spatial subset of the SIOSE database whose geographical limits agree on a particular administrative boundary (e.g. a municipality, a province, an autonomous region) or a Geohash extent. Moreover, this spatial subset will belong to either the 2005, 2011 or 2014 SIOSE edition. At the time of writing, SIOSE subsets are distributed as Docker images based on Alpine Linux 3.9. Hence, as part of the `pg_runoff` image multistage build process, the `runoff` binary is compiled using the `golang:1.12.14-alpine3.9`.
In order to modify the source code and compile the CLI in development, you should run the following Docker container:

```
~/gitrepos/pg_runoff$ docker run -it --rm -v $(pwd)/go:/go --workdir /go/src golang:1.12.14-alpine3.9 /bin/sh
```

This will leave you in an interactive shell where you can install the needed build dependencies:

```
/go/src # apk add --no-cache gcc git libc-dev
```

You may also install a terminal based Emacs for source code editing:

```
/go/src # apk add --no-cache emacs-nox
```

Finally install the Cobra Go package, which is `pg_runoff` CLI's only third-party dependency (you may skip this step in case `go/bin/cobra` and `github.com/spf13/cobra/cobra` already exists in your repo directory):

```
/go/src # go get -u github.com/spf13/cobra/cobra 
```

Remember you may use the `cobra` binary to automatically generate the scaffolding for new subcommands (e.g. `cobra add` or `cobra add -p`).

### Static binary compilation

To obtain a static binary file, run the `go build` command from within the `runoff` directory as follows:

```
/go/src # cd runoff
/go/src/runoff # GOOS=linux GARCH=amd64 go build -ldflags "-linkmode external -extldflags -static" -installsuffix cgo -o runoff -a main.go 
```
Should there be no errors, you'll be left with a binary file named `/go/src/runoff/runoff`.

### Limitations

Most `runoff` subcommands spawn external processes such as `ogrinfo`, `ogr2ogr` and `psql` using Go's `exec.Command`. Using `go run` from within the golang container will yield runtime errors anytime one of these external processes is invoked since neither **GDAL** nor **PostgreSQL** will be available. Therefore, build the `pg_runoff` image with the modified sources, create a container and execute the updated `runoff` utilities in order to test the binary.
