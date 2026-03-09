Give the player a subtle hint about the current Wordle target word without revealing it.

Read `main.go` to understand the game and word list. Then ask the user which word they are stuck on (or if they share their current board state, infer what words have been ruled out).

Generate a clue that is:
- Subtle: never state the word directly or give away more than one property at a time
- Creative: use wordplay, analogies, definitions, or thematic associations
- Honest: the clue must be genuinely true about the word

Good clue styles (pick one that fits naturally):
- A cryptic-style definition: "Something a knight might carry, but also what you do with a pen"
- A category/association: "You'd find this in a kitchen or a forest"
- A wordplay angle: "Rhymes with X, means the opposite of Y"
- A usage example: "You might do this to a problem — or to a puzzle"

Bad clues (avoid):
- Directly stating the word or an obvious synonym
- Giving away the starting letter or letter count (already known)
- Being so vague the clue is useless

Present the hint as a single clean line prefixed with "Hint:".
