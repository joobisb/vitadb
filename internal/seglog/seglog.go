package seglog

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	defaultSegmentSize = 1000 // Number of entries per segment
	logFilePrefix      = "log-"
	logFileExt         = ".seg"
)

type SegmentedLog struct {
	mu            sync.RWMutex
	dir           string
	segmentSize   int
	activeSegment *LogSegment
	segments      []*LogSegment
}

type LogSegment struct {
	file       *os.File
	baseOffset int64
	nextOffset int64
}

func NewSegmentedLog(dir string, segmentSize int) (*SegmentedLog, error) {
	if segmentSize <= 0 {
		segmentSize = defaultSegmentSize
	}

	sl := &SegmentedLog{
		dir:         dir,
		segmentSize: segmentSize,
	}

	if err := sl.initialize(); err != nil {
		return nil, err
	}

	return sl, nil
}

func (sl *SegmentedLog) initialize() error {
	if err := os.MkdirAll(sl.dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	files, err := filepath.Glob(filepath.Join(sl.dir, logFilePrefix+"*"+logFileExt))
	if err != nil {
		return fmt.Errorf("failed to read log directory: %v", err)
	}

	// Sort files and load existing segments
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	for _, file := range files {
		segment := &LogSegment{}
		segment.file, err = os.OpenFile(file, os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("failed to open segment file: %v", err)
		}

		// Extract base offset from filename
		baseOffset, err := strconv.ParseInt(strings.TrimSuffix(strings.TrimPrefix(filepath.Base(file), logFilePrefix), logFileExt), 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse base offset from filename: %v", err)
		}
		segment.baseOffset = baseOffset

		// Calculate nextOffset by scanning the file
		nextOffset, err := sl.calculateNextOffset(segment)
		if err != nil {
			return fmt.Errorf("failed to calculate next offset: %v", err)
		}
		segment.nextOffset = nextOffset

		sl.segments = append(sl.segments, segment)
	}

	if len(sl.segments) == 0 {
		if err := sl.createNewSegment(0); err != nil {
			return err
		}
	} else {
		sl.activeSegment = sl.segments[len(sl.segments)-1]
	}

	return nil
}

func (sl *SegmentedLog) calculateNextOffset(segment *LogSegment) (int64, error) {
	scanner := bufio.NewScanner(segment.file)
	offset := segment.baseOffset

	for scanner.Scan() {
		// Assuming each entry is on a new line
		offset++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// Reset file pointer to the end
	if _, err := segment.file.Seek(0, io.SeekEnd); err != nil {
		return 0, fmt.Errorf("failed to reset file pointer: %w", err)
	}

	return offset, nil
}

func (sl *SegmentedLog) createNewSegment(baseOffset int64) error {
	segment := &LogSegment{
		baseOffset: baseOffset,
		nextOffset: baseOffset,
	}

	segmentFile, err := os.OpenFile(filepath.Join(sl.dir, fmt.Sprintf("%s%d%s", logFilePrefix, baseOffset, logFileExt)), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new segment file: %v", err)
	}

	segment.file = segmentFile
	sl.segments = append(sl.segments, segment)
	sl.activeSegment = segment

	return nil
}

func (sl *SegmentedLog) Append(entry []byte) (int64, error) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	if sl.activeSegment.nextOffset-sl.activeSegment.baseOffset >= int64(sl.segmentSize) {
		if err := sl.createNewSegment(sl.activeSegment.nextOffset); err != nil {
			return 0, err
		}
	}

	offset := sl.activeSegment.nextOffset
	if _, err := sl.activeSegment.file.Write(entry); err != nil {
		return 0, err
	}
	if _, err := sl.activeSegment.file.Write([]byte("\n")); err != nil {
		return 0, err
	}

	sl.activeSegment.nextOffset++
	return offset, nil
}

func (sl *SegmentedLog) Read(offset int64) ([]byte, error) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	// Find the correct segment
	var segment *LogSegment
	for _, seg := range sl.segments {
		if offset >= seg.baseOffset && offset < seg.nextOffset {
			segment = seg
			break
		}
	}

	if segment == nil {
		return nil, fmt.Errorf("offset %d not found", offset)
	}

	// Calculate the relative offset within the segment
	relativeOffset := offset - segment.baseOffset

	// Seek to the correct position in the file
	_, err := segment.file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("failed to seek to start of file: %w", err)
	}

	// Read line by line until we reach the desired offset
	scanner := bufio.NewScanner(segment.file)
	for i := int64(0); i < relativeOffset; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected end of file")
		}
	}

	// Read the entry at the desired offset
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read entry at offset %d", offset)
	}

	return scanner.Bytes(), nil
}

func (sl *SegmentedLog) Close() error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	for _, segment := range sl.segments {
		if err := segment.file.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Add this method to the SegmentedLog struct
func (sl *SegmentedLog) GetActiveSegmentPath() string {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.activeSegment.file.Name()
}

// Other methods: Append, Read, Close, etc.

func (sl *SegmentedLog) GetAllSegmentPaths() []string {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	paths := make([]string, len(sl.segments))
	for i, segment := range sl.segments {
		paths[i] = segment.file.Name()
	}
	return paths
}
