# rcstate

[![Go Report Card](https://goreportcard.com/badge/github.com/marintailor/rcstate)](https://goreportcard.com/report/github.com/marintailor/rcstate)

## List of contents:

1. [Overview](#overview)
2. [Installation](#installation)
3. [Requirements](#requirements)
4. [Usage](#usage)
4. [Server mode](#server-mode)

<a name="overview"></a>
## Overview

rcstate is a CLI app written in Go to manage the state of resources in Google Cloud.

<a name="installation"></a>
## Installation

```bash
go install github.com/marintailor/rcstate@latest
```

<a name="requirements"></a>
## Requirements

### Google Cloud services

An environment variable `GOOGLE_APPLICATION_CREDENTIALS` is required for authentication with Google Cloud services.

For more information check [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials#GAC) documentation.

### AWS Route 53

Currently the DNS record is created with Route 53 DNS service.

All requests are made using the AWS SDK for Go, and credentials should be stored in `~/.aws/credentials` file.

For more information check [AWS SDK for Go](https://github.com/aws/aws-sdk-go).

<a name="usage"></a>
## Usage

### Manage environments

An environment represents one or more groups of resources that are already present in Google Cloud.

The environments are declared in a YAML file, and it is provided as an option with flag `--env-file`.

The path to the environment file can be also set as an environment variable `RCSTATE_ENV_FILE`.

NOTE: Option flag `--env-file` has higher priority.


**Examples:**

* show an environment

```bash
rcstate env show \
  --name <environment_name> \
  --env-file <environment_file>
```

* show all environments

```bash
rcstate env show \
  --all \
  --env-file <environment_file>
```

* show an environments with label "api"

```bash
rcstate env show \
  --name <environment_name> \
  --label api \
  --env-file <environment_file>
```

* show all environments with label "api"

```bash
rcstate env show \
  --all \
  --label api \
  --env-file <environment_file>
```

* change state of an environment

```bash
rcstate env up/down \
  --name <environment_name> \
  --env-file <environment_file>
```

*  change state of an environment labeled with "dev" and "qa"

```bash
rcstate env up/down \
  --name <environment_name> \
  --label qa,api \
  --env-file <environment_file>
```

* change state of all environments

```bash
rcstate env up/down \
  --all \
  --env-file <environment_file>
```

* change state of all environments labeled with "dev" and "qa"

```bash
rcstate env up/down \
  --all \
  --label qa,api \
  --env-file <environment_file>
```

**Schema example of the environment file:**

```yaml
variables:    # Global variables that are accessible from all environments
  APP_NAME:  test-app
  DNS_DOMAIN: example.com
  DNS_TYPE: A
  SSH_KEY: /home/user/.ssh/private_key
  SSH_PORT: 22
  SSH_USER: user
environment:    # List of the environments
  - name: dev    # Environment name
    label: dev    # Environment label(s)
    group:    # List of groups where resource are grouped
      - name: group-dev-1    # Group name
        project: project-dev-1    # GCP Project ID
        zone: us-central1-a    # GCP Zone name
        resource:    # List of different types of resources are specified per group
          vm:    # Virtual Machines
            script:    # Script at resource level will be run on all instance
              ssh:    # SSH configuration
                key: "{{ .SSH_KEY }}"
                port: "{{ .SSH_PORT }}"
                user: "{{ .SSH_USER }}"
              up:    # Shell commands to be executed AFTER instance is started
                - sudo shutdown -h +30
              down:    # Shell commands to be executed BEFORE instance is stopped
                - ~/clean-up.sh
            instance:    # List of the Virtual Machine instances
              - name: vm-dev-1    # Instance name
                record:    # Instance DNS record
                  domain: "{{ .DNS_DOMAIN }}"
                  external_ip: true    # Use instance's external IP for the DNS record
                  ip:    # List of ip addresses for the DNS record
                    - 123.123.123.123
                    - 145.145.145.145
                  type: "{{ .DNS_TYPE }}"    # The type of the DNS record
                  zone: "{{ .APP_NAME }}.dev-1.{{ .DNS_DOMAIN }}"    # The DNS record
                script:    # Script at instance level will be run per instance
                  ssh:
                    key: "{{ .SSH_KEY }}"
                    port: "{{ .SSH_PORT }}"
                    user: "{{ .SSH_USER }}"
                  up:
                    - cd /data/{{ .APP_NAME }} && docker-compose up -d
                  down:
                    - cd /data/{{ .APP_NAME }} && docker-compose down
      - name: group-dev-2
        project: project-dev-2
        zone: us-central1-a
        resource:
          vm:
            script: shutdown +10
            instance:
              - name: vm-dev-1
                script:
                  ssh:
                    key: "{{ .SSH_KEY }}"
                    port: "{{ .SSH_PORT }}"
                    user: "{{ .SSH_USER }}"
                  up:
                    - curl "https://{{ .APP_NAME }}.{{ .DOMAIN }}/health" \
              - name: vm-dev-2
                record:
                  domain: "{{ .DNS_DOMAIN }}"
                  external_ip: true
                  type: "{{ .DNS_TYPE }}"
                  zone: "{{ .APP_NAME }}.dev-2.{{ .DNS_DOMAIN }}"
                script:
                  ssh:
                    key: "{{ .SSH_KEY }}"
                    port: "{{ .SSH_PORT }}"
                    user: "{{ .SSH_USER }}"
                  up:
                    - curl "https://{{ .APP_NAME }}.{{ .DOMAIN }}/api/v1/start"
                  down:
                    - curl "https://{{ .APP_NAME }}.{{ .DOMAIN }}/api/v1/stop"
  - name: qa
    label: qa
    group:
      project: project-qa-1
      zone: us-central1-a
      - name: group-qa-1
        resource:
          vm:
            instance:
              - name: vm-qa-1
                record:
                  domain: "{{ .DNS_DOMAIN }}"
                  external_ip: true
                  type: "{{ .DNS_TYPE }}"
                  zone: "{{ .APP_NAME }}.qa.{{ .DNS_DOMAIN }}"
                script:
                  ssh:
                    key: "{{ .SSH_KEY }}"
                    port: "{{ .SSH_PORT }}"
                    user: "{{ .SSH_USER }}"
                  up:
                    - wget -O - https://{{ .APP_NAME }}.{{ .DOMAIN }}/init.sh | bash
```

### Manage virtual machine (Google Cloud Engine)

* list all virtual machine instances in specific project and zone

```bash
rcstate vm list \
  --project <project_id> \
  --zone <zone_name>
```

* start an instance in specific project and zone

```bash
rcstate vm start \
  --name <instance_name> \
  --project <project_id> \
  --zone <zone_name>
```

* start an instance in specific project and zone, and create a DNS record

```bash
rcstate vm start \
  --name <instance_name> \
  --project <project_id> \
  --zone <zone_name> \
  --domain <dns_domain> \
  --dns-record-name <record_name> \
  --dns-record-type <record_type>
```

* start an instance and run shell commands AFTER the instance is started

```bash
rcstate vm start \
  --name <instance_name> \
  --project <project_id> \
  --zone <zone_name> \
  --script "echo TEST > test-file" \
  --external-ip \
  --ssh-key <path_to_key> \
  --ssh-port <port_number> \
  --ssh-user <username>
```

* stop an instance and run shell commands BEFORE the instance is stopped

```bash
rcstate vm stop \
  --name <instance_name> \
  --project <project_id> \
  --zone <zone_name> \
  --script "echo TEST > test-file" \
  --external-ip \
  --ssh-key <path_to_key> \
  --ssh-port <port_number> \
  --ssh-user <username>
```

* show status of an instance in specific project and zone

```bash
rcstate vm status \
  --name <instance_name> \
  --project <project_id> \
  --zone <zone_name>
```

* stop an instance in specific project and zone

```bash
rcstate vm stop \
  --name <instance_name> \
  --project <project_id> \
  --zone <zone_name>
```

<a name="server-mode"></a>
## Server mode

The application can be started in the server mode to serve API requests.

To start in server mode provide only the port number on which the server will listening:

`rcstate --server 8080`

The requests can be sent using the application in CLI mode by providing the host flag `--host`.

```bash
rcstate env up \
  --name <environment_name> \
  --env-file <environment_file> \
  --host <remote_host>:<port_number>
```

```bash
rcstate vm start \
  --name <instance_name> \
  --project <project_id> \
  --zone <zone_name> \
  --host <remote_host>:<port_number>
  ```

Also, the request can be sent as a POST API request to a specific endpoint.

To bring an environment in the UP state send the request to `v1/env/up` endpoint:

```bash
curl -X POST http://<remove_host>:<port_number>/v1/env/up -H "Content-Type: application/json" -d '{
  "all": true,
  "name": "",
  "label": "",
  "data": {
    "Envs": [
      <environment_details>
      ...
    ],
    "Vars": {
      "APP_NAME": "foo",
      ...
    }
  }
}'
```

To start a virtual machine send the request to `v1/vm/start` endpoint:

```bash
curl -X POST http://<remove_host>:<port_number>/v1/vm/start -H "Content-Type: application/json" -d '{
    "domain": "example.com",
    "name": <instance_name>,
    "project": <project_id>,
    "zone": <zone_name>
}'
```

The `--format` flag can be used to get the JSON representation of the request sent to the server.

The boolean flag `--dry` can be used to run the command without executing the logic.

```bash
rcstate env list \
  --name <environment_name> \
  --env-file <environment_file> \
  --host <remote_host>:<port_number> \
  --format json \
  --dry
```

```bash
rcstate vm start \
  --name <instance_name> \
  --project <project_id> \
  --zone <zone_name> \
  --host <remote_host>:<port_number> \
  --format json \
  --dry
  ```

List of endpoints for environment management:

* v1/env/down
* v1/env/show
* v1/env/up

List of endpoints for virtual machine management:

* v1/vm/list
* v1/vm/start
* v1/vm/status
* v1/vm/stop
