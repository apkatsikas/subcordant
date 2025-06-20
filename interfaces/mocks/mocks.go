// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"context"
	"io"
	"net/url"

	"github.com/apkatsikas/go-subsonic"
	"github.com/apkatsikas/subcordant/interfaces"
	"github.com/apkatsikas/subcordant/types"
	"github.com/diamondburned/arikawa/v3/discord"
	mock "github.com/stretchr/testify/mock"
)

// NewICommandHandler creates a new instance of ICommandHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICommandHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICommandHandler {
	mock := &ICommandHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// ICommandHandler is an autogenerated mock type for the ICommandHandler type
type ICommandHandler struct {
	mock.Mock
}

type ICommandHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *ICommandHandler) EXPECT() *ICommandHandler_Expecter {
	return &ICommandHandler_Expecter{mock: &_m.Mock}
}

// Disconnect provides a mock function for the type ICommandHandler
func (_mock *ICommandHandler) Disconnect() {
	_mock.Called()
	return
}

// ICommandHandler_Disconnect_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Disconnect'
type ICommandHandler_Disconnect_Call struct {
	*mock.Call
}

// Disconnect is a helper method to define mock.On call
func (_e *ICommandHandler_Expecter) Disconnect() *ICommandHandler_Disconnect_Call {
	return &ICommandHandler_Disconnect_Call{Call: _e.mock.On("Disconnect")}
}

func (_c *ICommandHandler_Disconnect_Call) Run(run func()) *ICommandHandler_Disconnect_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ICommandHandler_Disconnect_Call) Return() *ICommandHandler_Disconnect_Call {
	_c.Call.Return()
	return _c
}

func (_c *ICommandHandler_Disconnect_Call) RunAndReturn(run func()) *ICommandHandler_Disconnect_Call {
	_c.Run(run)
	return _c
}

// Play provides a mock function for the type ICommandHandler
func (_mock *ICommandHandler) Play(albumId string, guildId discord.GuildID, channelId discord.ChannelID) (types.PlaybackState, error) {
	ret := _mock.Called(albumId, guildId, channelId)

	if len(ret) == 0 {
		panic("no return value specified for Play")
	}

	var r0 types.PlaybackState
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string, discord.GuildID, discord.ChannelID) (types.PlaybackState, error)); ok {
		return returnFunc(albumId, guildId, channelId)
	}
	if returnFunc, ok := ret.Get(0).(func(string, discord.GuildID, discord.ChannelID) types.PlaybackState); ok {
		r0 = returnFunc(albumId, guildId, channelId)
	} else {
		r0 = ret.Get(0).(types.PlaybackState)
	}
	if returnFunc, ok := ret.Get(1).(func(string, discord.GuildID, discord.ChannelID) error); ok {
		r1 = returnFunc(albumId, guildId, channelId)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// ICommandHandler_Play_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Play'
type ICommandHandler_Play_Call struct {
	*mock.Call
}

// Play is a helper method to define mock.On call
//   - albumId
//   - guildId
//   - channelId
func (_e *ICommandHandler_Expecter) Play(albumId interface{}, guildId interface{}, channelId interface{}) *ICommandHandler_Play_Call {
	return &ICommandHandler_Play_Call{Call: _e.mock.On("Play", albumId, guildId, channelId)}
}

func (_c *ICommandHandler_Play_Call) Run(run func(albumId string, guildId discord.GuildID, channelId discord.ChannelID)) *ICommandHandler_Play_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(discord.GuildID), args[2].(discord.ChannelID))
	})
	return _c
}

func (_c *ICommandHandler_Play_Call) Return(playbackState types.PlaybackState, err error) *ICommandHandler_Play_Call {
	_c.Call.Return(playbackState, err)
	return _c
}

func (_c *ICommandHandler_Play_Call) RunAndReturn(run func(albumId string, guildId discord.GuildID, channelId discord.ChannelID) (types.PlaybackState, error)) *ICommandHandler_Play_Call {
	_c.Call.Return(run)
	return _c
}

// Reset provides a mock function for the type ICommandHandler
func (_mock *ICommandHandler) Reset() {
	_mock.Called()
	return
}

// ICommandHandler_Reset_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Reset'
type ICommandHandler_Reset_Call struct {
	*mock.Call
}

// Reset is a helper method to define mock.On call
func (_e *ICommandHandler_Expecter) Reset() *ICommandHandler_Reset_Call {
	return &ICommandHandler_Reset_Call{Call: _e.mock.On("Reset")}
}

