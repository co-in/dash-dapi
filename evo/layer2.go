package evo

import (
	"errors"
	proto "github.com/co-in/dash-dapi/evo/protobuf"
	"github.com/co-in/dash-dapi/evo/structures"
)

func (c *connection) ApplyStateTransition(stateTransition []byte) (*proto.ApplyStateTransitionResponse, error) {
	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewPlatformClient(c.conn)
	request := &proto.ApplyStateTransitionRequest{
		StateTransition: stateTransition,
	}

	response, err := layer1.ApplyStateTransition(c.ctx, request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetIdentity(id string) (*proto.GetIdentityResponse, error) {
	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewPlatformClient(c.conn)
	request := &proto.GetIdentityRequest{
		Id: id,
	}

	response, err := layer1.GetIdentity(c.ctx, request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetDataContract(id string) (*proto.GetDataContractResponse, error) {
	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewPlatformClient(c.conn)
	request := &proto.GetDataContractRequest{
		Id: id,
	}

	response, err := layer1.GetDataContract(c.ctx, request)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *connection) GetDocuments(dataContractId string, documentType string, filter structures.GetDocumentsRequest) (*proto.GetDocumentsResponse, error) {
	if filter.StartAfter != nil && filter.StartAt != nil {
		return nil, errors.New("only one of fields (StartAfter, StartAt)")
	}

	err := c.LazyConnection()

	if err != nil {
		return nil, err
	}

	layer1 := proto.NewPlatformClient(c.conn)
	request := &proto.GetDocumentsRequest{
		DataContractId: dataContractId,
		DocumentType:   documentType,
	}

	if filter.Where != nil {
		request.Where = *filter.Where
	}

	if filter.Limit != nil {
		request.Limit = uint32(*filter.Limit)
	}

	if filter.OrderBy != nil {
		request.OrderBy = *filter.OrderBy
	}

	if filter.StartAfter != nil {
		request.Start = &proto.GetDocumentsRequest_StartAfter{
			StartAfter: uint32(*filter.StartAfter),
		}
	}

	if filter.StartAt != nil {
		request.Start = &proto.GetDocumentsRequest_StartAt{
			StartAt: uint32(*filter.StartAt),
		}
	}

	response, err := layer1.GetDocuments(c.ctx, request)

	if err != nil {
		return nil, err
	}

	return response, nil
}
