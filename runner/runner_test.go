package runner_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/interfaces/mocks"
	"github.com/apkatsikas/subcordant/runner"
	"github.com/apkatsikas/subcordant/types"
	flagutil "github.com/apkatsikas/subcordant/util/flag"
	"github.com/diamondburned/arikawa/v3/discord"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

const albumId = "foobar"

var guildId discord.GuildID = discord.NullGuildID

const albumName = "foobar album"
const dontSwitchChannels discord.ChannelID = discord.NullChannelID

type nopWriter struct{}

func (nopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

var fakeWriter = nopWriter{}

var anyUrl = mock.AnythingOfType("*url.URL")
var anyString = mock.AnythingOfType("string")
var anyCancelContext = mock.AnythingOfType("*context.cancelCtx")

var _ = DescribeTableSubtree("runner init and play",
	func(songCount int) {
		var songs = getSongs(songCount)

		var subcordantRunner *runner.SubcordantRunner
		var discordClient *mocks.IDiscordClient
		var subsonicClient *mocks.ISubsonicClient
		var streamer *mocks.IStreamer

		var initError error
		var playError error
		var playState types.PlaybackState

		BeforeEach(func() {
			discordClient = getDiscordClient([]string{albumName})
			discordClient.EXPECT().GetVoice().Return(fakeWriter).Times(songCount)
			streamer = getStreamer(len(songs))
			subsonicClient = getSubsonicClient(songs, true)
			subcordantRunner = &runner.SubcordantRunner{}

			initError = subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
			playState, playError = subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
		})

		It("should not error", func() {
			Expect(initError).NotTo(HaveOccurred())
			Expect(playError).NotTo(HaveOccurred())
			Expect(playState).To(Equal(types.PlaybackComplete))
		})

		It("should show complete playback", func() {
			Expect(playState).To(Equal(types.PlaybackComplete))
		})

		It("should complete all tracks", func() {
			Expect(subcordantRunner.GetPlaylist()).To(HaveLen(0))
		})
	},
	Entry("1 song", 1),
	Entry("2 songs", 2),
)

var _ = DescribeTableSubtree("runner init and play",
	func(songCount int) {
		var songs = getSongs(songCount)

		var subcordantRunner *runner.SubcordantRunner
		var discordClient *mocks.IDiscordClient
		var subsonicClient *mocks.ISubsonicClient
		var streamer *mocks.IStreamer

		var initError error
		var playError error
		var playState types.PlaybackState

		BeforeEach(func() {
			discordClient = getDiscordClient([]string{albumName})
			discordClient.EXPECT().GetVoice().Return(fakeWriter).Times(songCount)
			streamer = getStreamerFromFile(len(songs))
			subsonicClient = getSubsonicClient(songs, false)
			subcordantRunner = &runner.SubcordantRunner{}

			initError = subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromFile)
			playState, playError = subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
		})

		It("should not error", func() {
			Expect(initError).NotTo(HaveOccurred())
			Expect(playError).NotTo(HaveOccurred())
			Expect(playState).To(Equal(types.PlaybackComplete))
		})

		It("should show complete playback", func() {
			Expect(playState).To(Equal(types.PlaybackComplete))
		})

		It("should complete all tracks", func() {
			Expect(subcordantRunner.GetPlaylist()).To(HaveLen(0))
		})
	},
	Entry("1 song, stream from set to file", 1),
	Entry("1 song, stream from set to file", 2),
)

var _ = Describe("runner init and play resulting in a channel change during playback", func() {
	const album1Name = "album1"
	const album2Name = "album2"
	const albumSongCount = 2
	var album1Songs = getSongs(albumSongCount)
	var album2Songs = getSongs(albumSongCount)
	var songCount = albumSongCount
	var firstSongFromAlbum1 = []*subsonic.Child{album1Songs[0]}

	var sf = discord.NewSnowflake(time.Now())
	var switchToChannel = discord.ChannelID(sf)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var streamer *mocks.IStreamer

	BeforeEach(func() {
		discordClient = getDiscordClient([]string{album1Name, album2Name})
		discordClient.EXPECT().GetVoice().Return(fakeWriter).Times(songCount)

		discordClient.EXPECT().JoinVoiceChat(guildId, switchToChannel).Return(switchToChannel, nil).Once()
		discordClient.EXPECT().SwitchVoiceChannel(switchToChannel).Return(nil).Once()

		streamer = getStreamerDelay(songCount)
		subsonicClient = getMultipleAlbumSubsonicClient(map[string][]*subsonic.Child{
			album1Name: firstSongFromAlbum1,
			album2Name: album2Songs,
		})
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("should return playback complete state when invoked twice while playback is underway", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(album1Name, guildId, dontSwitchChannels)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.PlaybackComplete))
		}()
		time.Sleep(time.Millisecond * 1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(album2Name, guildId, switchToChannel)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.PlaybackComplete))
		}()
		wg.Wait()
	})
})

