# KEYBINDINGS

**Normal mode**:

- `h`,`l`: coarse decrease/increase the current slider by 5%
- `H`,`L`: fine decrease/increase the current slider by 1
- `j`,`k`: select the slider below/above
- `<Tab>`,`<S-Tab>`: move to the next/previous tab
- `f`,`b`: copy the color as an ANSI escape code for the foreground/background
- `x`,`r`,`s`,`c`: copy the color as a hex/rgb/hsl/cmyk
- `?`: expand/shrink the help menu
- `i`,`<cmd>`: enter Insert mode
- `q`/`<C-c>`: quit the application

**Insert mode**:

Manually type a color. Pressing <Esc> will cancel/leave insert mode.
Anything in the following formats will be used as a color input when
pressing enter:

- Hex values: `#rrggbb`
- RGB values: `rgb( r, g, b)`
- CMYK values: `cmyk(c, m, y, k)`
- HSL values: `hsl(h, s, l)`
