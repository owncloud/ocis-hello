Bugfix: Fix api path

The server path coming from ownCloud Web now has an enforced trailing slash.
Concatenating the api path to the server path resulted in a path containing a double slash.

https://github.com/owncloud/ocis-hello/pull/99
