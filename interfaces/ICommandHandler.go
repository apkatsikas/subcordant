package interfaces

type ICommandHandler interface {
	Queue(albumId string) error
	Play() error
	IsPlaying() bool
}
