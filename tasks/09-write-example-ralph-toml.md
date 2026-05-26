# Write example ralph.toml

Goal: Provide a working example config with two plan agents plus user, and loop executor/reviewer.

Context:
- This file doubles as documentation of the format.

Do:
1. Create ralph.toml at repo root.
2. Define agents: claude (cli=claude), codex (cli=codex), both enabled, each with a short system_prompt.
3. plan.turn_order = ["claude","codex","user"]; plan.plan_agent = "claude".
4. loop.executor = "claude"; loop.reviewer = "codex"; loop.max_retries = 2; loop.task_timeout = "15m".
5. ntfy.server = "https://ntfy.sh"; ntfy.topic = "CHANGE-ME"; ntfy.token = "".
6. paths.session_dir = ".ralph/sessions"; paths.log_dir = ".ralph/logs"; paths.out_dir = ".".

Done when:
- ralph.toml exists and config.Load + Validate accept it after replacing the ntfy topic.
- Comments explain each section.

Files: ralph.toml
