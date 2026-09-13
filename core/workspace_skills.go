package core

// skillsForAgent takes a fresh directory snapshot. Registries are request-local:
// changing a work directory or rebinding a channel cannot leave stale skill
// contents in another channel's cache. No model session is started here.
func (e *Engine) skillsForAgent(agent Agent) *SkillRegistry {
	if sp, ok := agent.(SkillProvider); ok {
		registry := NewSkillRegistry()
		registry.SetDirs(sp.SkillDirs())
		return registry
	}
	if agent == e.agent {
		return e.skills
	}
	return NewSkillRegistry()
}

func (e *Engine) skillsForMessage(p Platform, msg *Message) (*SkillRegistry, error) {
	agent, _, _, _, err := e.commandContextWithWorkspace(p, msg)
	if err != nil {
		return nil, err
	}
	return e.skillsForAgent(agent), nil
}

func (e *Engine) skillsCardForSession(sessionKey string) *Card {
	p := e.platformForName(extractPlatformName(sessionKey))
	registry, err := e.skillsForMessage(p, &Message{SessionKey: sessionKey, Platform: p.Name()})
	if err != nil {
		return e.simpleCard(e.i18n.T(MsgCardTitleSkills), "purple", e.i18n.Tf(MsgWsResolutionError, err))
	}
	return e.renderSkillsCard(registry)
}
