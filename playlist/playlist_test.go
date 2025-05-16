package playlist_test

import (
	"github.com/apkatsikas/subcordant/playlist"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("playlist service", func() {
	var playlistService *playlist.PlaylistService

	BeforeEach(func() {
		playlistService = &playlist.PlaylistService{}
	})

	It("will return a playlist with the added song after adding a song", func() {
		var trackId = "foobar"
		playlistService.Add(trackId)
		playlist := playlistService.GetPlaylist()

		Expect(len(playlist)).To(Equal(1))
		Expect(playlist[0]).To(Equal(trackId))
	})

	It("will return a playlist with the final song after removing a song", func() {
		var trackId = "foobar"
		var trackId2 = "cool"
		playlistService.Add(trackId)
		playlistService.Add(trackId2)
		playlist := playlistService.GetPlaylist()

		Expect(len(playlist)).To(Equal(2))

		playlistService.FinishTrack()
		newPlaylist := playlistService.GetPlaylist()

		Expect(newPlaylist).To(Equal([]string{trackId2}))
	})
})
