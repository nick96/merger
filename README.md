# merger

A super simple utility for merging github PRs that pass all their checks and are
mergeable.

`merger` is intended to be used on some sort of schedule. This can easily be
done using something like GitHub workflows (see
[.github/workflows/merger.yml](.github/workflows/merger.yml)).

:warning: Only use merger if you have CI setup so that your required checks are
run on PRs :warning:

## Usage

```
Usage of merger:
  -label string
    	Label to filter pull requests by. Only PRs with this label will be checked and merged.
  -repository string
    	GitHub repository to check issues on. Should be of the for <owner>/<repo>. Uses GITHUB_REPOSITORY if not provided.
  -token string
    	GitHub token used for authentication. Uses GITHUB_TOKEN if not provided.
```

If you're using `merger` in a GitHub workflow `GITHUB_REPOSITORY` is an
environment variable provided by default and `GITHUB_TOKEN` can be provided as a
secret. You can then call `merger` like so:

``` bash
merger -label dependencies
```

Where any PR with the `dependencies` label (e.g. dependabot) will be merged if
its checks are passing and it is mergeable.

## License

Licensed under

 * MIT license
   ([LICENSE-MIT](LICENSE-MIT) or http://opensource.org/licenses/MIT)

## Contribution

Any contribution shall be licensed under the MIT license without any additional
terms or conditions.