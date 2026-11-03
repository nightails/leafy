# Terminal UI

Code: `app/tui` (Bubble Tea model) and `internal/style` (Lip Gloss styles and helpers).

Libraries:
- [`bubbletea`](https://github.com/charmbracelet/bubbletea) v1.3: the program loop (Model / Update / View)
- [`bubbles`](https://github.com/charmbracelet/bubbles) v0.21: the `list` and `textinput` components
- [`lipgloss`](https://github.com/charmbracelet/lipgloss) v1.1: styling

## Screen layout

`View()` builds the screen top to bottom. Sections pile up as the workflow moves forward:

```
┌ HeaderStyle ─────────────┐
│ Leafy <version>          │      always
└──────────────────────────┘

Found media:                      step >= media
<media checklist>

Destination Path:                 step >= destination
<text input>

Copying Progress:                 step >= copying
<name copied/total (NN%)> ...

<help line>                       always (HelpTextStyle)
```

- If `m.err` is set, only the header and the error text are shown, with no help line.
- The help line reads `[j/k] up/down | [space] select | [enter] confirm | [ctrl+c] quit`. While copying, it changes
  to `[copying in progress: quitting disabled]`.
- When the transfer completes or fails, the program quits on its own. The last frame (completed progress lines or the
  error text) stays in the terminal.

## Media checklist

- Built on `list.Model` with a custom `mediaItemDelegate`: one line per item, no spacing.
- Each row is rendered as `"> "` (cursor only) + `"[x] "` or `"[ ] "` + `name (size)`.
- The title, status bar, pagination, built-in help, filtering and quit bindings are all turned off.
- The list height equals the number of items, so every file is visible. The width follows the terminal width
  (`tea.WindowSizeMsg`).
- `mediumItem.FilterValue()` returns the file name. It would be used if filtering were turned on.

## Destination input

- `textinput.Model` with the placeholder `Destination Path`.
- It is focused when you enter the destination step and blurred when copying starts.

## Styles (`internal/style/style.go`)

| Name            | Definition                                         | Used for              |
|-----------------|----------------------------------------------------|-----------------------|
| `HeaderStyle`   | ASCII border, foreground color `255`, bold         | `Leafy <version>` box |
| `TextStyle`     | Foreground `255` (white)                           | Main body             |
| `HelpTextStyle` | Foreground `244` (grey)                            | Help line             |

## Helpers not yet used (`internal/style`)

| Helper                    | Purpose                                                                                          |
|---------------------------|--------------------------------------------------------------------------------------------------|
| `NewLineSpinner()`        | Returns a Bubbles `spinner.Model` using the `Line` animation, meant for a loading indicator      |
| `MinDuration{Min}`        | `StartNow()`, `EnsureStarted()`, `Remaining()`: keeps a loading state visible for a minimum time  |
| `LoadDelay`, `QuitDelay`  | 1 s constants for the intended minimum loading and quit screens                                  |

Other unused assets:
- `assets/leafy_ascii_art.txt`: a boxed ASCII-art "LEAFY" logo.
- `assets/tui_layout.txt`: an empty placeholder for a future layout design.
- `deleteDoneMsg` and the `Selected media has been deleted.` text in `View()` hint at a planned "delete source after
  copy" feature. Nothing sends the message yet.
