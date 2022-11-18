# Design Documents

This is a location to record all high-level design decisions in the Ojo
project.

A design document should provide:

- Context on the current state
- Proposed changes to achieve the goals
- Detailed reasoning
- Example scenarios
- Discussions of pros, cons, hazards and alternatives

[Template](./TEMPLATE.md)

Note the distinction between a design document and a spec below.

## Rules

The current process for design docs is:

- A design document is drafted and discussed in a dedicated pull request.
- A design document, once merged, should not be significantly modified.
- When a design document's decision is superseded, a reference to the new design should be added to its text.
- We do NOT require that all features have a design document.

Meanwhile the Readme file of each module should be a living document that is kept up to date. Design changes should update this Readme in the same PR as their implementation, and the Readme as a whole should serve as a reliable, complete source of truth (for example, for onboarding new engineers).
