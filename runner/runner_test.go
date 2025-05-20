package runner_test

import (
	"crypto/rand"
	"encoding/hex"
	"net/url"

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

const albumId = "foobar"

var anyCancelFunc = mock.AnythingOfType("context.CancelFunc")
var fakeWriter = nopWriter{}

var _ = Describe("runner", func() {
	var songs = getSongs(1)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var streamer *mocks.IStreamer

	BeforeEach(func() {
		discordClient = getDiscordClient()
		streamer = getStreamer(len(songs))
		subsonicClient = getSubsonicClient(songs)
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("will Init, Queue and Play without error", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient, streamer)
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
	var streamer *mocks.IStreamer

	BeforeEach(func() {
		discordClient = getDiscordClient()
		streamer = getStreamer(len(songs))
		subsonicClient = getSubsonicClient(songs)
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("will Init, Queue and Play without error with a playlist of more than 1 song", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient, streamer)
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
	discordClient.EXPECT().Init(mock.AnythingOfType("*runner.SubcordantRunner")).Return(nil)
	discordClient.EXPECT().JoinVoiceChat(anyCancelFunc).Return(fakeWriter, nil)
	return discordClient
}

func getStreamer(songCount int) *mocks.IStreamer {
	streamer := mocks.NewIStreamer(GinkgoT())
	for i := 0; i < songCount; i++ {
		streamer.EXPECT().PrepStream(
			mock.Anything, mock.AnythingOfType("*url.URL"), anyCancelFunc).Return(nil)
		streamer.EXPECT().Stream(fakeWriter, anyCancelFunc).Return(nil)
	}
	return streamer
}

func getSubsonicClient(songs []*subsonic.Child) *mocks.ISubsonicClient {
	subsonicClient := mocks.NewISubsonicClient(GinkgoT())
	subsonicClient.EXPECT().Init().Return(nil)
	subsonicClient.EXPECT().GetAlbum(albumId).Return(&subsonic.AlbumID3{
		Song: songs,
	}, nil)
	for _, song := range songs {
		subsonicClient.EXPECT().Stream(song.ID).Return(&url.URL{}, nil)
	}
	return subsonicClient
}

func getSongs(n uint) []*subsonic.Child {
	songs := make([]*subsonic.Child, int(n))
	for i := range songs {
		b := make([]byte, 8) // 8 random bytes → 16 hex characters
		if _, err := rand.Read(b); err != nil {
			// fallback to something deterministic, or handle error
			b = []byte("darnbeefbebefunk")
		}
		songs[i] = &subsonic.Child{
			ID: hex.EncodeToString(b),
		}
	}
	return songs
}