func (_c *ICommandHandler_Reset_Call) Run(run func()) *ICommandHandler_Reset_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ICommandHandler_Reset_Call) Return() *ICommandHandler_Reset_Call {
	_c.Call.Return()
	return _c
}

func (_c *ICommandHandler_Reset_Call) RunAndReturn(run func()) *ICommandHandler_Reset_Call {
	_c.Run(run)
	return _c
}

// NewIDiscordClient creates a new instance of IDiscordClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIDiscordClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *IDiscordClient {
	mock := &IDiscordClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// IDiscordClient is an autogenerated mock type for the IDiscordClient type
type IDiscordClient struct {
	mock.Mock
}

type IDiscordClient_Expecter struct {
	mock *mock.Mock
}

func (_m *IDiscordClient) EXPECT() *IDiscordClient_Expecter {
	return &IDiscordClient_Expecter{mock: &_m.Mock}
}

// GetVoice provides a mock function for the type IDiscordClient
func (_mock *IDiscordClient) GetVoice() io.Writer {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetVoice")
	}

	var r0 io.Writer
	if returnFunc, ok := ret.Get(0).(func() io.Writer); ok {
		r0 = returnFunc()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.Writer)
		}
	}
	return r0
}

// IDiscordClient_GetVoice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetVoice'
type IDiscordClient_GetVoice_Call struct {
	*mock.Call
}

// GetVoice is a helper method to define mock.On call
func (_e *IDiscordClient_Expecter) GetVoice() *IDiscordClient_GetVoice_Call {
	return &IDiscordClient_GetVoice_Call{Call: _e.mock.On("GetVoice")}
}

func (_c *IDiscordClient_GetVoice_Call) Run(run func()) *IDiscordClient_GetVoice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IDiscordClient_GetVoice_Call) Return(writer io.Writer) *IDiscordClient_GetVoice_Call {
	_c.Call.Return(writer)
	return _c
}

func (_c *IDiscordClient_GetVoice_Call) RunAndReturn(run func() io.Writer) *IDiscordClient_GetVoice_Call {
	_c.Call.Return(run)
	return _c
}

// Init provides a mock function for the type IDiscordClient
func (_mock *IDiscordClient) Init(commandHandler interfaces.ICommandHandler) error {
	ret := _mock.Called(commandHandler)

	if len(ret) == 0 {
		panic("no return value specified for Init")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(interfaces.ICommandHandler) error); ok {
		r0 = returnFunc(commandHandler)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// IDiscordClient_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type IDiscordClient_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - commandHandler
func (_e *IDiscordClient_Expecter) Init(commandHandler interface{}) *IDiscordClient_Init_Call {
	return &IDiscordClient_Init_Call{Call: _e.mock.On("Init", commandHandler)}
}

func (_c *IDiscordClient_Init_Call) Run(run func(commandHandler interfaces.ICommandHandler)) *IDiscordClient_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(interfaces.ICommandHandler))
	})
	return _c
}

func (_c *IDiscordClient_Init_Call) Return(err error) *IDiscordClient_Init_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *IDiscordClient_Init_Call) RunAndReturn(run func(commandHandler interfaces.ICommandHandler) error) *IDiscordClient_Init_Call {
	_c.Call.Return(run)
	return _c
}

// JoinVoiceChat provides a mock function for the type IDiscordClient
func (_mock *IDiscordClient) JoinVoiceChat(guildId discord.GuildID, channelId discord.ChannelID) (discord.ChannelID, error) {
	ret := _mock.Called(guildId, channelId)

	if len(ret) == 0 {
		panic("no return value specified for JoinVoiceChat")
	}

	var r0 discord.ChannelID
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(discord.GuildID, discord.ChannelID) (discord.ChannelID, error)); ok {
		return returnFunc(guildId, channelId)
	}
	if returnFunc, ok := ret.Get(0).(func(discord.GuildID, discord.ChannelID) discord.ChannelID); ok {
		r0 = returnFunc(guildId, channelId)
	} else {
		r0 = ret.Get(0).(discord.ChannelID)
	}
	if returnFunc, ok := ret.Get(1).(func(discord.GuildID, discord.ChannelID) error); ok {
		r1 = returnFunc(guildId, channelId)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// IDiscordClient_JoinVoiceChat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'JoinVoiceChat'
type IDiscordClient_JoinVoiceChat_Call struct {
	*mock.Call
}

// JoinVoiceChat is a helper method to define mock.On call
//   - guildId
//   - channelId
func (_e *IDiscordClient_Expecter) JoinVoiceChat(guildId interface{}, channelId interface{}) *IDiscordClient_JoinVoiceChat_Call {
	return &IDiscordClient_JoinVoiceChat_Call{Call: _e.mock.On("JoinVoiceChat", guildId, channelId)}
}

func (_c *IDiscordClient_JoinVoiceChat_Call) Run(run func(guildId discord.GuildID, channelId discord.ChannelID)) *IDiscordClient_JoinVoiceChat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(discord.GuildID), args[1].(discord.ChannelID))
	})
	return _c
}

