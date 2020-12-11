---
title: Hello
geekdocRepo: https://github.com/owncloud/ocis-hello
geekdocEditPath: edit/master/docs
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

[![GitHub](https://img.shields.io/github/license/owncloud/ocis-hello)](https://github.com/owncloud/ocis-hello/blob/master/LICENSE)

## Abstract

When getting started with ocis development developers need to learn about the building blocks of ocis extensions.
Without guidance or orientation of the why and what of an extension they may start feeling lost.
The `ocis-hello` repository serves as a blueprint for ocis extensions.
It allows developers to get started with ocis extension development by looking at the code, configuration and documentation.

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

It can be integrated into [ocis web](https://github.com/owncloud/phoenix) as documented in the [extensions docs](https://owncloud.github.io/ocis/extensions/#external-phoenix-apps).

## Table of Contents

{{< toc-tree >}}
