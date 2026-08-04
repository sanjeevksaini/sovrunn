---
doc_type: decision_record
title: CloudEnrollment Joins Customer to Platform
status: Accepted
phase: 2R
ai_load_priority: important
---

# CloudEnrollment Joins Customer to Platform

## Status

Accepted

## Context

Customer relationship to a cloud platform and provider participation in that platform are conflated when one generic Provider owns both.

## Decision

CloudEnrollment joins a customer Organization to a CloudPlatform. Provider participation remains separate through CloudProviderParticipation.

## Consequences

Enrollment and participation have distinct endpoints, ownership, and lifecycle. Customer does not need to know provider identity. Provider does not own customer enrollment.

## Migration

Existing provider-customer links map to CloudEnrollment where the customer relationship exists.

## Conformance Evidence

Drift checks verify CloudEnrollment targets CloudPlatform and that CloudProviderParticipation is a separate boundary.

## Source ADH Reference

ADH-2026-021