func (_c *IDiscordClient_JoinVoiceChat_Call) Return(channelID discord.ChannelID, err error) *IDiscordClient_JoinVoiceChat_Call {
	_c.Call.Return(channelID, err)
	return _c
}

func (_c *IDiscordClient_JoinVoiceChat_Call) RunAndReturn(run func(guildId discord.GuildID, channelId discord.ChannelID) (discord.ChannelID, error)) *IDiscordClient_JoinVoiceChat_Call {
	_c.Call.Return(run)
	return _c
}

// LeaveVoiceSession provides a mock function for the type IDiscordClient
func (_mock *IDiscordClient) LeaveVoiceSession() {
	_mock.Called()
	return
}

// IDiscordClient_LeaveVoiceSession_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'LeaveVoiceSession'
type IDiscordClient_LeaveVoiceSession_Call struct {
	*mock.Call
}

// LeaveVoiceSession is a helper method to define mock.On call
func (_e *IDiscordClient_Expecter) LeaveVoiceSession() *IDiscordClient_LeaveVoiceSession_Call {
	return &IDiscordClient_LeaveVoiceSession_Call{Call: _e.mock.On("LeaveVoiceSession")}
}

func (_c *IDiscordClient_LeaveVoiceSession_Call) Run(run func()) *IDiscordClient_LeaveVoiceSession_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *IDiscordClient_LeaveVoiceSession_Call) Return() *IDiscordClient_LeaveVoiceSession_Call {
	_c.Call.Return()
	return _c
}

func (_c *IDiscordClient_LeaveVoiceSession_Call) RunAndReturn(run func()) *IDiscordClient_LeaveVoiceSession_Call {
	_c.Run(run)
	return _c
}

// SendMessage provides a mock function for the type IDiscordClient
func (_mock *IDiscordClient) SendMessage(message string) {
	_mock.Called(message)
	return
}

// IDiscordClient_SendMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMessage'
type IDiscordClient_SendMessage_Call struct {
	*mock.Call
}

// SendMessage is a helper method to define mock.On call
//   - message
func (_e *IDiscordClient_Expecter) SendMessage(message interface{}) *IDiscordClient_SendMessage_Call {
	return &IDiscordClient_SendMessage_Call{Call: _e.mock.On("SendMessage", message)}
}

func (_c *IDiscordClient_SendMessage_Call) Run(run func(message string)) *IDiscordClient_SendMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IDiscordClient_SendMessage_Call) Return() *IDiscordClient_SendMessage_Call {
	_c.Call.Return()
	return _c
}

func (_c *IDiscordClient_SendMessage_Call) RunAndReturn(run func(message string)) *IDiscordClient_SendMessage_Call {
	_c.Run(run)
	return _c
}

// SwitchVoiceChannel provides a mock function for the type IDiscordClient
func (_mock *IDiscordClient) SwitchVoiceChannel(channelId discord.ChannelID) error {
	ret := _mock.Called(channelId)

	if len(ret) == 0 {
		panic("no return value specified for SwitchVoiceChannel")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(discord.ChannelID) error); ok {
		r0 = returnFunc(channelId)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// IDiscordClient_SwitchVoiceChannel_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SwitchVoiceChannel'
type IDiscordClient_SwitchVoiceChannel_Call struct {
	*mock.Call
}

// SwitchVoiceChannel is a helper method to define mock.On call
//   - channelId
func (_e *IDiscordClient_Expecter) SwitchVoiceChannel(channelId interface{}) *IDiscordClient_SwitchVoiceChannel_Call {
	return &IDiscordClient_SwitchVoiceChannel_Call{Call: _e.mock.On("SwitchVoiceChannel", channelId)}
}

func (_c *IDiscordClient_SwitchVoiceChannel_Call) Run(run func(channelId discord.ChannelID)) *IDiscordClient_SwitchVoiceChannel_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(discord.ChannelID))
	})
	return _c
}

func (_c *IDiscordClient_SwitchVoiceChannel_Call) Return(err error) *IDiscordClient_SwitchVoiceChannel_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *IDiscordClient_SwitchVoiceChannel_Call) RunAndReturn(run func(channelId discord.ChannelID) error) *IDiscordClient_SwitchVoiceChannel_Call {
	_c.Call.Return(run)
	return _c
}

