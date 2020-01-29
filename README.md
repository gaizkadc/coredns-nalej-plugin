# coredns-nalej-plugin

CoreDNS Plugin for Nalej. This plugin will be used by the `coredns` instance deployed on the platform. This instance serves as an external DNS for the endpoint of the applications on the platform.

## Getting Started

### Prerequisites

* [`system-model`](https://github.com/nalej/system-model)

### Build and compile

In order to build and compile this repository use the provided Makefile:

```
make all
```

This operation generates the binaries for this repo, downloads the required dependencies, runs existing tests and generates ready-to-deploy Kubernetes files.

### Enabling debug mode

To enable the debug mode on the plugin, configure the `debug` flag into the plugin config file:

```
    .:53 {
      corednsnalejplugin {
        systemModelAddress system-model.__NPH_NAMESPACE:8800
        debug
      }
    }
```

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```
make test
```

No test files are available for this repo at the moment.

### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```
make dep
```

In order to have all dependencies up-to-date run:

```
dep ensure -update -v
```

## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the available versions, see the [tags on this repository](https://github.com/nalej/coredns-nalej-plugin/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/coredns-nalej-plugin/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.
