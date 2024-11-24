# Termpicker

A simple Color Picker Designed for your Terminal

This is very much a work in progress, but the end goal is to be able to
generate and copy colors from the terminal.

Here is a quick demo of what has been done so far:
<div align="center">
  <img src="./assets/demo.gif" width="600" alt="Termpicker Demo"><br>
</div>

Here is my roadmap to reach what I would consider a finished state:

- [ ] Add a "cmd" mode to manually input colors
- [ ] Allow to pass a starting color as an argument when launching the program
- [ ] Make the tabs interface prettier with [lipgloss][1] (similar to tabs in [soft-serve][2])
- [ ] Notify user of successful copy to clipboard (or failure)
- [ ] Unit-test color conversions near edge case colors
- [ ] Warn the user if the terminal is too small (and refuse to render)
- [x] Add a [help bubble][3] at the bottom of the interface to show available keybindings
- [x] Add Box-drawing to the picker and the previewer
- [x] Implement copying to clipboard for various formats (rgb, hex, hsl, cymk, etc...)
- [X] Make sliders reach the correct length on init/tab without pressing `j`,`k`
- [x] Make the preview windows prettier (perhaps same width as the sliders)


[1]: https://github.com/charmbracelet/lipgloss
[2]: https://github.com/charmbracelet/soft-serve
[3]: https://github.com/charmbracelet/bubbles#help
