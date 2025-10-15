package progress

// Forked from the original github.com/charmbracelet/bubbles/progress as
// my PR for the improvement is in review

import (
	"fmt"
	"image/color"
	"math"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ChausseBenjamin/termpicker/internal/colors"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/lucasb-eyer/go-colorful"
)

// ProgressFillFunc defines a function that returns which color
// to use for a given *current* part of the progress bar.
// totalCompletion: Overall completion of the progress bar (0.0 to 1.0)
// current: Position within the progress bar where the color is needed (0.0 to 1.0)
type ProgressFillFunc func(totalCompletion, current float64) color.Color

// Internal ID management. Used during animating to assure that frame messages
// can only be received by progress components that sent them.
var lastID int64

func nextID() int {
	return int(atomic.AddInt64(&lastID, 1))
}

const (
	fps              = 60
	defaultWidth     = 40
	defaultFrequency = 18.0
	defaultDamping   = 1.0
)

// FillStep is the only thing I added, you must refactor
// to make use of it
type FillStep struct {
	rune       rune
	completion float64 // 0% to 100% of that particular block
}

func defaultFillSteps() []FillStep {
	return []FillStep{
		{' ', 0.0},
		{'▏', 1.0 / 8.0},
		{'▎', 2.0 / 8.0},
		{'▍', 3.0 / 8.0},
		{'▌', 4.0 / 8.0},
		{'▋', 5.0 / 8.0},
		{'▊', 6.0 / 8.0},
		{'▉', 7.0 / 8.0},
		{'█', 1.0},
	}
}

// Option is used to set options in New. For example:
//
//	    progress := New(
//		       WithRamp("#ff0000", "#0000ff"),
//		       WithoutPercentage(),
//	    )
type Option func(*Model)

// WithDefaultGradient sets a gradient fill with default colors.
func WithDefaultGradient() Option {
	return WithGradientFunc(createLinearGradient("#5A56E0", "#EE6FF8", false))
}

// WithGradient sets a gradient fill blending between two colors.
func WithGradient(colorA, colorB string) Option {
	return WithGradientFunc(createLinearGradient(colorA, colorB, false))
}

// WithDefaultScaledGradient sets a gradient with default colors, and scales the
// gradient to fit the filled portion of the ramp.
func WithDefaultScaledGradient() Option {
	return WithGradientFunc(createLinearGradient("#5A56E0", "#EE6FF8", true))
}

// WithScaledGradient scales the gradient to fit the width of the filled portion of
// the progress bar.
func WithScaledGradient(colorA, colorB string) Option {
	return WithGradientFunc(createLinearGradient(colorA, colorB, true))
}

// WithGradientFunc sets a custom gradient function for the progress bar.
func WithGradientFunc(gradientFunc ProgressFillFunc) Option {
	return func(m *Model) {
		m.gradientFunc = gradientFunc
		m.useGradientFunc = true
	}
}

// WithStretchedGradient creates a linear gradient that stretches with progress.
// ColorB is reached at the totalCompletion mark (gradient stretches with progress).
func WithStretchedGradient(colorA, colorB string) Option {
	return WithGradientFunc(createStretchedLinearGradient(colorA, colorB))
}

// WithDefaultStretchedGradient creates a stretched gradient with default colors.
func WithDefaultStretchedGradient() Option {
	return WithStretchedGradient("#5A56E0", "#EE6FF8")
}

// WithHueGradient creates a hue-based gradient cycling through the color wheel.
// 0% -> 0° hue, 100% -> 360° hue. Saturation and lightness are fixed values (0.0-1.0).
func WithHueGradient(saturation, lightness float64) Option {
	return WithGradientFunc(createHueGradient(saturation, lightness))
}

// WithDefaultHueGradient creates a hue gradient with default saturation and lightness.
func WithDefaultHueGradient() Option {
	return WithHueGradient(0.8, 0.6) // High saturation, medium lightness
}

