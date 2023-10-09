config = {
    "binaryReleases": {
        "os": ["linux", "darwin", "windows"],
    },
    "dockerReleases": {
        "architectures": ["arm64", "amd64"],
    },
}

# volume for steps to cache Go dependencies between steps of a pipeline
# GOPATH must be set to /go inside the image, which is the case
stepVolumeGo = \
    {
        "name": "gopath",
        "path": "/go",
    }

# volume for pipeline to cache Go dependencies between steps of a pipeline
# to be used in combination with stepVolumeGo
pipelineVolumeGo = \
    {
        "name": "gopath",
        "temp": {},
    }

stepVolumeOC10Tests = \
    {
        "name": "oC10Tests",
        "path": "/srv/app",
    }

pipelineVolumeOC10Tests = \
    {
        "name": "oC10Tests",
        "temp": {},
    }

def pipelineDependsOn(pipeline, dependant_pipelines):
    pipeline["depends_on"] = getPipelineNames(dependant_pipelines)
    return pipeline

def pipelinesDependsOn(pipelines, dependant_pipelines):
    pipes = []
    for pipeline in pipelines:
        pipeline["depends_on"] = getPipelineNames(dependant_pipelines)
        pipes.append(pipeline)

    return pipes

def getPipelineNames(pipelines = []):
    """getPipelineNames returns names of pipelines as a string array

    Args:
      pipelines: array of drone pipelines

    Returns:
      names of the given pipelines as string array
    """
    names = []
    for pipeline in pipelines:
        names.append(pipeline["name"])
    return names

def main(ctx):
    """main is the entrypoint for drone

    Args:
        ctx: drone passes a context with information which the pipeline can be adapted to

    Returns:
        none
    """
    test_pipelines = [
        testHello(ctx),
        UITests(ctx),
    ]

    build_release_pipelines = pipelinesDependsOn(
        dockerReleases(ctx) + binaryReleases(ctx),
        test_pipelines,
    )

    build_release_helpers = pipelinesDependsOn(
        [changelog(ctx), docs(ctx)],
        build_release_pipelines,
    )

    if (ctx.build.event == "pull_request" and "[docs-only]" in ctx.build.title) or \
       (ctx.build.event != "pull_request" and "[docs-only]" in (ctx.build.title + ctx.build.message)):
        # [docs-only] is not taken from PR messages, but from commit messages
        pipelines = [docs(ctx)]

    else:
        pipelines = test_pipelines + build_release_pipelines + build_release_helpers

    pipelineSanityChecks(ctx, pipelines)
    return pipelines + checkStarlark()

def testHello(ctx):
    sonar_env = {
        "SONAR_TOKEN": {
            "from_secret": "sonar_token",
        },
    }
    if ctx.build.event == "pull_request":
        sonar_env.update({
            "SONAR_PULL_REQUEST_BASE": "%s" % (ctx.build.target),
            "SONAR_PULL_REQUEST_BRANCH": "%s" % (ctx.build.source),
            "SONAR_PULL_REQUEST_KEY": "%s" % (ctx.build.ref.replace("refs/pull/", "").split("/")[0]),
        })

    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "testing",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": makeGenerate() + [
            {
                "name": "golangci-lint",
                "image": "owncloudci/golang:1.20",
                "pull": "always",
                "commands": [
                    "make ci-golangci-lint",
                ],
                "volumes": [stepVolumeGo],
            },
            {
                "name": "test",
                "image": "owncloudci/golang:1.20",
                "pull": "always",
                "commands": [
                    "make test",
                ],
                "volumes": [stepVolumeGo],
            },
            {
                "name": "sonarcloud",
                "image": "sonarsource/sonar-scanner-cli:4.7.0",
                "pull": "always",
                "environment": sonar_env,
                "volumes": [stepVolumeGo],
            },
        ],
        "volumes": [pipelineVolumeGo],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
                "refs/pull/**",
            ],
        },
    }

def makeGenerate():
    return [
        {
            "name": "generate nodejs",
            "image": "owncloudci/nodejs:16",
            "pull": "always",
            "commands": [
                "make ci-node-generate",
            ],
            "volumes": [stepVolumeGo],
        },
        {
            "name": "generate go",
            "image": "owncloudci/golang:1.20",
            "pull": "always",
            "environment": {
                "HTTP_PROXY": {
                    "from_secret": "drone_http_proxy",
                },
                "HTTPS_PROXY": {
                    "from_secret": "drone_http_proxy",
                },
            },
            "commands": [
                "apk add protoc",
                "make ci-go-generate",
            ],
            "volumes": [stepVolumeGo],
        },
    ]

