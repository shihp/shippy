package main
import (
	"fmt"
	"log"

	pb "shippy/consignment-service/proto/consignment"
	vesselProto "shippy/vessel-service/proto/vessel"
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)
type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}
type ConsignmentRepository struct {
	consignments []*pb.Consignment
}
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}
func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}
type service struct {
	repo Repository
	// 请注意，我们在这里记录了一个货船服务的客户端对象
	vesselClient vesselProto.VesselServiceClient
}
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	// 这里，我们通过货船服务的客户端对象，向货船服务发出了一个请求
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity: int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}
	req.VesselId = vesselResponse.Vessel.Id
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}
func main() {
	repo := &ConsignmentRepository{}
	srv := micro.NewService(
		micro.Name("consignment"),
		micro.Version("latest"),
	)
	// 我们在这里使用预置的方法生成了一个货船服务的客户端对象
	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	srv.Init()
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}