// WithDefaultOKLCHHueGradient creates a hue gradient in OKLCH color space with default chroma and lightness.
func WithDefaultOKLCHHueGradient() Option {
	return WithGradientFunc(createOKLCHHueGradient(0.3, 0.6)) // Moderate chroma, medium lightness
}

// WithSolidFill sets the progress to use a solid fill with the given color.
func WithSolidFill(color color.Color) Option {
	return func(m *Model) {
		m.gradientFunc = createSolidGradient(color)
		m.useGradientFunc = true
	}
}

// WithFillCharacters sets the characters used to construct the full and empty components of the progress bar.
func WithFillCharacters(steps []FillStep) Option {
	sort.Slice(steps, func(i, j int) bool {
		return steps[i].completion < steps[j].completion
	})
	return func(m *Model) {
		m.FillSteps = steps
	}
}

// WithBinaryFill results in a less granular but possible more widely compatible
// progress bar as only two characters are used to represent completion of a
// single block (full/complete and empty/incomplete).
func WithBinaryFill() Option {
	return func(m *Model) {
		m.FillSteps = []FillStep{
			{' ', 0.0},
			{'█', 1.0},
		}
	}
}

// WithoutPercentage hides the numeric percentage.
func WithoutPercentage() Option {
	return func(m *Model) {
		m.ShowPercentage = false
	}
}

// WithWidth sets the initial width of the progress bar. Note that you can also
// set the width via the SetWidth method, which can come in handy if you're
// waiting for a tea.WindowSizeMsg.
func WithWidth(w int) Option {
	return func(m *Model) {
		m.width = w
	}
}

// WithSpringOptions sets the initial frequency and damping options for the
// progress bar's built-in spring-based animation. Frequency corresponds to
// speed, and damping to bounciness. For details see:
//
// https://github.com/charmbracelet/harmonica
func WithSpringOptions(frequency, damping float64) Option {
	return func(m *Model) {
		m.SetSpringOptions(frequency, damping)
		m.springCustomized = true
	}
}

// FrameMsg indicates that an animation step should occur.
type FrameMsg struct {
	id  int
	tag int
}

// Model stores values we'll use when rendering the progress bar.
type Model struct {
	// An identifier to keep us from receiving messages intended for other
	// progress bars.
	id int

	// An identifier to keep us from receiving frame messages too quickly.
	tag int

	// Total width of the progress bar, including percentage, if set.
	width int

	FillSteps []FillStep

	// "Filled" sections of the progress bar.
	Full      rune
	FullColor color.Color

	// "Empty" sections of the progress bar.
	Empty      rune
	EmptyColor color.Color

	// Settings for rendering the numeric percentage.
	ShowPercentage  bool
	PercentFormat   string // a fmt string for a float
	PercentageStyle lipgloss.Style

	// Members for animated transitions.
	spring           harmonica.Spring
	springCustomized bool
	percentShown     float64 // percent currently displaying
	targetPercent    float64 // percent to which we're animating
	velocity         float64

	// Gradient settings (legacy)
	useRamp    bool
	rampColorA colorful.Color
	rampColorB colorful.Color

	// New gradient function
	useGradientFunc bool
	gradientFunc    ProgressFillFunc

	// When true, we scale the gradient to fit the width of the filled section
	// of the progress bar. When false, the width of the gradient will be set
	// to the full width of the progress bar.
	scaleRamp bool
}

// New returns a model with default values.
func New(opts ...Option) Model {
	m := Model{
		id:              nextID(),
		width:           defaultWidth,
		FillSteps:       defaultFillSteps(),
		Full:            '█',
		FullColor:       lipgloss.Color("#7571F9"),
		Empty:           '░',
		EmptyColor:      lipgloss.Color("#606060"),
		ShowPercentage:  true,
		PercentFormat:   " %3.0f%%",
		gradientFunc:    createSolidGradient(lipgloss.Color("#7571F9")),
		useGradientFunc: true,
	}

	for _, opt := range opts {
		opt(&m)
	}

	if !m.springCustomized {
		m.SetSpringOptions(defaultFrequency, defaultDamping)
	}

	return m
}

