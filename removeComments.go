package oconf

// removeComments removes C-style comments from CJSON (JSON with comments) data.
// It handles both single-line comments (//) and multi-line comments (/* */).
// Comments inside quoted strings are preserved.
//
// The algorithm maintains state machines to track:
//   - Whether we're inside a quoted string
//   - Whether we're inside a single-line comment
//   - Whether we're inside a multi-line comment
//
// This implementation is careful to handle edge cases like:
//   - Quotes inside comments
//   - Comment markers inside strings
//   - Nested comment-like patterns
func removeComments(data []byte) []byte {
	// State flags
	var (
		inSingleLineComment = false // Currently inside a // comment
		inMultiLineComment  = false // Currently inside a /* */ comment
		inString            = false // Currently inside a quoted string
	)

	// Sliding windows to detect comment markers and quotes
	var (
		// For detecting single-line comment start (//)
		singleLineCommentStart = [3]byte{0, 0, 0}

		// For detecting single-line comment end (newline)
		singleLineCommentEnd = [2]byte{0, 0}

		// For detecting multi-line comment start (/*)
		multiLineCommentStart = [3]byte{0, 0, 0}

		// For detecting multi-line comment end (*/)
		multiLineCommentEnd = [3]byte{0, 0, 0}

		// For detecting quote characters
		quoteWindow = [2]byte{0, 0}
	)

	// Allocate output buffer (will be <= input size)
	output := make([]byte, len(data))
	outputPos := 0

	for _, currentByte := range data {
		// Shift quote window
		quoteWindow[0] = quoteWindow[1]
		quoteWindow[1] = currentByte

		// Track string state: toggle when we see a non-escaped quote
		// Only toggle if we're not already in a comment
		if quoteWindow[0] != '\\' && quoteWindow[1] == '"' && !inSingleLineComment && !inMultiLineComment {
			inString = !inString
		}

		// Only process comment markers if we're not in a string
		if !inString {
			// Check for single-line comment start (//)
			singleLineCommentStart[0] = singleLineCommentStart[1]
			singleLineCommentStart[1] = singleLineCommentStart[2]
			singleLineCommentStart[2] = currentByte

			if !inSingleLineComment && !inMultiLineComment &&
				singleLineCommentStart[1] == '/' && singleLineCommentStart[2] == '/' {
				inSingleLineComment = true
				// Remove the two '/' characters we just added
				outputPos -= 2
				if outputPos < 0 {
					outputPos = 0
				}
			}

			// Check for single-line comment end (newline)
			singleLineCommentEnd[0] = singleLineCommentEnd[1]
			singleLineCommentEnd[1] = currentByte

			if inSingleLineComment && (singleLineCommentEnd[1] == '\n' || singleLineCommentEnd[1] == '\r') {
				inSingleLineComment = false
				// Keep the newline character
				output[outputPos] = currentByte
				outputPos++
				continue
			}

			// Check for multi-line comment start (/*)
			multiLineCommentStart[0] = multiLineCommentStart[1]
			multiLineCommentStart[1] = multiLineCommentStart[2]
			multiLineCommentStart[2] = currentByte

			if !inMultiLineComment && !inSingleLineComment &&
				multiLineCommentStart[1] == '/' && multiLineCommentStart[2] == '*' {
				inMultiLineComment = true
				// Remove the '/' and '*' characters
				outputPos -= 2
				if outputPos < 0 {
					outputPos = 0
				}
			}

			// Check for multi-line comment end (*/)
			multiLineCommentEnd[0] = multiLineCommentEnd[1]
			multiLineCommentEnd[1] = multiLineCommentEnd[2]
			multiLineCommentEnd[2] = currentByte

			if inMultiLineComment &&
				multiLineCommentEnd[1] == '*' && multiLineCommentEnd[2] == '/' {
				inMultiLineComment = false
				continue
			}
		}

		// Copy byte to output if we're not in a comment
		if !inSingleLineComment && !inMultiLineComment {
			output[outputPos] = currentByte
			outputPos++
		}
	}

	return output[:outputPos]
}
