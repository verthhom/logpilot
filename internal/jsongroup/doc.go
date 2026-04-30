// Package jsongroup provides a Grouper that partitions JSON log lines by
// the value of a chosen field and produces compact summary objects on flush.
//
// Typical usage:
//
//	g, err := jsongroup.New("service")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		g.Feed(line)
//	}
//	for _, summary := range g.Flush() {
//		fmt.Println(summary)
//	}
//
// Non-JSON lines and lines that lack the target field are silently skipped.
package jsongroup
