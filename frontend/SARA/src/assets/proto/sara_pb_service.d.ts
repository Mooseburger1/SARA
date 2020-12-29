// package: server
// file: proto/sara.proto

import * as proto_sara_pb from "../proto/sara_pb";
import {grpc} from "@improbable-eng/grpc-web";

type pictureLinksServicepictureLinks = {
  readonly methodName: string;
  readonly service: typeof pictureLinksService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof proto_sara_pb.pictureLinksRequest;
  readonly responseType: typeof proto_sara_pb.pictureLinksResponse;
};

export class pictureLinksService {
  static readonly serviceName: string;
  static readonly pictureLinks: pictureLinksServicepictureLinks;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class pictureLinksServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  pictureLinks(
    requestMessage: proto_sara_pb.pictureLinksRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: proto_sara_pb.pictureLinksResponse|null) => void
  ): UnaryResponse;
  pictureLinks(
    requestMessage: proto_sara_pb.pictureLinksRequest,
    callback: (error: ServiceError|null, responseMessage: proto_sara_pb.pictureLinksResponse|null) => void
  ): UnaryResponse;
}

