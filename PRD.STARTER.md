# Product Requirements Document (PRD) - Starter Template
## Project: openCenter-base (baseline)
## Purpose: starting point for any openCenter webapp and backend

## Executive Summary
This document defines the baseline requirements for any openCenter web application and backend. It prescribes the high-performance, reactive stack and the architectural contracts that all openCenter projects must follow. Product-specific features and domain rules should be added as extensions to this baseline.

## Baseline Tech Stack (Non-Negotiable)
- Frontend: SolidJS with SolidStart (SSR/ISR)
- State management: Solid signals only (`createSignal`, `createStore`)
- Data fetching: TanStack Query (Solid Query)
- Styling: UnoCSS
- Backend: Go with Fiber (or high-performance stdlib routing)
- AuthN: Native SAML (crewjam/saml) and OIDC (coreos/go-oidc)
- AuthZ: Cedar + AVP for ReBAC/PBAC
- Database: PostgreSQL with pgx
- Data access: sqlc or raw pgx; no ORM

## System Architecture Strategy

### Frontend Architecture
- SolidStart routes prefetch data on the server using TanStack Query.
- Hydrate on the client for low-latency interactivity.
- Keep UI state in Solid signals; reserve TanStack Query for async server state.
- Required directories:
  - `src/routes/`
  - `src/components/`
  - `src/lib/api/`
  - `uno.config.ts`

### Backend Architecture
- Hexagonal architecture (Ports & Adapters) with:
  - Domain use cases in the core
  - Interfaces for storage, policy, identity, and GitOps systems
  - Adapters for HTTP (Fiber), PostgreSQL (pgx/sqlc), and AuthN/AuthZ providers

### Auth Flow (SAML)
1. `/auth/saml/login` issues AuthnRequest.
2. IdP authenticates user.
3. `/auth/saml/acs` consumes SAML Response.
4. Attributes map to identity and Cedar + AVP policy context.
5. Session/token issued for UI access.

## Baseline Functional Requirements
- AuthN with SAML and OIDC, including IdP group mapping.
- AuthZ with Cedar + AVP and support for read-only roles.
- Audit logging for authentication and authorization changes.
- API contracts must be REST/JSON with consistent error schemas.

## Data Strategy (PostgreSQL)
- Use sqlc where possible for compile-time SQL safety.
- Use B-Tree indexes for identifiers and foreign keys.
- Use GIN indexes for JSONB configuration fields.
- Prefer explicit schemas per bounded context where applicable.

## Non-Functional Requirements
- Frontend Lighthouse >= 90 for performance and best practices.
- Backend API p50 < 50ms, p95 < 200ms for core endpoints.
- TLS everywhere, least privilege access, and auditability as defaults.

## Extension Points (Project-Specific)
- Domain-specific workflows and UI modules.
- Integration targets (GitOps, infra providers, observability stacks).
- Data models and API endpoints beyond the baseline.
- Operational runbooks and SLOs.
