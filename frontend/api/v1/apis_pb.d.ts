// package: v1
// file: api/v1/apis.proto

import * as jspb from "google-protobuf";

export class PlayerState extends jspb.Message {
  getButtonposition(): number;
  setButtonposition(value: number): void;

  getBettinground(): number;
  setBettinground(value: number): void;

  getLeadplayer(): number;
  setLeadplayer(value: number): void;

  getActingplayer(): number;
  setActingplayer(value: number): void;

  clearPotsList(): void;
  getPotsList(): Array<number>;
  setPotsList(value: Array<number>): void;
  addPots(value: number, index?: number): number;

  clearPlayercardsList(): void;
  getPlayercardsList(): Array<number>;
  setPlayercardsList(value: Array<number>): void;
  addPlayercards(value: number, index?: number): number;

  clearCommunitycardsList(): void;
  getCommunitycardsList(): Array<number>;
  setCommunitycardsList(value: Array<number>): void;
  addCommunitycards(value: number, index?: number): number;

  clearSeatsList(): void;
  getSeatsList(): Array<Seat>;
  setSeatsList(value: Array<Seat>): void;
  addSeats(value?: Seat, index?: number): Seat;

  getIshandactive(): boolean;
  setIshandactive(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PlayerState.AsObject;
  static toObject(includeInstance: boolean, msg: PlayerState): PlayerState.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: PlayerState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PlayerState;
  static deserializeBinaryFromReader(message: PlayerState, reader: jspb.BinaryReader): PlayerState;
}

export namespace PlayerState {
  export type AsObject = {
    buttonposition: number,
    bettinground: number,
    leadplayer: number,
    actingplayer: number,
    potsList: Array<number>,
    playercardsList: Array<number>,
    communitycardsList: Array<number>,
    seatsList: Array<Seat.AsObject>,
    ishandactive: boolean,
  }
}

export class Seat extends jspb.Message {
  getIndex(): number;
  setIndex(value: number): void;

  getOccupied(): boolean;
  setOccupied(value: boolean): void;

  hasPlayer(): boolean;
  clearPlayer(): void;
  getPlayer(): Player | undefined;
  setPlayer(value?: Player): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Seat.AsObject;
  static toObject(includeInstance: boolean, msg: Seat): Seat.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Seat, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Seat;
  static deserializeBinaryFromReader(message: Seat, reader: jspb.BinaryReader): Seat;
}

export namespace Seat {
  export type AsObject = {
    index: number,
    occupied: boolean,
    player?: Player.AsObject,
  }
}

export class Player extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getStack(): number;
  setStack(value: number): void;

  getSeatnumber(): number;
  setSeatnumber(value: number): void;

  getFolded(): boolean;
  setFolded(value: boolean): void;

  getIsallin(): boolean;
  setIsallin(value: boolean): void;

  getSittingout(): boolean;
  setSittingout(value: boolean): void;

  hasLastaction(): boolean;
  clearLastaction(): void;
  getLastaction(): Action | undefined;
  setLastaction(value?: Action): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Player.AsObject;
  static toObject(includeInstance: boolean, msg: Player): Player.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Player, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Player;
  static deserializeBinaryFromReader(message: Player, reader: jspb.BinaryReader): Player;
}

export namespace Player {
  export type AsObject = {
    name: string,
    stack: number,
    seatnumber: number,
    folded: boolean,
    isallin: boolean,
    sittingout: boolean,
    lastaction?: Action.AsObject,
  }
}

export class Action extends jspb.Message {
  getActiontype(): number;
  setActiontype(value: number): void;

  getValue(): number;
  setValue(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Action.AsObject;
  static toObject(includeInstance: boolean, msg: Action): Action.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Action, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Action;
  static deserializeBinaryFromReader(message: Action, reader: jspb.BinaryReader): Action;
}

export namespace Action {
  export type AsObject = {
    actiontype: number,
    value: number,
  }
}

export class GetPlayerStateRequest extends jspb.Message {
  getPlayerid(): number;
  setPlayerid(value: number): void;

  getGameid(): number;
  setGameid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPlayerStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPlayerStateRequest): GetPlayerStateRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: GetPlayerStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPlayerStateRequest;
  static deserializeBinaryFromReader(message: GetPlayerStateRequest, reader: jspb.BinaryReader): GetPlayerStateRequest;
}

