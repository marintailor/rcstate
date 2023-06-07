# rcstate

## Overview

rcstate is a CLI app written in Go to manage the state of resources in Google Cloud.

## Installation

```bash
go install github.com/marintailor/rcstate@latest
```

## Usage

### Manage environments

An environment represents one or more groups of resources that are already present in Google Cloud.

The environments are declared in a YAML file, and it is provided as an option with flag `--env-file`.

```bash
rcstate env show --all --env-file env-dev.yaml
```

**Examples:**

```bash
# show an environment
rcstate env show \
  --name <environment_name> \
  --env-file <environment_file>

# show all environments
rcstate env show \
  --all \
  --env-file <environment_file>

# show an environments with label "api"
rcstate env show \
  --name <environment_name> \
  --label api \
  --env-file <environment_file>

# show all environments with label "api"
rcstate env show \
  --all \
  --label api \
  --env-file <environment_file>

# change state of an environment
rcstate env up/down \
  --name <environment_name> \
  --env-file <environment_file>

#  change state of an environment labeled with "dev" and "qa"
rcstate env up/down \
  --name <environment_name> \
  --label qa,api \
  --env-file <environment_file>

# change state of all environments
rcstate env up/down \
  --all \
  --env-file <environment_file>

# change state of all environments labeled with "dev" and "qa"
rcstate env up/down \
  --all \
  --label qa,api \
  --env-file <environment_file>
```

**Schema example of the environment YAML file:**

```yaml
provider: gcp    # Acronym of the Public Cloud provider
variables:    # Global variables that are accessible from all environments
  APP_NAME:  test-app
  DOMAIN:  test-domain.com
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
              up:    # Shell commands to be executed after instance is started
                - sudo shutdown -h +30
              down:    # Shell commands to be executed before instance is stopped
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
                    - wget -O - https://g{{ .APP_NAME }}.{{ .DOMAIN }}/apps/{{ .APP_NAME }}/-/raw/main/init.sh | bash
```

### Google Cloud Engine virtual machine

```bash
# list all virtual machine instances in specific project and zone
rcstate vm list \
  --project <project_name> \
  --zone <zone_name>

# start an instance in specific project and zone
rcstate vm start \
  --name <instance_name> \
  --project <project_name> \
  --zone <zone_name>

# start an instance in specific project and zone, and create a DNS record
rcstate vm start \
  --name <instance_name> \
  --project <project_name> \
  --zone <zone_name> \
  --domain <dns_domain> \
  --dns-record-name <record_name> \
  --dns-record-type <record_type>

# start an instance and run shell commands AFTER the instance is started
rcstate vm start \
  --name <instance_name> \
  --project <project_name> \
  --zone <zone_name> \
  --script "echo TEST > test-file" \
  --ip <ip_addr> \
  --ssh-key <path_to_key> \
  --ssh-port <port_number> \
  --ssh-user <username>

# stop an instance and run shell commands BEFORE the instance is stopped
rcstate vm stop \
  --name <instance_name> \
  --project <project_name> \
  --zone <zone_name> \
  --script "echo TEST > test-file" \
  --ip <ip_addr> \
  --ssh-key <path_to_key> \
  --ssh-port <port_number> \
  --ssh-user <username>

# show status of an instance in specific project and zone
rcstate vm status \
  --name <instance_name> \
  --project <project_name> \
  --zone <zone_name>

# stop an instance in specific project and zone
rcstate vm stop \
  --name <instance_name> \
  --project <project_name> \
  --zone <zone_name>
```
