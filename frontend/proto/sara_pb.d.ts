// package: server
// file: proto/sara.proto

import * as jspb from "google-protobuf";

export class pictureLinksRequest extends jspb.Message {
  getPageNumber(): number;
  setPageNumber(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): pictureLinksRequest.AsObject;
  static toObject(includeInstance: boolean, msg: pictureLinksRequest): pictureLinksRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: pictureLinksRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): pictureLinksRequest;
  static deserializeBinaryFromReader(message: pictureLinksRequest, reader: jspb.BinaryReader): pictureLinksRequest;
}

export namespace pictureLinksRequest {
  export type AsObject = {
    pageNumber: number,
  }
}

export class pictureLinksResponse extends jspb.Message {
  clearUrlsList(): void;
  getUrlsList(): Array<string>;
  setUrlsList(value: Array<string>): void;
  addUrls(value: string, index?: number): string;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): pictureLinksResponse.AsObject;
  static toObject(includeInstance: boolean, msg: pictureLinksResponse): pictureLinksResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: pictureLinksResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): pictureLinksResponse;
  static deserializeBinaryFromReader(message: pictureLinksResponse, reader: jspb.BinaryReader): pictureLinksResponse;
}

export namespace pictureLinksResponse {
  export type AsObject = {
    urlsList: Array<string>,
  }
}

