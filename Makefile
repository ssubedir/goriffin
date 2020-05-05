.PHONY: protos

protos:
	 protoc protos/service.proto --go_out=plugins=grpc:protos/service