package processor

import (
	"context"
	"image"
	"os"
	"path/filepath"

	"github.com/chai2010/webp"
)

type Processor struct {
	cfg *ProcessorConfig
}

type ProcessorConfig struct {
	OutputDir string
}

func NewProcessorConfig(outputDir string) *ProcessorConfig {
	return &ProcessorConfig{
		OutputDir: outputDir,
	}
}

func NewProcessor(cfg *ProcessorConfig) *Processor {
	return &Processor{cfg: cfg}
}

func (p *Processor) Process(ctx context.Context, image image.Image) error {
	outFile, err := p.prepareFiles()

	if err != nil {
		return err
	}

	defer outFile.Close()

	if err := webp.Encode(outFile, image, &webp.Options{Lossless: true}); err != nil {
		return err
	}

	return nil
}

func (p *Processor) prepareFiles() (*os.File, error) {
	absOutputDir, err := filepath.Abs(p.cfg.OutputDir)

	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(absOutputDir, 0755); err != nil {
		return nil, err
	}

	imagePath := filepath.Join(absOutputDir, "image.webp")

	outFile, err := os.Create(imagePath)

	if err != nil {
		return nil, err
	}

	return outFile, nil
}
