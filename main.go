package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jstncnnr/go-hyprland/hypr"
	"github.com/jstncnnr/go-hyprland/hypr/event"
	"os"
	"os/signal"
	"slices"
	"syscall"
)

func main() {
	client, err := events.NewClient()
	if err != nil {
		fmt.Printf("Error creating event client: %v\n", err)
		os.Exit(1)
	}

	client.RegisterListener(EventListener)

	// Setup interrupt handler so we can cleanly close the event client
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		interrupt, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		<-interrupt.Done()

		cancel()
	}()

	if err := client.Listen(ctx); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Printf("Error running event client: %v\n", err)
		os.Exit(1)
	}
}

func EventListener(event events.Event) {
	switch event.(type) {
	case events.OpenWindowEvent,
		events.CloseWindowEvent,
		events.MoveWindowEvent,
		events.MoveWindowV2Event,
		events.ChangeFloatingModeEvent,
		events.WorkspaceEvent,
		events.WorkspaceV2Event:
		CheckWorkspace()
		break

	case events.ActiveSpecialEvent:
		event := event.(events.ActiveSpecialEvent)
		if event.WorkspaceName == "" {
			// WorkspaceName is empty when we close a special workspace
			CheckWorkspace()
		} else {
			// Keep the special workspace full width, always
			_ = RemoveReservedSpace(event.MonitorName)
		}
		break

	case events.ActiveSpecialV2Event:
		event := event.(events.ActiveSpecialV2Event)
		if event.WorkspaceName == "" {
			// WorkspaceName is empty when we close a special workspace
			CheckWorkspace()
		} else {
			// Keep the special workspace full width, always
			_ = RemoveReservedSpace(event.MonitorName)
		}
		break
	}
}

func CheckWorkspace() {
	workspace, err := hypr.GetActiveWorkspace()
	if err != nil {
		fmt.Printf("Error getting active workspace: %v\n", err)
		return
	}

	windows, err := hypr.GetWindows()
	if err != nil {
		fmt.Printf("Error getting windows: %v\n", err)
		return
	}

	// Get the number of non-floating windows in this workspace
	windows = slices.DeleteFunc(windows, func(window hypr.Window) bool {
		return window.Workspace.Id != workspace.Id || window.Floating
	})

	if len(windows) == 1 {
		//If there is only 1 window, add reserved space on the left and right side
		//this is useful for ultrawide monitors so you don't have one stretched window
		_ = AddReservedSpace(workspace.Monitor, 0, 0, 865, 865)
	} else {
		_ = RemoveReservedSpace(workspace.Monitor)
	}
}

func AddReservedSpace(monitor string, top, bottom, left, right int) error {
	return hypr.Keyword(fmt.Sprintf("monitor %s,addreserved,%d,%d,%d,%d", monitor, top, bottom, left, right))
}

func RemoveReservedSpace(monitor string) error {
	return hypr.Keyword(fmt.Sprintf("monitor %s,addreserved,0,0,0,0", monitor))
}
