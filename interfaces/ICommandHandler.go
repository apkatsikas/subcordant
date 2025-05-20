package interfaces

type ICommandHandler interface {
	Play(albumId string) error
}
