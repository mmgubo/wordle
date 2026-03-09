package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[42;30m"
	colorYellow = "\033[43;30m"
	colorGray   = "\033[100;37m"
)

type tileState int

const (
	stateUnknown tileState = iota
	stateAbsent
	statePresent
	stateCorrect
)

// evaluate scores a guess against the target word.
// Handles duplicate letters correctly: a letter is only marked present/correct
// as many times as it appears in the target.
func evaluate(guess, target string) [5]tileState {
	res := [5]tileState{}
	freq := [26]int{}

	// First pass: find exact matches and count remaining target letters.
	for i := 0; i < 5; i++ {
		if guess[i] == target[i] {
			res[i] = stateCorrect
		} else {
			freq[target[i]-'a']++
		}
	}

	// Second pass: find letters present but in wrong positions.
	for i := 0; i < 5; i++ {
		if res[i] == stateCorrect {
			continue
		}
		idx := guess[i] - 'a'
		if freq[idx] > 0 {
			res[i] = statePresent
			freq[idx]--
		} else {
			res[i] = stateAbsent
		}
	}

	return res
}

func renderTile(ch byte, s tileState) string {
	letter := string(rune(ch - 32)) // ASCII lowercase to uppercase
	switch s {
	case stateCorrect:
		return colorGreen + " " + letter + " " + colorReset
	case statePresent:
		return colorYellow + " " + letter + " " + colorReset
	case stateAbsent:
		return colorGray + " " + letter + " " + colorReset
	}
	return "   "
}

