package vtscfg

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Parse parses raw config bytes into a tree of Nodes. The returned root
// Node is synthetic (Name == "") and its Fields are the file's top-level
// blocks. Most callers should use Unmarshal instead; Parse is exposed for
// callers who want to walk or manipulate the tree directly.
func Parse(data []byte) (*Node, error) {
	root := &Node{}
	stack := []*Node{root}

	pendingName := ""
	pendingLine := 0

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)

	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimRight(scanner.Text(), "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if trimmed == "{" {
			if pendingName == "" {
				return nil, &SyntaxError{Line: lineNo, Msg: "unexpected '{' with no preceding block name"}
			}
			openBlock(&stack, pendingName)
			pendingName = ""
			continue
		}

		if trimmed == "}" {
			if pendingName != "" {
				return nil, &SyntaxError{Line: pendingLine, Msg: fmt.Sprintf("block %q was never opened with '{'", pendingName)}
			}
			if len(stack) <= 1 {
				return nil, &SyntaxError{Line: lineNo, Msg: "unexpected '}' with no matching '{'"}
			}
			stack = stack[:len(stack)-1]
			continue
		}

		// "key = value" - only if it's not actually an inline "Name {" block
		// opener that happens to contain a literal '=' further along (that
		// never occurs in real files, but we guard defensively by checking
		// the '=' comes before any trailing '{').
		if idx := strings.IndexByte(trimmed, '='); idx >= 0 && !strings.HasSuffix(trimmed, "{") {
			if pendingName != "" {
				return nil, &SyntaxError{Line: pendingLine, Msg: fmt.Sprintf("block %q was never opened with '{'", pendingName)}
			}
			key := strings.TrimSpace(trimmed[:idx])
			val := strings.TrimSpace(trimmed[idx+1:])
			top := stack[len(stack)-1]
			top.Fields = append(top.Fields, Field{Key: key, Value: &val})
			continue
		}

		// Inline "Name {" form, in case a file (or hand-written test fixture)
		// puts the brace on the same line as the block name.
		if strings.HasSuffix(trimmed, "{") {
			if pendingName != "" {
				return nil, &SyntaxError{Line: pendingLine, Msg: fmt.Sprintf("block %q was never opened with '{'", pendingName)}
			}
			name := strings.TrimSpace(strings.TrimSuffix(trimmed, "{"))
			openBlock(&stack, name)
			continue
		}

		// Otherwise this line is a bare block name; the '{' is expected on
		// one of the following lines.
		if pendingName != "" {
			return nil, &SyntaxError{Line: pendingLine, Msg: fmt.Sprintf("block %q was never opened with '{'", pendingName)}
		}
		pendingName = trimmed
		pendingLine = lineNo
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if pendingName != "" {
		return nil, &SyntaxError{Line: pendingLine, Msg: fmt.Sprintf("block %q was never opened with '{'", pendingName)}
	}
	if len(stack) != 1 {
		return nil, &SyntaxError{Line: lineNo, Msg: "unexpected end of input: missing closing '}'"}
	}
	return root, nil
}

func openBlock(stack *[]*Node, name string) {
	child := &Node{Name: name}
	top := (*stack)[len(*stack)-1]
	top.Fields = append(top.Fields, Field{Key: name, Child: child})
	*stack = append(*stack, child)
}

// Disclaimer: All code in this document has been taken from LLM outputs.
// TODO: Consider rewriting eventually.