def build():
    return [
        {
            "name": "build",
            "image": "owncloudci/golang:1.20",
            "pull": "always",
            "environment": {
                "HTTP_PROXY": {
                    "from_secret": "drone_http_proxy",
                },
                "HTTPS_PROXY": {
                    "from_secret": "drone_http_proxy",
                },
            },
            "commands": [
                "apk add protoc",
                "make build",
            ],
            "volumes": [stepVolumeGo],
        },
    ]

def dockerReleases(ctx):
    pipelines = []
    for arch in config["dockerReleases"]["architectures"]:
        pipelines.append(dockerRelease(ctx, arch))

    pipelines.append(
        pipelineDependsOn(
            releaseDockerManifest(ctx),
            pipelines,
        ),
    )
    pipelines.append(
        pipelineDependsOn(
            releaseDockerReadme(ctx),
            pipelines,
        ),
    )

    return pipelines

def dockerRelease(ctx, arch):
    build_args = [
        "REVISION=%s" % (ctx.build.commit),
        "VERSION=%s" % (ctx.build.ref.replace("refs/tags/", "") if ctx.build.event == "tag" else "latest"),
    ]

    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "docker-%s" % (arch),
        "platform": {
            "os": "linux",
            "arch": arch,
        },
        "steps": makeGenerate() + [
            {
                "name": "build",
                "image": "owncloudci/golang:1.20",
                "pull": "always",
                "environment": {
                    "HTTP_PROXY": {
                        "from_secret": "drone_http_proxy",
                    },
                    "HTTPS_PROXY": {
                        "from_secret": "drone_http_proxy",
                    },
                },
                "commands": [
                    "make release-linux-docker",
                ],
            },
            {
                "name": "dryrun",
                "image": "plugins/docker:latest",
                "pull": "always",
                "settings": {
                    "dry_run": True,
                    "tags": "linux-%s" % (arch),
                    "dockerfile": "docker/Dockerfile.linux.%s" % (arch),
                    "repo": ctx.repo.slug,
                    "build_args": build_args,
                },
                "when": {
                    "ref": {
                        "include": [
                            "refs/pull/**",
                        ],
                    },
                },
            },
            {
                "name": "docker",
                "image": "plugins/docker:latest",
                "pull": "always",
                "settings": {
                    "username": {
                        "from_secret": "docker_username",
                    },
                    "password": {
                        "from_secret": "docker_password",
                    },
                    "auto_tag": True,
                    "auto_tag_suffix": "linux-%s" % (arch),
                    "dockerfile": "docker/Dockerfile.linux.%s" % (arch),
                    "repo": ctx.repo.slug,
                    "build_args": build_args,
                },
                "when": {
                    "ref": {
                        "exclude": [
                            "refs/pull/**",
                        ],
                    },
                },
            },
        ],
        "volumes": [pipelineVolumeGo],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
                "refs/pull/**",
            ],
        },
    }

def ocisServer(storage, accounts_hash_difficulty = 4, volumes = []):
    environment = {
        "IDM_ADMIN_PASSWORD": "admin",  # override the random admin password from `ocis init`
        "OCIS_URL": "https://ocis-server:9200",
        "OCIS_INSECURE": "true",
        "STORAGE_HOME_DRIVER": "%s" % (storage),
        "STORAGE_USERS_DRIVER": "%s" % (storage),
        "PROXY_ENABLE_BASIC_AUTH": True,
        "WEB_UI_CONFIG": "/drone/src/ui/tests/config/drone/web-config.json",
        "OCIS_LOG_LEVEL": "error",
        "OCIS_JWT_SECRET": "aaaa",
        "OCIS_ADMIN_USER_ID": "c59a6ae9-5f5e-4eef-b82e-0e5c34f93e52",
        "MICRO_REGISTRY": "mdns",
        "SETTINGS_GRPC_ADDR": "ocis-server:9191",
    }

    # Pass in "default" accounts_hash_difficulty to not set this environment variable.
    # That will allow OCIS to use whatever its built-in default is.
    # Otherwise pass in a value from 4 to about 11 or 12 (default 4, for making regular tests fast)
    # The high values cause lots of CPU to be used when hashing passwords, and really slow down the tests.
    if (accounts_hash_difficulty != "default"):
        environment["ACCOUNTS_HASH_DIFFICULTY"] = accounts_hash_difficulty

    return [
        {
            "name": "ocis-server",
            "image": "owncloud/ocis:3.0.0-rc.4",
            "pull": "always",
            "detach": True,
            "environment": environment,
            "volumes": volumes,
            "commands": [
                "/usr/bin/ocis init",
                "cp /drone/src/ui/tests/config/drone/proxy.yaml /etc/ocis/",
                "/usr/bin/ocis server",
            ],
        },
        {
            "name": "wait-for-ocis-server",
            "image": "owncloudci/wait-for:latest",
            "pull": "always",
            "commands": [
                "wait-for -it ocis-server:9200 -t 300",
            ],
        },
    ]

