# Quick reference

- Maintained by [Digital Plumbers Union]
- [Open Issues]

# Tags and respective `Dockerfile` links

- [`edge`]

# What is Taskwarrior Server?

> [Taskwarrior] is Free and Open Source Software that manages your TODO list from the command line.

Taskwarrior Server (a.k.a Taskserver or `taskd`) is a remote service that syncs your tasks across multiple devices. With
a personal Taskwarrior Server you can more efficiently keep track of your TODOs, backup your tasks and have complete
control of your configuration and privacy.

# How to use this image

The required certificate configuration files for `taskd` should be mounted them at the following locations:

- `client.cert` → `/var/lib/taskd/certs/client.cert.pem`
- `client.key` → `/var/lib/taskd/certs/client.key.pem`
- `server.cert` → `/var/lib/taskd/certs/server.cert.pem`
- `server.key` → `/var/lib/taskd/certs/server.key.pem`
- `server.crl` → `/var/lib/taskd/certs/server.crl.pem`
- `ca.cert` → `/var/lib/taskd/certs/ca.cert.pem`

The recommended approach is to combine the files in a local folder `certs` and use a read-only volume mount.

### Starting server

``` sh
$ docker run -d --name taskd \
    -p 53589:53589 \
    --mount source=certs,destination=/var/lib/taskd/certs,readonly \
    dpush/taskwarrior-server:edge
```

# Image Variants

The `dpush/taskwarrior-server` image currently only comes in one flavor.

## `dpush/taskwarrior-server:edge`

This is the bleeding edge image that is regularly iterated on and provides no guarantees in terms of stability or
support. It is based on the [`Alpine` official image] which is about ~5MB in size and uses [musl libc] as its [C
standard library]. This image may be used as the base of other images but is discouraged due to rate at which it may
change.

# License

This work is [dual-licensed] as `Apache-2.0 OR MIT` to achieve compatibility with GPLv2 and to be friendly for both
individuals and enterprise.

As with all container images, these images probably contain software which may be under other licenses; this includes
licenses for direct or indirect dependencies of the main program. It is the image user's responsibility to ensure that
any use of this image complies with all of the relevant licenses of software contained within.

[Digital Plumbers Union]:https://github.com/digital-plumbers-union
[Open Issues]:https://github.com/digital-plumbers-union/taskwarrior-operator/issues
[Taskwarrior]:https://taskwarrior.org
[dual-licensed]:https://github.com/digital-plumbers-union/taskwarrior-operator/blob/master/LICENSE
[`edge`]:https://github.com/digital-plumbers/digital-plumbers-union/taskwarrior-server/blob/master/docker-images/taskwarrior-server/Dockerfile
[`Alpine` official image]:https://hub.docker.com/_/alpine 
[musl libc]:https://musl.libc.org
[C standard library]:https://www.etalabs.net/compare_libcs.html
