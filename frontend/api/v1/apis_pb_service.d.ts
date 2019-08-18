// package: v1
// file: api/v1/apis.proto

import * as api_v1_apis_pb from "../../api/v1/apis_pb";
import {grpc} from "@improbable-eng/grpc-web";

type PokerServiceStartGame = {
  readonly methodName: string;
  readonly service: typeof PokerService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof api_v1_apis_pb.StartGameRequest;
  readonly responseType: typeof api_v1_apis_pb.OperationResponse;
};

type PokerServiceAddPlayer = {
  readonly methodName: string;
  readonly service: typeof PokerService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof api_v1_apis_pb.AddPlayerRequest;
  readonly responseType: typeof api_v1_apis_pb.OperationResponse;
};

type PokerServiceSitPlayer = {
  readonly methodName: string;
  readonly service: typeof PokerService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof api_v1_apis_pb.SitPlayerRequest;
  readonly responseType: typeof api_v1_apis_pb.OperationResponse;
};

type PokerServiceStandPlayer = {
  readonly methodName: string;
  readonly service: typeof PokerService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof api_v1_apis_pb.StandPlayerRequest;
  readonly responseType: typeof api_v1_apis_pb.OperationResponse;
};

type PokerServiceTakeAction = {
  readonly methodName: string;
  readonly service: typeof PokerService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof api_v1_apis_pb.TakeActionRequest;
  readonly responseType: typeof api_v1_apis_pb.OperationResponse;
};

type PokerServiceDealHand = {
  readonly methodName: string;
  readonly service: typeof PokerService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof api_v1_apis_pb.DealHandRequest;
  readonly responseType: typeof api_v1_apis_pb.OperationResponse;
};

type PokerServiceGetPlayerState = {
  readonly methodName: string;
  readonly service: typeof PokerService;
  readonly requestStream: false;
  readonly responseStream: false;
  readonly requestType: typeof api_v1_apis_pb.GetPlayerStateRequest;
  readonly responseType: typeof api_v1_apis_pb.PlayerState;
};

export class PokerService {
  static readonly serviceName: string;
  static readonly StartGame: PokerServiceStartGame;
  static readonly AddPlayer: PokerServiceAddPlayer;
  static readonly SitPlayer: PokerServiceSitPlayer;
  static readonly StandPlayer: PokerServiceStandPlayer;
  static readonly TakeAction: PokerServiceTakeAction;
  static readonly DealHand: PokerServiceDealHand;
  static readonly GetPlayerState: PokerServiceGetPlayerState;
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

export class PokerServiceClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  startGame(
    requestMessage: api_v1_apis_pb.StartGameRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  startGame(
    requestMessage: api_v1_apis_pb.StartGameRequest,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  addPlayer(
    requestMessage: api_v1_apis_pb.AddPlayerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  addPlayer(
    requestMessage: api_v1_apis_pb.AddPlayerRequest,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  sitPlayer(
    requestMessage: api_v1_apis_pb.SitPlayerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  sitPlayer(
    requestMessage: api_v1_apis_pb.SitPlayerRequest,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  standPlayer(
    requestMessage: api_v1_apis_pb.StandPlayerRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  standPlayer(
    requestMessage: api_v1_apis_pb.StandPlayerRequest,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  takeAction(
    requestMessage: api_v1_apis_pb.TakeActionRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  takeAction(
    requestMessage: api_v1_apis_pb.TakeActionRequest,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  dealHand(
    requestMessage: api_v1_apis_pb.DealHandRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  dealHand(
    requestMessage: api_v1_apis_pb.DealHandRequest,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.OperationResponse|null) => void
  ): UnaryResponse;
  getPlayerState(
    requestMessage: api_v1_apis_pb.GetPlayerStateRequest,
    metadata: grpc.Metadata,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.PlayerState|null) => void
  ): UnaryResponse;
  getPlayerState(
    requestMessage: api_v1_apis_pb.GetPlayerStateRequest,
    callback: (error: ServiceError|null, responseMessage: api_v1_apis_pb.PlayerState|null) => void
  ): UnaryResponse;
}

