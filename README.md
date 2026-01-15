# Find GitHub paired PRs

This tool helps you identify when a PR to a _base repo_ depends on a "paired" PR in a _dependency repo_, and compile the base REPO with the PR in the dependency.

An an example, [mit-pdos/perennial](https://github.com/mit-pdos/perennial) (our base repo) depends on [goose-lang/goose](https://github.com/goose-lang/goose). We run:

```sh
go run . -base mit-pdos/perennial -pr 313 -dependency goose-lang/goose -verbose
```

This tells us that perennial PR 313 depends on goose PR 119, which is branch `tchajed/gen-wp-globals-alloc`. It allows us in CI to compile perennial with that branch of goose, anticipating that the dependency PR will be merged together with the base PR. (In this case, a URL is not output since the goose PR is already merged.)
