// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package geninit

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// SelectVersion prompts user to select a version interactively
func SelectVersion(ctx context.Context, versions []string, modulePath string) (string, error) {
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions available for selection")
	}

	if len(versions) == 1 {
		mlog.Printf("Only one version available: %s", versions[0])
		return versions[0], nil
	}

	// Display available versions
	fmt.Printf("\nAvailable versions for %s:\n", modulePath)
	fmt.Println(strings.Repeat("-", 40))

	// Show versions with index (newest first)
	maxDisplay := 20 // Limit display to avoid overwhelming output
	displayCount := len(versions)
	if displayCount > maxDisplay {
		displayCount = maxDisplay
	}

	for i := 0; i < displayCount; i++ {
		marker := ""
		if i == 0 {
			marker = " (latest)"
		}
		fmt.Printf("  [%2d] %s%s\n", i+1, versions[i], marker)
	}

	if len(versions) > maxDisplay {
		fmt.Printf("  ... and %d more versions\n", len(versions)-maxDisplay)
	}

	fmt.Println(strings.Repeat("-", 40))

	// Prompt for selection
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("Select version [1-%d] (default: 1 for latest): ", displayCount)

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		// Default to latest
		if input == "" {
			fmt.Printf("Selected: %s (latest)\n", versions[0])
			return versions[0], nil
		}

		// Parse selection
		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > displayCount {
			fmt.Printf("Invalid selection. Please enter a number between 1 and %d.\n", displayCount)
			continue
		}

		selected := versions[idx-1]
		fmt.Printf("Selected: %s\n", selected)
		return selected, nil
	}
}
