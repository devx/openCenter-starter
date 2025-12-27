# Product Requirements Document (PRD)
## Project: openCenter-base
## Domain: Enterprise Kubernetes platform delivery and operations

## Executive Summary
openCenter-base is a high-performance web platform for bootstrapping and operating production-grade Kubernetes clusters on private and hybrid infrastructure (OpenStack/VMware). It delivers an opinionated GitOps model and a curated baseline of enterprise services (identity, RBAC, ingress, certificates, observability, backup/restore, operators). The frontend uses SolidStart for reactive SSR/ISR, while the backend leverages Go + Fiber for low-latency APIs. AuthN is implemented natively using SAML and OIDC to maintain control over identity flows; AuthZ uses Cedar + AVP for fine-grained access policies.

## Goals & Success Criteria
- Standardize cluster delivery and day-2 operations across environments.
- Reduce time-to-production for Kubernetes clusters.
- Provide reliable, auditable access control and change tracking.
- Achieve top-tier performance for UI responsiveness and API latency.

## Personas
- Platform Owner / Platform Engineering Lead
- Cluster Admin / Kubernetes Ops
- SRE / Operations Engineer
- Security / IAM / Compliance Engineer
- Observability Engineer
- Infrastructure / Cloud Engineer
- Application Team Lead / Developer
- Read-only / Auditor

## System Architecture Strategy

### Frontend Architecture (SolidStart + TanStack Query + UnoCSS)
- SolidStart SSR/ISR routes prefetch critical data on the server with TanStack Query, then hydrate on the client for instant interactivity.
- Use Solid signals for local UI state only (`createSignal`, `createStore`).
- Directory conventions:
  - `src/routes/` for page routes
  - `src/components/` for UI components
  - `src/lib/api/` for API client and query hooks
  - `uno.config.ts` for UnoCSS rules, presets, and theme tokens

### Backend Architecture (Go + Fiber + Hexagonal)
- Ports & Adapters pattern:
  - Core domain: use cases (cluster bootstrap, upgrades, drift detection, RBAC mapping)
  - Ports: interfaces for storage, identity, policy engine, external GitOps tools
  - Adapters: HTTP handlers (Fiber), PostgreSQL via pgx/sqlc, Cedar + AVP, SAML/OIDC providers
- Keep HTTP, DB, and policy engines as replaceable adapters to preserve testability and performance.

### Auth Flow (SAML)
Sequence description:
1. User initiates login at `/auth/saml/login`.
2. Backend generates SAML AuthnRequest using `crewjam/saml`.
3. User is redirected to the IdP (e.g., Entra ID, Keycloak).
4. IdP authenticates and posts SAML Response to `/auth/saml/acs`.
5. Backend validates assertion, extracts attributes/groups.
6. Backend maps attributes to internal identity and Cedar + AVP policy context.
7. Backend issues session/token and redirects to frontend.

## Functional Requirements

### AuthN/AuthZ
- SAML SP (crewjam/saml) with enterprise IdPs (Keycloak, Entra ID).
- OIDC (coreos/go-oidc) for token validation and federated sign-in.
- Group-to-role mapping with external IdP claims.
- ReBAC/PBAC via Cedar + AVP:
  - Example relationship: `user:alice` is `editor` of `cluster:prod-1`.
  - Support read-only auditor role for non-mutating access.
- Audit trail for authentication events and access changes.

### Core Platform Features
- Cluster bootstrap workflows:
  - Standardized base cluster definition plus overlays.
  - GitOps reconciliation integration.
- Day-2 operations:
  - Upgrade orchestration (control plane and workers).
  - Drift detection and remediation.
  - Break/fix workflows with approvals.
- Enterprise service baseline:
  - Identity and RBAC
  - Ingress/Gateway
  - Certificates
  - Observability stack
  - Backup/restore
  - Operator lifecycle management
- Environment support:
  - OpenStack, VMware, and cloud VM targets
  - Infrastructure requirements validation

## Data Strategy (PostgreSQL + pgx + sqlc)

### Schema Approach
- Use `sqlc` for type-safe queries and strongly typed models.
- Domain tables (examples):
  - `clusters`, `cluster_versions`, `cluster_events`
  - `gitops_repos`, `reconciliation_runs`
  - `identities`, `idp_mappings`
  - `policies`, `relationship_tuples`
  - `audit_logs`

### Indexing Strategy
- B-Tree for primary and lookup keys (e.g., `cluster_id`, `org_id`).
- GIN for JSONB metadata fields (e.g., cluster config, reconciliation status).
- Partial indexes for status filters (e.g., `WHERE status = 'failed'`).

## Frontend/Backend Interface

### API Contract (REST/JSON)
- Example endpoints:
  - `GET /clusters`
  - `POST /clusters`
  - `POST /clusters/:id/upgrade`
  - `GET /audit/logs`
- All APIs return structured JSON with a consistent error schema.

### Cache Invalidation (TanStack Query)
- Queries keyed by domain entity (`['clusters']`, `['cluster', id]`).
- Mutations call `queryClient.invalidateQueries(...)` on affected keys.
- Background refetch for status-heavy resources (e.g., upgrade progress).

## Non-Functional Requirements

### Performance
- Frontend: Lighthouse >= 90 (Performance, Best Practices).
- Backend API: p50 < 50ms, p95 < 200ms for core endpoints.

### Security & Compliance
- Secure session handling, TLS enforced end-to-end.
- Audit logs for auth and mutation actions.
- Principle of least privilege via ReBAC/PBAC.
- Configurable retention for compliance.

## Open Decisions
- Decide baseline GitOps tool integration (Argo CD vs Flux).
