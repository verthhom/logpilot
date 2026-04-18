// Package source provides the Source interface and built-in implementations
// for reading structured JSON log lines from different origins.
//
// # Available sources
//
//   - FileSource – tails a regular file from the beginning to EOF.
//   - StdinSource – reads from standard input until EOF or context cancellation.
//
// # Usage
//
//	fs := source.NewFile("/var/log/app.log")
//	ch, err := fs.Tail(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for line := range ch {
//		fmt.Println(string(line.Raw))
//	}
//
// Each Line carries the raw bytes of a single log entry together with the
// source name, which is useful when multiplexing several sources.
package source
