# Specification Rules

This document defines the mandatory structure and semantics of all specification files (*.md).

---

## 1. File Location Rules

- Each Go file MUST have a specification file with the same name and path.
- Example:
    - usecases/create_user.go
    - usecases/create_user.md

- Specifications MUST be colocated with code.

---

## 2. Mandatory Sections

Every specification MUST contain the following sections in this order:

1. # Title
2. ## Responsibility
3. ## Inputs
4. ## Outputs
5. ## Business Rules
6. ## Flow
7. ## Dependencies
8. ## Errors
9. ## Notes (optional)

---

## 3. Responsibility

- Describes WHAT the component does.
- Must NOT describe implementation details.
- One responsibility only.

Example:
"Creates a new user according to business rules."

---

## 4. Inputs / Outputs

- Describe logical inputs and outputs.
- Do NOT describe transport-level details (HTTP, JSON).

---

## 5. Business Rules

- Enumerated list.
- Each rule must be testable.
- No technical language.

Example:
1. Email must be unique.
2. Password length must be at least 8 characters.

---

## 6. Flow

- Ordered list of steps.
- Each step may reference another spec.

Example:
1. Validate input
2. Check email uniqueness  
   → uses: ../repositories/user_repository.md#existsbyemail
3. Hash password  
   → uses: ../services/password_service.md#hash

---

## 7. Dependencies

- Explicit list of all referenced components.
- Each dependency MUST be a markdown link.

Example:
- [UserRepository](../repositories/user_repository.md)
- [PasswordService](../services/password_service.md)

---

## 8. Errors

- Enumerate all domain-level errors.
- No technical errors (SQL, HTTP, etc).

---

## 9. Links Semantics

All links MUST be explicit and meaningful.

Allowed link annotations:
- uses
- reads
- writes
- calls
- validates

Example:
→ writes: ../repositories/user_repository.md#create

---

## 10. Forbidden Practices

- Describing SQL queries
- Describing HTTP handlers inside usecases
- Referring to code line numbers
- Implicit dependencies

---

## 11. Change Management

- Any change in behavior requires updating:
    - Business Rules
    - Flow
    - Dependencies (if applicable)

- Specs MUST be updated before code.

---

This document is authoritative.
If a spec violates these rules, it MUST be fixed.

---

## 12. Language Rules

- All specification content MUST be written in Russian.
- Specifications are business-facing documents.
- English is allowed ONLY for:
    - Go identifiers
    - File names
    - Code references inside links

Correct:
- "Создаёт пользователя согласно бизнес-правилам."
- "→ uses: ../services/password_service.md#Hash"

Incorrect:
- "Creates a user"
- "Validates input parameters"

If a specification is written partially or fully in English, it MUST be rewritten in Russian.
