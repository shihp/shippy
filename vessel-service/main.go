package main
import (
	"context"
	"errors"
	"fmt"
	pb "github.com/EwanValentine/shippy/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
)
type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}
type VesselRepository struct {
	vessels []*pb.Vessel
}
// FindAvailable - 根据Specification，从若干货船中挑选出合适的货船来运送货物
// 如果货物的数量和重量都没有超过一个货船的数量和重量上限，
// 那么我们就返回这个货船
func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, vessel := range repo.vessels {
		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
			return vessel, nil
		}
	}
	return nil, errors.New("No vessel found by that spec")
}
type service struct {
	repo Repository
}
func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}
	res.Vessel = vessel
	return nil
}
func main() {
	vessels := []*pb.Vessel{
		&pb.Vessel{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},
	}
	repo := &VesselRepository{vessels}
	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)
	srv.Init()
	pb.RegisterVesselServiceHandler(srv.Server(), &service{repo})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
