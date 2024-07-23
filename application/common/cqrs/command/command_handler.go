package command

type ICommandHandler[Command any, Result any] interface {
	Handle(command Command) (Result, error)
}
