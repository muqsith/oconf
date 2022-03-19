package oconf

/*
removeComments removes comments from cjson file
*/
func removeComments(data []byte) []byte {
	// stop for single-line comment
	stopforslc := false
	// stop for multi-line comment
	stopformlc := false

	quotestart := false
	quoteend := true

	// single-line comment start
	slcstart := []byte{0, 0, 0}
	// single-line commment end
	slcend := []byte{0, 0}
	// multi-line comment start
	mlcstart := []byte{0, 0, 0}
	// multi-line comment end
	mlcend := []byte{0, 0, 0}
	// quote
	quotes := []byte{0, 0}

	dataSize := len(data)
	dataCopy := make([]byte, dataSize)

	j := 0

	for _, c := range data {

		quotes[0] = quotes[1]
		quotes[1] = c

		if quotes[0] == '"' && (!stopforslc && !stopformlc) {
			quotestart = quoteend
			quoteend = !quotestart
		}

		slcstart[0] = slcstart[1]
		slcstart[1] = slcstart[2]
		slcstart[2] = c
		if !stopforslc && !quotestart && quoteend && slcstart[0] == '/' && slcstart[1] == '/' {
			stopforslc = true
			j -= 2
		}
		slcend[0] = slcend[1]
		slcend[1] = c
		if stopforslc && !quotestart && quoteend && (slcend[0] == '\n' || slcend[0] == '\r') {
			stopforslc = false
		}

		mlcstart[0] = mlcstart[1]
		mlcstart[1] = mlcstart[2]
		mlcstart[2] = c
		if !stopformlc && !quotestart && quoteend && mlcstart[0] == '/' && mlcstart[1] == '*' {
			stopformlc = true
			j -= 2
		}
		mlcend[0] = mlcend[1]
		mlcend[1] = mlcend[2]
		mlcend[2] = c
		if stopformlc && !quotestart && quoteend && mlcend[0] == '*' && mlcend[1] == '/' {
			stopformlc = false
		}

		if !stopforslc && !stopformlc {
			dataCopy[j] = c
			j++
		}
	}

	return dataCopy[0:j]
}
