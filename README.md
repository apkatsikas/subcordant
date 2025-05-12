# subcordant

TODO
ffmpeg will be a dependency, as we need it to create opus stream:
https://github.com/Gimzie/submeister/blob/015218a906599f9abe208f7cd6685b8209147f4d/player.py#L69

where we provide the subsonic stream to the audio source for ffmpeg to stream (streaming the stream)
compare to golang stream:
https://github.com/diamondburned/arikawa/blob/8a78eb04430cfd0f4997c8bf206cf36c0c2e604d/0-examples/voice/main.go#L75

architechture

main (main.go), app, discord, subsonic, ffmpeg

discord is entry point (via commands), but all components are started in main

discord has app as dependency (injected), so does app have discord... circular dep?

in my case, my discord client is the entry point of my application. its where we get commands from users via discord. i want these commands to get piped to an Application component, which in turn can call my subsonic client/component. i want the application to be the orchestrator, but it depends on being able to check the state of the discord client, but the discord client needs to call the application when it gets a command
ChatGPT said:

Your scenario involves a natural circular dependency: the Discord client must send commands to the Application, and the Application needs to check the state of the Discord client. This can be resolved by using patterns that avoid direct circular dependencies. Here's a clean way to structure this:

Solution: Use Interfaces and Inversion of Control

    Define Interfaces for Dependencies:

        Create interfaces for the functionalities each component needs from the other.

        The Discord client and the Application interact via these interfaces, breaking the direct dependency.

    Example:

// ICommandHandler defines what the Application can do with commands
type ICommandHandler interface {
    HandleCommand(command string)
}

// IDiscordClient defines what the Application can query about the Discord client
type IDiscordClient interface {
    IsConnected() bool
}

Implement the Interfaces:

    The DiscordClient implements IDiscordClient for the Application to query its state.

    The Application implements ICommandHandler to process commands sent from the Discord client.

type DiscordClient struct {
    CommandHandler ICommandHandler
}

func (d *DiscordClient) IsConnected() bool {
    // Check if the Discord client is connected
    return true
}

func (d *DiscordClient) ReceiveCommand(command string) {
    if d.CommandHandler != nil {
        d.CommandHandler.HandleCommand(command)
    }
}

type Application struct {
    DiscordClient IDiscordClient
}

func (a *Application) HandleCommand(command string) {
    // Handle the command
    fmt.Printf("Processing command: %s\n", command)
}

Wire Everything Together in the Entry Point:

    Use the main function to wire the dependencies.

    Pass the Discord client to the Application and set the CommandHandler of the Discord client.

    package main

    import "fmt"

    func main() {
        app := &Application{}
        discord := &DiscordClient{}

        // Set up mutual dependencies via interfaces
        app.DiscordClient = discord
        discord.CommandHandler = app

        // Simulate a command coming from Discord
        discord.ReceiveCommand("play music")
    }

Advantages of This Approach

    Decoupling:

        The Application and DiscordClient depend only on interfaces, not on concrete implementations. This avoids circular dependencies.

    Testability:

        You can easily mock the ICommandHandler or IDiscordClient in tests.

    Flexibility:

        You can replace the Discord client or application logic without impacting the other.

Alternative: Use an Event-Driven Architecture

If the interactions between components become more complex, consider using an event-driven pattern where:

    The Discord client emits events (e.g., OnCommandReceived).

    The Application subscribes to these events and orchestrates the logic.

This decouples the components further and makes it easier to add new features.

Let me know if you'd like to explore this further or need a more tailored example!
