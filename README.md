# rcstate

## Overview

rcstate is a CLI app written in Go to manage the state of resources in the Google Cloud.

## Installation

```bash
go install github.com/marintailor/rcstate@latest
```

## Usage

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

# show status of an instance in specific project and zone
rcstate vm status \
  --name <instance_name> \
  --project <project_name> \
  --zone <zone_name>
```
