package main

import (
	"fmt"
	"os"
)

func Write(target string, clippings map[string][]*Clipping) error {

	if _, err := os.Stat(target); os.IsNotExist(err) {
		return fmt.Errorf("Target directory must exists")
	}

	for book, clips := range clippings {
		path := fmt.Sprintf("%s/%s", target, book)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			err := os.Mkdir(path, 0755)
			if err != nil {
				return err
			}
		}

		files := make(map[ClippingType][]string)

		for _, clip := range clips {

			line := fmt.Sprintf("\n%d - %d \n%s \n ----", clip.start, clip.end, clip.content)

			switch clip.clippingType {
			case Note:
				// Ensure that we're appending correctly
				files[Note] = append(files[Note], line)
			case Highlight:
				files[Highlight] = append(files[Highlight], line)
			case Marker:
				files[Marker] = append(files[Marker], line)
			default:
				files[Unknown] = append(files[Unknown], line)
			}
		}

		for clipType, lines := range files {
			typePath := fmt.Sprintf("%s/%s.md", path, clipType)
			file, err := os.Create(typePath)
			if err != nil {
				return err
			}
			defer file.Close()

			for _, line := range lines {
				_, err = file.WriteString(line)
				if err != nil {
					return err
				}
				file.Sync()
			}
		}
	}

	return nil
}