var _ = Describe("runner init and play resulting in a failed channel change during playback", func() {
	const album1Name = "album1"
	const album2Name = "album2"
	const albumSongCount = 2
	var album1Songs = getSongs(albumSongCount)
	var songCount = 1
	var firstSongFromAlbum1 = []*subsonic.Child{album1Songs[0]}

	var sf = discord.NewSnowflake(time.Now())
	var switchToChannel = discord.ChannelID(sf)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var streamer *mocks.IStreamer

	BeforeEach(func() {
		// Only the first album will produce a message about queuing - only pass in 1
		discordClient = getDiscordClient([]string{album1Name})
		discordClient.EXPECT().GetVoice().Return(fakeWriter).Times(songCount)
		discordClient.EXPECT().SendMessage("Failed to switch channels, error is Failed to switch voice channel").Once()

		discordClient.EXPECT().JoinVoiceChat(guildId, switchToChannel).Return(switchToChannel, nil).Once()
		discordClient.EXPECT().SwitchVoiceChannel(switchToChannel).Return(
			fmt.Errorf("Failed to switch voice channel")).Once()

		streamer = getStreamerDelay(songCount)
		subsonicClient = getMultipleAlbumSubsonicClient(map[string][]*subsonic.Child{
			album1Name: firstSongFromAlbum1,
		})
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("should return an invalid state when invoked twice - on the second time - while playback is underway", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(album1Name, guildId, dontSwitchChannels)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.PlaybackComplete))
		}()
		time.Sleep(time.Millisecond * 1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(album2Name, guildId, switchToChannel)
			Expect(err).To(HaveOccurred())
			Expect(state).To(Equal(types.Invalid))
		}()
		wg.Wait()
	})
})

var _ = Describe("runner", func() {
	const album1Name = "album1"
	const album2Name = "album2"
	const albumSongCount = 2
	var album1Songs = getSongs(albumSongCount)
	var album2Songs = getSongs(albumSongCount)
	var songCount = len(album1Songs) + len(album2Songs)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var streamer *mocks.IStreamer

	BeforeEach(func() {
		discordClient = getDiscordClient([]string{album1Name, album2Name})
		discordClient.EXPECT().GetVoice().Return(fakeWriter).Times(songCount)
		// Add an additional JoinVoiceChat expectation
		discordClient.EXPECT().JoinVoiceChat(guildId, dontSwitchChannels).Return(dontSwitchChannels, nil).Once()
		streamer = getStreamerDelay(songCount)
		subsonicClient = getMultipleAlbumSubsonicClient(map[string][]*subsonic.Child{
			album1Name: album1Songs,
			album2Name: album2Songs,
		})
		subcordantRunner = &runner.SubcordantRunner{}
	})

	It("should return already playing state when invoked twice while playback is underway", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(album1Name, guildId, dontSwitchChannels)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.PlaybackComplete))
		}()
		time.Sleep(time.Millisecond * 1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(album2Name, guildId, dontSwitchChannels)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.AlreadyPlaying))
		}()
		wg.Wait()
	})
})

var _ = Describe("runner", func() {
	const songCount = 3
	var songs = getSongs(songCount)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var streamer *mocks.IStreamer

	BeforeEach(func() {
		discordClient = getDiscordClient([]string{albumName})
		discordClient.EXPECT().GetVoice().Return(fakeWriter).Once()
		// We only want the first song to build our expectations, as the rest will be skipped
		streamer = getStreamerDelay(1)
		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(nil).Once()
		subsonicClient.EXPECT().GetAlbum(albumId).Return(&subsonic.AlbumID3{
			Name: albumName,
			Song: songs,
		}, nil).Once()
		subsonicClient.EXPECT().StreamUrl(songs[0].ID).Return(&url.URL{}, nil).Once()
		subcordantRunner = &runner.SubcordantRunner{}
		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should clear the playlist when reset during playback", func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.PlaybackComplete))
		}()
		time.Sleep(time.Millisecond * 1)

		Expect(subcordantRunner.PlaylistService.GetPlaylist()).To(HaveLen(len(songs)))
		subcordantRunner.Reset()

		Expect(subcordantRunner.PlaylistService.GetPlaylist()).To(HaveLen(0))

		wg.Wait()
	})
})

