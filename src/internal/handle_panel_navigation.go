package internal

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	variable "github.com/yorukot/superfile/src/config"
)

// Pinned directory
func (m *model) pinnedDirectory() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	unPinned := false

	dirs := getPinnedDirectories()
	for i, other := range dirs {
		if other.location == panel.location {
			dirs = append(dirs[:i], dirs[i+1:]...)
			unPinned = true
		}
	}

	if !unPinned {
		dirs = append(dirs, directory{
			location: panel.location,
			name:     filepath.Base(panel.location),
		})
	}

	// Todo : This anonymous struct is defined at 3 places. Too much duplication. Need to fix.
	type pinnedDir struct {
		Location string `json:"location"`
		Name     string `json:"name"`
	}
	var pinnedDirs []pinnedDir
	for _, dir := range dirs {
		pinnedDirs = append(pinnedDirs, pinnedDir{Location: dir.location, Name: dir.name})
	}

	updatedData, err := json.Marshal(pinnedDirs)
	if err != nil {
		slog.Error("Error while pinned folder function updatedData superfile data", "error", err)
	}

	err = os.WriteFile(variable.PinnedFile, updatedData, 0644)
	if err != nil {
		slog.Error("Error while pinned folder function updatedData superfile data", "error", err)
	}

	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Create new file panel
func (m *model) createNewFilePanel() {
	if len(m.fileModel.filePanels) == m.fileModel.maxFilePanel {
		return
	}

	m.fileModel.filePanels = append(m.fileModel.filePanels, filePanel{
		location:        newFilePanelDir,
		sortOptions:     m.fileModel.filePanels[m.filePanelFocusIndex].sortOptions,
		panelMode:       browserMode,
		focusType:       secondFocus,
		directoryRecord: make(map[string]directoryRecord),
		searchBar:       generateSearchBar(),
	})

	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - Config.SidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {
			m.fileModel.filePreview.width = (m.fullWidth - Config.SidebarWidth) / Config.FilePreviewWidth
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	m.fileModel.filePanels[m.filePanelFocusIndex+1].focusType = returnFocusType(m.focusPanel)
	m.fileModel.width = (m.fullWidth - Config.SidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.filePanelFocusIndex++

	m.fileModel.maxFilePanel = (m.fullWidth - Config.SidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
}

// Close current focus file panel
func (m *model) closeFilePanel() {
	if len(m.fileModel.filePanels) == 1 {
		return
	}

	m.fileModel.filePanels = append(m.fileModel.filePanels[:m.filePanelFocusIndex], m.fileModel.filePanels[m.filePanelFocusIndex+1:]...)

	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - Config.SidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {

			m.fileModel.filePreview.width = (m.fullWidth - Config.SidebarWidth) / Config.FilePreviewWidth
		}
	}

	if m.filePanelFocusIndex != 0 {
		m.filePanelFocusIndex--
	}

	m.fileModel.width = (m.fullWidth - Config.SidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)

	m.fileModel.maxFilePanel = (m.fullWidth - Config.SidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
}

func (m *model) toggleFilePreviewPanel() {
	m.fileModel.filePreview.open = !m.fileModel.filePreview.open
	m.fileModel.filePreview.width = 0
	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - Config.SidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {
			m.fileModel.filePreview.width = (m.fullWidth - Config.SidebarWidth) / Config.FilePreviewWidth
		}
	}

	m.fileModel.width = (m.fullWidth - Config.SidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)

	m.fileModel.maxFilePanel = (m.fullWidth - Config.SidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}

}

// Focus on next file panel
func (m *model) nextFilePanel() {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == (len(m.fileModel.filePanels) - 1) {
		m.filePanelFocusIndex = 0
	} else {
		m.filePanelFocusIndex++
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	newFilePanelDir = m.fileModel.filePanels[m.filePanelFocusIndex].location

}

// Focus on previous file panel
func (m *model) previousFilePanel() {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == 0 {
		m.filePanelFocusIndex = (len(m.fileModel.filePanels) - 1)
	} else {
		m.filePanelFocusIndex--
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
	newFilePanelDir = m.fileModel.filePanels[m.filePanelFocusIndex].location
}

// Focus on sidebar
func (m *model) focusOnSideBar() {
	if Config.SidebarWidth == 0 {
		return
	}
	if m.focusPanel == sidebarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = sidebarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
}

// Focus on processbar
func (m *model) focusOnProcessBar() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == processBarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = processBarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
}

// focus on metadata
func (m *model) focusOnMetadata() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == metadataFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = metadataFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
}
