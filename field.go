package walk

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type field struct {
	*Model

	value *FilePath
	key   string

	// error handling
	validate func(FilePath) error
	err      error

	// state
	isFocused  bool
	isFiltered bool
	filter     textinput.Model

	// customization
	heading string
	caption string

	// options
	accessible bool
	showAll    bool
	theme      *huh.Theme
}

// Theme returns an Option that sets the theme of a field.
func Theme(theme *huh.Theme) Option[*field] {
	return func(f *field) *field { return f.WithTheme(theme).(*field) }
}

// Accessible returns an Option that sets the accessibility flag of a field.
func Accessible(accessible bool) Option[*field] {
	return func(f *field) *field { return f.WithAccessible(accessible).(*field) }
}

// KeyMap returns an Option that sets the key bindings of a field.
func KeyMap(keys *huh.KeyMap) Option[*field] {
	return func(f *field) *field { return f.WithKeyMap(keys).(*field) }
}

// Width returns an Option that sets the width of a field.
func Width(width int) Option[*field] {
	return func(f *field) *field { return f.WithWidth(width).(*field) }
}

// Height returns an Option that sets the height of a field.
func Height(height int) Option[*field] {
	return func(f *field) *field { return f.WithHeight(height).(*field) }
}

// Value returns an Option that sets the value of a field.
func Value(value string) Option[*field] {
	return func(f *field) *field { return f.WithValue(value).(*field) }
}

// Key returns an Option that sets the key of a field.
func Key(key string) Option[*field] {
	return func(f *field) *field { return f.WithKey(key).(*field) }
}

// Heading returns an Option that sets the heading of a field.
func Heading(heading string) Option[*field] {
	return func(f *field) *field { return f.WithHeading(heading).(*field) }
}

// Caption returns an Option that sets the caption of a field.
func Caption(caption string) Option[*field] {
	return func(f *field) *field { return f.WithCaption(caption).(*field) }
}

// Validate returns an Option that sets the validation function of a field.
func Validate(validate func(FilePath) error) Option[*field] {
	return func(f *field) *field { return f.WithValidate(validate).(*field) }
}

// Prompt returns an Option that sets the prompt of a field.
func Prompt(prompt string) Option[*field] {
	return func(f *field) *field { return f.WithPrompt(prompt).(*field) }
}

// Init initializes the internal state of a field.
func (f *field) Init() tea.Cmd {
	return f.Model.Init()
}

// Update processes and manages the internal state of a field.
func (f *field) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	_, cmd := f.Model.Update(msg)
	return f, cmd
}

// View renders the field according to its internal state.
func (f *field) View() string {
	return f.Model.View()
}

// Blur blurs the field.
func (f *field) Blur() tea.Cmd {
	f.isFocused = false
	f.err = f.validate(*f.value)
	return nil
}

// Focus focuses the field.
func (f *field) Focus() tea.Cmd {
	f.isFocused = true
	return nil
}

// Skip returns whether the textarea should be skipped or should be blocking.
func (f *field) Skip() bool {
	return false
}

// Zoom returns whether the note should be zoomed.
func (f *field) Zoom() bool {
	return false
}

// Error returns the error of the field.
func (f *field) Error() error {
	return f.err
}

// Run runs the field.
func (f *field) Run() error {
	if f.accessible {
		return f.runAccessible()
	}
	return newRunError(huh.Run(f))
}

// KeyBinds returns the keybindings for the field.
func (f *field) KeyBinds() []key.Binding {
	return []key.Binding{} // f.keys.bindings()
}

// HelpKeyBinds returns the help keybindings for the field.
func (f *field) HelpKeyBinds(huh.HelpFormat) []key.Binding {
	return []key.Binding{}
}

// With returns the receiver with the given options applied.
func (f *field) With(options ...Option[*field]) huh.Field {
	for _, option := range options {
		f = option(f)
	}
	return f
}

// WithTheme sets the theme of the field.
func (f *field) WithTheme(theme *huh.Theme) huh.Field {
	f.theme = theme
	f.filter.Cursor.Style = f.theme.Focused.TextInput.Cursor
	f.filter.PromptStyle = f.theme.Focused.TextInput.Prompt
	return f
}

// WithAccessible sets the accessible mode of the field.
func (f *field) WithAccessible(accessible bool) huh.Field {
	f.accessible = accessible
	return f
}

// WithKeyMap sets the keymap on a field.
func (f *field) WithKeyMap(keys *huh.KeyMap) huh.Field {
	// TBD
	return f
}

// WithWidth sets the width of the field.
func (f *field) WithWidth(width int) huh.Field {
	f.width = width
	return f
}

// WithHeight sets the height of the field.
func (f *field) WithHeight(height int) huh.Field {
	f.height = height
	return f
}

// GetKey returns the key of the field.
func (f *field) GetKey() string {
	return f.key
}

// GetValue returns the value of the field.
func (f *field) GetValue() any {
	return f.value.path()
}

// WithValue sets the value of the field.
func (f *field) WithValue(value string) huh.Field {
	f.value = f.value.init(value)
	return f
}

// WithKey sets the key of the field.
func (f *field) WithKey(key string) huh.Field {
	f.key = key
	return f
}

// WithHeading sets the heading of the field.
func (f *field) WithHeading(heading string) huh.Field {
	f.heading = heading
	return f
}

// WithCaption sets the caption of the field.
func (f *field) WithCaption(caption string) huh.Field {
	f.caption = caption
	return f
}

// WithValidate sets the validation function of the field.
func (f *field) WithValidate(validate func(FilePath) error) huh.Field {
	f.validate = validate
	return f
}

// WithPrompt sets the prompt of the field.
func (f *field) WithPrompt(prompt string) huh.Field {
	f.filter.Prompt = prompt
	return f
}

// WithPosition sets the position information of the text field.
func (f *field) WithPosition(p huh.FieldPosition) huh.Field {
        f.keys.Prev.SetEnabled(!p.IsFirst())
        f.keys.Next.SetEnabled(!p.IsLast())
        f.keys.Submit.SetEnabled(p.IsLast())
        return f
}

func (f *field) runAccessible() error {
	var sb strings.Builder
	sb.WriteString(f.theme.Focused.Title.Render(f.heading) + "\n")

	// for i, option := range t.option {
	// 	sb.WriteString(fmt.Sprintf("%d. %s", i+1, option.Key))
	// 	sb.WriteString("\n")
	// }
	//
	// fmt.Println(t.theme.Blurred.Base.Render(sb.String()))
	//
	// for {
	//	choice := accessibility.PromptInt("Choose: ", 1, len(t.options))
	//	option := t.options[choice-1]
	//	if err := t.validate(option.Value); err != nil {
	//		fmt.Println(err.Error())
	//		continue
	//	}
	//	fmt.Println(t.theme.Focused.SelectedOption.Render("Chose: " + option.Key + "\n"))
	//	*t.value = option.Value
	//	break
	//}

	return nil
}

func (f *field) setIsFiltered(isFiltered bool) {
	f.isFiltered = isFiltered
}

func (f *field) filterFunc(option string) bool {
	// XXX: remove diacritics or allow customization of filter function.
	return strings.Contains(
		strings.ToLower(option),
		strings.ToLower(f.filter.Value()),
	)
}
