package orchestrator

import (
	"context"
	"fmt"
	"strings"

	"github.com/archsinit/ralph-go/internal/agent"
	"github.com/archsinit/ralph-go/internal/config"
	"github.com/archsinit/ralph-go/internal/plan"
	"github.com/archsinit/ralph-go/internal/session"
)

// TurnOrder represents a sequence of agent turns for orchestration.
type TurnOrder struct {
	// Order lists agent names (or "user") in the order they should take turns.
	Order []string
}

// NewTurnOrder creates a new turn order with the given sequence.
func NewTurnOrder(order []string) *TurnOrder {
	return &TurnOrder{
		Order: order,
	}
}

// UIBridge is the interface for communicating with the user interface.
type UIBridge interface {
	// RequestUserInput blocks until the user submits text, then returns it.
	// Returns context.Canceled if the user cancels the current turn.
	RequestUserInput(ctx context.Context) (string, error)

	// RequestCancelSignal returns a channel that signals when the user presses cancel (ctrl-c).
	// Multiple signals on the channel may occur for repeated ctrl-c.
	RequestCancelSignal(ctx context.Context) <-chan struct{}

	// StreamStart notifies the UI that an agent stream is starting.
	StreamStart(author string) error

	// StreamToken sends a token to the UI.
	StreamToken(text string) error

	// StreamEnd notifies the UI that a stream has ended.
	StreamEnd() error

	// AddMessage adds a finalized message to the transcript.
	AddMessage(author, text string) error
}

// Engine drives turns over a turn order, blocking on user slots and invoking agents.
type Engine struct {
	cfg       *config.Config
	sess      *session.Session
	turnOrder *TurnOrder
	adapters  map[string]agent.Adapter
	bridge    UIBridge
	nextTurn  string // Override the cyclic next turn with a prefix-targeted agent
}

// NewEngine creates a new Engine.
func NewEngine(cfg *config.Config, sess *session.Session, turnOrder *TurnOrder, adapters map[string]agent.Adapter, bridge UIBridge) *Engine {
	return &Engine{
		cfg:       cfg,
		sess:      sess,
		turnOrder: turnOrder,
		adapters:  adapters,
		bridge:    bridge,
	}
}

// Run executes the plan's turn order until completion or error.
func (e *Engine) Run(ctx context.Context) error {
	turnIdx := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Determine next turn: use override if set, otherwise use cyclic order
		var slot string
		if e.nextTurn != "" {
			slot = e.nextTurn
			e.nextTurn = ""
		} else {
			if len(e.turnOrder.Order) == 0 {
				return fmt.Errorf("empty turn order")
			}
			slot = e.turnOrder.Order[turnIdx%len(e.turnOrder.Order)]
			turnIdx++
		}

		// Handle user turn
		if slot == "user" {
			err := e.runUserTurn(ctx)
			if err != nil && err != context.Canceled {
				return err
			}
			// On context cancellation, propagate it up
			if err == context.Canceled {
				return err
			}
			continue
		}

		// Handle agent turn
		err := e.runAgentTurn(ctx, slot)
		if err != nil && err != context.Canceled {
			return err
		}
		// On context cancellation, propagate it up
		if err == context.Canceled {
			return err
		}
	}
}

// runUserTurn handles a user input turn.
func (e *Engine) runUserTurn(ctx context.Context) error {
	cancelSignal := e.bridge.RequestCancelSignal(ctx)
	turnCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	cancelPressed := false
	go func() {
		for {
			select {
			case <-cancelSignal:
				if cancelPressed {
					// Second ctrl-c: request full app quit
					cancel()
					return
				}
				cancelPressed = true
			case <-turnCtx.Done():
				return
			}
		}
	}()

	userInput, err := e.bridge.RequestUserInput(turnCtx)

	// Check for context cancellation (ctrl-c)
	if err == context.Canceled {
		// User pressed cancel on an input turn - just retry this turn
		return nil
	}
	if err != nil {
		return err
	}

	// Check for /plan command
	if userInput == "/plan" {
		return e.runPlanTurn(ctx)
	}

	// Check for /quit command
	if userInput == "/quit" {
		return context.Canceled // Signal app quit
	}

	// Parse prefix to see if user is targeting a specific agent
	target, body := ParsePrefix(userInput, e.getAgentNames())
	if target != "" {
		e.nextTurn = target
	}

	// Add the user message to the session
	msg := session.Message{
		Author: "user",
		Role:   "user",
		Text:   body,
	}
	if err := e.sess.Append(msg); err != nil {
		return err
	}

	// Display the user message in the UI
	if err := e.bridge.AddMessage("user", body); err != nil {
		return err
	}

	return nil
}