def UITests(ctx):
    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "UiTests",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": makeGenerate() +
                 build() + ocisServer("ocis", 4, [stepVolumeOC10Tests]) + [
            {
                "name": "hello-server",
                "image": "owncloudci/alpine:latest",
                "pull": "always",
                "detach": True,
                "environment": {
                    "HELLO_ADMIN_USER_ID": "c59a6ae9-5f5e-4eef-b82e-0e5c34f93e52",
                    "HELLO_JWT_SECRET": "aaaa",
                    "MICRO_REGISTRY": "mdns",
                },
                "commands": [
                    "bin/hello server",
                ],
                "volumes": [
                    {
                        "name": "uploads",
                        "path": "/uploads",
                    },
                ],
            },
            {
                "name": "wait-for-hello-server",
                "image": "owncloudci/wait-for:latest",
                "pull": "always",
                "commands": [
                    "wait-for -it hello-server:9105 -t 300",
                ],
            },
            {
                "name": "WebUIAcceptanceTests",
                "image": "owncloudci/nodejs:16",
                "pull": "always",
                "environment": {
                    "SERVER_HOST": "https://ocis-server:9200",
                    "FEATURE_PATH": "/drone/src/ui/tests/acceptance/features",
                },
                "commands": [
                    "npx playwright install",
                    "make test-acceptance-webui",
                ],
                "volumes": [
                    stepVolumeOC10Tests,
                    {
                        "name": "uploads",
                        "path": "/uploads",
                    },
                ],
            },
        ],
        "volumes": [
            pipelineVolumeGo,
            pipelineVolumeOC10Tests,
            {
                "name": "uploads",
                "temp": {},
            },
        ],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/v*",
                "refs/pull/**",
            ],
        },
    }

def binaryReleases(ctx):
    pipelines = []
    for os in config["binaryReleases"]["os"]:
        pipelines.append(binaryRelease(ctx, os))

    return pipelines

