# Specification Rules

This document defines mandatory requirements for all specification files (`*.md`).

---

## 1. File Location Rules

- Each Go file MUST have a specification file with the same name and path.
- Example:
  - `usecases/create_user.go`
  - `usecases/create_user.md`
- Specifications MUST be colocated with code.

---

## 2. Spec Marker Block (mandatory)

Every specification file MUST start with this exact marker block (first lines of the file):

```md
<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=usecases/create_user -->
<!-- SPEC:KIND=usecase -->
<!-- SPEC:MENU=true -->
<!-- SPEC:END -->
```

Rules:
- `SPEC:FILE=true` is mandatory and identifies a file as a specification.
- `SPEC:ID` is mandatory, unique in the project, and must match the relative logical path.
- `SPEC:KIND` is mandatory. Allowed values: `controller`, `command`, `usecase`, `service`, `repository`, `other`.
- `SPEC:MENU` is mandatory. Allowed values: `true` or `false`.

---

## 3. Documentation Build and Menu Rules

- Only files containing `SPEC:FILE=true` MUST be included into documentation HTML build.
- Markdown files without this marker MUST be excluded from HTML spec visualization (for example `README.md`).
- Only files with `SPEC:MENU=true` MUST appear in the generated menu/table of contents.
- Files with `SPEC:MENU=false` MAY be rendered and reachable only by internal links.
- Recommended usage for menu:
  - include: `controller`, `command`
  - exclude from menu (but keep linkable): `service`, `usecase`, `repository`

---

## 4. Mandatory Section Template

After marker block, every specification MUST contain sections in this exact order:

1. `# Title`
2. `## Responsibility`
3. `## Inputs`
4. `## Outputs`
5. `## Business Logic`
6. `## Flow`
7. `## Links`
8. `## Dependencies`
9. `## Errors`
10. `## Notes` (optional)

If any mandatory section is missing, the file is invalid.

---

## 5. Business Logic Block Rules

Section `## Business Logic` is mandatory and MUST:
- contain numbered rules (`1.`, `2.`, ...);
- describe only domain rules;
- keep each rule testable and unambiguous;
- avoid transport/infrastructure details (HTTP, SQL, broker specifics).

Example:
1. Email must be unique.
2. Password length must be at least 8 characters.

---

## 6. Flow Block Rules

Section `## Flow` is mandatory and MUST:
- contain ordered steps (`1.`, `2.`, ...);
- describe behavioral sequence only;
- reference other specs only through the link format defined below.

---

## 7. Links Block Rules

Section `## Links` is mandatory and contains only explicit cross-spec links.

Each line MUST match this format:

```md
- <relation>: [<Name>](<relative/path/to/spec.md#anchor>)
```

Allowed `<relation>` values:
- `uses`
- `reads`
- `writes`
- `calls`
- `validates`

Examples:
- `- uses: [PasswordService](../services/password_service.md#hash)`
- `- writes: [UserRepository](../repositories/user_repository.md#create)`

Rules:
- only relative paths;
- only links to files with `SPEC:FILE=true`;
- anchor after `#` is mandatory;
- free-text references without markdown links are forbidden.

---

## 8. Dependencies Block Rules

Section `## Dependencies` is mandatory.
- Explicit list of dependent components.
- Every dependency MUST be a markdown link.
- Dependencies MUST be consistent with `## Links`.

---

## 9. Responsibility / Inputs / Outputs / Errors Rules

- `Responsibility`: describes WHAT component does; one responsibility only.
- `Inputs/Outputs`: describe logical contract, not transport-level payloads.
- `Errors`: domain-level errors only; no SQL/HTTP/internal framework errors.

---

## 10. Forbidden Practices

- Missing spec marker block.
- `SPEC:FILE` absent or not `true` in spec files.
- Including non-marked markdown files in HTML build.
- Showing non-`SPEC:MENU=true` files in menu.
- Describing SQL queries or HTTP handler internals in usecase/service specs.
- Referring to code line numbers.
- Implicit dependencies.

---

## 11. Change Management

Any change in behavior requires updating:
- `## Business Logic`
- `## Flow`
- `## Links`
- `## Dependencies` (if affected)

Specs MUST be updated before code.

---

## 12. Language Rules

- All specification content MUST be written in Russian.
- Specifications are business-facing documents.
- English is allowed ONLY for:
  - Go identifiers
  - File names
  - Link/code references

Correct:
- `Создаёт пользователя согласно бизнес-правилам.`
- `- uses: [PasswordService](../services/password_service.md#Hash)`

Incorrect:
- `Creates a user`
- `Validates input parameters`

If a specification is partially or fully in English, it MUST be rewritten in Russian.

---

## 13. Canonical Spec Skeleton

Use this skeleton for every new specification:

```md
<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=path/to/file_without_ext -->
<!-- SPEC:KIND=controller|command|usecase|service|repository|other -->
<!-- SPEC:MENU=true|false -->
<!-- SPEC:END -->

# <Название>

## Responsibility
<Краткое описание ответственности>

## Inputs
- ...

## Outputs
- ...

## Business Logic
1. ...
2. ...

## Flow
1. ...
2. ...

## Links
- uses: [Name](../path/spec.md#anchor)

## Dependencies
- [Name](../path/spec.md)

## Errors
- ...

## Notes
- ...
```

---

This document is authoritative.
If a spec violates these rules, it MUST be fixed.
