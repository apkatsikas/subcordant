package interfaces

type IDiscordClient interface {
	Init(commandHandler ICommandHandler) error
}
