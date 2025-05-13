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
})
