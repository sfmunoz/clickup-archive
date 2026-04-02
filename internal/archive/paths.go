package archive

import "path/filepath"

func WorkspaceDir(base, workspaceID string) string {
	return filepath.Join(base, workspaceID)
}

func SpaceDir(workspaceDir, spaceID string) string {
	return filepath.Join(workspaceDir, spaceID)
}

func FolderDir(spaceDir, folderID string) string {
	return filepath.Join(spaceDir, folderID)
}

func ListDir(parentDir, listID string) string {
	return filepath.Join(parentDir, listID)
}

func TaskDir(listDir, taskID string) string {
	return filepath.Join(listDir, taskID)
}

func CommentsDir(taskDir string) string {
	return filepath.Join(taskDir, "comments")
}

func CommentDir(taskDir, commentID string) string {
	return filepath.Join(CommentsDir(taskDir), commentID)
}

func IndexFile(dir string) string {
	return filepath.Join(dir, "index.json")
}

func DoneFile(taskDir string) string {
	return filepath.Join(taskDir, "comments.done")
}
