package createbento

import (
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/kafka"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/ui/components/layout"
	"github.com/openmdd/koicha/internal/tui/ui/components/metarows"
	"github.com/openmdd/koicha/internal/tui/ui/components/notice"
	"github.com/openmdd/koicha/internal/tui/ui/components/statusbar"
	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

type wizardStep int

const (
	stepMetadata wizardStep = iota
	stepKafka
	stepRules
	stepReview
)

const (
	metadataFieldName = iota
	metadataFieldProfile
	metadataFieldDisplayName
	metadataFieldCount
)

const (
	kafkaStepFieldBootstrapServers = iota
	kafkaStepFieldSecurityProtocol
	kafkaStepFieldCount
)

const kafkaStepInputCount = 1

const (
	rulesFieldResource = iota
	rulesFieldAction
	rulesFieldMatch
	rulesFieldValue
	rulesFieldAllowOutOfScope
	rulesFieldAddRule
	rulesFieldDone
	rulesFieldCount
)

type selectOption struct {
	Title    string
	Detail   string
	Disabled bool
}

type ruleResource int

const (
	ruleResourceTopics ruleResource = iota
	ruleResourceConsumerGroups
)

type ruleAction int

const (
	ruleActionInclude ruleAction = iota
	ruleActionExclude
)

type Screen struct {
	layout    layout.Model
	statusbar statusbar.Model
	styles    theme.Styles
	notice    notice.Model
	store     bento.Store

	currentStep     wizardStep
	editMode        bool
	status          string
	statusError     bool
	saveSuccessPath string

	metadataFocus   int
	metadataInputs  [metadataFieldCount]textinput.Model
	metadataTouched [metadataFieldCount]bool
	metadataErrors  [metadataFieldCount]string

	kafkaStepFocus   int
	kafkaStepInputs  [kafkaStepInputCount]textinput.Model
	kafkaStepTouched [kafkaStepFieldCount]bool
	protocolOptions  []selectOption
	protocolIndex    int

	rulesFocus             int
	ruleValueInput         textinput.Model
	ruleValueTouched       bool
	ruleResourceOptions    []selectOption
	ruleResourceIndex      int
	ruleActionOptions      []selectOption
	ruleActionIndex        int
	ruleMatchOptions       []selectOption
	ruleMatchIndex         int
	allowOutOfScopeOptions []selectOption
	allowOutOfScopeIndex   int

	draft bento.Bento
}

func New(styles theme.Styles, store bento.Store) *Screen {
	metadataInputs := [metadataFieldCount]textinput.Model{
		newInput("local-dev"),
		newInput("dev"),
		newInput("Local Dev"),
	}
	metadataInputs[metadataFieldName].Focus()

	kafkaStepInputs := [kafkaStepInputCount]textinput.Model{
		newInput("localhost:9092"),
	}

	ruleValueInput := newInput("payments.events or payments.")

	return &Screen{
		layout:          layout.New(styles),
		statusbar:       statusbar.New(styles),
		styles:          styles,
		notice:          notice.New(styles),
		store:           store,
		metadataInputs:  metadataInputs,
		kafkaStepInputs: kafkaStepInputs,
		protocolOptions: []selectOption{
			{Title: string(kafka.SecurityProtocolPlaintext), Detail: "available"},
			{Title: "SASL_SSL", Detail: "WIP", Disabled: true},
			{Title: "SSL", Detail: "WIP", Disabled: true},
			{Title: "SASL_PLAINTEXT", Detail: "WIP", Disabled: true},
		},
		ruleValueInput: ruleValueInput,
		ruleResourceOptions: []selectOption{
			{Title: "topics"},
			{Title: "consumer groups"},
		},
		ruleActionOptions: []selectOption{
			{Title: "include"},
			{Title: "exclude"},
		},
		ruleMatchOptions: []selectOption{
			{Title: string(bento.ResourcePatternExact)},
			{Title: string(bento.ResourcePatternPrefix)},
			{Title: string(bento.ResourcePatternRegex)},
		},
		allowOutOfScopeOptions: []selectOption{
			{Title: "yes", Detail: "recommended"},
			{Title: "no"},
		},
		draft: bento.Bento{
			SchemaVersion: bento.CurrentSchemaVersion,
			Spec: bento.Spec{
				Kafka: kafka.Config{
					ClientID: kafka.DefaultClientID,
					Auth:     kafka.Auth{Protocol: kafka.SecurityProtocolPlaintext},
				},
				Resources: bento.ResourceView{AllowOutOfScope: true},
			},
		},
	}
}

func (s *Screen) ID() nav.ScreenID { return nav.ScreenCreateBento }

func (s *Screen) Init() tea.Cmd { return nil }

func (s *Screen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.layout.SetSize(msg.Width, msg.Height)
		s.resizeInputs()
	case saveBentoMsg:
		return s.handleSaveBentoMsg(msg)
	case tea.KeyPressMsg:
		s.status = ""
		s.statusError = false
		switch s.currentStep {
		case stepMetadata:
			return s.updateMetadataStep(msg)
		case stepKafka:
			return s.updateKafkaStep(msg)
		case stepRules:
			return s.updateRulesStep(msg)
		case stepReview:
			return s.updateReviewStep(msg)
		}
	}
	return s, nil
}

