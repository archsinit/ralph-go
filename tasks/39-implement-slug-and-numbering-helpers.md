# Implement slug and numbering helpers

Goal: Generate stable file numbers and kebab-case slugs for task files.

Context:
- The agent: prefix lives in the index title, not the filename.

Do:
1. In internal/plan/slug.go add func Slug(title string) string (lowercase, non-alnum -> hyphen, collapse repeats, trim).
2. Add func Number(i int) string -> zero-padded width 2 (01..99; widen if more).
3. Ensure collision-safe filenames across a plan (append -2 etc. if slugs repeat).

Done when:
- go build ./... succeeds.
- Slug strips an 'executor:'/'reviewer:' prefix from the title text used in the filename.
- Collisions resolved deterministically.

Files: internal/plan/slug.go
