# PHP Specification Rules

This document defines mandatory requirements for PHP specification files (`*.md`).

---

## 1. File Location Rules

- Each PHP file MUST have a specification file with the same name and path.
- Example:
  - `app/Services/PaymentService.php`
  - `app/Services/PaymentService.md`
- Specifications MUST be colocated with code.

---

## 2. Spec Marker Block (mandatory)

Every specification file MUST start with this marker block:

```md
<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=app/Services/PaymentService -->
<!-- SPEC:KIND=controller -->
<!-- SPEC:MENU=true -->
<!-- SPEC:END -->
```

Rules:
- `SPEC:FILE=true` is mandatory.
- `SPEC:ID` is mandatory and unique.
- `SPEC:KIND` allowed values: `controller`, `command`, `usecase`, `service`, `repository`, `job`, `listener`, `other`.
- `SPEC:MENU` allowed values: `true` or `false`.

---

## 3. Documentation Build Rules

- HTML build includes only files with `SPEC:FILE=true`.
- Menu includes only files with `SPEC:MENU=true`.
- Non-marked markdown files MUST be excluded.

---

## 4. Mandatory Sections

Each specification MUST contain, in order:
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

---

## 5. Links Rules

Link format:

```md
- <relation>: [<Name>](<relative/path/to/spec.md#anchor>)
```

Allowed relations: `uses`, `reads`, `writes`, `calls`, `validates`.

---

## 6. Change Management

Behavior changes require updates in:
- `## Business Logic`
- `## Flow`
- `## Links`
- `## Dependencies`

Specs MUST be updated before code.

---

## 7. Language Rules

- Specification text MUST be in Russian.
- English is allowed only for identifiers, file names, links.
