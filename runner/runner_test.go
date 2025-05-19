package runner_test

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/interfaces/mocks"
	"github.com/apkatsikas/subcordant/runner"
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

const albumId = "foobar"

var anyCancelFunc = mock.AnythingOfType("context.CancelFunc")
var anySubcordantRunner = mock.AnythingOfType("*runner.SubcordantRunner")
var anyString = mock.AnythingOfType("string")

var fakeWriter = nopWriter{}
var fakeReadCloser nopReadCloser = nopReadCloser{}

var _ = Describe("runner", func() {
	var songs = getSongs(1)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var execCommander *mocks.IExecCommander

	BeforeEach(func() {
		discordClient = getDiscordClient()
		execCommander = getExecCommander(len(songs))
		subsonicClient = getSubsonicClient(songs)
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("will Init, Queue and Play without error", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient, execCommander)
		Expect(err).NotTo(HaveOccurred())

		err = subcordantRunner.Queue(albumId)
		Expect(err).NotTo(HaveOccurred())

		err = subcordantRunner.Play()
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("runner", func() {
	var songs = getSongs(2)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var execCommander *mocks.IExecCommander

	BeforeEach(func() {
		discordClient = getDiscordClient()
		execCommander = getExecCommander(len(songs))
		subsonicClient = getSubsonicClient(songs)
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("will Init, Queue and Play without error with a playlist of more than 1 song", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient, execCommander)
		Expect(err).NotTo(HaveOccurred())

		err = subcordantRunner.Queue(albumId)
		Expect(err).NotTo(HaveOccurred())

		err = subcordantRunner.Play()
		Expect(err).NotTo(HaveOccurred())
	})
})

// TODO - should i add ginkgo helper to these functions below?
func getDiscordClient() *mocks.IDiscordClient {
	discordClient := mocks.NewIDiscordClient(GinkgoT())
	discordClient.EXPECT().Init(anySubcordantRunner).Return(nil)
	discordClient.EXPECT().JoinVoiceChat(anyCancelFunc).Return(fakeWriter, nil)
	return discordClient
}

func getExecCommander(songCount int) *mocks.IExecCommander {
	execCommander := mocks.NewIExecCommander(GinkgoT())
	for i := 0; i < songCount; i++ {
		execCommander.EXPECT().Start(
			mock.Anything, fakeReadCloser, anyString, anyCancelFunc).Return(nil)
		execCommander.EXPECT().Stream(fakeWriter, anyCancelFunc).Return(nil)
	}
	return execCommander
}

func getSubsonicClient(songs []*subsonic.Child) *mocks.ISubsonicClient {
	subsonicClient := mocks.NewISubsonicClient(GinkgoT())
	subsonicClient.EXPECT().Init().Return(nil)
	subsonicClient.EXPECT().GetAlbum(albumId).Return(&subsonic.AlbumID3{
		Song: songs,
	}, nil)
	for _, song := range songs {
		subsonicClient.EXPECT().Stream(song.ID).Return(fakeReadCloser, nil)
	}
	return subsonicClient
}

func getSongs(n uint) []*subsonic.Child {
	songs := make([]*subsonic.Child, int(n))
	for i := range songs {
		b := make([]byte, 8) // 8 random bytes → 16 hex characters
		if _, err := rand.Read(b); err != nil {
			// fallback to something deterministic, or handle error
			b = []byte("deadbeefcafebabe")
		}
		songs[i] = &subsonic.Child{
			ID: hex.EncodeToString(b),
		}
	}
	return songs
}
