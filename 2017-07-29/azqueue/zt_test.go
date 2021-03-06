package azqueue_test

import (
    "testing"
    chk "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { chk.TestingT(t) }

type QueueSuite struct{}

var _ = chk.Suite(&QueueSuite{})

/*
import (
	"context"
	"fmt"

	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)*/
/*
Add 204 to Create Queue success status codes

Call delete on non-existant queue
Set access condition and try op that always fails
//Bad credentials for SAS or shared key
test failing ACL
time to live - test that we send this properly
visibility timeout
 */

/*
type ContainerURLSuite struct{}

var _ = chk.Suite(&ContainerURLSuite{})

const (
	containerPrefix = "azblobstest"
)

func generateContainerName(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, newUUID().String())
}

func getContainer(c *chk.C) ContainerURL {
	name := generateContainerName(containerPrefix)
	sa := getStorageAccount(c)
	container := sa.NewContainerURL(name)

	cResp, err := container.Create(context.Background(), nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.Response().StatusCode, chk.Equals, 201)

	return container
}

func getContainerWithPrefix(c *chk.C, prefix string) ContainerURL {
	name := generateContainerName(prefix)
	sa := getStorageAccount(c)
	container := sa.NewContainerURL(name)

	cResp, err := container.Create(context.Background(), nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.Response().StatusCode, chk.Equals, 201)

	return container
}

func delContainer(c *chk.C, container ContainerURL) {
	resp, err := container.Delete(context.Background(), ContainerAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 202)
}

func (s *ContainerURLSuite) TestCreateDelete(c *chk.C) {
	containerName := generateContainerName(containerPrefix)
	sa := getStorageAccount(c)
	container := sa.NewContainerURL(containerName)

	cResp, err := container.Create(context.Background(), nil, PublicAccessNone)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp.Response().StatusCode, chk.Equals, 201)
	c.Assert(cResp.Date().IsZero(), chk.Equals, false)
	c.Assert(cResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(cResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(cResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(cResp.Version(), chk.Not(chk.Equals), "")

	containers, err := sa.ListContainers(context.Background(), Marker{}, ListContainersOptions{Prefix: containerPrefix})
	c.Assert(err, chk.IsNil)
	c.Assert(containers.Containers, chk.HasLen, 1)
	c.Assert(containers.Containers[0].Name, chk.Equals, containerName)

	dResp, err := container.Delete(context.Background(), ContainerAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(dResp.Response().StatusCode, chk.Equals, 202)
	c.Assert(dResp.Date().IsZero(), chk.Equals, false)
	c.Assert(dResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(dResp.Version(), chk.Not(chk.Equals), "")

	containers, err = sa.ListContainers(context.Background(), Marker{}, ListContainersOptions{Prefix: containerPrefix})
	c.Assert(err, chk.IsNil)
	c.Assert(containers.Containers, chk.HasLen, 0)
}

/*func (s *ContainerURLSuite) TestGetProperties(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	props, err := container.GetProperties(context.Background(), LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(props.Response().StatusCode, chk.Equals, 200)
	c.Assert(props.BlobPublicAccess().IsZero(), chk.Equals, true)
	c.Assert(props.ETag(), chk.Not(chk.Equals), ETagNone)
	verifyDateResp(c, props.LastModified, false)
}
/*
func (s *ContainerURLSuite) TestGetSetPermissions(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	now := time.Now().UTC().Truncate(10000 * time.Millisecond) // Enough resolution
	permissions := []SignedIdentifier{{
		ID: "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=",
		AccessPolicy: AccessPolicy{
			Start:      now,
			Expiry:     now.Add(5 * time.Minute).UTC(),
			Permission: "rw",
		},
	}}
	sResp, err := container.SetPermissions(context.Background(), PublicAccessNone, permissions, ContainerAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(sResp.Response().StatusCode, chk.Equals, 200)
	c.Assert(sResp.Date().IsZero(), chk.Equals, false)
	c.Assert(sResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(sResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(sResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(sResp.Version(), chk.Not(chk.Equals), "")

	gResp, err := container.GetPermissions(context.Background(), LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(gResp.Response().StatusCode, chk.Equals, 200)
	c.Assert(gResp.BlobPublicAccess(), chk.Equals, PublicAccessNone)
	c.Assert(gResp.Date().IsZero(), chk.Equals, false)
	c.Assert(gResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(gResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(gResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(gResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(gResp.Value, chk.HasLen, 1)
	c.Assert(gResp.Value[0], chk.DeepEquals, permissions[0])
}

func (s *ContainerURLSuite) TestGetSetMetadata(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	// TODO: add test case ensuring that we get back case-sensitive keys
	md := Metadata{
		"foo": "FooValuE",
		"bar": "bArvaLue",
	}
	sResp, err := container.SetMetadata(context.Background(), md, ContainerAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(sResp.Response().StatusCode, chk.Equals, 200)
	c.Assert(sResp.Date().IsZero(), chk.Equals, false)
	c.Assert(sResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(sResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(sResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(sResp.Version(), chk.Not(chk.Equals), "")

	gResp, err := container.GetMetadata(context.Background(), LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(gResp.Response().StatusCode, chk.Equals, 200)
	c.Assert(gResp.Date().IsZero(), chk.Equals, false)
	c.Assert(gResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(gResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(gResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(gResp.Version(), chk.Not(chk.Equals), "")
	nmd := gResp.NewMetadata()
	c.Assert(nmd, chk.DeepEquals, md)
}

func (s *ContainerURLSuite) TestListBlobs(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blobs, err := container.ListBlobs(context.Background(), Marker{}, ListBlobsOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(blobs.Response().StatusCode, chk.Equals, 200)
	c.Assert(blobs.ContentType(), chk.Not(chk.Equals), "")
	c.Assert(blobs.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(blobs.Version(), chk.Not(chk.Equals), "")
	c.Assert(blobs.Date().IsZero(), chk.Equals, false)
	c.Assert(blobs.Blobs.Blob, chk.HasLen, 0)
	c.Assert(blobs.ServiceEndpoint, chk.NotNil)
	c.Assert(blobs.ContainerName, chk.NotNil)
	c.Assert(blobs.Prefix, chk.Equals, "")
	c.Assert(blobs.Marker, chk.Equals, "")
	c.Assert(blobs.MaxResults, chk.Equals, int32(0))
	c.Assert(blobs.Delimiter, chk.Equals, "")

	blob := container.NewBlockBlobURL(generateBlobName())

	_, err = blob.PutBlob(context.Background(), nil, BlobHTTPHeaders{}, nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	blobs, err = container.ListBlobs(context.Background(), Marker{}, ListBlobsOptions{})
	c.Assert(err, chk.IsNil)
	c.Assert(blobs.Blobs.BlobPrefix, chk.HasLen, 0)
	c.Assert(blobs.Blobs.Blob, chk.HasLen, 1)
	c.Assert(blobs.Blobs.Blob[0].Name, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Snapshot.IsZero(), chk.Equals, true)
	c.Assert(blobs.Blobs.Blob[0].Metadata, chk.HasLen, 0)
	c.Assert(blobs.Blobs.Blob[0].Properties, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.LastModified, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.Etag, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.ContentLength, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.ContentType, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.ContentEncoding, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.ContentLanguage, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.ContentMD5, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.ContentDisposition, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.CacheControl, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.BlobSequenceNumber, chk.IsNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.BlobType, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.LeaseStatus, chk.Equals, LeaseStatusUnlocked)
	c.Assert(blobs.Blobs.Blob[0].Properties.LeaseState, chk.Equals, LeaseStateAvailable)
	c.Assert(string(blobs.Blobs.Blob[0].Properties.LeaseDuration), chk.Equals, "")
	c.Assert(blobs.Blobs.Blob[0].Properties.CopyID, chk.IsNil)
	c.Assert(string(blobs.Blobs.Blob[0].Properties.CopyStatus), chk.Equals, "")
	c.Assert(blobs.Blobs.Blob[0].Properties.CopySource, chk.IsNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.CopyProgress, chk.IsNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.CopyCompletionTime, chk.IsNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.CopyStatusDescription, chk.IsNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.ServerEncrypted, chk.NotNil)
	c.Assert(blobs.Blobs.Blob[0].Properties.IncrementalCopy, chk.IsNil)
}

func (s *ContainerURLSuite) TestLeaseAcquireRelease(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	leaseID := newUUID().String()
	resp, err := container.AcquireLease(context.Background(), leaseID, 15, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, leaseID)
	c.Assert(resp.LeaseTime(), chk.Equals, int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")

	resp, err = container.ReleaseLease(context.Background(), leaseID, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, "")
	c.Assert(resp.LeaseTime(), chk.Equals, int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
}

func (s *ContainerURLSuite) TestLeaseRenewChangeBreak(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	leaseID := newUUID().String()
	resp, err := container.AcquireLease(context.Background(), leaseID, 15, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)

	newID := newUUID().String()
	resp, err = container.ChangeLease(context.Background(), leaseID, newID, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, newID)
	c.Assert(resp.LeaseTime(), chk.Equals, int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")

	resp, err = container.RenewLease(context.Background(), newID, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, newID)
	c.Assert(resp.LeaseTime(), chk.Equals, int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")

	resp, err = container.BreakLease(context.Background(), newID, 5, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 202)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, "")
	c.Assert(resp.LeaseTime(), chk.Not(chk.Equals), int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")

	resp, err = container.ReleaseLease(context.Background(), newID, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
}

func (s *ContainerURLSuite) TestListBlobsPaged(c *chk.C) {
	container := getContainer(c)

	const numBlobs = 4
	const maxResultsPerPage = 2

	blobs := make([]BlockBlobURL, numBlobs)
	for i := 0; i < numBlobs; i++ {
		blobs[i] = getBlob(c, container)
	}

	defer delContainer(c, container)

	marker := Marker{}
	iterations := numBlobs / maxResultsPerPage

	for i := 0; i < iterations; i++ {
		resp, err := container.ListBlobs(context.Background(), marker, ListBlobsOptions{MaxResults: maxResultsPerPage})
		c.Assert(err, chk.IsNil)
		c.Assert(resp.Blobs.Blob, chk.HasLen, maxResultsPerPage)

		hasMore := i < iterations-1
		c.Assert(resp.NextMarker.NotDone(), chk.Equals, hasMore)
		marker = resp.NextMarker
	}
}

func (s *ContainerURLSuite) TestSetMetadataCondition(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)
	time.Sleep(time.Second * 3)
	currTime := time.Now()
	rResp, err := container.SetMetadata(context.Background(), Metadata{"foo": "bar"},
		ContainerAccessConditions{HTTPAccessConditions: HTTPAccessConditions{IfModifiedSince: currTime}})
	c.Assert(err, chk.NotNil)
	c.Assert(rResp, chk.IsNil)
	se, ok := err.(StorageError)
	c.Assert(ok, chk.Equals, true)
	c.Assert(se.Response().StatusCode, chk.Equals, http.StatusPreconditionFailed)
	gResp, err := container.GetMetadata(context.Background(), LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	md := gResp.NewMetadata()
	c.Assert(md, chk.HasLen, 0)
}

func (s *ContainerURLSuite) TestListBlobsWithPrefix(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	prefixes := []string{
		"one/",
		"one/",
		"one/",
		"two/",
		"three/",
		"three/",
	}

	for _, prefix := range prefixes {
		blob := container.NewBlockBlobURL(generateBlobNameWithPrefix(prefix))

		_, err := blob.PutBlob(context.Background(), nil, BlobHTTPHeaders{}, nil, BlobAccessConditions{})
		c.Assert(err, chk.IsNil)
	}

	blobs, err := container.ListBlobs(context.Background(), Marker{}, ListBlobsOptions{Delimiter: "/"})
	c.Assert(err, chk.IsNil)
	c.Assert(blobs.Blobs.BlobPrefix, chk.HasLen, 3)
	c.Assert(blobs.Blobs.Blob, chk.HasLen, 0)
}
*/