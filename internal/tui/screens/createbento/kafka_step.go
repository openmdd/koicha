package createbento

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/kafka"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/ui/components/layout"
)

func (s *Screen) updateKafkaStep(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	switch msg.String() {
	case "left", "esc":
		return s.previousStep()
	case "shift+tab", "up":
		s.moveKafkaFocus(-1)
		return s, nil
	case "down":
		s.moveKafkaFocus(1)
		return s, nil
	case "tab", "enter":
		if !s.advanceKafka() {
			return s, nil
		}
		return s, nil
	case "ctrl+n":
		if s.kafkaStepFocus == kafkaStepFieldSecurityProtocol {
			s.protocolIndex = moveOption(s.protocolOptions, s.protocolIndex, 1)
		}
		return s, nil
	case "ctrl+p":
		if s.kafkaStepFocus == kafkaStepFieldSecurityProtocol {
			s.protocolIndex = moveOption(s.protocolOptions, s.protocolIndex, -1)
		}
		return s, nil
	}

	if s.kafkaStepFocus == kafkaStepFieldSecurityProtocol {
		return s, nil
	}

	var cmd tea.Cmd
	s.kafkaStepInputs[s.kafkaStepFocus], cmd = s.kafkaStepInputs[s.kafkaStepFocus].Update(msg)
	return s, cmd
}

func (s *Screen) viewKafkaStep() string {
	body := s.layout.Stack(
		s.layout.Panel(layout.PanelOptions{
			Body: s.notice.Render(
				"Step 2/4: Kafka connection",
				"Tell Koicha how to reach your Kafka cluster. Known values are selected, free-form values are typed.",
			),
			BodyAlign: layout.AlignLeft,
		}),
		s.layout.Panel(layout.PanelOptions{
			Title:     "Required connection",
			Subtitle:  "At least one bootstrap server is required.",
			Body:      strings.Join([]string{s.renderKafkaField(kafkaStepFieldBootstrapServers, "Bootstrap servers", true), "", s.renderSecurityProtocol()}, "\n"),
			BodyAlign: layout.AlignLeft,
			Active:    s.kafkaStepFocus == kafkaStepFieldBootstrapServers || s.kafkaStepFocus == kafkaStepFieldSecurityProtocol,
		}),
	)

	help := s.statusbar.Help("tab/enter next", "shift+tab previous", "ctrl+n/ctrl+p choose", "left/esc back")
	if s.editMode {
		help = s.statusbar.Help("tab/enter done", "shift+tab previous", "ctrl+n/ctrl+p choose", "left/esc review")
	}
	return s.layout.Render("", body, help)
}

func (s *Screen) renderKafkaField(index int, label string, required bool) string {
	return s.renderField(
		s.kafkaStepInputs[index],
		label,
		required,
		s.isKafkaInvalid(index),
	)
}

func (s *Screen) renderSecurityProtocol() string {
	parts := []string{s.fieldLabel("Security protocol", true)}
	if s.kafkaStepTouched[kafkaStepFieldSecurityProtocol] && selectedTitle(s.protocolOptions, s.protocolIndex) == "" {
		parts = append(parts, s.styles.Danger.Render("Choose one available protocol."))
	}
	parts = append(parts, renderSelector(
		s.styles,
		s.protocolOptions,
		s.protocolIndex,
		s.kafkaStepFocus == kafkaStepFieldSecurityProtocol,
	))
	return strings.Join(parts, "\n")
}

func (s *Screen) advanceKafka() bool {
	if !s.validateKafkaFocused() {
		return false
	}
	if s.kafkaStepFocus < kafkaStepFieldCount-1 {
		s.moveKafkaFocus(1)
		return true
	}

	s.commitKafka()
	s.completeStep(stepRules)
	return true
}

func (s *Screen) validateKafkaFocused() bool {
	s.kafkaStepTouched[s.kafkaStepFocus] = true
	if s.kafkaStepFocus == kafkaStepFieldBootstrapServers {
		return len(splitCommaValues(s.kafkaStepInputs[kafkaStepFieldBootstrapServers].Value())) > 0
	}
	if s.kafkaStepFocus == kafkaStepFieldSecurityProtocol {
		return !s.protocolOptions[s.protocolIndex].Disabled
	}
	return true
}

func (s *Screen) isKafkaInvalid(index int) bool {
	return index == kafkaStepFieldBootstrapServers &&
		s.kafkaStepTouched[index] &&
		len(splitCommaValues(s.kafkaStepInputs[index].Value())) == 0
}

func (s *Screen) moveKafkaFocus(delta int) {
	next := s.kafkaStepFocus + delta
	if next < 0 {
		next = 0
	}
	if next >= kafkaStepFieldCount {
		next = kafkaStepFieldCount - 1
	}
	s.focusKafka(next)
}

func (s *Screen) focusKafka(index int) {
	s.blurAllInputs()
	s.kafkaStepFocus = index
	if index < len(s.kafkaStepInputs) {
		s.kafkaStepInputs[index].Focus()
	}
}

func (s *Screen) commitKafka() {
	s.draft.Spec.Kafka.BootstrapServers = splitCommaValues(s.kafkaStepInputs[kafkaStepFieldBootstrapServers].Value())
	s.draft.Spec.Kafka.ClientID = kafka.DefaultClientID
	s.draft.Spec.Kafka.Auth.Protocol = kafka.SecurityProtocol(selectedTitle(s.protocolOptions, s.protocolIndex))
}
