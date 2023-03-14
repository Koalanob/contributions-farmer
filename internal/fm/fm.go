package fm

import "context"

type FileManager interface {
	CreateAndEditFile(ctx context.Context, filename string) error
}
