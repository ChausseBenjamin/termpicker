# Termpicker

A simple Color Picker Designed for your Terminal

This is very much a work in progress, but the end goal is to be able to
generate and copy colors from the terminal.

Here is a quick demo of what has been done so far:
<div align="center">
  <img src="./assets/demo.gif" width="600" alt="Termpicker Demo"><br>
</div>

Here is my roadmap to reach what I would consider a finished state:

- [ ] Implement copying to clipboard for various formats (rgb, hex, hsl, cymk, etc...)
- [ ] Make the tabs interface prettier with [lipgloss][1] (similar to tabs in [soft-serve][2])
- [ ] Add a [help bubble][3] at the bottom of the interface to show available keybindings
- [ ] Add some form of stdout cli flag to output to stdout instead of copying colors
- [ ] Auto-adjust geometry on terminal resize (+ warn the user if the terminal is too small)
- [ ] Make the preview windows prettier (perhaps same width as the sliders)
- [ ] Add Box-drawing to the picker and the previewer
- [ ] Add more color conversion unit-tests around edge case colors
- [X] Make sliders reach the correct length on init/tab without pressing `j`,`k`


[1]: https://github.com/charmbracelet/lipgloss
[2]: https://github.com/charmbracelet/soft-serve
[3]: https://github.com/charmbracelet/bubbles#help
