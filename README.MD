This repository contains functionality for upgrading [okctl](https://github.com/oslokommune/okctl).

# Developing principles

Developers contributing to this project [MUST](https://www.ietf.org/rfc/rfc2119.txt) follow these principles:

* Any code in a migration SHOULD NOT import code from another migration.
  Upgrading is a complex domain, and this is to ensure that upgrades stays simple and isolated. If downstream migrations
  depend on an upstream migrations, we might break any of the downstream migrations, which we want to avoid. If you want
  to reuse logic, import it from somewhere common outside the migrations. However, make sure changes to reuse logic
  doesn't break any of the migrations using the common logic (by having tests, keeping the API stable, etc).
* Upgrades SHOULD be idempotent. How to implement is this is up to the update itself, but the straight forward way is
  using the same logic as okctl uses to check if an upgrade has been run, which is checking the cluster's state storage.
* If you want okctl upgrade not to possibly re-run an upgrade that has been run manually, that is, outside of okctl
  upgrade, then you MUST ensure that the upgrade updates okctl's state, marking the upgrade as run.

# Outputs

This repository will produce the following outputs.

## Binaries

This project creates binaries whose file name will follow this naming convention:

```
okctl-upgrade_<okctl target version>_<os>_<arch>
```

* `okctl target version` is the version of okctl that this upgrade will support after running the upgrade.
*  `os` is `Linux` or `Darwin`
* `arch` is for instance `amd64`
* `description` can in theory be anything, but we will by default name it with the version of okctl we are upgrading
to, plus possible a textual description if that gives value in the given situation.

Examples:
* `okctl-upgrade_0.0.63_Linux_amd64`
* `okctl-upgrade_0.0.64_Linux_amd64`
* `okctl-upgrade_0.0.64_Linux_amd64`
* `okctl-upgrade_0.0.65_Linux_amd64`

## Releases

Every binary will be put in its own release, and tagged with `<okctl target version>`, for instance `0.0.64`.

See https://github.com/oslokommune/okctl-upgrade/releases.
