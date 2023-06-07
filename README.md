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
