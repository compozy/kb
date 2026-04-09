package greeter

import "fmt"

// Hello formats a deterministic greeting for the fixture program.
func Hello(name string) string {
	return fmt.Sprintf("hello, %s", name)
}