func draw(guesses []string, results [][5]tileState) {
	fmt.Print("\033[H\033[2J")
	fmt.Println(colorBold + "   W O R D L E  [HARD]" + colorReset)
	fmt.Println()

	// Board: 6 rows of 5 tiles
	for i := 0; i < 6; i++ {
		fmt.Print("  ")
		if i < len(guesses) {
			for j := 0; j < 5; j++ {
				fmt.Print(renderTile(guesses[i][j], results[i][j]))
			}
		} else {
			fmt.Print("\033[90m[ ][ ][ ][ ][ ]\033[0m")
		}
		fmt.Println()
	}

	// On-screen keyboard showing letter states
	best := make(map[byte]tileState)
	for g, guess := range guesses {
		for j := 0; j < 5; j++ {
			ch := guess[j]
			if s := results[g][j]; s > best[ch] {
				best[ch] = s
			}
		}
	}

	rows := [3]string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}
	pads := [3]string{"", " ", "   "}
	fmt.Println()
	for i, row := range rows {
		fmt.Print("  " + pads[i])
		for j := 0; j < len(row); j++ {
			ch := row[j]
			letter := string(rune(ch - 32))
			switch best[ch] {
			case stateCorrect:
				fmt.Print(colorGreen + letter + colorReset + " ")
			case statePresent:
				fmt.Print(colorYellow + letter + colorReset + " ")
			case stateAbsent:
				fmt.Print(colorGray + letter + colorReset + " ")
			default:
				fmt.Print(letter + " ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// printHint reveals a progressively more specific clue about the target.
// Hints never directly expose the word; they give structural or positional info.
func printHint(target string, n int) {
	switch n {
	case 1:
		vowelCount := 0
		for _, ch := range target {
			if ch == 'a' || ch == 'e' || ch == 'i' || ch == 'o' || ch == 'u' {
				vowelCount++
			}
		}
		fmt.Printf("  Hint: The word contains %d vowel(s).\n", vowelCount)
	case 2:
		fmt.Printf("  Hint: The word ends with '%s'.\n", strings.ToUpper(target[4:]))
	default:
		fmt.Printf("  Hint: The word begins with '%s'.\n", strings.ToUpper(target[:1]))
	}
}

func main() {
	enableANSI()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	target := answers[r.Intn(len(answers))]

	valid := make(map[string]bool, len(answers))
	for _, w := range answers {
		valid[w] = true
	}

	guesses := make([]string, 0, 6)
	results := make([][5]tileState, 0, 6)
	hintCount := 0
	scanner := bufio.NewScanner(os.Stdin)

	// Hard mode constraints: greens fix a position, yellows must reappear.
	fixed := [5]byte{}
	mustHave := map[byte]bool{}

	for {
		draw(guesses, results)

		n := len(guesses)

		// Check win/loss after redrawing so the final board is shown.
		if n > 0 && guesses[n-1] == target {
			messages := [7]string{"", "Genius!", "Magnificent!", "Impressive!", "Splendid!", "Great!", "Phew!"}
			fmt.Println("  " + colorBold + messages[n] + colorReset)
			break
		}
		if n == 6 {
			fmt.Printf("  The word was: %s%s%s\n", colorBold, strings.ToUpper(target), colorReset)
			break
		}

		// Prompt for input, re-asking on invalid entries.
		var guess string
		for {
			fmt.Printf("  Guess %d/6 (or /hint): ", n+1)
			if !scanner.Scan() {
				os.Exit(0)
			}
			g := strings.ToLower(strings.TrimSpace(scanner.Text()))

			if g == "/hint" {
				hintCount++
				printHint(target, hintCount)
				continue
			}

			if len(g) != 5 {
				fmt.Println("  Enter a 5-letter word.")
				continue
			}
			allAlpha := true
			for i := 0; i < 5; i++ {
				if g[i] < 'a' || g[i] > 'z' {
					allAlpha = false
					break
				}
			}
			if !allAlpha {
				fmt.Println("  Letters only.")
				continue
			}
			if !valid[g] {
				fmt.Println("  Not in word list.")
				continue
			}

			// Hard mode: enforce previously revealed clues.
			badPos := -1
			for i := 0; i < 5; i++ {
				if fixed[i] != 0 && g[i] != fixed[i] {
					badPos = i
					break
				}
			}
			if badPos >= 0 {
				fmt.Printf("  Position %d must be %s.\n", badPos+1, strings.ToUpper(string(fixed[badPos])))
				continue
			}
			missingLetter := byte(0)
			for ch := range mustHave {
				if !strings.ContainsRune(g, rune(ch)) {
					missingLetter = ch
					break
				}
			}
			if missingLetter != 0 {
				fmt.Printf("  Guess must contain %s.\n", strings.ToUpper(string(missingLetter)))
				continue
			}

			guess = g
			break
		}

		res := evaluate(guess, target)
		guesses = append(guesses, guess)
		results = append(results, res)

		// Update hard mode constraints from this result.
		for i := 0; i < 5; i++ {
			switch res[i] {
			case stateCorrect:
				fixed[i] = guess[i]
			case statePresent:
				mustHave[guess[i]] = true
			}
		}
	}
}

var answers = []string{
	"about", "above", "abuse", "actor", "acute",
	"admit", "adult", "after", "again", "agent",
	"agree", "ahead", "alarm", "album", "alert",
	"alike", "alive", "alley", "allow", "alone",
	"along", "alter", "angel", "anger", "angle",
	"angry", "ankle", "apple", "apply", "arena",
	"argue", "arise", "armor", "aroma", "array",
	"aside", "asset", "atlas", "attic", "avoid",
	"awake", "award", "aware", "awful", "basic",
	"beach", "beard", "beast", "began", "being",
	"below", "bench", "birth", "black", "blade",
	"blame", "blast", "blaze", "bleed", "blend",
	"bless", "blind", "block", "blood", "bloom",
	"blown", "blues", "blunt", "board", "bonus",
	"boost", "booth", "bound", "boxer", "brain",
	"brave", "bread", "break", "breed", "brick",
	"bride", "brief", "bring", "broad", "broke",
	"broom", "brown", "brush", "build", "built",
	"burst", "buyer", "cable", "candy", "carry",
	"catch", "cause", "chain", "chair", "chaos",
	"chase", "cheap", "check", "cheek", "chess",
	"chest", "chief", "child", "chose", "civil",
	"claim", "class", "clean", "clear", "clerk",
	"click", "cliff", "climb", "clock", "clone",
	"close", "cloth", "cloud", "coach", "coast",
	"comet", "comic", "coral", "count", "court",
	"cover", "craft", "crane", "crash", "cream",
	"crime", "crisp", "cross", "crowd", "crown",
	"cruel", "crush", "curve", "cycle", "dance",
	"death", "delta", "dense", "depth", "devil",
	"digit", "dirty", "dizzy", "dodge", "doubt",
	"dough", "draft", "drain", "drama", "drawn",
	"dream", "dress", "drink", "drive", "drove",
	"drums", "drunk", "dying", "eagle", "early",
	"earth", "eight", "elite", "empty", "enemy",
	"enjoy", "enter", "entry", "equal", "error",
	"essay", "event", "exact", "exist", "extra",
	"fable", "falls", "false", "fancy", "fatal",
	"favor", "feast", "fence", "fever", "final",
	"first", "fixed", "flame", "flash", "fleet",
	"flesh", "float", "flood", "floor", "flour",
	"focus", "force", "forge", "forum", "found",
	"frame", "fresh", "front", "frost", "funny",
	"fuzzy", "giant", "given", "glass", "globe",
	"gloom", "glove", "going", "grace", "grade",
	"grain", "grand", "grant", "grape", "grasp",
	"grass", "grave", "great", "green", "grill",
	"grind", "gross", "grown", "guard", "guess",
	"guide", "guilt", "habit", "happy", "harsh",
	"heart", "heavy", "herbs", "hinge", "hotel",
	"house", "human", "hurry", "ideal", "image",
	"imply", "index", "inner", "input", "irony",
	"ivory", "jewel", "joker", "joint", "juice",
	"juicy", "jumbo", "jumpy", "knife", "kneel",
	"known", "label", "labor", "large", "laser",
	"later", "laugh", "layer", "learn", "lease",
	"legal", "lemon", "level", "light", "limit",
	"liver", "lodge", "logic", "loose", "lower",
	"lucky", "lusty", "magic", "major", "maker",
	"maple", "march", "marsh", "match", "mayor",
	"media", "mercy", "merge", "merit", "metal",
	"might", "minor", "minus", "mirth", "model",
	"money", "month", "moose", "moral", "motor",
	"mount", "mouse", "mouth", "muddy", "music",
	"musty", "naive", "ninja", "noble", "noise",
	"north", "novel", "ocean", "often", "olive",
	"opera", "orbit", "order", "other", "outer",
	"oxide", "paint", "panic", "paper", "party",
	"pasta", "patch", "pause", "peace", "peach",
	"pearl", "penny", "perch", "piano", "piece",
	"pilot", "pinch", "pizza", "place", "plain",
	"plane", "plant", "plead", "pluck", "plumb",
	"plume", "plush", "point", "porch", "power",
	"press", "price", "pride", "prime", "prior",
	"prize", "probe", "prone", "prose", "proud",
	"prove", "pulse", "pupil", "purse", "queen",
	"quest", "queue", "quick", "quiet", "quota",
	"quote", "radar", "raise", "ranch", "range",
	"rapid", "raven", "reach", "react", "relay",
	"repay", "rider", "ridge", "rifle", "right",
	"risky", "rival", "river", "rocky", "rough",
	"round", "route", "royal", "ruler", "rusty",
	"sadly", "saint", "sandy", "sauce", "savvy",
	"scale", "scene", "scope", "score", "scout",
	"screw", "seize", "sense", "serum", "serve",
	"setup", "seven", "shard", "share", "shark",
	"sharp", "sheer", "shelf", "shell", "shift",
	"shirt", "shock", "shore", "short", "shout",
	"sight", "silly", "since", "sixth", "sixty",
	"skill", "skull", "slash", "slave", "sleep",
	"slice", "slide", "slope", "smart", "smell",
	"smile", "smoke", "snack", "snake", "solar",
	"solid", "solve", "sorry", "south", "space",
	"spare", "spark", "speak", "spell", "spend",
	"spice", "spine", "spoke", "spoon", "sport",
	"spray", "stack", "staff", "stage", "stain",
	"stale", "stalk", "stand", "stark", "start",
	"state", "steak", "steel", "steep", "steer",
	"stern", "stick", "stiff", "still", "stock",
	"stone", "stood", "store", "storm", "story",
	"stout", "strap", "straw", "stray", "stuck",
	"study", "stuff", "stunt", "style", "sugar",
	"suite", "super", "surge", "sword", "swamp",
	"swear", "sweet", "swept", "swift", "swirl",
	"table", "talon", "tango", "taste", "tasty",
	"taunt", "tense", "theft", "theme", "there",
	"these", "third", "thorn", "those", "three",
	"thumb", "tiger", "tired", "title", "toast",
	"today", "token", "total", "touch", "tough",
	"tower", "toxic", "track", "trade", "trail",
	"train", "trait", "trash", "treat", "tribe",
	"trick", "troop", "trout", "truce", "truck",
	"truly", "trunk", "trust", "truth", "tulip",
	"twice", "twist", "ultra", "uncle", "under",
	"union", "unity", "until", "upper", "upset",
	"urban", "usual", "utter", "vague", "valid",
	"value", "valve", "vapor", "vault", "verse",
	"vigor", "viper", "viral", "virus", "visor",
	"vista", "vital", "vivid", "vocal", "voice",
	"voter", "vowel", "wager", "waste", "watch",
	"water", "weary", "wedge", "weird", "whale",
	"wheat", "wheel", "where", "which", "while",
	"white", "whole", "whose", "woman", "women",
	"world", "worry", "worst", "worth", "would",
	"wound", "wrath", "wrist", "wrong", "yacht",
	"yearn", "yield", "young", "youth", "zesty",
}
