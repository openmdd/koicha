package createbento

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/ui/components/layout"
)

func (s *Screen) updateMetadataStep(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	switch msg.String() {
	case "left", "esc":
		return s.previousStep()
	case "shift+tab", "up":
		s.moveMetadataFocus(-1)
		return s, nil
	case "down":
		s.moveMetadataFocus(1)
		return s, nil
	case "tab", "enter":
		if !s.advanceMetadata() {
			return s, nil
		}
		return s, nil
	}

	var cmd tea.Cmd
	s.metadataErrors[s.metadataFocus] = ""
	s.metadataInputs[s.metadataFocus], cmd = s.metadataInputs[s.metadataFocus].Update(msg)
	return s, cmd
}

func (s *Screen) viewMetadataStep() string {
	required := s.layout.Panel(layout.PanelOptions{
		Title:     "Step 1/4: Required metadata",
		Subtitle:  "These values identify the Bento profile.",
		Body:      strings.Join([]string{s.renderMetadataField(metadataFieldName, "Bento name", true), "", s.renderMetadataField(metadataFieldProfile, "Profile", true)}, "\n"),
		BodyAlign: layout.AlignLeft,
		Active:    s.metadataFocus == metadataFieldName || s.metadataFocus == metadataFieldProfile,
	})

	optional := s.layout.Panel(layout.PanelOptions{
		Title:     "Optional metadata",
		Subtitle:  "You can leave display name empty for now.",
		Body:      s.renderMetadataField(metadataFieldDisplayName, "Display name", false),
		BodyAlign: layout.AlignLeft,
		Active:    s.metadataFocus == metadataFieldDisplayName,
	})

	body := s.layout.Stack(
		s.layout.Panel(layout.PanelOptions{
			Body: s.notice.Render(
				"Create Bento",
				"Start with the metadata Koicha needs. Kafka connection and Bento rules come next.",
			),
			BodyAlign: layout.AlignLeft,
		}),
		required,
		optional,
	)
	help := s.statusbar.Help("tab/enter next", "shift+tab previous", "left/esc back")
	if s.editMode {
		help = s.statusbar.Help("tab/enter done", "shift+tab previous", "left/esc review")
	}
	return s.layout.Render("", body, help)
}

func (s *Screen) renderMetadataField(index int, label string, required bool) string {
	parts := []string{s.fieldLabel(label, required)}
	if s.metadataErrors[index] != "" {
		parts = append(parts, s.styles.Danger.Render(s.metadataErrors[index]))
	}
	parts = append(parts, s.metadataInputs[index].View())
	return strings.Join(parts, "\n")
}

func (s *Screen) advanceMetadata() bool {
	if !s.validateMetadataFocused() {
		return false
	}
	if s.metadataFocus < metadataFieldCount-1 {
		s.moveMetadataFocus(1)
		return true
	}

	s.commitMetadata()
	s.completeStep(stepKafka)
	return true
}

func (s *Screen) validateMetadataFocused() bool {
	if !isRequiredMetadataField(s.metadataFocus) {
		return true
	}

	s.metadataTouched[s.metadataFocus] = true
	s.metadataErrors[s.metadataFocus] = ""

	value := strings.TrimSpace(s.metadataInputs[s.metadataFocus].Value())
	if value == "" {
		s.metadataErrors[s.metadataFocus] = "Required value. Fill it in to continue."
		return false
	}
	if s.metadataFocus != metadataFieldName {
		return true
	}
	if err := bento.ValidateName(value); err != nil {
		s.metadataErrors[s.metadataFocus] = "Use lowercase letters, digits, and hyphens only."
		return false
	}

	exists, err := s.store.Exists(value)
	if err != nil {
		s.metadataErrors[s.metadataFocus] = "Could not check this Bento name. Try again."
		return false
	}
	if exists {
		s.metadataErrors[s.metadataFocus] = "Bento with this name already exists. Choose another name."
		return false
	}
	return true
}

func (s *Screen) isMetadataInvalid(index int) bool {
	return s.metadataErrors[index] != ""
}

func (s *Screen) moveMetadataFocus(delta int) {
	next := s.metadataFocus + delta
	if next < 0 {
		next = 0
	}
	if next >= metadataFieldCount {
		next = metadataFieldCount - 1
	}
	s.focusMetadata(next)
}

func (s *Screen) focusMetadata(index int) {
	s.blurAllInputs()
	s.metadataFocus = index
	s.metadataInputs[index].Focus()
}

func (s *Screen) commitMetadata() {
	s.draft.SchemaVersion = bento.CurrentSchemaVersion
	s.draft.Metadata.Name = strings.TrimSpace(s.metadataInputs[metadataFieldName].Value())
	s.draft.Metadata.DisplayName = strings.TrimSpace(s.metadataInputs[metadataFieldDisplayName].Value())
	s.draft.Spec.Profile.Environment = strings.TrimSpace(s.metadataInputs[metadataFieldProfile].Value())
}

func isRequiredMetadataField(index int) bool {
	return index == metadataFieldName || index == metadataFieldProfile
}
