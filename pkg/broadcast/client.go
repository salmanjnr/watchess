// Client defines a subscriber to a certain "broker". A broker organizes the communication between round struct and clients.
package broadcast

// A client defines an interface of a broker subscriber
type Client interface {
	Updates() chan<- GameUpdate
	Close()
}

// A subscriber of a game broker
type GameClient struct {
	broker  *gameBroker
	updates chan GameUpdate
	done    chan struct{}
}

// A subscriber of a round broker
type RoundClient struct {
	broker  *roundBroker
	updates chan GameUpdate
	done    chan struct{}
}

// Get channel over which relevant updates should be sent. Updates channel is only used by game broker to send and client to receive
func (c *GameClient) Updates() chan<- GameUpdate {
	return c.updates
}

// Close client
func (c *GameClient) Close() {
	close(c.updates)
	c.done <- struct{}{}
}

func (c *RoundClient) Updates() chan<- GameUpdate {
	return c.updates
}

func (c *RoundClient) Close() {
	close(c.updates)
	c.done <- struct{}{}
}
