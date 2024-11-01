# gear

An infrastructure automation tool inspired by JetPorch, using YAML-based playbooks for task definition and a compiled Go execution engine for running automation over SSH.

## Overview

`gear` is an experiment in replacing Ansible-style configuration management with a compiled Go CLI. Tasks are defined in YAML — similar to Ansible playbooks — and executed against remote hosts over SSH. The difference is that the execution engine itself is a single, dependency-free Go binary rather than a Python runtime.

> **Heavily inspired by [JetPorch](https://github.com/jhawkesworth/jetporch_docs)** — a general-purpose IT automation platform designed by Michael DeHaan, the original creator of Ansible and Cobbler. Gear borrows many of JetPorch's ideas around language simplicity and agentless SSH-based execution, reimplementing them in Go.

**Why gear over Ansible?**

- **No Python runtime dependency** on target hosts — all orchestration runs from the `gear` binary
- **YAML for tasks, Go for execution** — familiar declarative task syntax backed by a compiled engine
- **Type safety** — catch configuration errors at parse time, not halfway through a playbook run
- **Single binary** — easy to distribute and version alongside your infrastructure code

## Project Status

Early-stage / experimental. Current version: `0.0.1`.

## Repository Structure

```
gear/
├── cmd/
│   └── gear/          # CLI entrypoint
├── internal/
│   └── cli/           # CLI commands and version metadata
├── examples/          # Example configuration files
├── vagrant/           # Vagrant environment for local testing
├── Makefile           # Build tasks
├── go.mod
└── go.sum
```

## Prerequisites

- Go 1.22.3 or later
- SSH access to target hosts (key-based auth recommended)

## Building

Generate the version file and build:

```bash
make gensrc
go build -o gear ./cmd/gear
```

Or install directly:

```bash
go install github.com/smahm006/gear/cmd/gear@latest
```

## Usage

Define your tasks in a YAML file:

```yaml
# example.yaml
hosts:
  - user: ubuntu
    address: 192.168.0.10

tasks:
  - name: Ensure /tmp/hello exists
    # ... task definition
```

Then run:

```bash
gear apply -f example.yaml
```

See the [`examples/`](./examples) directory for working configuration samples.

## Local Development

A Vagrant environment is included for testing against a local VM:

```bash
cd vagrant
vagrant up
```

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/pkg/sftp` | SFTP file transfer to remote hosts |
| `golang.org/x/crypto` | SSH client authentication |
| `gopkg.in/yaml.v3` | Parsing task configuration files |

## Roadmap

- [ ] Core task primitives (file, command, template, package)
- [ ] Idempotent execution (check mode)
- [ ] Inventory file support
- [ ] Parallel host execution
- [ ] Role/module system analogous to Ansible roles

## License

MIT
