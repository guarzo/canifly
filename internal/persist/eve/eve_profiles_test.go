package eve_test

import (
	"archive/tar"
	"compress/gzip"
	"github.com/guarzo/canifly/internal/persist/eve"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEveProfilesStore_ListSettingsFiles(t *testing.T) {
	logger := &testutil.MockLogger{}
	store := eve.NewEveProfilesStore(logger)

	baseDir := t.TempDir()
	subDir := "profile1"
	settingsDir := filepath.Join(baseDir, subDir)
	require.NoError(t, os.MkdirAll(settingsDir, 0755))

	// Create some files
	// Valid char file
	charFile := filepath.Join(settingsDir, "core_char_12345.dat")
	require.NoError(t, os.WriteFile(charFile, []byte("char"), 0644))

	// Valid user file
	userFile := filepath.Join(settingsDir, "core_user_6789.dat")
	require.NoError(t, os.WriteFile(userFile, []byte("user"), 0644))

	// Invalid files
	require.NoError(t, os.WriteFile(filepath.Join(settingsDir, "core_char_invalid.dat"), []byte{}, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(settingsDir, "core_user_abc.dat"), []byte{}, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(settingsDir, "some_other_file.txt"), []byte{}, 0644))

	results, err := store.ListSettingsFiles(subDir, baseDir)
	assert.NoError(t, err)
	assert.Len(t, results, 2) // Only the valid char and user files

	// Check details
	var charFound, userFound bool
	for _, r := range results {
		if r.IsChar && r.CharOrUserID == "12345" {
			charFound = true
		}
		if !r.IsChar && r.CharOrUserID == "6789" {
			userFound = true
		}
	}

	assert.True(t, charFound)
	assert.True(t, userFound)
}

func TestEveProfilesStore_GetSubDirectories(t *testing.T) {
	logger := &testutil.MockLogger{}
	store := eve.NewEveProfilesStore(logger)

	baseDir := t.TempDir()

	// Create some dirs
	require.NoError(t, os.MkdirAll(filepath.Join(baseDir, "settings_alpha"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(baseDir, "settings_beta"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(baseDir, "normaldir"), 0755))

	dirs, err := store.GetSubDirectories(baseDir)
	assert.NoError(t, err)
	assert.Len(t, dirs, 2)
	assert.Contains(t, dirs, "settings_alpha")
	assert.Contains(t, dirs, "settings_beta")
}

func TestEveProfilesStore_BackupDirectory(t *testing.T) {
	logger := &testutil.MockLogger{}
	store := eve.NewEveProfilesStore(logger)

	targetDir := t.TempDir()
	backupDir := t.TempDir()

	// No settings_ dir initially
	err := store.BackupDirectory(targetDir, backupDir)
	assert.Error(t, err)

	// Create some settings_ dirs and files
	dir1 := filepath.Join(targetDir, "settings_1")
	dir2 := filepath.Join(targetDir, "settings_2")

	require.NoError(t, os.MkdirAll(dir1, 0755))
	require.NoError(t, os.MkdirAll(dir2, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir1, "file1.txt"), []byte("content1"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir2, "file2.txt"), []byte("content2"), 0644))

	err = store.BackupDirectory(targetDir, backupDir)
	assert.NoError(t, err)

	// Check that a .tar.gz file was created
	files, err := os.ReadDir(backupDir)
	assert.NoError(t, err)
	assert.Len(t, files, 1)
	backupFile := filepath.Join(backupDir, files[0].Name())
	assert.FileExists(t, backupFile)
	assert.True(t, strings.HasSuffix(files[0].Name(), ".bak.tar.gz"))

	// (Optional) Check tar contents if desired
	f, err := os.Open(backupFile)
	require.NoError(t, err)
	defer f.Close()
	gz, err := gzip.NewReader(f)
	require.NoError(t, err)
	defer gz.Close()
	tr := tar.NewReader(gz)

	var foundFile1, foundFile2 bool
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if filepath.Base(hdr.Name) == "file1.txt" {
			foundFile1 = true
		}
		if filepath.Base(hdr.Name) == "file2.txt" {
			foundFile2 = true
		}
	}
	assert.True(t, foundFile1)
	assert.True(t, foundFile2)
}

func TestEveProfilesStore_SyncSubdirectory(t *testing.T) {
	logger := &testutil.MockLogger{}
	store := eve.NewEveProfilesStore(logger)

	settingsDir := t.TempDir()
	subDir := "settings_base"
	subDirPath := filepath.Join(settingsDir, subDir)
	require.NoError(t, os.MkdirAll(subDirPath, 0755))

	userId := "111"
	charId := "222"
	userFileName := "core_user_111.dat"
	charFileName := "core_char_222.dat"

	// Create the required user and char files
	require.NoError(t, os.WriteFile(filepath.Join(subDirPath, userFileName), []byte("userdata"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(subDirPath, charFileName), []byte("chardata"), 0644))

	// Create another char/user file that should be replaced
	require.NoError(t, os.WriteFile(filepath.Join(subDirPath, "core_user_333.dat"), []byte("olduser"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(subDirPath, "core_char_444.dat"), []byte("oldchar"), 0644))

	userCopied, charCopied, err := store.SyncSubdirectory(subDir, userId, charId, settingsDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, userCopied) // only one user file replaced
	assert.Equal(t, 1, charCopied) // only one char file replaced

	// Check that the replaced files now contain the "userdata" and "chardata"
	replacedUserData, err := os.ReadFile(filepath.Join(subDirPath, "core_user_333.dat"))
	require.NoError(t, err)
	assert.Equal(t, []byte("userdata"), replacedUserData)

	replacedCharData, err := os.ReadFile(filepath.Join(subDirPath, "core_char_444.dat"))
	require.NoError(t, err)
	assert.Equal(t, []byte("chardata"), replacedCharData)
}

func TestEveProfilesStore_SyncAllSubdirectories(t *testing.T) {
	logger := &testutil.MockLogger{}
	store := eve.NewEveProfilesStore(logger)

	settingsDir := t.TempDir()

	baseSubDir := "settings_base"
	baseSubDirPath := filepath.Join(settingsDir, baseSubDir)
	require.NoError(t, os.MkdirAll(baseSubDirPath, 0755))

	otherSubDir := "settings_other"
	otherSubDirPath := filepath.Join(settingsDir, otherSubDir)
	require.NoError(t, os.MkdirAll(otherSubDirPath, 0755))

	userId := "999"
	charId := "888"
	baseUserFile := "core_user_999.dat"
	baseCharFile := "core_char_888.dat"

	// Create base files
	require.NoError(t, os.WriteFile(filepath.Join(baseSubDirPath, baseUserFile), []byte("masterUser"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(baseSubDirPath, baseCharFile), []byte("masterChar"), 0644))

	// In the other directory, create files that should be replaced
	require.NoError(t, os.WriteFile(filepath.Join(otherSubDirPath, "core_user_777.dat"), []byte("oldUserData"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(otherSubDirPath, "core_char_666.dat"), []byte("oldCharData"), 0644))

	userCopied, charCopied, err := store.SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir)
	assert.NoError(t, err)
	// Should update one user and one char file in the 'otherSubDir'
	assert.Equal(t, 1, userCopied)
	assert.Equal(t, 1, charCopied)

	// Check updated content
	newUserData, err := os.ReadFile(filepath.Join(otherSubDirPath, "core_user_777.dat"))
	require.NoError(t, err)
	assert.Equal(t, []byte("masterUser"), newUserData)

	newCharData, err := os.ReadFile(filepath.Join(otherSubDirPath, "core_char_666.dat"))
	require.NoError(t, err)
	assert.Equal(t, []byte("masterChar"), newCharData)
}

func TestEveProfilesStore_SubdirNotExist(t *testing.T) {
	logger := &testutil.MockLogger{}
	store := eve.NewEveProfilesStore(logger)

	settingsDir := t.TempDir()

	_, _, err := store.SyncSubdirectory("missing_subdir", "u", "c", settingsDir)
	assert.Error(t, err)

	_, _, err = store.SyncAllSubdirectories("missing_base", "u", "c", settingsDir)
	assert.Error(t, err)
}

func TestEveProfilesStore_ReadFilesError(t *testing.T) {
	logger := &testutil.MockLogger{}
	store := eve.NewEveProfilesStore(logger)

	settingsDir := t.TempDir()
	baseSubDir := "settings_base"
	baseSubDirPath := filepath.Join(settingsDir, baseSubDir)
	require.NoError(t, os.MkdirAll(baseSubDirPath, 0755))

	userId := "999"
	charId := "888"
	// We do not create the user and char files, so reading them should fail
	_, _, err := store.SyncSubdirectory(baseSubDir, userId, charId, settingsDir)
	assert.Error(t, err)

	// For SyncAllSubdirectories, also should fail if base files don't exist
	_, _, err = store.SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir)
	assert.Error(t, err)
}
