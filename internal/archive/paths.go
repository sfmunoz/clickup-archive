package archive

import "path/filepath"

func workspaceDir(base, workspaceID string) string {
	return filepath.Join(base, workspaceID)
}

func spaceDir(workspaceDir, spaceID string) string {
	return filepath.Join(workspaceDir, spaceID)
}

func folderDir(spaceDir, folderID string) string {
	return filepath.Join(spaceDir, folderID)
}

func listDir(parentDir, listID string) string {
	return filepath.Join(parentDir, listID)
}

func taskDir(listDir, taskID string) string {
	return filepath.Join(listDir, taskID)
}

func commentsDir(taskDir string) string {
	return filepath.Join(taskDir, "comments")
}

func commentDir(taskDir, commentID string) string {
	return filepath.Join(commentsDir(taskDir), commentID)
}

func indexFile(dir string) string {
	return filepath.Join(dir, "index.json")
}

func doneFile(taskDir string) string {
	return filepath.Join(taskDir, "comments.done")
}
