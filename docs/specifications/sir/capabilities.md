# Capabilities

Document
- ID: sir-capabilities
- Version: 1.0
- Status: Stable

Purpose
- Define canonical capability catalog
- Define capability taxonomy
- Define capability identifiers

Rules

MUST
- Use canonical identifiers
- Preserve identifier stability
- Preserve backward compatibility

MUST NOT
- Use display names as identifiers
- Rename published identifiers
- Duplicate capabilities

Capability Taxonomy

Data

- data.read
- data.create
- data.update
- data.delete
- data.merge

Query

- query.scan
- query.filter
- query.project
- query.aggregate
- query.sort
- query.join

Transaction

- transaction.atomic
- transaction.consistent
- transaction.isolated
- transaction.durable
- transaction.savepoint

Analytics

- analytics.window
- analytics.statistics
- analytics.timeseries

Search

- search.full-text
- search.vector
- search.graph

Streaming

- stream.publish
- stream.subscribe
- stream.consume

Object Storage

- object.read
- object.write
- object.delete
- object.multipart-upload

Security

- security.authentication
- security.authorization
- security.encryption
- security.audit

Capability Definition

Every Capability

MUST define

- Identifier
- Category
- Name

MAY define

- Description
- Version
- Aliases

References

- capability-model.md
