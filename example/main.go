package main

import (
	"fmt"
	"math"

	"github.com/schollz/goaubio-onset"
)

func main() {
	// Parameters
	bufSize := uint(512)
	hopSize := uint(256)
	samplerate := uint(44100)

	// Create onset detector with "hfc" method (High Frequency Content)
	o := onset.NewOnset("hfc", bufSize, hopSize, samplerate)

	// Create input buffer and onset output buffer
	input := onset.NewFvec(hopSize)
	onsetOut := onset.NewFvec(1)

	fmt.Println("Onset Detection Example")
	fmt.Println("=======================")
	fmt.Printf("Buffer size: %d\n", bufSize)
	fmt.Printf("Hop size: %d\n", hopSize)
	fmt.Printf("Sample rate: %d Hz\n", samplerate)
	fmt.Printf("Method: HFC (High Frequency Content)\n\n")

	// Simulate processing audio frames
	// In a real application, you would read audio data from a file or stream
	totalFrames := 20
	fmt.Println("Processing audio frames...")

	for frame := 0; frame < totalFrames; frame++ {
		// Generate a simple test signal
		// In real usage, this would be actual audio data
		for i := uint(0); i < hopSize; i++ {
			// Simple sine wave with some variation
			t := float64(frame*int(hopSize)+int(i)) / float64(samplerate)
			frequency := 440.0 // A4 note

			// Add an "onset" every 5 frames by changing amplitude
			amplitude := 0.5
			if frame%5 == 0 && i < hopSize/4 {
				amplitude = 1.0 // Simulated onset
			}

			input.Data[i] = amplitude * math.Sin(2*math.Pi*frequency*t)
		}

		// Process the input
		o.Do(input, onsetOut)

		// Check if an onset was detected
		if onsetOut.Data[0] > 0 {
			onsetTime := o.GetLastMs()
			fmt.Printf("Frame %3d: ONSET DETECTED at %.2f ms (value: %.3f)\n",
				frame, onsetTime, onsetOut.Data[0])
		} else {
			// Uncomment to see all frames
			// fmt.Printf("Frame %3d: No onset (descriptor: %.3f)\n",
			//    frame, o.GetDescriptor())
		}
	}

	fmt.Println("\nDetection complete!")
	fmt.Printf("Total frames processed: %d\n", totalFrames)
	fmt.Printf("Last onset detected at: %.2f ms\n", o.GetLastMs())
}