// NewModel returns a model with default values.
//
// Deprecated: use [New] instead.
var NewModel = New

// Init exists to satisfy the tea.Model interface.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update is used to animate the progress bar during transitions. Use
// SetPercent to create the command you'll need to trigger the animation.
//
// If you're rendering with ViewAs you won't need this.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case FrameMsg:
		if msg.id != m.id || msg.tag != m.tag {
			return m, nil
		}

		// If we've more or less reached equilibrium, stop updating.
		if !m.IsAnimating() {
			return m, nil
		}

		m.percentShown, m.velocity = m.spring.Update(m.percentShown, m.velocity, m.targetPercent)
		return m, m.nextFrame()

	default:
		return m, nil
	}
}

// SetSpringOptions sets the frequency and damping for the current spring.
// Frequency corresponds to speed, and damping to bounciness. For details see:
//
// https://github.com/charmbracelet/harmonica
func (m *Model) SetSpringOptions(frequency, damping float64) {
	m.spring = harmonica.NewSpring(harmonica.FPS(fps), frequency, damping)
}

// Percent returns the current visible percentage on the model. This is only
// relevant when you're animating the progress bar.
//
// If you're rendering with ViewAs you won't need this.
func (m Model) Percent() float64 {
	return m.targetPercent
}

// SetPercent sets the percentage state of the model as well as a command
// necessary for animating the progress bar to this new percentage.
//
// If you're rendering with ViewAs you won't need this.
func (m *Model) SetPercent(p float64) tea.Cmd {
	m.targetPercent = math.Max(0, math.Min(1, p))
	m.tag++
	return m.nextFrame()
}

// IncrPercent increments the percentage by a given amount, returning a command
// necessary to animate the progress bar to the new percentage.
//
// If you're rendering with ViewAs you won't need this.
func (m *Model) IncrPercent(v float64) tea.Cmd {
	return m.SetPercent(m.Percent() + v)
}

// DecrPercent decrements the percentage by a given amount, returning a command
// necessary to animate the progress bar to the new percentage.
//
// If you're rendering with ViewAs you won't need this.
func (m *Model) DecrPercent(v float64) tea.Cmd {
	return m.SetPercent(m.Percent() - v)
}

// View renders an animated progress bar in its current state. To render
// a static progress bar based on your own calculations use ViewAs instead.
func (m Model) View() string {
	return m.ViewAs(m.percentShown)
}

// ViewAs renders the progress bar with a given percentage.
func (m Model) ViewAs(percent float64) string {
	b := strings.Builder{}
	percentView := m.percentageView(percent)
	m.barView(&b, percent, ansi.StringWidth(percentView))
	b.WriteString(percentView)
	return b.String()
}

// SetWidth sets the width of the progress bar.
func (m *Model) SetWidth(w int) {
	m.width = w
}

// Width returns the width of the progress bar.
func (m Model) Width() int {
	return m.width
}

func (m *Model) nextFrame() tea.Cmd {
	return tea.Tick(time.Second/time.Duration(fps), func(time.Time) tea.Msg {
		return FrameMsg{id: m.id, tag: m.tag}
	})
}

