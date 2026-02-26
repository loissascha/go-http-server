package server

import "fmt"

func (s *Server) setupExportInterfaces() error {
	for _, paths := range s.Paths {
		for _, p := range paths {
			fmt.Println("Path: ", p)
		}
	}
	return nil
}
