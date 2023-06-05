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
```
