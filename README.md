# Polyroll

Simple CLI utility to create Elasticsearch index templates with ILM policies attached.

## Install

```shell
echo -e "TBD\n"
```

## Usage

To create index templates with ILM policies:

1. Create a yaml config file
2. Run command `polyroll <path-to-config-file>`

## Config

```yaml
elasticsearch:
  host: "elasticsearch-host"
  basicAuthToken: "basic-auth-token"

policies:
  index-policy-foo:
    phases:
      warm: 1
      cold: 30
      delete: 60

templates:
  index-template-foo:
    policy: "index-policy-foo"
    patterns: [ "index-foo-*" ]
```

> **Note**:
> `basicAuthToken` value must be a base64 encoded string `username:password`

### Config rules

Required parameters:

- `elasticsearch.host` - ELK host
- `elasticsearch.basicAuthToken` - ELK basic auth token 

Optional parameters:

- `policies` - map of `<policy-name>: <phases>`
  - `policies.*.phases.{warm|cold|delete}` - optional, any integer value greater than `0` 
- `templates` - map of `<template-name>: <settings>`
  - `templates.*.policy` - required, a valid policy name from `policies` list
  - `templates.*.paterns` - required, non-empty list of strings