func (s *Screen) View() string {
	switch s.currentStep {
	case stepKafka:
		return s.viewKafkaStep()
	case stepRules:
		return s.viewRulesStep()
	case stepReview:
		return s.viewReviewStep()
	default:
		return s.viewMetadataStep()
	}
}

func (s *Screen) completeStep(next wizardStep) {
	if s.editMode {
		s.editMode = false
		s.currentStep = stepReview
		s.blurAllInputs()
		return
	}
	s.currentStep = next
	s.focusStepStart(next)
}

func (s *Screen) editStep(step wizardStep) {
	s.editMode = true
	s.currentStep = step
	s.focusStepStart(step)
}

func (s *Screen) previousStep() (nav.Screen, tea.Cmd) {
	if s.editMode {
		s.editMode = false
		s.currentStep = stepReview
		s.blurAllInputs()
		return s, nil
	}

	switch s.currentStep {
	case stepMetadata:
		return s, func() tea.Msg { return nav.Back() }
	case stepKafka:
		s.currentStep = stepMetadata
		s.focusStepEnd(stepMetadata)
	case stepRules:
		s.currentStep = stepKafka
		s.focusStepEnd(stepKafka)
	case stepReview:
		s.currentStep = stepRules
		s.focusStepEnd(stepRules)
	}
	return s, nil
}

func (s *Screen) focusStepStart(step wizardStep) {
	switch step {
	case stepMetadata:
		s.focusMetadata(metadataFieldName)
	case stepKafka:
		s.focusKafka(kafkaStepFieldBootstrapServers)
	case stepRules:
		s.focusRules(rulesFieldResource)
	default:
		s.blurAllInputs()
	}
}

func (s *Screen) focusStepEnd(step wizardStep) {
	switch step {
	case stepMetadata:
		s.focusMetadata(metadataFieldDisplayName)
	case stepKafka:
		s.focusKafka(kafkaStepFieldSecurityProtocol)
	case stepRules:
		s.focusRules(rulesFieldDone)
	default:
		s.blurAllInputs()
	}
}

func (s *Screen) resizeInputs() {
	width := s.layout.InnerWidth()
	for i := range s.metadataInputs {
		s.metadataInputs[i].SetWidth(width)
	}
	for i := range s.kafkaStepInputs {
		s.kafkaStepInputs[i].SetWidth(width)
	}
	s.ruleValueInput.SetWidth(width)
}

func (s *Screen) blurAllInputs() {
	for i := range s.metadataInputs {
		s.metadataInputs[i].Blur()
	}
	for i := range s.kafkaStepInputs {
		s.kafkaStepInputs[i].Blur()
	}
	s.ruleValueInput.Blur()
}

func (s *Screen) renderField(input textinput.Model, label string, required bool, invalid bool) string {
	parts := []string{s.fieldLabel(label, required)}
	if invalid {
		parts = append(parts, s.styles.Danger.Render("Required value. Fill it in to continue."))
	}
	parts = append(parts, input.View())
	return strings.Join(parts, "\n")
}

func (s *Screen) fieldLabel(label string, required bool) string {
	if required {
		return s.styles.Error.Render(label + " *")
	}
	return s.styles.Subtle.Render(label)
}

func (s *Screen) renderMetadataLine(label, value string) string {
	return metarows.Render(s.styles, []metarows.Row{
		{Label: label, Value: optionalValue(value)},
	})
}

func newInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Prompt = "> "
	input.CharLimit = 240
	input.SetWidth(40)
	return input
}

func optionalValue(value string) string {
	if strings.TrimSpace(value) == "" {
		return "not set"
	}
	return value
}

func splitCommaValues(value string) []string {
	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}