var _ = Describe("runner", func() {
	const songCount = 3
	var songs = getSongs(songCount)

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient
	var streamer *mocks.IStreamer

	BeforeEach(func() {
		discordClient = getDiscordClient([]string{albumName})
		discordClient.EXPECT().GetVoice().Return(fakeWriter).Once()
		discordClient.EXPECT().LeaveVoiceSession().Return().Once()
		// We only want the first song to build our expectations, as the rest will be skipped
		streamer = getStreamerDelay(1)
		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(nil).Once()
		subsonicClient.EXPECT().GetAlbum(albumId).Return(&subsonic.AlbumID3{
			Name: albumName,
			Song: songs,
		}, nil).Once()
		subsonicClient.EXPECT().StreamUrl(songs[0].ID).Return(&url.URL{}, nil).Once()
		subcordantRunner = &runner.SubcordantRunner{}
		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should clear the playlist when disconnected during playback", func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer GinkgoRecover()
			state, err := subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
			Expect(err).NotTo(HaveOccurred())
			Expect(state).To(Equal(types.PlaybackComplete))
		}()
		time.Sleep(time.Millisecond * 1)

		Expect(subcordantRunner.PlaylistService.GetPlaylist()).To(HaveLen(len(songs)))
		subcordantRunner.Disconnect()

		Expect(subcordantRunner.PlaylistService.GetPlaylist()).To(HaveLen(0))

		wg.Wait()
	})
})

var _ = Describe("runner", func() {
	var subcordantRunner *runner.SubcordantRunner
	var subsonicClient *mocks.ISubsonicClient

	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(fmt.Errorf("init error")).Once()
	})

	It("should error on init if subsonic init errors", func() {
		err := subcordantRunner.Init(subsonicClient, mocks.NewIDiscordClient(GinkgoT()),
			mocks.NewIStreamer(GinkgoT()), flagutil.StreamFromStream)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("runner", func() {
	var subcordantRunner *runner.SubcordantRunner
	var subsonicClient *mocks.ISubsonicClient
	var discordClient *mocks.IDiscordClient

	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(nil).Once()
		discordClient = mocks.NewIDiscordClient(GinkgoT())
		discordClient.EXPECT().Init(subcordantRunner).Return(fmt.Errorf("init error")).Once()
	})

	It("should error on init if discord init errors", func() {
		err := subcordantRunner.Init(subsonicClient, discordClient,
			mocks.NewIStreamer(GinkgoT()), flagutil.StreamFromStream)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("runner play if get album errors", func() {
	var subcordantRunner *runner.SubcordantRunner
	var subsonicClient *mocks.ISubsonicClient
	var discordClient *mocks.IDiscordClient

	var playError error
	var playbackState types.PlaybackState

	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(nil).Once()
		subsonicClient.EXPECT().GetAlbum(albumId).Return(nil, fmt.Errorf("get album error")).Once()
		discordClient = mocks.NewIDiscordClient(GinkgoT())
		discordClient.EXPECT().Init(subcordantRunner).Return(nil).Once()
		discordClient.EXPECT().SendMessage(fmt.Sprintf("Could not find album with ID of %v", albumId)).Once()
		discordClient.EXPECT().JoinVoiceChat(guildId, dontSwitchChannels).Return(dontSwitchChannels, nil).Once()
		err := subcordantRunner.Init(subsonicClient, discordClient,
			mocks.NewIStreamer(GinkgoT()), flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())

		playbackState, playError = subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
	})

	It("should return an invalid state", func() {
		Expect(playbackState).To(Equal(types.Invalid))
	})

	It("should error", func() {
		Expect(playError).To(HaveOccurred())
	})
})

var _ = Describe("runner play if stream url errors", func() {
	const songCount = 1
	var songs = getSongs(songCount)

	var subcordantRunner *runner.SubcordantRunner
	var subsonicClient *mocks.ISubsonicClient
	var discordClient *mocks.IDiscordClient

	var playError error
	var playbackState types.PlaybackState

	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(nil).Once()
		subsonicClient.EXPECT().GetAlbum(albumId).Return(&subsonic.AlbumID3{
			Name: albumName,
			Song: songs,
		}, nil).Once()
		subsonicClient.EXPECT().StreamUrl(songs[0].ID).Return(nil, fmt.Errorf("stream url error")).Once()
		discordClient = getDiscordClient([]string{albumName})
		err := subcordantRunner.Init(subsonicClient, discordClient,
			mocks.NewIStreamer(GinkgoT()), flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())

		playbackState, playError = subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
	})

	It("should return an invalid state", func() {
		Expect(playbackState).To(Equal(types.Invalid))
	})

	It("should error", func() {
		Expect(playError).To(HaveOccurred())
	})
})

