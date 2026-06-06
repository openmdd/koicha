package createbento

import (
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/openmdd/koicha/internal/bento"
	"github.com/openmdd/koicha/internal/tui/nav"
	"github.com/openmdd/koicha/internal/tui/ui/components/layout"
	"github.com/openmdd/koicha/internal/tui/ui/theme"
)

func (s *Screen) updateRulesStep(msg tea.KeyPressMsg) (nav.Screen, tea.Cmd) {
	switch msg.String() {
	case "left", "esc":
		return s.previousStep()
	case "shift+tab", "up":
		s.moveRulesFocus(-1)
		return s, nil
	case "down":
		s.moveRulesFocus(1)
		return s, nil
	case "tab":
		return s.advanceRules()
	case "right":
		s.moveRulesFocus(1)
		return s, nil
	case "enter":
		return s.activateRulesFocused()
	case "space", "ctrl+n":
		s.moveRulesOption(1)
		return s, nil
	case "ctrl+p":
		s.moveRulesOption(-1)
		return s, nil
	}

	if s.rulesFocus != rulesFieldValue {
		return s, nil
	}

	var cmd tea.Cmd
	s.ruleValueInput, cmd = s.ruleValueInput.Update(msg)
	return s, cmd
}

func (s *Screen) viewRulesStep() string {
	// TODO: Replace manual resource value entry with a broker-assisted picker
	// that fetches topics and consumer groups, supports search, and inserts
	// selected resources into include/exclude rules.
	body := s.layout.Stack(
		s.layout.Panel(layout.PanelOptions{
			Body: s.notice.Render(
				"Step 3/4: Bento rules",
				"Choose what Koicha should show by default. Empty rules mean show everything.",
			),
			BodyAlign: layout.AlignLeft,
		}),
		s.layout.Panel(layout.PanelOptions{
			Title:     "Rule builder",
			Subtitle:  "Use space or ctrl+n/ctrl+p on choices. Add several rules if needed.",
			Body:      s.renderRuleBuilder(),
			BodyAlign: layout.AlignLeft,
			Active:    true,
		}),
		s.layout.Panel(layout.PanelOptions{
			Title:     "Current rules",
			Body:      s.renderCurrentRules(),
			BodyAlign: layout.AlignLeft,
		}),
	)

	help := s.statusbar.Help("tab next", "enter activate", "space choose", "left/esc back")
	if s.editMode {
		help = s.statusbar.Help("tab next", "enter activate", "space choose", "left/esc review")
	}
	return s.layout.Render("", body, help)
}

func (s *Screen) renderRuleBuilder() string {
	return strings.Join([]string{
		s.renderRuleSelector("Resource", s.ruleResourceOptions, s.ruleResourceIndex, rulesFieldResource),
		"",
		s.renderRuleSelector("Action", s.ruleActionOptions, s.ruleActionIndex, rulesFieldAction),
		"",
		s.renderRuleSelector("Match", s.ruleMatchOptions, s.ruleMatchIndex, rulesFieldMatch),
		"",
		s.renderRuleValueField(),
		"",
		s.renderRuleSelector("Allow browsing outside Bento scope", s.allowOutOfScopeOptions, s.allowOutOfScopeIndex, rulesFieldAllowOutOfScope),
		"",
		s.renderRuleButton(rulesFieldAddRule, "Add rule"),
		s.renderRuleButton(rulesFieldDone, "Done"),
		s.renderRulesStatus(),
	}, "\n")
}

func (s *Screen) renderRuleSelector(label string, options []selectOption, selected int, field int) string {
	return strings.Join([]string{
		s.styles.Subtle.Render(label),
		renderInlineChoice(s.styles, options, selected, s.rulesFocus == field),
	}, "\n")
}

func (s *Screen) renderRuleValueField() string {
	parts := []string{s.fieldLabel("Value", false)}
	if s.ruleValueTouched && strings.TrimSpace(s.ruleValueInput.Value()) == "" {
		parts = append(parts, s.styles.Danger.Render("Enter a topic or consumer group value before adding a rule."))
	}
	parts = append(parts, s.ruleValueInput.View())
	return strings.Join(parts, "\n")
}

func (s *Screen) renderRuleButton(field int, label string) string {
	prefix := "  "
	style := s.styles.Subtle
	if s.rulesFocus == field {
		prefix = "> "
		style = s.styles.Title
	}
	return style.Render(prefix + label)
}

func (s *Screen) renderRulesStatus() string {
	if s.status == "" {
		return ""
	}
	return "\n" + s.styles.NoticeBody.Render(s.status)
}

func (s *Screen) renderCurrentRules() string {
	lines := []string{
		s.renderMetadataLine("allowOutOfScope", selectedTitle(s.allowOutOfScopeOptions, s.allowOutOfScopeIndex)),
		"",
		s.styles.Title.Render("topics"),
		renderScope(s.styles, s.draft.Spec.Resources.Topics),
		"",
		s.styles.Title.Render("consumer groups"),
		renderScope(s.styles, s.draft.Spec.Resources.ConsumerGroups),
	}
	return strings.Join(lines, "\n")
}

