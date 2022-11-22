package timings

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestSinceMonitor(t *testing.T) {

	ctx := context.Background()

	m, err := NewMonitor(ctx, "since://")

	if err != nil {
		t.Fatalf("Failed to create monitor, %v", err)
	}

	r, wr := io.Pipe()

	scanner := bufio.NewScanner(r)

	err = m.Start(ctx, wr)

	if err != nil {
		t.Fatalf("Failed to start monitor, %v", err)
	}

	err_ch := make(chan error)
	done_ch := make(chan bool)

	ticker := time.NewTicker(time.Second * 2)
	after := time.After(10 * time.Second)

	go func() {

		for scanner.Scan() {

			br := bytes.NewReader(scanner.Bytes())

			var rsp *SinceResponse

			dec := json.NewDecoder(br)
			err := dec.Decode(&rsp)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to decoder since response, %w", err)
				return
			}

			_, err = time.ParseDuration(rsp.Duration)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to parse duration, %w", err)
				return
			}

			_, err = br.Seek(0, 0)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to rewind byte reader, %w", err)
				return
			}

			_, err = io.Copy(os.Stdout, br)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to copy response output, %w", err)
				return
			}

		}
	}()

	go func() {

		defer func() {
			done_ch <- true
		}()

		for {
			select {
			case t := <-ticker.C:

				msg := fmt.Sprintf("Send signal at %v", t)
				err = m.Signal(ctx, msg)

				if err != nil {
					err_ch <- err
					return
				}

			case <-after:
				fmt.Println("DONE 1")
				return
			}
		}
	}()

	working := true

	for {
		select {
		case <-done_ch:
			working = false
		case err := <-err_ch:
			t.Fatalf("Monitor reported an error, %v", err)
		}

		if !working {
			break
		}
	}

	err = m.Stop(ctx)

	if err != nil {
		t.Fatalf("Failed to stop monitor, %v", err)
	}

	err = scanner.Err()

	if err != nil {
		t.Fatalf("Scanner error, %v", err)
	}
}