// NewIStreamer creates a new instance of IStreamer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIStreamer(t interface {
	mock.TestingT
	Cleanup(func())
}) *IStreamer {
	mock := &IStreamer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// IStreamer is an autogenerated mock type for the IStreamer type
type IStreamer struct {
	mock.Mock
}

type IStreamer_Expecter struct {
	mock *mock.Mock
}

func (_m *IStreamer) EXPECT() *IStreamer_Expecter {
	return &IStreamer_Expecter{mock: &_m.Mock}
}

// PrepStreamFromFile provides a mock function for the type IStreamer
func (_mock *IStreamer) PrepStreamFromFile(file string) error {
	ret := _mock.Called(file)

	if len(ret) == 0 {
		panic("no return value specified for PrepStreamFromFile")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(string) error); ok {
		r0 = returnFunc(file)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// IStreamer_PrepStreamFromFile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PrepStreamFromFile'
type IStreamer_PrepStreamFromFile_Call struct {
	*mock.Call
}

// PrepStreamFromFile is a helper method to define mock.On call
//   - file
func (_e *IStreamer_Expecter) PrepStreamFromFile(file interface{}) *IStreamer_PrepStreamFromFile_Call {
	return &IStreamer_PrepStreamFromFile_Call{Call: _e.mock.On("PrepStreamFromFile", file)}
}

func (_c *IStreamer_PrepStreamFromFile_Call) Run(run func(file string)) *IStreamer_PrepStreamFromFile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *IStreamer_PrepStreamFromFile_Call) Return(err error) *IStreamer_PrepStreamFromFile_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *IStreamer_PrepStreamFromFile_Call) RunAndReturn(run func(file string) error) *IStreamer_PrepStreamFromFile_Call {
	_c.Call.Return(run)
	return _c
}

// PrepStreamFromStream provides a mock function for the type IStreamer
func (_mock *IStreamer) PrepStreamFromStream(streamUrl *url.URL) error {
	ret := _mock.Called(streamUrl)

	if len(ret) == 0 {
		panic("no return value specified for PrepStreamFromStream")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(*url.URL) error); ok {
		r0 = returnFunc(streamUrl)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// IStreamer_PrepStreamFromStream_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PrepStreamFromStream'
type IStreamer_PrepStreamFromStream_Call struct {
	*mock.Call
}

// PrepStreamFromStream is a helper method to define mock.On call
//   - streamUrl
func (_e *IStreamer_Expecter) PrepStreamFromStream(streamUrl interface{}) *IStreamer_PrepStreamFromStream_Call {
	return &IStreamer_PrepStreamFromStream_Call{Call: _e.mock.On("PrepStreamFromStream", streamUrl)}
}

func (_c *IStreamer_PrepStreamFromStream_Call) Run(run func(streamUrl *url.URL)) *IStreamer_PrepStreamFromStream_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*url.URL))
	})
	return _c
}

func (_c *IStreamer_PrepStreamFromStream_Call) Return(err error) *IStreamer_PrepStreamFromStream_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *IStreamer_PrepStreamFromStream_Call) RunAndReturn(run func(streamUrl *url.URL) error) *IStreamer_PrepStreamFromStream_Call {
	_c.Call.Return(run)
	return _c
}

// Stream provides a mock function for the type IStreamer
func (_mock *IStreamer) Stream(ctx context.Context, voice io.Writer) error {
	ret := _mock.Called(ctx, voice)

	if len(ret) == 0 {
		panic("no return value specified for Stream")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, io.Writer) error); ok {
		r0 = returnFunc(ctx, voice)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// IStreamer_Stream_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stream'
type IStreamer_Stream_Call struct {
	*mock.Call
}

// Stream is a helper method to define mock.On call
//   - ctx
//   - voice
func (_e *IStreamer_Expecter) Stream(ctx interface{}, voice interface{}) *IStreamer_Stream_Call {
	return &IStreamer_Stream_Call{Call: _e.mock.On("Stream", ctx, voice)}
}

func (_c *IStreamer_Stream_Call) Run(run func(ctx context.Context, voice io.Writer)) *IStreamer_Stream_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(io.Writer))
	})
	return _c
}

func (_c *IStreamer_Stream_Call) Return(err error) *IStreamer_Stream_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *IStreamer_Stream_Call) RunAndReturn(run func(ctx context.Context, voice io.Writer) error) *IStreamer_Stream_Call {
	_c.Call.Return(run)
	return _c
}

// NewISubsonicClient creates a new instance of ISubsonicClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewISubsonicClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *ISubsonicClient {
	mock := &ISubsonicClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// ISubsonicClient is an autogenerated mock type for the ISubsonicClient type
