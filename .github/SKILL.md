# Skill: learning-center

Purpose
- Provide a reproducible, workspace-scoped workflow for implementing the `learning-center` feature.
- Capture the step-by-step implementation pattern, decision points, and quality checks so agents and contributors follow a consistent process.

Scope
- Workspace-scoped: stores conventions and a checklist for the `learning-center` feature in this repo.
- Assumes API endpoints and data models will be provided later; this skill focuses on process and guardrails.

When to use
- Starting a new learning-center backend or frontend task.
- Adding or updating APIs, handlers, routes, components, or tests for learning-center.

Inputs
- API contract (to be added later by the developer).
- Any auth/permission requirements.
- Data model shapes (course, lesson, enrollment, progress).

Step-by-step workflow
1. Define API contract
   - List endpoints (e.g., GET /courses, GET /courses/:id, POST /enrollments).
   - Define request/response payloads and errors.
   - Decide pagination, filtering, and sorting semantics.
2. Model design and DB migration
   - Map API fields to DB models and types.
   - Add migrations and seed data for dev/testing.
3. Backend implementation
   - Implement models, repository/db layer, and business logic.
   - Add handlers (HTTP) and wire routes with appropriate authorization.
   - Add unit tests for services and handler-level integration tests.
4. Frontend implementation
   - Create typed API client types and hooks (use `lib/api.ts` pattern).
   - Build pages/components under `features/learning-center/` matching routes.
   - Add accessibility and responsiveness checks.
5. End-to-end and integration tests
   - Add integration tests covering major user flows (list, view, enroll, progress).
6. Docs and API examples
   - Add API examples for frontend and external integrators.
7. Review and release
   - Code review checklist, run linters/formatters, and confirm tests pass.

Decision points / branching
- Auth: anonymous read-only vs. authenticated enrollment flows.
- Pagination strategy: offset-based vs. cursor-based.
- Caching: server-side cache for course lists vs. client caching only.
- Media hosting: local vs. external CDN for lesson assets.

Quality criteria (done when)
- API contract exists and is documented.
- Unit and integration tests cover critical logic with CI passing.
- Linting and type checks pass.
- Frontend pages meet basic accessibility standards.
- API performance considerations (pagination, limits) are addressed.

Example prompts
- "Create API spec for learning-center course list with pagination and filters."
- "Generate Go handler skeletons and DB models for learning-center courses."
- "Add a React page under features/learning-center/pages to display course list."

Customization notes
- When APIs arrive, update the `Inputs` section with precise endpoint names and payloads.
- Add references to any shared design tokens, auth roles, or feature flags used by learning-center.

Contact
- Add the primary owner or team contact here when known.
