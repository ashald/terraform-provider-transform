# Terraform Transform Data Sources

## Overview

This provider defines a set of data sources providing data transformation primitives missing from core Terraform.

## Data Sources

This plugin defines following data sources: 
* `transform_group_by_value` - group map keys by map values

## Reference

### transform_group_by_value

#### Arguments

The following arguments are supported:

* `input` - (Required) A map with both keys and values as strings
* `extract` - (Required) A key (one of values in the `input`) which value should be returned in output after the grouping.
Because of Terraform limitation it's not possible to return the full grouped map (which would be a map of strings to lists).  

#### Attributes

The following attribute is exported:

* `output` - A list of keys from `input` that contain `extract` value. Sorted in lexicographical order.

## Installation

> Terraform automatically discovers the Providers when it parses configuration files.
> This only occurs when the init command is executed.

Currently Terraform is able to automatically download only
[official plugins distributed by HashiCorp](https://github.com/terraform-providers).

[All other plugins](https://www.terraform.io/docs/providers/type/community-index.html) should be installed manually.

> Terraform will search for matching Providers via a
> [Discovery](https://www.terraform.io/docs/extend/how-terraform-works.html#discovery) process, **including the current
> local directory**.

This means that the plugin should either be placed into current working directory where Terraform will be executed from
or it can be [installed system-wide](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

## Usage

### main.tf
```hcl
locals {
  input = {
    "aaa/bbb/111" = "val1"
    "aaa/ccc/111" = "val1"
    "aaa/ddd/222" = "val2"
  }
}

data "transform_group_by_value" "data" { input="${local.input}" extract = "val1" }

data "transform_glob_map" "include"         { input="${local.input}" pattern="aaa/*/111"                }
data "transform_glob_map" "exclude"         { input="${local.input}" pattern="aaa/*/111" exclude = true }
data "transform_glob_map" "include_w_sep"   { input="${local.input}" pattern="aaa/*"     separator="/"  }
data "transform_glob_map" "include_wo_sep"  { input="${local.input}" pattern="aaa/*"                    }

output "result" {
  value = {
    grouped = "${data.transform_group_by_value.data.items}"

    glob_include        = "${data.transform_glob_map.include.output}"
    glob_exclude        = "${data.transform_glob_map.exclude.output}"

    glob_include_w_sep  = "${data.transform_glob_map.include_w_sep.output}"
    glob_include_wo_sep = "${data.transform_glob_map.include_wo_sep.output}"
  }
}
```

### Download
```bash
$ wget "https://github.com/ashald/terraform-provider-transform/releases/download/v1.1.0/terraform-provider-transform_v1.1.0-$(uname -s | tr '[:upper:]' '[:lower:]')-amd64"
$ chmod +x ./terraform-provider-transform*
```

### Init
```bash
$ ls -1
  main.tf
  terraform-provider-transform_v1.1.0-linux-amd64

$ terraform init
  
  Initializing provider plugins...
  
  The following providers do not have any version constraints in configuration,
  so the latest version was installed.
  
  To prevent automatic upgrades to new major versions that may contain breaking
  changes, it is recommended to add version = "..." constraints to the
  corresponding provider blocks in configuration, with the constraint strings
  suggested below.
  
  * provider.transform: version = "~> 1.1"
  
  Terraform has been successfully initialized!
  
  You may now begin working with Terraform. Try running "terraform plan" to see
  any changes that are required for your infrastructure. All Terraform commands
  should now work.
  
  If you ever set or change modules or backend configuration for Terraform,
  rerun this command to reinitialize your working directory. If you forget, other
  commands will detect it and remind you to do so if necessary.
```

### Apply

```bash
$ terraform apply
  data.transform_glob_map.include: Refreshing state...
  data.transform_glob_map.include_wo_sep: Refreshing state...
  data.transform_glob_map.include_w_sep: Refreshing state...
  data.transform_glob_map.exclude: Refreshing state...
  data.transform_group_by_value.data: Refreshing state...
  
  Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
  
  Outputs:
  
  result = {
    glob_exclude = map[aaa/ddd/222:val2]
    glob_include = map[aaa/ccc/111:val1 aaa/bbb/111:val1]
    glob_include_w_sep = map[]
    glob_include_wo_sep = map[aaa/bbb/111:val1 aaa/ddd/222:val2 aaa/ccc/111:val1]
    grouped = [aaa/ccc/111 aaa/bbb/111]
  }

```


## Development

### Go

In order to work on the provider, [Go](http://www.golang.org) should be installed first (version 1.8+ is *required*).
[goenv](https://github.com/syndbg/goenv) and [gvm](https://github.com/moovweb/gvm) are great utilities that can help a
lot with that and simplify setup tremendously. 
[GOPATH](http://golang.org/doc/code.html#GOPATH) should be setup correctly and as long as `$GOPATH/bin` should be
added `$PATH`.

### Source Code

Source code can be retrieved either with `go get`

```bash
$ go get -u -d github.com/ashald/terraform-provider-transform
```

or with `git`
```bash
$ mkdir -p ${GOPATH}/src/github.com/ashald/terraform-provider-transform
$ cd ${GOPATH}/src/github.com/ashald/terraform-provider-transform
$ git clone git@github.com:ashald/terraform-provider-transform.git .
```

### Dependencies

This project uses `govendor` to manage its dependencies. When adding a dependency on a new package it should be fetched
with:
```bash
$ govendor fetch +o
```

### Test

```bash
$ make test
  go test -v ./...
  ?   	github.com/ashald/terraform-provider-transform	[no test files]
  === RUN   TestGlobMapDataSource
  --- PASS: TestGlobMapDataSource (0.09s)
  === RUN   TestGroupByValueDataSource
  --- PASS: TestGroupByValueDataSource (0.09s)
  === RUN   TestProvider
  --- PASS: TestProvider (0.00s)
  PASS
  ok  	github.com/ashald/terraform-provider-transform/transform	(cached)
  go vet ./...
```

### Build
In order to build plugin for the current platform use [GNU]make:
```bash
$ make build
  go build -o terraform-provider-transform_v1.1.0

```

it will build provider from sources and put it into current working directory.

If Terraform was installed (as a binary) or via `go get -u github.com/hashicorp/terraform` it'll pick up the plugin if 
executed against a configuration in the same directory.

### Release

In order to prepare provider binaries for all platforms:
```bash
$ make release
  GOOS=darwin GOARCH=amd64 go build -o './release/terraform-provider-transform_v1.1.0-darwin-amd64'
  GOOS=linux GOARCH=amd64 go build -o './release/terraform-provider-transform_v1.1.0-linux-amd64'
```

### Versioning

This project follow [Semantic Versioning](https://semver.org/)

### Changelog

This project follows [keep a changelog](https://keepachangelog.com/en/1.0.0/) guidelines for changelog.

### Contributors

Please see [CONTRIBUTORS.md](./CONTRIBUTORS.md)

## License

This is free and unencumbered software released into the public domain. See [LICENSE](./LICENSE)
