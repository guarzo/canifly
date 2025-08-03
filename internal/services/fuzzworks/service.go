package fuzzworks

import (
	"compress/bzip2"
	"context"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/guarzo/canifly/internal/services/interfaces"
)

const (
	FuzzworkInvTypesURL     = "https://www.fuzzwork.co.uk/dump/latest/invTypes.csv.bz2"
	FuzzworkSolarSystemsURL = "https://www.fuzzwork.co.uk/dump/latest/mapSolarSystems.csv.bz2"
	MaxRetries              = 3
	RequestTimeout          = 60 * time.Second
	MetadataFile            = "fuzzworks_metadata.json"
)

type DataType string

const (
	InvTypes     DataType = "invTypes"
	SolarSystems DataType = "solarSystems"
)

type FileMetadata struct {
	URL          string    `json:"url"`
	DownloadTime time.Time `json:"download_time"`
	FileSize     int64     `json:"file_size"`
	MD5Hash      string    `json:"md5_hash"`
	ETag         string    `json:"etag"`
}

type Metadata struct {
	InvTypes     *FileMetadata `json:"inv_types"`
	SolarSystems *FileMetadata `json:"solar_systems"`
}

type Service struct {
	logger      interfaces.Logger
	httpClient  *http.Client
	dataPath    string
	metadata    Metadata
	metadataMux sync.RWMutex
	forceUpdate bool
}

func New(logger interfaces.Logger, basePath string, forceUpdate bool) *Service {
	dataPath := filepath.Join(basePath, "config", "fuzzworks")

	return &Service{
		logger: logger,
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
		dataPath:    dataPath,
		forceUpdate: forceUpdate,
	}
}

func (s *Service) Initialize(ctx context.Context) error {
	// Ensure directory exists
	if err := os.MkdirAll(s.dataPath, 0755); err != nil {
		return fmt.Errorf("failed to create fuzzworks data directory: %w", err)
	}

	// Load existing metadata
	if err := s.loadMetadata(); err != nil {
		s.logger.Warnf("Failed to load metadata: %v", err)
		// Not a critical error, continue with empty metadata
	}

	// Check if we have local files
	if !s.hasLocalFiles() {
		s.logger.Warnf("No local Fuzzworks data found, downloading now...")
		// Force download if no local files exist
		s.forceUpdate = true
		if err := s.UpdateData(ctx); err != nil {
			return fmt.Errorf("failed to download required EVE data: %w", err)
		}
	} else {
		// Try to update existing files
		if err := s.UpdateData(ctx); err != nil {
			s.logger.Warnf("Failed to update data files: %v", err)
			// Validate existing files are still good
			if err := s.validateLocalFiles(); err != nil {
				s.logger.Errorf("Local files validation failed: %v", err)
				// Force re-download
				s.forceUpdate = true
				if updateErr := s.UpdateData(ctx); updateErr != nil {
					return fmt.Errorf("failed to re-download EVE data after validation failure: %w", updateErr)
				}
			} else {
				s.logger.Infof("Using existing local data files")
			}
		}
	}

	return nil
}

func (s *Service) UpdateData(ctx context.Context) error {
	s.logger.Infof("Checking for Fuzzworks data updates")

	var wg sync.WaitGroup
	var errMux sync.Mutex
	var errors []error

	// Download invTypes
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.updateFile(ctx, InvTypes, FuzzworkInvTypesURL, "invTypes.csv"); err != nil {
			errMux.Lock()
			errors = append(errors, fmt.Errorf("invTypes update failed: %w", err))
			errMux.Unlock()
		}
	}()

	// Download solar systems
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.updateFile(ctx, SolarSystems, FuzzworkSolarSystemsURL, "mapSolarSystems.csv"); err != nil {
			errMux.Lock()
			errors = append(errors, fmt.Errorf("solarSystems update failed: %w", err))
			errMux.Unlock()
		}
	}()

	wg.Wait()

	if len(errors) > 0 {
		return fmt.Errorf("update errors: %v", errors)
	}

	// Save metadata after successful updates
	if err := s.saveMetadata(); err != nil {
		s.logger.Warnf("Failed to save metadata: %v", err)
	}

	return nil
}

func (s *Service) updateFile(ctx context.Context, dataType DataType, url, filename string) error {
	filePath := filepath.Join(s.dataPath, filename)

	// Check if update is needed
	if !s.forceUpdate && !s.needsUpdate(dataType, url) {
		s.logger.Infof("%s is up to date", filename)
		return nil
	}

	s.logger.Infof("Downloading %s from %s", filename, url)

	// Download with retries
	data, metadata, err := s.downloadWithRetries(ctx, url)
	if err != nil {
		return err
	}

	// If data is nil, file hasn't changed (304 response)
	if data == nil {
		s.logger.Infof("%s unchanged (ETag match), skipping update", filename)
		return nil
	}

	// Validate data based on type
	if err := s.validateData(dataType, data); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Create backup if file exists
	if err := s.createBackup(filePath); err != nil {
		s.logger.Warnf("Failed to create backup: %v", err)
	}

	// Write to temporary file first
	tempPath := filePath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Atomic replace
	if err := os.Rename(tempPath, filePath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to replace file: %w", err)
	}

	// Update metadata
	s.updateMetadata(dataType, metadata)

	s.logger.Infof("Successfully updated %s (%d bytes)", filename, len(data))
	return nil
}

