// Package client provides a client from the hub. Clients are created by calling Dial
// and each client implements the text.Editor interface. A frame (graphical text box)
// must be provided to Dial, and updates the frame when a remote or the local client
// performs an editing operation.
//
// Note: This is not an ideal architecture and it will be modified when time permits
//
// The workflow will be changed such that a screen.Window can optionally be provided
// (instead of a frame) and the graphical synchronization point will rely on the shiny
// event pump (thereby freeing this package and ../hub from brokering graphical operations)
package client
