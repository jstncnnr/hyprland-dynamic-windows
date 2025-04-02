# hyprland-dynamic-windows
This script is used as the example in https://github.com/jstncnnr/go-hyprland

I use an ultrawide monitor, and hate that one window on a workspace will stretch across
the entire screen. I do not like the master layout, and thought the problem was easier to solve in
code.

## How it works
We listen for events from Hyprland, and anytime the below events are received, we count all
of the eligible windows on the active workspace. If there is only 1 we add reserved space on the
sides to effectively center the window, and if there is 0 or more than 1 we remove the reserved space.

I like to keep my special workspaces full width, so we remove the reserved space anytime the special workspace
is opened.

## Events Listened To
- Window opened/closed
- Window moved to workspace
- Workspace changed
- Special workspace changed
- Window floating mode changed
