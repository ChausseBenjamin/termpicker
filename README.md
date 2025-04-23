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
- Copy the color to your clipboard in various formats ([RGB][4], [HEX][5], [HSL][6], [CMYK][7], [ANSI truecolor][8])

## Usage:

The keybindings are pretty simple and shown in the UI. Their description can
be expanded by pressing the `?` key. The exhaustive list is available if you
RTFM :P (either `termpicker --help` or `man termpicker` if the manpage is
installed)

## Installation

**Via Go**:

```sh
go install github.com/ChausseBenjamin/termpicker@latest
```

**From the aur**:

Termpicker is on the AUR! :tada: Just install it with you favourite
aur package manager (yay, paru, yaourt, etc...)

```sh
yay -S termpicker
```

**Manual Installation**:

Just grab the latest release for your platform and install the binary
somewhere in your `PATH`. Releases also include a manpage which you can
install to your `$XDG_DATA_HOME/man/man1/`.

## Roadmap

- [ ] Publish release to more mainstream repositories (Homebrew, nix, etc...)
- [ ] Unit-test color conversions near edge case colors
- [ ] Migrate to bubbletea/V2 once it comes out of beta
- [X] Warn the user if the terminal is too small (and refuse to render)

[1]: https://github.com/charmbracelet/lipgloss
[2]: https://github.com/charmbracelet/soft-serve
[3]: https://github.com/charmbracelet/bubbles#help
[4]: https://www.w3schools.com/colors/colors_rgb.asp
[5]: https://www.w3schools.com/colors/colors_hexadecimal.asp
[6]: https://www.w3schools.com/colors/colors_hsl.asp
[7]: https://www.w3schools.com/colors/colors_cmyk.asp
[8]: https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797#rgb-colors
