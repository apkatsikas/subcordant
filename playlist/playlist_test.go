package playlist_test

import (
	"github.com/apkatsikas/subcordant/playlist"
	"github.com/apkatsikas/subcordant/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var track1 = types.Track{
	ID:   "foo",
	Path: "/a/path",
}

var track2 = types.Track{
	ID:   "foo",
	Path: "/a/path",
}

var _ = Describe("playlist service", func() {
	var playlistService *playlist.PlaylistService

	BeforeEach(func() {
		playlistService = &playlist.PlaylistService{}
	})

	It("will return a playlist with the added song after adding a song", func() {
		playlistService.Add(track1)
		playlist := playlistService.GetPlaylist()

		Expect(len(playlist)).To(Equal(1))
		Expect(playlist[0]).To(Equal(track1))
	})

	It("will return a playlist with the final song after removing a song", func() {
		playlistService.Add(track1)
		playlistService.Add(track2)
		playlist := playlistService.GetPlaylist()

		Expect(len(playlist)).To(Equal(2))

		playlistService.FinishTrack()
		newPlaylist := playlistService.GetPlaylist()

		Expect(newPlaylist).To(Equal([]types.Track{track2}))
	})

	It("will return an empty playlist after removing a song from an empty playlist", func() {
		playlist := playlistService.GetPlaylist()

		Expect(len(playlist)).To(Equal(0))

		playlistService.FinishTrack()
		newPlaylist := playlistService.GetPlaylist()

		Expect(len(newPlaylist)).To(Equal(0))
	})

	It("will clear the playlist when clear is called", func() {
		playlistService.Add(track1)
		playlistService.Add(track2)

		playlistService.Clear()

		Expect(len(playlistService.GetPlaylist())).To(Equal(0))
	})
})