// runAgentTurn invokes an agent and streams its response.
func (e *Engine) runAgentTurn(ctx context.Context, agentName string) error {
	adapter, ok := e.adapters[agentName]
	if !ok {
		return fmt.Errorf("unknown agent: %s", agentName)
	}

	// Get the agent's system prompt from config
	systemPrompt := ""
	if e.cfg != nil {
		for _, cfgAgent := range e.cfg.Agents {
			if cfgAgent.Name == agentName && cfgAgent.SystemPrompt != "" {
				systemPrompt = cfgAgent.SystemPrompt
				break
			}
		}
	}

	// Convert session messages to agent messages
	var messages []agent.Message
	for _, sm := range e.sess.Messages {
		messages = append(messages, agent.Message{
			Author: sm.Author,
			Role:   sm.Role,
			Text:   sm.Text,
		})
	}

	// Populate NewMessage with the latest message in the transcript
	newMessage := ""
	if len(messages) > 0 {
		newMessage = messages[len(messages)-1].Text
	}

	resumeID := e.sess.ResumeID(agentName)

	req := agent.Request{
		SystemPrompt:    systemPrompt,
		Transcript:      messages,
		NewMessage:      newMessage,
		ResumeSessionID: resumeID,
	}

	// Create a per-turn context for cancel monitoring
	turnCtx, turnCancel := context.WithCancel(ctx)
	defer turnCancel()

	// Monitor cancel signals during the agent turn
	cancelSignal := e.bridge.RequestCancelSignal(ctx)
	streamCanceled := false

	go func() {
		select {
		case <-cancelSignal:
			streamCanceled = true
			turnCancel()
		case <-turnCtx.Done():
		}
	}()

	// Invoke the agent
	tokenChan, resultChan := adapter.Invoke(turnCtx, req)

	// Notify UI that we're starting a stream
	if err := e.bridge.StreamStart(agentName); err != nil {
		return err
	}

	// Accumulate the full response as we stream tokens
	var fullResponse string
	for {
		select {
		case token, ok := <-tokenChan:
			if !ok {
				tokenChan = nil
				break
			}
			fullResponse += token.Text
			if err := e.bridge.StreamToken(token.Text); err != nil {
				return err
			}
		case <-turnCtx.Done():
			// Context was canceled; drain remaining tokens safely
			for {
				select {
				case token, ok := <-tokenChan:
					if !ok {
						tokenChan = nil
						break
					}
					fullResponse += token.Text
				default:
					tokenChan = nil
				}
				if tokenChan == nil {
					break
				}
			}
			break
		}

		if tokenChan == nil {
			break
		}
	}

	// Notify UI that stream ended (even if canceled)
	if err := e.bridge.StreamEnd(); err != nil {
		return err
	}

	// If the turn was canceled, return the cancellation error
	if streamCanceled {
		// Append cancellation note to session
		if err := e.sess.Append(session.Message{
			Author: agentName,
			Role:   "assistant",
			Text:   "(turn cancelled)",
		}); err != nil {
			return err
		}
		// Display the cancellation note in the UI
		if err := e.bridge.AddMessage(agentName, "(turn cancelled)"); err != nil {
			return err
		}
		// Return to user's turn
		e.nextTurn = "user"
		return context.Canceled
	}

	// Get the result (should be available after token channel closes)
	result := <-resultChan
	if result.Err != nil {
		return result.Err
	}

	// Append the agent's message to the session
	msg := session.Message{
		Author: agentName,
		Role:   "assistant",
		Text:   fullResponse,
	}
	if err := e.sess.Append(msg); err != nil {
		return err
	}

	// Update the agent session ID if provided
	if result.SessionID != "" {
		if err := e.sess.SetAgentSession(agentName, result.SessionID); err != nil {
			return err
		}
	}

	// Display the agent message in the UI
	return e.bridge.AddMessage(agentName, fullResponse)
}

// runPlanTurn handles the /plan command, requesting a JSON task list from the plan agent.
func (e *Engine) runPlanTurn(ctx context.Context) error {
	// Get the plan agent from config
	planAgent := ""
	if e.cfg != nil && e.cfg.Plan.PlanAgent != "" {
		planAgent = e.cfg.Plan.PlanAgent
	}

	if planAgent == "" {
		// No plan agent configured; silently return to user
		return e.bridge.AddMessage("system", "Plan agent not configured in settings")
	}

	// Get the adapter for the plan agent
	adapter, ok := e.adapters[planAgent]
	if !ok {
		return fmt.Errorf("plan agent %q not found", planAgent)
	}

	// Convert session messages to agent messages
	var messages []agent.Message
	for _, sm := range e.sess.Messages {
		messages = append(messages, agent.Message{
			Author: sm.Author,
			Role:   sm.Role,
			Text:   sm.Text,
		})
	}

	// Create a special request for plan generation
	req := agent.Request{
		SystemPrompt: `You are a task planning assistant. Given the conversation context, generate a JSON array of tasks to complete. Respond with ONLY valid JSON, no markdown or extra text. The JSON should be an array of task objects, each with "id", "title", and "description" fields.`,
		Transcript:   messages,
		NewMessage:   "Please generate a plan based on our discussion",
	}

	// Notify UI that we're generating a plan
	if err := e.bridge.AddMessage("system", "Generating plan..."); err != nil {
		return err
	}

	// Invoke the plan agent
	tokenChan, resultChan := adapter.Invoke(ctx, req)

	// Accumulate the full response
	var planJSON string
	for token := range tokenChan {
		planJSON += token.Text
	}

	// Get the result
	result := <-resultChan
	if result.Err != nil {
		return result.Err
	}

	// Parse the JSON plan
	decodedPlan, err := plan.Decode(planJSON)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to parse plan: %v", err)
		if err := e.bridge.AddMessage("system", errMsg); err != nil {
			return err
		}
		e.nextTurn = "user"
		return nil
	}

	// Write the plan files
	outDir := "."
	if e.cfg != nil && e.cfg.Paths.OutDir != "" {
		outDir = e.cfg.Paths.OutDir
	}

	written, err := decodedPlan.Write(outDir)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to write plan files: %v", err)
		if err := e.bridge.AddMessage("system", errMsg); err != nil {
			return err
		}
		e.nextTurn = "user"
		return nil
	}

	// Report success
	successMsg := fmt.Sprintf("Plan generated successfully. Files written:\n%s", strings.Join(written, "\n"))
	if err := e.bridge.AddMessage("system", successMsg); err != nil {
		return err
	}

	// Return to user turn
	e.nextTurn = "user"
	return nil
}

// getAgentNames returns the list of configured agent names.
func (e *Engine) getAgentNames() []string {
	var names []string
	for name := range e.adapters {
		names = append(names, name)
	}
	return names
}
