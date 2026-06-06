# koicha

**koicha** is a Kafka TUI written in Go with [Bubble Tea](https://charm.land/bubbletea).

The project is currently in early development. Most features are still being designed, prototyped, or implemented.

## What is koicha?

`koicha` is a terminal UI for working with Apache Kafka.

The goal of the project is to make routine Kafka tasks simpler, faster, and more pleasant for everyone who works with Kafka.

Kafka tooling can often feel scattered, noisy, or overloaded.  
`koicha` is meant to be the opposite: simple, focused, and friendly to use.

Working with Kafka is often a daily routine - and `koicha` aims to make that routine feel a little more like tea time: familiar, calm, and enjoyable.

## Why the name?

**Koicha** is a thick and concentrated style of matcha tea.

The name reflects the idea behind the project: a concentrated Kafka experience inside the terminal.

A dense little cup of Kafka.

## Bento

One of the core ideas behind `koicha` is the **bento**.

A bento is a portable set of Kafka connection and workspace configuration that can be imported, exported, and shared.

Bento profiles are designed to make it easy to:

- keep separate profiles for different environments, such as `prod`, `test`, and `dev`
- switch between Kafka clusters quickly
- share ready-to-use configurations with a team
- limit the visible workspace to the topics and consumer groups that matter for a specific team or project

For example, a development team may own only a specific set of topics and consumer groups. In that case, a bento can define a focused workspace containing only those resources.

This helps reduce noise and makes Kafka easier to use for people who don't need to see the entire cluster.

When needed, users should still be able to leave the bento scope and inspect other resources available on the same bootstrap server.

## Planned features

`koicha` is still in the design and prototyping stage. All listed features are currently **WIP**.

| Feature | Status |
|---|---|
| Bento profiles | WIP |
| Broker overview | WIP |
| Topic listing and metadata | WIP |
| Partition and offset inspection | WIP |
| Consumer group listing and metadata | WIP |
| Convenient offset reset flows | WIP |
| Reading messages from topics | WIP |
| Producing messages to topics | WIP |
| Quick switching between environments and profiles | WIP |
| Common Kafka authentication methods | WIP |

## Authentication

Must have:
- PLAINTEXT
- SASL_SSL with username/password
- SASL_SSL with SCRAM
- SSL/TLS with optional CA certificate

Later:
- mTLS with client certificate
- SASL_SSL + client certificate
- OAuth / Kerberos

## Platform support

The goal is to make `koicha` available on major platforms:

- Linux
- macOS
- Windows

## Status

🚧 **Work in progress**

`koicha` is not production-ready yet.

The project is currently mostly in the design and prototyping phase. APIs, UI, configuration format, behavior, and feature scope may change at any time.

## Tech stack

- Go
- Bubble Tea v2 (`charm.land/bubbletea/v2`)
- Bubbles v2 (`charm.land/bubbles/v2`)
- Lip Gloss v2 (`charm.land/lipgloss/v2`)
- franz-go