func (m Model) barView(b *strings.Builder, percent float64, textWidth int) {
	var (
		tw = max(0, m.width-textWidth) // total width of the progress bar
		fw = percent * float64(tw)     // filled width in exact units
	)

	for i := range tw {
		cellPercent := float64(i) / float64(tw) // percentage of each cell
		if cellPercent < percent {
			// Filled cell: calculate the closest FillStep
			step := interpolateFillStep(m.FillSteps, fw-float64(i))

			// For full blocks, use half-block technique for smoother gradients
			if step.rune == '█' {
				// Use half-block character with foreground and background colors
				// for 2x resolution gradient effect
				leftPos := float64(i) / float64(tw)
				rightPos := float64(i) + 0.5
				if rightPos < float64(tw) {
					rightPos = rightPos / float64(tw)
				} else {
					rightPos = leftPos
				}

				var leftColor, rightColor color.Color
				if m.useGradientFunc {
					leftColor = m.gradientFunc(percent, leftPos)
					rightColor = m.gradientFunc(percent, rightPos)
				} else if m.useRamp {
					// Legacy gradient support
					p1 := leftPos
					p2 := rightPos
					if m.scaleRamp {
						p1 = leftPos
						p2 = rightPos
					}
					leftColor = m.rampColorA.BlendLuv(m.rampColorB, p1)
					rightColor = m.rampColorA.BlendLuv(m.rampColorB, p2)
				} else {
					// Legacy solid color support
					leftColor = m.FullColor
					rightColor = m.FullColor
				}

				// Use half-block character: foreground = left half, background = right half
				b.WriteString(lipgloss.NewStyle().Foreground(leftColor).Background(rightColor).Render("▌"))
			} else {
				// For edge blocks (using the eights), keep original logic
				if m.useGradientFunc {
					current := float64(i) / float64(tw)
					c := m.gradientFunc(percent, current)
					b.WriteString(lipgloss.NewStyle().Foreground(c).Render(string(step.rune)))
				} else if m.useRamp {
					// Legacy gradient support
					p := float64(i) / float64(tw-1)
					if m.scaleRamp {
						p = float64(i) / float64(tw-1)
					}
					c := m.rampColorA.BlendLuv(m.rampColorB, p)
					b.WriteString(lipgloss.NewStyle().Foreground(c).Render(string(step.rune)))
				} else {
					// Legacy solid color support
					b.WriteString(lipgloss.NewStyle().Foreground(m.FullColor).Render(string(step.rune)))
				}
			}
		} else {
			// Empty cell - always use static color, no gradient
			emptyStep := m.FillSteps[0]
			b.WriteString(lipgloss.NewStyle().Foreground(m.EmptyColor).Render(string(emptyStep.rune)))
		}
	}
}

// Helper: Interpolate between FillSteps
func interpolateFillStep(steps []FillStep, remaining float64) FillStep {
	for i := len(steps) - 1; i >= 0; i-- {
		if remaining >= steps[i].completion {
			return steps[i]
		}
	}
	return steps[0]
}

// createLinearGradient creates a linear gradient function between two colors.
func createLinearGradient(colorA, colorB string, scaled bool) ProgressFillFunc {
	a, _ := colorful.Hex(colorA)
	b, _ := colorful.Hex(colorB)

	return func(totalCompletion, current float64) color.Color {
		var p float64
		if scaled {
			// Scale gradient to fit only the filled portion
			if totalCompletion > 0 {
				p = current / totalCompletion
			} else {
				p = 0
			}
		} else {
			// Gradient spans the entire bar width
			p = current
		}
		p = math.Max(0, math.Min(1, p))
		return a.BlendLuv(b, p)
	}
}

// createStretchedLinearGradient creates a linear gradient where colorB is reached at totalCompletion.
// The gradient stretches with progress - when bar is 30% full, colorB is at the 30% mark.
func createStretchedLinearGradient(colorA, colorB string) ProgressFillFunc {
	a, _ := colorful.Hex(colorA)
	b, _ := colorful.Hex(colorB)

	return func(totalCompletion, current float64) color.Color {
		if totalCompletion <= 0 {
			return a
		}

		// Map current position to the stretched gradient
		// When totalCompletion=0.3, current=0.3 should give us colorB (p=1.0)
		p := current / totalCompletion
		p = math.Max(0, math.Min(1, p))
		return a.BlendLuv(b, p)
	}
}

// createHueGradient creates a hue-based gradient where colors cycle through the hue wheel.
// 0% completion -> 0° hue, 100% completion -> 360° hue
func createHueGradient(saturation, lightness float64) ProgressFillFunc {
	// Clamp saturation and lightness to valid ranges
	sat := math.Max(0, math.Min(1, saturation))
	light := math.Max(0, math.Min(1, lightness))

	return func(totalCompletion, current float64) color.Color {
		// Map current position to hue angle (0-360 degrees)
		hue := current * 360.0

		// Create HSL color and convert to RGB
		hslColor := colorful.Hsl(hue, sat, light)
		return hslColor
	}
}

