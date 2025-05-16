package runner_test

import (
	"io"

	"github.com/apkatsikas/subcordant/interfaces/mocks"
	"github.com/apkatsikas/subcordant/runner"
	"github.com/delucks/go-subsonic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

type nopReadCloser struct{}

func (nopReadCloser) Read(p []byte) (int, error) {
	// Return EOF immediately, simulating an empty reader
	return 0, io.EOF
}

func (nopReadCloser) Close() error {
	// No-op for Close
	return nil
}

var _ = Describe("runner", func() {
	const albumId = "foobar"

	var songs = []*subsonic.Child{
		{
			ID: "bloop",
		},
		{
			ID: "scoop",
		},
	}

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var ffmpegCommander *mocks.IFfmpegCommander
	var fakeWriter nopWriter
	var fakeReadCloser nopReadCloser
	var initError error

	var cancelFunc = mock.AnythingOfType("context.CancelFunc")

	BeforeEach(func() {
		discordClient = mocks.NewIDiscordClient(GinkgoT())
		discordClient.EXPECT().Init(mock.AnythingOfType("*runner.SubcordantRunner")).Return(nil)
		discordClient.EXPECT().JoinVoiceChat(cancelFunc).Return(fakeWriter, nil)

		ffmpegCommander = mocks.NewIFfmpegCommander(GinkgoT())
		ffmpegCommander.EXPECT().Start(
			mock.Anything, fakeReadCloser, mock.AnythingOfType("string"), cancelFunc).Return(nil)
		ffmpegCommander.EXPECT().Stream(fakeWriter, cancelFunc).Return(nil)

		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(nil)
		subsonicClient.EXPECT().GetAlbum(albumId).Return(&subsonic.AlbumID3{
			Song: songs,
		}, nil)
		subsonicClient.EXPECT().Stream(songs[0].ID).Return(fakeReadCloser, nil)
		subcordantRunner = &runner.SubcordantRunner{}

		initError = subcordantRunner.Init(subsonicClient, discordClient, ffmpegCommander)
	})

	It("will Init and HandlePlay without error", func() {
		Expect(initError).NotTo(HaveOccurred())
		err := subcordantRunner.HandlePlay(albumId)
		Expect(err).NotTo(HaveOccurred())
	})
})
