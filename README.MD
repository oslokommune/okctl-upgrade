This repository contains functionality for upgrading [okctl](https://github.com/oslokommune/okctl).

# Developing principles

Developers contributing to this project MUST follow these principles:

* Any code in a migration [SHOULD NOT](https://www.ietf.org/rfc/rfc2119.txt) import code from another migration.
  Upgrading is a complex domain, and this is to ensure that upgrades stays simple and isolated. If downstream migrations
  depend on an upstream migrations, we might break any of the downstream migrations, which we want to avoid. If you want
  to reuse logic, import it from somewhere common outside the migrations. However, make sure changes to reuse logic
  doesn't break any of the migrations using the common logic (by having tests, keeping the API stable, etc).

# Outputs

This repository will produce the following outputs.

## Binaries

Output of this project will

This project creates migrations just like Flyway. Each migration is a binary, whose file name will follow this naming
convention:

```
okctl_<migration number>_<okctl target version>_<os>_<arch>
```

* `migration number` will be used to determine execution order, and must thus be a parsable integer.
*  
* `description` can in theory be anything, but we will by default name it with the version of okctl we are upgrading
to, plus possible a textual description if that gives value in the given situation.

Examples:
* `okctl-upgrade_1_0.0.63_Linux_amd64`
* `okctl-upgrade_2_0.0.64_Linux_amd64`
* `okctl-upgrade_3_0.0.64_Linux_amd64`
* `okctl-upgrade_4_0.0.65_Linux_amd64`

## Releases

Every binary will be put in its own release, and tagged with

```
<migration number>-<okctl target version>
```

Examples:
* `2_0.0.64`
* `3_0.0.64`

See https://github.com/oslokommune/okctl-upgrade/releases.
