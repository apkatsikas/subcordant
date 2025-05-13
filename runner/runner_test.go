package runner_test

import (
	"github.com/apkatsikas/subcordant/interfaces/mocks"
	"github.com/apkatsikas/subcordant/runner"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("runner", func() {
	const albumId = "foobar"

	var subcordantRunner *runner.SubcordantRunner
	var discordClient *mocks.IDiscordClient
	var subsonicClient *mocks.ISubsonicClient

	BeforeEach(func() {
		discordClient = mocks.NewIDiscordClient(GinkgoT())
		discordClient.EXPECT().Init(mock.AnythingOfType("*runner.SubcordantRunner")).Return(nil)

		subsonicClient = mocks.NewISubsonicClient(GinkgoT())
		subsonicClient.EXPECT().Init().Return(nil)
		//subsonicClient.EXPECT().GetAlbum(albumId).Return(nil)
		subcordantRunner = &runner.SubcordantRunner{}
		subcordantRunner.Init(subsonicClient, discordClient)
	})

	It("will run", func() {
		subcordantRunner.HandlePlay(albumId)
		Expect(1).To(Equal(1))
	})
})
