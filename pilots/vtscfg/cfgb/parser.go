package cfgb

import (
	"bufio"
	"fmt"
	"strings"
)

type Node struct {
	Name     string
	Values   map[string]string
	Children []*Node
}

func Parse(text string) (*Node, error) {
	scanner := bufio.NewScanner(strings.NewReader(text))

	root := &Node{
		Values: make(map[string]string),
	}

	stack := []*Node{root}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		// Closing brace.
		if line == "}" {
			if len(stack) <= 1 {
				return nil, fmt.Errorf("unexpected }")
			}

			stack = stack[:len(stack)-1]
			continue
		}

		// Opening a block:
		//
		// BAN
		// {
		//
		// or:
		//
		// USER
		// {
		if line != "{" && scanner.Scan() {
			next := strings.TrimSpace(scanner.Text())

			if next == "{" {
				node := &Node{
					Name:   line,
					Values: make(map[string]string),
				}

				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, node)

				stack = append(stack, node)
				continue
			}

			// It wasn't an opening brace, so treat the line
			// we consumed as a key/value line too.
			line += "\n" + next
		}

		// key = value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			stack[len(stack)-1].Values[key] = value
			continue
		}

		return nil, fmt.Errorf("cannot parse line: %q", line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return root, nil
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