// createOKLCHHueGradient creates a hue-based gradient in OKLCH color space.
// 0% completion -> 0° hue, 100% completion -> 360° hue
func createOKLCHHueGradient(chroma, lightness float64) ProgressFillFunc {
	// Clamp chroma and lightness to valid ranges
	c := math.Max(0, math.Min(0.5, chroma))
	l := math.Max(0, math.Min(1, lightness))

	return func(totalCompletion, current float64) color.Color {
		// Map current position to hue angle (0-360 degrees)
		hue := current * 360.0

		// Create OKLCH color and convert to RGB
		oklch := colors.OKLCH{L: l, C: c, H: hue}
		precise := oklch.ToPrecise()
		return color.RGBA{
			R: uint8(math.Round(precise.R * 255)),
			G: uint8(math.Round(precise.G * 255)),
			B: uint8(math.Round(precise.B * 255)),
			A: 255,
		}
	}
}

// createSolidGradient creates a gradient function that always returns the same color.
func createSolidGradient(c color.Color) ProgressFillFunc {
	return func(totalCompletion, current float64) color.Color {
		return c
	}
}

// CreateStripedLuminanceGradient creates a moving striped progress bar using HSL luminance variations.
// The stripes move as the progress advances, creating an animated effect.
// hue: base hue (0-360), saturation: base saturation (0-1), stripeWidth: width of each stripe (0-1)
func CreateStripedLuminanceGradient(hue, saturation, stripeWidth float64) ProgressFillFunc {
	// Clamp parameters to valid ranges
	h := math.Mod(hue, 360)
	if h < 0 {
		h += 360
	}
	sat := math.Max(0, math.Min(1, saturation))
	width := math.Max(0.01, math.Min(1, stripeWidth)) // Minimum stripe width to avoid division by zero

	return func(totalCompletion, current float64) color.Color {
		// Create moving stripe pattern based on current position and total completion
		// The totalCompletion acts as the "x offset" making stripes appear to move
		stripePosition := (current - totalCompletion) / width

		// Use modulo to create repeating pattern (0 to 1)
		stripe := math.Mod(stripePosition, 1.0)

		// Create luminance variation: oscillate between 0.2 and 0.8
		// Using sine wave for smooth transitions
		luminance := 0.5 + 0.3*math.Sin(stripe*2*math.Pi)

		// Clamp luminance to valid range
		luminance = math.Max(0.1, math.Min(0.9, luminance))

		// Create HSL color and return
		return colorful.Hsl(h, sat, luminance)
	}
}

// WithStripedLuminanceGradient creates a moving striped progress bar option.
// As progress advances, the stripes appear to move, creating a dynamic visual effect.
func WithStripedLuminanceGradient(hue, saturation, stripeWidth float64) Option {
	return WithGradientFunc(CreateStripedLuminanceGradient(hue, saturation, stripeWidth))
}

// WithDefaultStripedLuminanceGradient creates a striped gradient with default blue hue and moderate stripe width.
func WithDefaultStripedLuminanceGradient() Option {
	return WithStripedLuminanceGradient(240, 0.8, 0.1) // Blue hue, high saturation, narrow stripes
}

func (m Model) percentageView(percent float64) string {
	if !m.ShowPercentage {
		return ""
	}
	percent = math.Max(0, math.Min(1, percent))
	percentage := fmt.Sprintf(m.PercentFormat, percent*100) //nolint:mnd
	percentage = m.PercentageStyle.Inline(true).Render(percentage)
	return percentage
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// IsAnimating returns false if the progress bar reached equilibrium and is no longer animating.
func (m *Model) IsAnimating() bool {
	dist := math.Abs(m.percentShown - m.targetPercent)
	return !(dist < 0.001 && m.velocity < 0.01)
}
