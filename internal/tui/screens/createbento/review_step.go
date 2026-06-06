package createbento

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/ui/components/layout"
)

type saveBentoMsg struct {
	path string
	err  error
}

func (s *Screen) updateReviewStep(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	if s.saveSuccessPath != "" {
		switch msg.String() {
		case "enter", "space":
			return s, func() tea.Msg { return nav.SelectBento(s.draft) }
		}
		return s, nil
	}

	switch msg.String() {
	case "left", "esc":
		return s.previousStep()
	case "1":
		s.editStep(stepMetadata)
	case "2":
		s.editStep(stepKafka)
	case "3":
		s.editStep(stepRules)
	case "s":
		return s, s.saveBentoCmd()
	}
	return s, nil
}

func (s *Screen) viewReviewStep() string {
	yamlPreview, validation := s.reviewYAML()
	actions := strings.Join([]string{
		s.styles.Title.Render("Edit shortcuts"),
		s.styles.Subtle.Render("1 metadata"),
		s.styles.Subtle.Render("2 kafka connection"),
		s.styles.Subtle.Render("3 bento rules"),
		s.styles.Subtle.Render("s save"),
	}, "\n")

	sections := []string{
		s.layout.Panel(layout.PanelOptions{
			Body: s.notice.Render(
				"Step 4/4: Review Bento YAML",
				"Check the generated config. Jump directly to any step if something needs fixing.",
			),
			BodyAlign: layout.AlignLeft,
		}),
		s.layout.Panel(layout.PanelOptions{
			Title:     "Generated Bento",
			Subtitle:  "This is the YAML Koicha will save later.",
			Body:      s.styles.Subtle.Render(yamlPreview),
			BodyAlign: layout.AlignLeft,
		}),
		s.layout.Panel(layout.PanelOptions{
			Title:     "Actions",
			Body:      actions,
			BodyAlign: layout.AlignLeft,
			Active:    s.saveSuccessPath == "",
		}),
	}

	if validation != "" {
		sections = append(sections, s.layout.Panel(layout.PanelOptions{
			Title:     "Needs attention",
			Body:      s.styles.Danger.Render(validation),
			BodyAlign: layout.AlignLeft,
		}))
	}
	if s.status != "" {
		style := s.styles.NoticeBody
		if s.statusError {
			style = s.styles.Danger
		}
		sections = append(sections, s.layout.Panel(layout.PanelOptions{
			Body:      style.Render(s.status),
			BodyAlign: layout.AlignLeft,
		}))
	}
	if s.saveSuccessPath != "" {
		sections = append(sections, s.renderSaveSuccessPopup())
	}

	body := s.layout.Stack(sections...)
	help := s.statusbar.Help("1/2/3 edit", "s save", "left/esc back")
	if s.saveSuccessPath != "" {
		help = s.statusbar.Help("enter/space ok")
	}
	return s.layout.Render("", body, help)
}

func (s *Screen) saveBentoCmd() tea.Cmd {
	draft := s.draft
	store := s.store
	return func() tea.Msg {
		path, err := store.Save(draft)
		return saveBentoMsg{path: path, err: err}
	}
}

func (s *Screen) handleSaveBentoMsg(msg saveBentoMsg) (nav.Screen, tea.Cmd) {
	if msg.err != nil {
		s.status = fmt.Sprintf("Could not save Bento: %v", msg.err)
		s.statusError = true
		return s, nil
	}
	s.saveSuccessPath = msg.path
	s.status = ""
	s.statusError = false
	return s, nil
}

func (s *Screen) renderSaveSuccessPopup() string {
	body := strings.Join([]string{
		s.styles.Subtle.Render("File saved to:"),
		s.styles.Error.Render(s.saveSuccessPath),
		"",
		s.layout.Center(s.styles.Error.Render("[ ok ]")),
	}, "\n")

	return s.layout.Panel(layout.PanelOptions{
		Title:     "Bento saved",
		Body:      body,
		BodyAlign: layout.AlignLeft,
		Active:    true,
	})
}

func (s *Screen) reviewYAML() (string, string) {
	if err := bento.Validate(s.draft); err != nil {
		return s.marshalDraft(), err.Error()
	}
	return s.marshalDraft(), ""
}

func (s *Screen) marshalDraft() string {
	data, err := bento.Marshal(s.draft)
	if err != nil {
		return "failed to render bento YAML: " + err.Error()
	}
	return strings.TrimSpace(string(data))
}
