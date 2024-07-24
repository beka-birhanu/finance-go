/*
Package command provides a generic interface for handling commands.

It includes the `ICommandHandler` interface for processing commands of any type and returning results of any type.
*/
package command

// ICommandHandler defines a generic interface for handling commands.
//
// Type Parameters:
// - Command: The type of the command to be handled.
// - Result: The type of the result returned after handling the command.
type ICommandHandler[Command any, Result any] interface {

	// Handle processes the provided command and returns the result or an error.
	Handle(command Command) (Result, error)
}