var _ = Describe("runner play if prep stream from file errors, using stream from stream", func() {
	const songCount = 1
	var songs = getSongs(songCount)

	var subcordantRunner *runner.SubcordantRunner
	var subsonicClient *mocks.ISubsonicClient
	var discordClient *mocks.IDiscordClient
	var streamer *mocks.IStreamer

	var playError error
	var playbackState types.PlaybackState

	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
		subsonicClient = getSubsonicClient(songs, true)
		discordClient = getDiscordClient([]string{albumName})
		streamer = mocks.NewIStreamer(GinkgoT())

		streamer.EXPECT().PrepStreamFromStream(anyUrl).Return(fmt.Errorf("prep stream error")).Once()

		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())

		playbackState, playError = subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
	})

	It("should return an invalid state", func() {
		Expect(playbackState).To(Equal(types.Invalid))
	})

	It("should error", func() {
		Expect(playError).To(HaveOccurred())
	})

	It("should finish the track", func() {
		Expect(subcordantRunner.GetPlaylist()).To(HaveLen(songCount - 1))
	})
})

var _ = Describe("runner play if prep stream from file errors, using stream from file", func() {
	const songCount = 1
	var songs = getSongs(songCount)

	var subcordantRunner *runner.SubcordantRunner
	var subsonicClient *mocks.ISubsonicClient
	var discordClient *mocks.IDiscordClient
	var streamer *mocks.IStreamer

	var playError error
	var playbackState types.PlaybackState

	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
		subsonicClient = getSubsonicClient(songs, false)
		discordClient = getDiscordClient([]string{albumName})
		streamer = mocks.NewIStreamer(GinkgoT())

		streamer.EXPECT().PrepStreamFromFile(anyString).Return(fmt.Errorf("prep stream error")).Once()

		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromFile)
		Expect(err).NotTo(HaveOccurred())

		playbackState, playError = subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
	})

	It("should return an invalid state", func() {
		Expect(playbackState).To(Equal(types.Invalid))
	})

	It("should error", func() {
		Expect(playError).To(HaveOccurred())
	})

	It("should finish the track", func() {
		Expect(subcordantRunner.GetPlaylist()).To(HaveLen(songCount - 1))
	})
})

var _ = Describe("runner play if stream errors", func() {
	const songCount = 1
	var songs = getSongs(songCount)

	var subcordantRunner *runner.SubcordantRunner
	var subsonicClient *mocks.ISubsonicClient
	var discordClient *mocks.IDiscordClient
	var streamer *mocks.IStreamer

	var playError error
	var playbackState types.PlaybackState

	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
		subsonicClient = getSubsonicClient(songs, true)
		discordClient = getDiscordClient([]string{albumName})
		discordClient.EXPECT().GetVoice().Return(fakeWriter).Once()
		streamer = mocks.NewIStreamer(GinkgoT())

		streamer.EXPECT().PrepStreamFromStream(anyUrl).Return(nil).Once()
		streamer.EXPECT().Stream(anyCancelContext, fakeWriter).Return(fmt.Errorf("stream error")).Once()

		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())

		playbackState, playError = subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
	})

	It("should return an invalid state", func() {
		Expect(playbackState).To(Equal(types.Invalid))
	})

	It("should error", func() {
		Expect(playError).To(HaveOccurred())
	})

	It("should finish the track", func() {
		Expect(subcordantRunner.GetPlaylist()).To(HaveLen(songCount - 1))
	})
})

