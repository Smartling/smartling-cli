# User Documentation

See the [Wiki](https://github.com/Smartling/smartling-cli/wiki) page for this repository.

# Development
For developers interested in modifying the tool.

## Building package

```
make <target>
```

Where target is:

* `deb` for building Debian packages:
   ```
   make deb
   ```

* `rpm` for building Fedora packages:
   ```
   make rpm
   ```

Specific settings can be set in built-time:

*VERSION*:

```
make VERSION=<version-string> <target>
```

*MAINTAINER*:

```
make MAINTAINER=<maintainer> <target>
```

An executable named `smartling-cli` should become available in your
`$GOPATH/bin`.
