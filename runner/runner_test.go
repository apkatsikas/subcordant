package runner_test

import (
	"io"
	"net/url"
	"sync"
	"time"

	"github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/interfaces/mocks"
	"github.com/apkatsikas/subcordant/runner"
	"github.com/apkatsikas/subcordant/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

const albumId = "foobar"

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

var fakeWriter = nopWriter{}

var _ = DescribeTableSubtree("runner",
	func(songCount int) {
		var songs = getSongs(songCount)

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

		It("should Init and Play without error", func() {
			err := subcordantRunner.Init(subsonicClient, discordClient, streamer)
			Expect(err).NotTo(HaveOccurred())

			state, err := subcordantRunner.Play(albumId)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.PlaybackComplete))
		})
	},
	Entry("1 song", 1),
	Entry("2 songs", 2),
)

var _ = Describe("runner", func() {
	var songs = getSongs(2)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var streamer *mocks.IStreamer

	BeforeEach(func() {
		discordClient = getDiscordClient()
		streamer = getStreamerDelay(len(songs))
		subsonicClient = getSubsonicClient(songs)
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("should return already playing state when invoked twice while playback is underway", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient, streamer)
		Expect(err).NotTo(HaveOccurred())

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(albumId)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.PlaybackComplete))
		}()
		time.Sleep(time.Millisecond * 1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(albumId)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.AlreadyPlaying))
		}()
		wg.Wait()
	})
})

// TODO - test that we finish track if there is an error with subsonic/streamer (either function)/discord

// TODO - should i add ginkgo helper to these functions below?
func getDiscordClient() *mocks.IDiscordClient {
	discordClient := mocks.NewIDiscordClient(GinkgoT())
	discordClient.EXPECT().Init(mock.AnythingOfType("*runner.SubcordantRunner")).Return(nil)
	discordClient.EXPECT().JoinVoiceChat().Return(fakeWriter, nil)
	return discordClient
}

func getStreamer(songCount int) *mocks.IStreamer {
	streamer := mocks.NewIStreamer(GinkgoT())
	for range songCount {
		streamer.EXPECT().PrepStream(mock.AnythingOfType("*url.URL")).Return(nil)
		streamer.EXPECT().Stream(fakeWriter).Return(nil)
	}
	return streamer
}

// Simulates the delay for Stream to return as if a song is playing
func getStreamerDelay(songCount int) *mocks.IStreamer {
	streamer := mocks.NewIStreamer(GinkgoT())
	for range songCount {
		streamer.EXPECT().PrepStream(mock.AnythingOfType("*url.URL")).Return(nil)
		streamer.EXPECT().Stream(fakeWriter).RunAndReturn(func(voice io.Writer) error {
			time.Sleep(time.Millisecond * 50)
			return nil
		})
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
		subsonicClient.EXPECT().StreamUrl(song.ID).Return(&url.URL{}, nil)
	}
	return subsonicClient
}

func getSongs(n int) []*subsonic.Child {
	songs := make([]*subsonic.Child, int(n))
	for i := range songs {
		songs[i] = &subsonic.Child{
			ID: time.Now().String(),
		}
	}
	return songs
}
