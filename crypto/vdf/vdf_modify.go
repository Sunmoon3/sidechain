package vdf

type VDFM struct {
	difficulty int
	input      []byte
	output     []byte
	outputChan chan []byte
	finished   bool
}

func NewVDFM(difficulty int, input []byte) *VDFM {
	return &VDFM{
		difficulty: difficulty,
		input:      input,
		outputChan: make(chan []byte),
	}
}


func (vdf *VDFM) GetOutputChannel() chan []byte {
	return vdf.outputChan
}


func (vdf *VDFM) Execute() {
	vdf.finished = false

	yBuf, proofBuf := GenerateVDF(vdf.input[:], vdf.difficulty, sizeInBits)

	vdf.output= append(vdf.output, yBuf...)
	vdf.output= append(vdf.output, proofBuf...)
	//copy(assist.output[:], yBuf)
	//copy(assist.output[len(yBuf):], proofBuf)

	go func() {
		vdf.outputChan <- vdf.output
	}()

	vdf.finished = true
}


// Verify runs the verification of generated proof
// currently on i7-6700K, verification takes about 350 ms
func (vdf *VDFM) Verify(proof []byte) bool {
	return VerifyVDF(vdf.input[:], proof[:], vdf.difficulty, sizeInBits)
}

// IsFinished returns whether the assist execution is finished or not.
func (vdf *VDFM) IsFinished() bool {
	return vdf.finished
}

// GetOutput returns the assist output, which can be bytes of 0s is the assist is not finished.
func (vdf *VDFM) GetOutput() []byte {
	return vdf.output
}
