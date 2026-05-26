# Define TaskSpec schema and strict decoder

Goal: Define the JSON contract the plan agent emits and a strict decoder for it.

Context:
- Agent is instructed to emit JSON only; defensive fence-stripping handles stray markdown.

Do:
1. In internal/plan/spec.go define TaskSpec { Title string; Agent *string; Goal string; Context, Do, DoNot, DoneWhen, Files []string } with json tags (do_not, done_when).
2. Define Plan { Tasks []TaskSpec }.
3. Add func Decode(b []byte) (*Plan, error): strip surrounding ``` fences if present, use json.Decoder with DisallowUnknownFields, then validate.
4. Validate: Tasks non-empty; each Title, Goal non-empty; Do and DoneWhen non-empty; Agent if set is "executor" or "reviewer".

Done when:
- go build ./... succeeds.
- Valid JSON decodes; unknown field rejected; missing required field rejected.
- Empty tasks rejected.

Files: internal/plan/spec.go
