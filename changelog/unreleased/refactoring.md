Enhancement: Streamline project structure

- We have aligned the project structure of ocis-hello with other repositories and improved error logging.
- When running this service through `make watch` it now regenerates embedded assets properly as soon as the web bundle
  is changed / saved.
- In the package.json file we're now declaring owncloud-design-system as peer dependency, since we're actively using it.
  It comes from ocis-web, so we don't need to bundle it.

https://github.com/owncloud/ocis-hello/pull/79
https://github.com/owncloud/ocis-hello/pull/80
