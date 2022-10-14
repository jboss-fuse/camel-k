# CI configuration

This directory contains any required configuration for product in order to automate the build and release work.

## Divergencies synchronization directory (ci/downstream)

In the `ci/downstream` directory you can provide all those configuration files that we know for sure are divergenging from upstream. In this way we simplify the managment of those changes and we guarantee they are always in sync with the conflict we solve. The files contained in this directory must reflect the structure of the project and they are used when initializing a new branch from upstream.

There is a `sync` process (available at https://github.com/squakez/sync-pin-divergencies) that is in charge to copy all the changes performed to the original sources, so, the developer must focus on adding the required feature. The process will take care to bring the changes back to the init directory.