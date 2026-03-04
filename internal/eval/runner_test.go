package eval

import "testing"

func TestIsAllowedPath(t *testing.T) {
	allowed := []string{"internal/calc/**", "spec_changes/**"}
	if !isAllowedPath("internal/calc/calc.go", allowed) {
		t.Fatalf("expected path to be allowed")
	}
	if isAllowedPath("README.md", allowed) {
		t.Fatalf("expected path to be denied")
	}
}

func TestValidateSpecContent(t *testing.T) {
	good := `<!-- SPEC:START -->
<!-- SPEC:FILE=true -->
<!-- SPEC:ID=usecases/create_user -->
<!-- SPEC:KIND=usecase -->
<!-- SPEC:MENU=false -->
<!-- SPEC:END -->

# Title

## Responsibility
text

## Inputs
- a

## Outputs
- b

## Business Logic
1. one

## Flow
1. one

## Links
- uses: [X](../x.md#id)

## Dependencies
- [X](../x.md)

## Errors
- e
`
	if got := validateSpecContent(good); got != 0 {
		t.Fatalf("expected 0 violations, got %d", got)
	}
}
