---
title: Hello (example extension)
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

When getting started with the oCIS development developers need to learn about the building blocks of oCIS extensions.
Without guidance or orientation of the why and what of an extension they may start feeling lost.
The `ocis-hello` repository serves as a blueprint for oCIS extensions.
It allows developers to get started with oCIS extension development by looking at the code, configuration and documentation.

{{< mermaid class="text-center">}}
graph TD
    subgraph ow[ocis-web]
        owh[ocis-web-hello]
    end
    owh ---|"greet()"| ows[ocis-hello-server]
{{< /mermaid >}}


`ocis-hello` provides a simple hello world example with
- a protobuf based greeter API
- a grpc service implementing the API
- a vue.js frontend using the API
- configuration options for admins

It can be integrated into [ownCloud Web](https://github.com/owncloud/web) as documented in the [extensions docs](https://owncloud.github.io/ocis/extensions/#external-web-apps).