func (s *Service) downloadWithRetries(ctx context.Context, url string) ([]byte, *FileMetadata, error) {
	var lastErr error

	// Get existing ETag if available
	existingETag := s.getExistingETag(url)

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		if attempt > 1 {
			s.logger.Infof("Download attempt %d/%d", attempt, MaxRetries)
		}

		data, metadata, err := s.download(ctx, url, existingETag)
		if err == nil {
			return data, metadata, nil
		}

		lastErr = err
		if attempt < MaxRetries {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
		}
	}

	return nil, nil, fmt.Errorf("download failed after %d attempts: %w", MaxRetries, lastErr)
}

func (s *Service) download(ctx context.Context, url string, existingETag string) ([]byte, *FileMetadata, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("User-Agent", "CanIFly/1.0 (EVE Online Tool)")

	// Add If-None-Match header if we have an existing ETag
	if existingETag != "" {
		req.Header.Set("If-None-Match", existingETag)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	// Handle 304 Not Modified
	if resp.StatusCode == http.StatusNotModified {
		s.logger.Infof("File not modified (ETag match): %s", url)
		return nil, nil, nil // No error, but no data either
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	// Read and decompress
	bzReader := bzip2.NewReader(resp.Body)
	data, err := io.ReadAll(bzReader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decompress: %w", err)
	}

	// Calculate hash
	hash := md5.Sum(data)
	hashStr := hex.EncodeToString(hash[:])

	metadata := &FileMetadata{
		URL:          url,
		DownloadTime: time.Now(),
		FileSize:     int64(len(data)),
		MD5Hash:      hashStr,
		ETag:         resp.Header.Get("ETag"),
	}

	return data, metadata, nil
}

func (s *Service) validateData(dataType DataType, data []byte) error {
	reader := csv.NewReader(strings.NewReader(string(data)))

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Validate based on type
	switch dataType {
	case InvTypes:
		return s.validateInvTypes(reader, header)
	case SolarSystems:
		return s.validateSolarSystems(reader, header)
	default:
		return fmt.Errorf("unknown data type: %s", dataType)
	}
}

func (s *Service) validateInvTypes(reader *csv.Reader, header []string) error {
	// Check required columns
	requiredCols := []string{"typeID", "typeName", "groupID"}
	if !hasRequiredColumns(header, requiredCols) {
		return fmt.Errorf("missing required columns in invTypes")
	}

	// Count rows and check for known items
	rowCount := 0
	foundEssential := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		rowCount++

		// Check for Tritanium (typeID 34) as a sanity check
		if len(record) > 0 && record[0] == "34" {
			foundEssential = true
		}
	}

	if rowCount < 10000 {
		return fmt.Errorf("insufficient data: only %d types found", rowCount)
	}

	if !foundEssential {
		return fmt.Errorf("validation failed: essential items not found")
	}

	return nil
}

func (s *Service) validateSolarSystems(reader *csv.Reader, header []string) error {
	// Check required columns
	requiredCols := []string{"solarSystemID", "solarSystemName"}
	if !hasRequiredColumns(header, requiredCols) {
		return fmt.Errorf("missing required columns in solarSystems")
	}

	// Count rows
	rowCount := 0
	foundJita := false

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		rowCount++

		// Check for Jita as a sanity check
		for _, field := range record {
			if field == "Jita" {
				foundJita = true
				break
			}
		}
	}

	if rowCount < 5000 {
		return fmt.Errorf("insufficient data: only %d systems found", rowCount)
	}

	if !foundJita {
		return fmt.Errorf("validation failed: known systems not found")
	}

	return nil
}

func (s *Service) needsUpdate(dataType DataType, url string) bool {
	s.metadataMux.RLock()
	defer s.metadataMux.RUnlock()

	var metadata *FileMetadata
	switch dataType {
	case InvTypes:
		metadata = s.metadata.InvTypes
	case SolarSystems:
		metadata = s.metadata.SolarSystems
	}

	if metadata == nil {
		return true // No metadata, needs update
	}

	// Check if file exists
	filename := s.getFilename(dataType)
	filePath := filepath.Join(s.dataPath, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return true
	}

	// For now, update daily
	return time.Since(metadata.DownloadTime) > 24*time.Hour
}