type ISubsonicClient struct {
	mock.Mock
}

type ISubsonicClient_Expecter struct {
	mock *mock.Mock
}

func (_m *ISubsonicClient) EXPECT() *ISubsonicClient_Expecter {
	return &ISubsonicClient_Expecter{mock: &_m.Mock}
}

// GetAlbum provides a mock function for the type ISubsonicClient
func (_mock *ISubsonicClient) GetAlbum(albumId string) (*subsonic.AlbumID3, error) {
	ret := _mock.Called(albumId)

	if len(ret) == 0 {
		panic("no return value specified for GetAlbum")
	}

	var r0 *subsonic.AlbumID3
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (*subsonic.AlbumID3, error)); ok {
		return returnFunc(albumId)
	}
	if returnFunc, ok := ret.Get(0).(func(string) *subsonic.AlbumID3); ok {
		r0 = returnFunc(albumId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*subsonic.AlbumID3)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(albumId)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// ISubsonicClient_GetAlbum_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAlbum'
type ISubsonicClient_GetAlbum_Call struct {
	*mock.Call
}

// GetAlbum is a helper method to define mock.On call
//   - albumId
func (_e *ISubsonicClient_Expecter) GetAlbum(albumId interface{}) *ISubsonicClient_GetAlbum_Call {
	return &ISubsonicClient_GetAlbum_Call{Call: _e.mock.On("GetAlbum", albumId)}
}

func (_c *ISubsonicClient_GetAlbum_Call) Run(run func(albumId string)) *ISubsonicClient_GetAlbum_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ISubsonicClient_GetAlbum_Call) Return(albumID3 *subsonic.AlbumID3, err error) *ISubsonicClient_GetAlbum_Call {
	_c.Call.Return(albumID3, err)
	return _c
}

func (_c *ISubsonicClient_GetAlbum_Call) RunAndReturn(run func(albumId string) (*subsonic.AlbumID3, error)) *ISubsonicClient_GetAlbum_Call {
	_c.Call.Return(run)
	return _c
}

// Init provides a mock function for the type ISubsonicClient
func (_mock *ISubsonicClient) Init() error {
	ret := _mock.Called()

	if len(ret) == 0 {
		panic("no return value specified for Init")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func() error); ok {
		r0 = returnFunc()
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// ISubsonicClient_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type ISubsonicClient_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
func (_e *ISubsonicClient_Expecter) Init() *ISubsonicClient_Init_Call {
	return &ISubsonicClient_Init_Call{Call: _e.mock.On("Init")}
}

func (_c *ISubsonicClient_Init_Call) Run(run func()) *ISubsonicClient_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ISubsonicClient_Init_Call) Return(err error) *ISubsonicClient_Init_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *ISubsonicClient_Init_Call) RunAndReturn(run func() error) *ISubsonicClient_Init_Call {
	_c.Call.Return(run)
	return _c
}

// StreamUrl provides a mock function for the type ISubsonicClient
func (_mock *ISubsonicClient) StreamUrl(trackId string) (*url.URL, error) {
	ret := _mock.Called(trackId)

	if len(ret) == 0 {
		panic("no return value specified for StreamUrl")
	}

	var r0 *url.URL
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(string) (*url.URL, error)); ok {
		return returnFunc(trackId)
	}
	if returnFunc, ok := ret.Get(0).(func(string) *url.URL); ok {
		r0 = returnFunc(trackId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url.URL)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(string) error); ok {
		r1 = returnFunc(trackId)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// ISubsonicClient_StreamUrl_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StreamUrl'
type ISubsonicClient_StreamUrl_Call struct {
	*mock.Call
}

// StreamUrl is a helper method to define mock.On call
//   - trackId
func (_e *ISubsonicClient_Expecter) StreamUrl(trackId interface{}) *ISubsonicClient_StreamUrl_Call {
	return &ISubsonicClient_StreamUrl_Call{Call: _e.mock.On("StreamUrl", trackId)}
}

func (_c *ISubsonicClient_StreamUrl_Call) Run(run func(trackId string)) *ISubsonicClient_StreamUrl_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ISubsonicClient_StreamUrl_Call) Return(uRL *url.URL, err error) *ISubsonicClient_StreamUrl_Call {
	_c.Call.Return(uRL, err)
	return _c
}

func (_c *ISubsonicClient_StreamUrl_Call) RunAndReturn(run func(trackId string) (*url.URL, error)) *ISubsonicClient_StreamUrl_Call {
	_c.Call.Return(run)
	return _c
}
