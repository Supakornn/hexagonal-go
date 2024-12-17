package filesUsecases

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Supakornn/hexagonal-go/config"
	"github.com/Supakornn/hexagonal-go/modules/files"
)

type IFilesUsecase interface {
	UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error)
	DeleteFromGCP(req []*files.DeleteFileReq) error
}

type filesUsecase struct {
	cfg config.IConfig
}

func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}

type filesPub struct {
	file        *files.FileRes
	bucket      string
	destination string
}

func (f *filesPub) makePublic(ctx context.Context, client *storage.Client) error {
	acl := client.Bucket(f.bucket).Object(f.destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %v", err)
	}

	fmt.Printf("File %v is now public.\n", f.destination)
	return nil
}

func (u *filesUsecase) uploadWorkers(ctx context.Context, client *storage.Client, jobs <-chan *files.FileReq, results chan<- *files.FileRes, errs chan<- error) {
	for job := range jobs {
		container, err := job.File.Open()
		if err != nil {
			errs <- err
			return
		}

		b, err := io.ReadAll(container)
		if err != nil {
			errs <- err
			return
		}

		buf := bytes.NewBuffer(b)

		wc := client.Bucket(u.cfg.App().GCPBucket()).Object(job.Destination).NewWriter(ctx)

		if _, err = io.Copy(wc, buf); err != nil {
			errs <- fmt.Errorf("io.Copy: %w", err)
			return
		}

		if err := wc.Close(); err != nil {
			errs <- fmt.Errorf("Writer.Close: %w", err)
			return
		}

		fmt.Printf("%v uploaded to %v.\n", job.FileName, job.Extension)
		newFile := &filesPub{
			file: &files.FileRes{
				FileName: job.FileName,
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.cfg.App().GCPBucket(), job.Destination),
			},
			bucket:      u.cfg.App().GCPBucket(),
			destination: job.Destination,
		}

		if err := newFile.makePublic(ctx, client); err != nil {
			errs <- err
		}

		errs <- nil
		results <- newFile.file
	}
}

func (u *filesUsecase) UploadToGCP(req []*files.FileReq) ([]*files.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewStorageClient: %v", err)
	}
	defer client.Close()

	jobsCh := make(chan *files.FileReq, len(req))
	resultsCh := make(chan *files.FileRes, len(req))
	errsCh := make(chan error, len(req))

	res := make([]*files.FileRes, 0)

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.uploadWorkers(ctx, client, jobsCh, resultsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		if err != nil {
			return nil, err
		}

		result := <-resultsCh
		res = append(res, result)
	}

	return res, nil
}

func (u *filesUsecase) deleteWorkers(ctx context.Context, client *storage.Client, jobs <-chan *files.DeleteFileReq, errs chan<- error) {
	for job := range jobs {
		o := client.Bucket(u.cfg.App().GCPBucket()).Object(job.Destination)
		attrs, err := o.Attrs(ctx)
		if err != nil {
			errs <- fmt.Errorf("object.Attrs: %w", err)
			return
		}

		o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

		if err := o.Delete(ctx); err != nil {
			errs <- fmt.Errorf("Object(%q).Delete: %w", job.Destination, err)
			return
		}
		fmt.Printf("Blob %v deleted.\n", job.Destination)

		errs <- nil
	}
}

func (u *filesUsecase) DeleteFromGCP(req []*files.DeleteFileReq) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}

	jobsCh := make(chan *files.DeleteFileReq, len(req))
	errsCh := make(chan error, len(req))

	for _, r := range req {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.deleteWorkers(ctx, client, jobsCh, errsCh)
	}

	for a := 0; a < len(req); a++ {
		err := <-errsCh
		if err != nil {
			return err
		}
	}

	return nil
}
