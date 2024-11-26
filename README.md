# Termpicker

A simple Color Picker Designed for your Terminal

Here is a quick demo of what has been done so far:
<div align="center">
  <img src="./assets/demo.gif" width="600" alt="Termpicker Demo"><br>
</div>

## Features:

- Preview any color using a truecolor terminal
- Create colors using sliders for RGB, HSL, and CMYK
- Seamlessly convert between color formats (RGB, HSL, CMYK) as you create
- Copy the color to your clipboard in various formats (RGB, HEX, HSL, CMYK)

## Usage:

The keybindings are pretty simple and shown in the UI. For a more exhaustive
list, pressing `?` expands the help section to show all available keybindings.

## Installation

Just grab the latest release for your platform and install the binary
somewhere in your `PATH`.

Alternatively, you can install it directly from go with:
```sh
go install github.com/ChausseBenjamin/termpicker@latest
```

## Roadmap

- [ ] Refactor the code to streamline and centralize lipgloss styles
- [ ] Unit-test color conversions near edge case colors
- [ ] Warn the user if the terminal is too small (and refuse to render)
- [ ] Publish release to mainstream repositories (AUR, Homebrew, etc...)
- [x] Add a "cmd" mode to manually input colors
- [x] Add a [help bubble][3] at the bottom of the interface to show available keybindings
- [x] Add an input flag to pass specific color as a starting value
- [x] Add Box-drawing to the picker and the previewer
- [x] Implement copying to clipboard for various formats (rgb, hex, hsl, cymk, etc...)
- [X] Make sliders reach the correct length on init/tab without pressing `j`,`k`
- [x] Make the preview windows prettier (perhaps same width as the sliders)
- [x] Make the tabs interface prettier with [lipgloss][1] (similar to tabs in [soft-serve][2])
- [x] Notify user of successful copy to clipboard (or failure)

[1]: https://github.com/charmbracelet/lipgloss
[2]: https://github.com/charmbracelet/soft-serve
[3]: https://github.com/charmbracelet/bubbles#help