def binaryRelease(ctx, name):
    # uploads binary to https://download.owncloud.com/ocis/ocis/testing/
    target = "/ocis/%s/testing" % (ctx.repo.name.replace("ocis-", ""))
    if ctx.build.event == "tag":
        # uploads binary to eg. https://download.owncloud.com/ocis/ocis/1.0.0-beta9/
        target = "/ocis/%s/%s" % (ctx.repo.name.replace("ocis-", ""), ctx.build.ref.replace("refs/tags/v", ""))

    settings = {
        "endpoint": {
            "from_secret": "upload_s3_endpoint",
        },
        "access_key": {
            "from_secret": "upload_s3_access_key",
        },
        "secret_key": {
            "from_secret": "upload_s3_secret_key",
        },
        "bucket": {
            "from_secret": "upload_s3_bucket",
        },
        "path_style": True,
        "strip_prefix": "dist/release/",
        "source": "dist/release/*",
        "target": target,
    }

    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "binaries-%s" % (name),
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": makeGenerate() + [
            {
                "name": "build",
                "image": "owncloudci/golang:1.20",
                "pull": "always",
                "commands": [
                    "make release-%s" % (name),
                ],
                "volumes": [stepVolumeGo],
            },
            {
                "name": "finish",
                "image": "owncloudci/golang:1.20",
                "pull": "always",
                "commands": [
                    "make release-finish",
                ],
                "volumes": [stepVolumeGo],
            },
            {
                "name": "upload",
                "image": "plugins/s3:1",
                "pull": "always",
                "settings": settings,
                "when": {
                    "ref": [
                        "refs/heads/master",
                        "refs/tags/**",
                    ],
                },
            },
            {
                "name": "changelog",
                "image": "owncloudci/golang:1.20",
                "pull": "always",
                "commands": [
                    "make changelog CHANGELOG_VERSION=%s" % ctx.build.ref.replace("refs/tags/v", "").split("-")[0],
                ],
                "volumes": [stepVolumeGo],
                "when": {
                    "ref": [
                        "refs/tags/v*",
                    ],
                },
            },
            {
                "name": "release",
                "image": "plugins/github-release:1",
                "pull": "always",
                "settings": {
                    "api_key": {
                        "from_secret": "github_token",
                    },
                    "files": [
                        "dist/release/*",
                    ],
                    "title": ctx.build.ref.replace("refs/tags/v", ""),
                    "note": "dist/CHANGELOG.md",
                    "overwrite": True,
                    "prerelease": len(ctx.build.ref.split("-")) > 1,
                },
                "when": {
                    "ref": [
                        "refs/tags/**",
                    ],
                },
            },
        ],
        "volumes": [pipelineVolumeGo],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/**",
                "refs/pull/**",
            ],
        },
    }

def releaseDockerManifest(ctx):
    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "manifest",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": [
            {
                "name": "execute",
                "image": "plugins/manifest:1",
                "pull": "always",
                "settings": {
                    "username": {
                        "from_secret": "docker_username",
                    },
                    "password": {
                        "from_secret": "docker_password",
                    },
                    "spec": "docker/manifest.tmpl",
                    "auto_tag": True,
                    "ignore_missing": True,
                },
            },
        ],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/v*",
            ],
        },
    }

def changelog(ctx):
    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "changelog",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": [
            {
                "name": "generate",
                "image": "owncloudci/golang:1.20",
                "pull": "always",
                "commands": [
                    "make changelog",
                ],
            },
            {
                "name": "diff",
                "image": "owncloudci/alpine:latest",
                "pull": "always",
                "commands": [
                    "git diff",
                ],
            },
            {
                "name": "output",
                "image": "owncloudci/alpine:latest",
                "pull": "always",
                "commands": [
                    "cat CHANGELOG.md",
                ],
            },
            {
                "name": "publish",
                "image": "plugins/git-action:1",
                "pull": "always",
                "settings": {
                    "actions": [
                        "commit",
                        "push",
                    ],
                    "message": "Automated changelog update [skip ci]",
                    "branch": "master",
                    "author_email": "devops@owncloud.com",
                    "author_name": "ownClouders",
                    "netrc_machine": "github.com",
                    "netrc_username": {
                        "from_secret": "github_username",
                    },
                    "netrc_password": {
                        "from_secret": "github_token",
                    },
                },
                "when": {
                    "ref": {
                        "exclude": [
                            "refs/pull/**",
                        ],
                    },
                },
            },
        ],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/pull/**",
            ],
        },
    }

def releaseDockerReadme(ctx):
    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "readme",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": [
            {
                "name": "execute",
                "image": "chko/docker-pushrm:1",
                "pull": "always",
                "environment": {
                    "DOCKER_USER": {
                        "from_secret": "docker_username",
                    },
                    "DOCKER_PASS": {
                        "from_secret": "docker_password",
                    },
                    "PUSHRM_TARGET": "owncloud/${DRONE_REPO_NAME}",
                    "PUSHRM_SHORT": "Docker images for %s" % (ctx.repo.name),
                    "PUSHRM_FILE": "README.md",
                },
            },
        ],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/tags/v*",
            ],
        },
    }

def docs(ctx):
    return {
        "kind": "pipeline",
        "type": "docker",
        "name": "docs",
        "platform": {
            "os": "linux",
            "arch": "amd64",
        },
        "steps": [
            {
                "name": "docs-generate",
                "image": "owncloudci/golang:1.20",
                "environment": {
                    "HTTP_PROXY": {
                        "from_secret": "drone_http_proxy",
                    },
                    "HTTPS_PROXY": {
                        "from_secret": "drone_http_proxy",
                    },
                },
                "commands": [
                    "apk add protoc",
                    "make docs-generate",
                ],
            },
            {
                "name": "prepare",
                "image": "owncloudci/golang:1.20",
                "environment": {
                    "HTTP_PROXY": {
                        "from_secret": "drone_http_proxy",
                    },
                    "HTTPS_PROXY": {
                        "from_secret": "drone_http_proxy",
                    },
                },
                "commands": [
                    "make -C docs docs-copy",
                ],
            },
            {
                "name": "test",
                "image": "owncloudci/golang:1.20",
                "environment": {
                    "HTTP_PROXY": {
                        "from_secret": "drone_http_proxy",
                    },
                    "HTTPS_PROXY": {
                        "from_secret": "drone_http_proxy",
                    },
                },
                "commands": [
                    "make -C docs test",
                ],
            },
            {
                "name": "publish",
                "image": "plugins/gh-pages:1",
                "pull": "always",
                "settings": {
                    "username": {
                        "from_secret": "github_username",
                    },
                    "password": {
                        "from_secret": "github_token",
                    },
                    "pages_directory": "docs/extensions/hello/",
                    "copy_contents": "true",
                    "target_branch": "docs",
                    "delete": "true",
                },
                "when": {
                    "ref": {
                        "exclude": [
                            "refs/pull/**",
                        ],
                    },
                },
            },
        ],
        "trigger": {
            "ref": [
                "refs/heads/master",
                "refs/pull/**",
            ],
        },
    }

def checkStarlark():
    return [{
        "kind": "pipeline",
        "type": "docker",
        "name": "check-starlark",
        "steps": [
            {
                "name": "format-check-starlark",
                "image": "owncloudci/bazel-buildifier",
                "pull": "always",
                "commands": [
                    "buildifier --mode=check .drone.star",
                ],
            },
            {
                "name": "show-diff",
                "image": "owncloudci/bazel-buildifier",
                "pull": "always",
                "commands": [
                    "buildifier --mode=fix .drone.star",
                    "git diff",
                ],
                "when": {
                    "status": [
                        "failure",
                    ],
                },
            },
        ],
        "depends_on": [],
        "trigger": {
            "ref": [
                "refs/pull/**",
            ],
        },
    }]

def pipelineSanityChecks(ctx, pipelines):
    """pipelineSanityChecks helps the CI developers to find errors before running it

    These sanity checks are only executed on when converting starlark to yaml.
    Error outputs are only visible when the conversion is done with the drone cli.

    Args:
      ctx: drone passes a context with information which the pipeline can be adapted to
      pipelines: pipelines to be checked, normally you should run this on the return value of main()

    Returns:
      none
    """

    # check if name length of pipeline and steps are exceeded.
    max_name_length = 50
    for pipeline in pipelines:
        pipeline_name = pipeline["name"]
        if len(pipeline_name) > max_name_length:
            print("Error: pipeline name %s is longer than 50 characters" % (pipeline_name))

        for step in pipeline["steps"]:
            step_name = step["name"]
            if len(step_name) > max_name_length:
                print("Error: step name %s in pipeline %s is longer than 50 characters" % (step_name, pipeline_name))

    # check for non existing depends_on
    possible_depends = []
    for pipeline in pipelines:
        possible_depends.append(pipeline["name"])

    for pipeline in pipelines:
        if "depends_on" in pipeline.keys():
            for depends in pipeline["depends_on"]:
                if not depends in possible_depends:
                    print("Error: depends_on %s for pipeline %s is not defined" % (depends, pipeline["name"]))

    # check for non declared volumes
    for pipeline in pipelines:
        pipeline_volumes = []
        if "volumes" in pipeline.keys():
            for volume in pipeline["volumes"]:
                pipeline_volumes.append(volume["name"])

        for step in pipeline["steps"]:
            if "volumes" in step.keys():
                for volume in step["volumes"]:
                    if not volume["name"] in pipeline_volumes:
                        print("Warning: volume %s for step %s is not defined in pipeline %s" % (volume["name"], step["name"], pipeline["name"]))

    # list used docker images
    print("")
    print("List of used docker images:")

    images = {}

    for pipeline in pipelines:
        for step in pipeline["steps"]:
            image = step["image"]
            if image in images.keys():
                images[image] = images[image] + 1
            else:
                images[image] = 1

    for image in images.keys():
        print(" %sx\t%s" % (images[image], image))
