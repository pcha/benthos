package output

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Jeffail/benthos/v3/lib/log"
	"github.com/Jeffail/benthos/v3/lib/message"
	"github.com/Jeffail/benthos/v3/lib/metrics"
	"github.com/Jeffail/benthos/v3/lib/processor"
	"github.com/Jeffail/benthos/v3/lib/types"
	"github.com/stretchr/testify/require"
)

func TestFallbackOutputBasic(t *testing.T) {
	dir := t.TempDir()

	outOne, outTwo, outThree := NewConfig(), NewConfig(), NewConfig()
	outOne.Type, outTwo.Type, outThree.Type = TypeHTTPClient, TypeFiles, TypeFile
	outOne.HTTPClient.URL = "http://localhost:11111111/badurl"
	outOne.HTTPClient.NumRetries = 1
	outOne.HTTPClient.Retry = "1ms"
	outTwo.Files.Path = filepath.Join(dir, "two", `bar-${!count("fallbacktofoo")}-${!count("fallbacktobar")}.txt`)
	outThree.File.Path = "/dev/null"

	procOne, procTwo, procThree := processor.NewConfig(), processor.NewConfig(), processor.NewConfig()
	procOne.Type, procTwo.Type, procThree.Type = processor.TypeText, processor.TypeText, processor.TypeText
	procOne.Text.Operator = "prepend"
	procOne.Text.Value = "this-should-never-appear ${!count(\"fallbacktofoo\")}"
	procTwo.Text.Operator = "prepend"
	procTwo.Text.Value = "two-"
	procThree.Text.Operator = "prepend"
	procThree.Text.Value = "this-should-never-appear ${!count(\"fallbacktobar\")}"

	outOne.Processors = append(outOne.Processors, procOne)
	outTwo.Processors = append(outTwo.Processors, procTwo)
	outThree.Processors = append(outThree.Processors, procThree)

	conf := NewConfig()
	conf.Type = TypeFallback
	conf.Fallback = append(conf.Fallback, outOne, outTwo, outThree)

	s, err := New(conf, nil, log.Noop(), metrics.Noop())
	require.NoError(t, err)

	sendChan := make(chan types.Transaction)
	resChan := make(chan types.Response)
	require.NoError(t, s.Consume(sendChan))

	t.Cleanup(func() {
		s.CloseAsync()
		require.NoError(t, s.WaitForClose(time.Second))
	})

	inputs := []string{
		"first", "second", "third", "fourth",
	}
	expFiles := map[string]string{
		"./two/bar-2-1.txt": "two-first",
		"./two/bar-4-2.txt": "two-second",
		"./two/bar-6-3.txt": "two-third",
		"./two/bar-8-4.txt": "two-fourth",
	}

	for _, input := range inputs {
		testMsg := message.New([][]byte{[]byte(input)})
		select {
		case sendChan <- types.NewTransaction(testMsg, resChan):
		case <-time.After(time.Second * 2):
			t.Fatal("Action timed out")
		}

		select {
		case res := <-resChan:
			if res.Error() != nil {
				t.Fatal(res.Error())
			}
		case <-time.After(time.Second * 2):
			t.Fatal("Action timed out")
		}
	}

	for k, exp := range expFiles {
		k = filepath.Join(dir, k)
		fileBytes, err := os.ReadFile(k)
		if err != nil {
			t.Errorf("Expected file '%v' could not be read: %v", k, err)
			continue
		}
		if act := string(fileBytes); exp != act {
			t.Errorf("Wrong contents for file '%v': %v != %v", k, act, exp)
		}
	}
}
