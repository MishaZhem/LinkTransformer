package tests

import (
	"context"
	"net"
	"testing"
	"time"

	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func getGrpcTestClient(t *testing.T) (context.Context, grpcPort.AdServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		_ = lis.Close()
	})

	srv := grpcPort.NewGRPCServer(app.NewApp(adrepo.New(), userrepo.New()))
	t.Cleanup(func() {
		srv.Stop()
	})

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	//nolint:staticcheck // SA1019: grpc.DialContext is needed for bufconn testing
	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		_ = conn.Close()
	})

	return ctx, grpcPort.NewAdServiceClient(conn)
}

func TestGRPCCreateUser(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	assert.Equal(t, "Oleg", res.Name)
}

func TestGRPCGetUser(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 0})
	assert.NoError(t, err, "client.GetUser")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "test@gmail.com", res.Email)
}

func TestGRPCDeleteUser(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: 0})
	assert.NoError(t, err, "client.DeleteUser")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "test@gmail.com", res.Email)

	_, err2 := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 0})
	assert.Error(t, err2, "client.GetUser")
}

func TestGRPCUpdateUser(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{UserId: 0, Nickname: "Oleg2", Email: "bob@mail.ru"})
	assert.NoError(t, err, "client.UpdateUser")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Oleg2", res.Name)
	assert.Equal(t, "bob@mail.ru", res.Email)
}

func TestGRPCFindUser(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.FindUser(ctx, &grpcPort.FindUserRequest{Nickname: "Oleg"})
	assert.NoError(t, err, "client.FindUser")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "test@gmail.com", res.Email)
}

func TestGRPCCreateAd(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "About", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "About", res.Text)
}

func TestGRPCUpdateAd(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "About", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: 0, Title: "Title2", Text: "About2", UserId: 0})
	assert.NoError(t, err, "client.UpdateAd")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Title2", res.Title)
	assert.Equal(t, "About2", res.Text)
}

func TestGRPCUpdateAdWrongAuthor(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "About", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: 0, Title: "Title2", Text: "About2", UserId: 1})
	assert.Error(t, err, "client.CreateAdWrongAuthor")
}

func TestGRPCGetAd(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "About", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.GetAdByID(ctx, &grpcPort.GetAdByIDRequest{AdId: 0})
	assert.NoError(t, err, "client.UpdateAd")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "About", res.Text)
}

func TestGRPCDeleteAd(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "About", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: 0, AuthorId: 0})
	assert.NoError(t, err, "client.DeleteAd")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "About", res.Text)

	_, err2 := client.GetAdByID(ctx, &grpcPort.GetAdByIDRequest{AdId: 0})
	assert.Error(t, err2, "client.GetAd")
}

func TestGRPCChangeAdStatus(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "About", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: 0, UserId: 0, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "About", res.Text)
	assert.Equal(t, true, res.Published)
}

func TestGRPCListAds(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "About", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res2, _ := client.ListAds(ctx, &grpcPort.ListAdsRequest{Method: ""})
	assert.Equal(t, 0, len(res2.List))

	res, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: 0, UserId: 0, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "About", res.Text)
	assert.Equal(t, true, res.Published)

	res2, _ = client.ListAds(ctx, &grpcPort.ListAdsRequest{Method: ""})
	assert.Equal(t, 1, len(res2.List))
}

func TestGRPCFindAd(t *testing.T) {
	ctx, client := getGrpcTestClient(t)
	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "test@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "About", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.FindAd(ctx, &grpcPort.FindAdRequest{Title: "Title"})
	assert.NoError(t, err, "client.ChangeAdStatus")
	assert.Equal(t, int64(0), res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "About", res.Text)
	assert.Equal(t, false, res.Published)
}