func (s *Screen) advanceRules() (nav.Screen, tea.Cmd) {
	if s.rulesFocus < rulesFieldCount-1 {
		s.moveRulesFocus(1)
		return s, nil
	}
	s.commitRulesOptions()
	s.completeStep(stepReview)
	return s, nil
}

func (s *Screen) activateRulesFocused() (nav.Screen, tea.Cmd) {
	switch s.rulesFocus {
	case rulesFieldAddRule:
		s.addResourceRule()
	case rulesFieldDone:
		s.commitRulesOptions()
		s.completeStep(stepReview)
	default:
		return s.advanceRules()
	}
	return s, nil
}

func (s *Screen) addResourceRule() {
	value := strings.TrimSpace(s.ruleValueInput.Value())
	if value == "" {
		s.ruleValueTouched = true
		return
	}

	pattern := bento.ResourcePattern{
		Kind:  bento.ResourcePatternKind(selectedTitle(s.ruleMatchOptions, s.ruleMatchIndex)),
		Value: value,
	}

	resource := s.selectedRuleResource()
	action := s.selectedRuleAction()
	switch {
	case resource == ruleResourceTopics && action == ruleActionInclude:
		s.draft.Spec.Resources.Topics.Include = append(s.draft.Spec.Resources.Topics.Include, pattern)
	case resource == ruleResourceTopics && action == ruleActionExclude:
		s.draft.Spec.Resources.Topics.Exclude = append(s.draft.Spec.Resources.Topics.Exclude, pattern)
	case resource == ruleResourceConsumerGroups && action == ruleActionInclude:
		s.draft.Spec.Resources.ConsumerGroups.Include = append(s.draft.Spec.Resources.ConsumerGroups.Include, pattern)
	case resource == ruleResourceConsumerGroups && action == ruleActionExclude:
		s.draft.Spec.Resources.ConsumerGroups.Exclude = append(s.draft.Spec.Resources.ConsumerGroups.Exclude, pattern)
	}

	s.commitRulesOptions()
	s.ruleValueInput.SetValue("")
	s.ruleValueTouched = false
	s.status = "Rule added. Add another one or choose Done."
	s.focusRules(rulesFieldValue)
}

func (s *Screen) commitRulesOptions() {
	s.draft.Spec.Resources.AllowOutOfScope = s.allowOutOfScopeIndex == 0
}

func (s *Screen) moveRulesFocus(delta int) {
	next := s.rulesFocus + delta
	if next < 0 {
		next = 0
	}
	if next >= rulesFieldCount {
		next = rulesFieldCount - 1
	}
	s.focusRules(next)
}

func (s *Screen) focusRules(index int) {
	s.blurAllInputs()
	s.rulesFocus = index
	if index == rulesFieldValue {
		s.ruleValueInput.Focus()
	}
}

func (s *Screen) moveRulesOption(delta int) {
	switch s.rulesFocus {
	case rulesFieldResource:
		s.ruleResourceIndex = moveOption(s.ruleResourceOptions, s.ruleResourceIndex, delta)
	case rulesFieldAction:
		s.ruleActionIndex = moveOption(s.ruleActionOptions, s.ruleActionIndex, delta)
	case rulesFieldMatch:
		s.ruleMatchIndex = moveOption(s.ruleMatchOptions, s.ruleMatchIndex, delta)
	case rulesFieldAllowOutOfScope:
		s.allowOutOfScopeIndex = moveOption(s.allowOutOfScopeOptions, s.allowOutOfScopeIndex, delta)
		s.commitRulesOptions()
	}
}

func (s *Screen) selectedRuleResource() ruleResource {
	if s.ruleResourceIndex == 1 {
		return ruleResourceConsumerGroups
	}
	return ruleResourceTopics
}

func (s *Screen) selectedRuleAction() ruleAction {
	if s.ruleActionIndex == 1 {
		return ruleActionExclude
	}
	return ruleActionInclude
}

func renderScope(styles theme.Styles, scope bento.ResourceScope) string {
	if len(scope.Include) == 0 && len(scope.Exclude) == 0 {
		return styles.Subtle.Render("  show all")
	}

	lines := make([]string, 0, len(scope.Include)+len(scope.Exclude))
	for _, pattern := range scope.Include {
		lines = append(lines, renderPattern(styles, "include", pattern))
	}
	for _, pattern := range scope.Exclude {
		lines = append(lines, renderPattern(styles, "exclude", pattern))
	}
	return strings.Join(lines, "\n")
}

func renderPattern(styles theme.Styles, action string, pattern bento.ResourcePattern) string {
	return strings.Join([]string{
		"  ",
		styles.Error.Render(action),
		styles.Subtle.Render(" " + string(pattern.Kind) + " "),
		styles.Title.Render(pattern.Value),
	}, "")
}