func (s *Service) getFilename(dataType DataType) string {
	switch dataType {
	case InvTypes:
		return "invTypes.csv"
	case SolarSystems:
		return "mapSolarSystems.csv"
	default:
		return ""
	}
}

func (s *Service) getExistingETag(url string) string {
	s.metadataMux.RLock()
	defer s.metadataMux.RUnlock()

	switch url {
	case FuzzworkInvTypesURL:
		if s.metadata.InvTypes != nil {
			return s.metadata.InvTypes.ETag
		}
	case FuzzworkSolarSystemsURL:
		if s.metadata.SolarSystems != nil {
			return s.metadata.SolarSystems.ETag
		}
	}
	return ""
}

func (s *Service) updateMetadata(dataType DataType, metadata *FileMetadata) {
	s.metadataMux.Lock()
	defer s.metadataMux.Unlock()

	switch dataType {
	case InvTypes:
		s.metadata.InvTypes = metadata
	case SolarSystems:
		s.metadata.SolarSystems = metadata
	}
}

func (s *Service) loadMetadata() error {
	metadataPath := filepath.Join(s.dataPath, MetadataFile)

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No metadata file yet
		}
		return err
	}

	s.metadataMux.Lock()
	defer s.metadataMux.Unlock()

	return json.Unmarshal(data, &s.metadata)
}

func (s *Service) saveMetadata() error {
	s.metadataMux.RLock()
	data, err := json.MarshalIndent(s.metadata, "", "  ")
	s.metadataMux.RUnlock()

	if err != nil {
		return err
	}

	metadataPath := filepath.Join(s.dataPath, MetadataFile)
	return os.WriteFile(metadataPath, data, 0644)
}

func (s *Service) createBackup(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // No file to backup
	}

	backupPath := filePath + ".backup"

	source, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

func (s *Service) hasLocalFiles() bool {
	invTypesPath := filepath.Join(s.dataPath, "invTypes.csv")
	solarSystemsPath := filepath.Join(s.dataPath, "mapSolarSystems.csv")

	_, err1 := os.Stat(invTypesPath)
	_, err2 := os.Stat(solarSystemsPath)

	return err1 == nil && err2 == nil
}

func (s *Service) GetInvTypesPath() string {
	return filepath.Join(s.dataPath, "invTypes.csv")
}

func (s *Service) GetSolarSystemsPath() string {
	return filepath.Join(s.dataPath, "mapSolarSystems.csv")
}

// ParseSolarSystemsCSV parses the downloaded solar systems CSV and returns ID->Name mapping
func (s *Service) ParseSolarSystemsCSV() (map[int64]string, map[string]int64, error) {
	filePath := s.GetSolarSystemsPath()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open solar systems file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Find column indices
	var sysIDIdx, sysNameIdx int = -1, -1
	for i, col := range header {
		switch col {
		case "solarSystemID":
			sysIDIdx = i
		case "solarSystemName":
			sysNameIdx = i
		}
	}

	if sysIDIdx == -1 || sysNameIdx == -1 {
		return nil, nil, fmt.Errorf("required columns not found")
	}

	idToName := make(map[int64]string)
	nameToId := make(map[string]int64)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.logger.Warnf("Error reading CSV record: %v", err)
			continue
		}

		if len(record) <= sysIDIdx || len(record) <= sysNameIdx {
			continue
		}

		sysID, err := strconv.ParseInt(record[sysIDIdx], 10, 64)
		if err != nil {
			s.logger.Warnf("Invalid system ID: %v", err)
			continue
		}

		sysName := record[sysNameIdx]
		idToName[sysID] = sysName
		nameToId[sysName] = sysID
	}

	return idToName, nameToId, nil
}

func (s *Service) validateLocalFiles() error {
	// Validate invTypes.csv
	invTypesPath := filepath.Join(s.dataPath, "invTypes.csv")
	invTypesData, err := os.ReadFile(invTypesPath)
	if err != nil {
		return fmt.Errorf("failed to read invTypes.csv: %w", err)
	}

	if err := s.validateData(InvTypes, invTypesData); err != nil {
		return fmt.Errorf("invTypes.csv validation failed: %w", err)
	}

	// Validate mapSolarSystems.csv
	solarSystemsPath := filepath.Join(s.dataPath, "mapSolarSystems.csv")
	solarSystemsData, err := os.ReadFile(solarSystemsPath)
	if err != nil {
		return fmt.Errorf("failed to read mapSolarSystems.csv: %w", err)
	}

	if err := s.validateData(SolarSystems, solarSystemsData); err != nil {
		return fmt.Errorf("mapSolarSystems.csv validation failed: %w", err)
	}

	return nil
}

func hasRequiredColumns(header []string, required []string) bool {
	headerMap := make(map[string]bool)
	for _, col := range header {
		headerMap[col] = true
	}

	for _, req := range required {
		if !headerMap[req] {
			return false
		}
	}

	return true
}
