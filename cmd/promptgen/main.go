package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SmirnovND/spec-agent/internal/promptgen"
)

func main() {
	source := flag.String("source", "internal/fs/assets/shared/prompt_specs", "директория с prompt.yaml")
	check := flag.Bool("check", false, "только проверить актуальность сгенерированных markdown")
	flag.Parse()

	if err := promptgen.GenerateAll(*source, *check); err != nil {
		fmt.Fprintf(os.Stderr, "promptgen error: %v\n", err)
		os.Exit(1)
	}

	if *check {
		fmt.Println("prompt markdown is up-to-date")
		return
	}
	fmt.Println("prompt markdown generated")
}