export namespace GetPlayerStateRequest {
  export type AsObject = {
    playerid: number,
    gameid: number,
  }
}

export class StartGameRequest extends jspb.Message {
  getGameid(): number;
  setGameid(value: number): void;

  getBigblind(): number;
  setBigblind(value: number): void;

  getSmallblind(): number;
  setSmallblind(value: number): void;

  getDeterministic(): boolean;
  setDeterministic(value: boolean): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StartGameRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StartGameRequest): StartGameRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StartGameRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StartGameRequest;
  static deserializeBinaryFromReader(message: StartGameRequest, reader: jspb.BinaryReader): StartGameRequest;
}

export namespace StartGameRequest {
  export type AsObject = {
    gameid: number,
    bigblind: number,
    smallblind: number,
    deterministic: boolean,
  }
}

export class AddPlayerRequest extends jspb.Message {
  getPlayerid(): number;
  setPlayerid(value: number): void;

  getName(): string;
  setName(value: string): void;

  getStack(): number;
  setStack(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddPlayerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddPlayerRequest): AddPlayerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: AddPlayerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddPlayerRequest;
  static deserializeBinaryFromReader(message: AddPlayerRequest, reader: jspb.BinaryReader): AddPlayerRequest;
}

export namespace AddPlayerRequest {
  export type AsObject = {
    playerid: number,
    name: string,
    stack: number,
  }
}

export class SitPlayerRequest extends jspb.Message {
  getPlayerid(): number;
  setPlayerid(value: number): void;

  getSeatnumber(): number;
  setSeatnumber(value: number): void;

  getGameid(): number;
  setGameid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SitPlayerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SitPlayerRequest): SitPlayerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: SitPlayerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SitPlayerRequest;
  static deserializeBinaryFromReader(message: SitPlayerRequest, reader: jspb.BinaryReader): SitPlayerRequest;
}

export namespace SitPlayerRequest {
  export type AsObject = {
    playerid: number,
    seatnumber: number,
    gameid: number,
  }
}

export class StandPlayerRequest extends jspb.Message {
  getPlayerid(): number;
  setPlayerid(value: number): void;

  getSeatnumber(): number;
  setSeatnumber(value: number): void;

  getGameid(): number;
  setGameid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): StandPlayerRequest.AsObject;
  static toObject(includeInstance: boolean, msg: StandPlayerRequest): StandPlayerRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: StandPlayerRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): StandPlayerRequest;
  static deserializeBinaryFromReader(message: StandPlayerRequest, reader: jspb.BinaryReader): StandPlayerRequest;
}

export namespace StandPlayerRequest {
  export type AsObject = {
    playerid: number,
    seatnumber: number,
    gameid: number,
  }
}

export class TakeActionRequest extends jspb.Message {
  getPlayerid(): number;
  setPlayerid(value: number): void;

  hasAction(): boolean;
  clearAction(): void;
  getAction(): Action | undefined;
  setAction(value?: Action): void;

  getGameid(): number;
  setGameid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TakeActionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: TakeActionRequest): TakeActionRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: TakeActionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TakeActionRequest;
  static deserializeBinaryFromReader(message: TakeActionRequest, reader: jspb.BinaryReader): TakeActionRequest;
}

export namespace TakeActionRequest {
  export type AsObject = {
    playerid: number,
    action?: Action.AsObject,
    gameid: number,
  }
}

export class OperationResponse extends jspb.Message {
  getSuccessful(): boolean;
  setSuccessful(value: boolean): void;

  getMessage(): string;
  setMessage(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): OperationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: OperationResponse): OperationResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: OperationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): OperationResponse;
  static deserializeBinaryFromReader(message: OperationResponse, reader: jspb.BinaryReader): OperationResponse;
}

export namespace OperationResponse {
  export type AsObject = {
    successful: boolean,
    message: string,
  }
}

export class DealHandRequest extends jspb.Message {
  getGameid(): number;
  setGameid(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DealHandRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DealHandRequest): DealHandRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: DealHandRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DealHandRequest;
  static deserializeBinaryFromReader(message: DealHandRequest, reader: jspb.BinaryReader): DealHandRequest;
}

export namespace DealHandRequest {
  export type AsObject = {
    gameid: number,
  }
}