var _ = Describe("runner play if join voice errors", func() {
	const errorMessage = "join voice error"

	var subcordantRunner *runner.SubcordantRunner
	var subsonicClient *mocks.ISubsonicClient
	var discordClient *mocks.IDiscordClient
	var streamer *mocks.IStreamer

	var playError error
	var playbackState types.PlaybackState

	BeforeEach(func() {
		subcordantRunner = &runner.SubcordantRunner{}
		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(nil).Once()
		discordClient = mocks.NewIDiscordClient(GinkgoT())
		discordClient.EXPECT().Init(subcordantRunner).Return(nil).Once()
		discordClient.EXPECT().JoinVoiceChat(guildId, dontSwitchChannels).Return(
			dontSwitchChannels, errors.New(errorMessage)).Once()
		discordClient.EXPECT().SendMessage(fmt.Sprintf("Could not join voice, error is %v", errorMessage))
		streamer = mocks.NewIStreamer(GinkgoT())

		err := subcordantRunner.Init(subsonicClient, discordClient, streamer, flagutil.StreamFromStream)
		Expect(err).NotTo(HaveOccurred())

		playbackState, playError = subcordantRunner.Play(albumId, guildId, dontSwitchChannels)
	})

	It("should return an invalid state", func() {
		Expect(playbackState).To(Equal(types.Invalid))
	})

	It("should error", func() {
		Expect(playError).To(HaveOccurred())
	})

	It("should clear the playlist", func() {
		Expect(subcordantRunner.GetPlaylist()).To(HaveLen(0))
	})
})

func getDiscordClient(albums []string) *mocks.IDiscordClient {
	discordClient := mocks.NewIDiscordClient(GinkgoT())
	discordClient.EXPECT().Init(mock.AnythingOfType("*runner.SubcordantRunner")).Return(nil).Once()
	discordClient.EXPECT().JoinVoiceChat(guildId, dontSwitchChannels).Return(dontSwitchChannels, nil).Once()
	for _, album := range albums {
		discordClient.EXPECT().SendMessage(getQueuedAlbumMessage(album)).Once()
	}

	return discordClient
}

func getQueuedAlbumMessage(albumName string) string {
	return fmt.Sprintf("Queued album: %v", albumName)
}

func getStreamer(songCount int) *mocks.IStreamer {
	streamer := mocks.NewIStreamer(GinkgoT())
	for range songCount {
		streamer.EXPECT().PrepStreamFromStream(anyUrl).Return(nil).Once()
		streamer.EXPECT().Stream(anyCancelContext, fakeWriter).Return(nil).Once()
	}
	return streamer
}

func getStreamerFromFile(songCount int) *mocks.IStreamer {
	streamer := mocks.NewIStreamer(GinkgoT())
	for range songCount {
		streamer.EXPECT().PrepStreamFromFile(anyString).Return(nil).Once()
		streamer.EXPECT().Stream(anyCancelContext, fakeWriter).Return(nil).Once()
	}
	return streamer
}

// Simulates the delay for Stream to return as if a song is playing
func getStreamerDelay(songCount int) *mocks.IStreamer {
	streamer := mocks.NewIStreamer(GinkgoT())
	streamer.EXPECT().PrepStreamFromStream(anyUrl).Return(nil).Times(songCount)
	streamer.EXPECT().Stream(anyCancelContext, fakeWriter).RunAndReturn(func(_ context.Context, _ io.Writer) error {
		time.Sleep(time.Millisecond * 50)
		return nil
	}).Times(songCount)
	return streamer
}

func getSubsonicClient(songs []*subsonic.Child, fromStream bool) *mocks.ISubsonicClient {
	subsonicClient := mocks.NewISubsonicClient(GinkgoT())
	subsonicClient.EXPECT().Init().Return(nil).Once()
	subsonicClient.EXPECT().GetAlbum(albumId).Return(&subsonic.AlbumID3{
		Name: albumName,
		Song: songs,
	}, nil).Once()

	if fromStream {
		for _, song := range songs {
			subsonicClient.EXPECT().StreamUrl(song.ID).Return(&url.URL{}, nil).Once()
		}
	}

	return subsonicClient
}

func getMultipleAlbumSubsonicClient(albumSongs map[string][]*subsonic.Child) *mocks.ISubsonicClient {
	subsonicClient := mocks.NewISubsonicClient(GinkgoT())
	subsonicClient.EXPECT().Init().Return(nil).Once()

	for albumName, songs := range albumSongs {
		subsonicClient.EXPECT().GetAlbum(albumName).Return(&subsonic.AlbumID3{
			Name: albumName,
			Song: songs,
		}, nil).Once()
		for _, song := range songs {
			subsonicClient.EXPECT().StreamUrl(song.ID).Return(&url.URL{}, nil).Once()
		}
	}

	return subsonicClient
}

func getSongs(n int) []*subsonic.Child {
	songs := make([]*subsonic.Child, int(n))
	for i := range songs {
		songs[i] = &subsonic.Child{
			ID: fmt.Sprintf("%v%v", time.Now().String(), i),
		}
	}
	return songs
}
