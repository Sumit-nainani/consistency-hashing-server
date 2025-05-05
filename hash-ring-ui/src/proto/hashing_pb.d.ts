import * as jspb from 'google-protobuf'



export class WebSocketMetadata extends jspb.Message {
  getType(): string;
  setType(value: string): WebSocketMetadata;

  getAction(): string;
  setAction(value: string): WebSocketMetadata;

  getNodeMetaData(): NodeMetaData | undefined;
  setNodeMetaData(value?: NodeMetaData): WebSocketMetadata;
  hasNodeMetaData(): boolean;
  clearNodeMetaData(): WebSocketMetadata;

  getRequestMetaData(): RequestMetaData | undefined;
  setRequestMetaData(value?: RequestMetaData): WebSocketMetadata;
  hasRequestMetaData(): boolean;
  clearRequestMetaData(): WebSocketMetadata;

  getDataCase(): WebSocketMetadata.DataCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebSocketMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: WebSocketMetadata): WebSocketMetadata.AsObject;
  static serializeBinaryToWriter(message: WebSocketMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebSocketMetadata;
  static deserializeBinaryFromReader(message: WebSocketMetadata, reader: jspb.BinaryReader): WebSocketMetadata;
}

export namespace WebSocketMetadata {
  export type AsObject = {
    type: string,
    action: string,
    nodeMetaData?: NodeMetaData.AsObject,
    requestMetaData?: RequestMetaData.AsObject,
  }

  export enum DataCase { 
    DATA_NOT_SET = 0,
    NODE_META_DATA = 3,
    REQUEST_META_DATA = 4,
  }
}

export class NodeMetaData extends jspb.Message {
  getNodeName(): string;
  setNodeName(value: string): NodeMetaData;

  getNodeIp(): string;
  setNodeIp(value: string): NodeMetaData;

  getNodeHash(): number;
  setNodeHash(value: number): NodeMetaData;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NodeMetaData.AsObject;
  static toObject(includeInstance: boolean, msg: NodeMetaData): NodeMetaData.AsObject;
  static serializeBinaryToWriter(message: NodeMetaData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NodeMetaData;
  static deserializeBinaryFromReader(message: NodeMetaData, reader: jspb.BinaryReader): NodeMetaData;
}

export namespace NodeMetaData {
  export type AsObject = {
    nodeName: string,
    nodeIp: string,
    nodeHash: number,
  }
}

export class RequestMetaData extends jspb.Message {
  getAssignedNodeName(): string;
  setAssignedNodeName(value: string): RequestMetaData;

  getAssignedNodeIp(): string;
  setAssignedNodeIp(value: string): RequestMetaData;

  getRequestHash(): number;
  setRequestHash(value: number): RequestMetaData;

  getAssignedNodeHash(): number;
  setAssignedNodeHash(value: number): RequestMetaData;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RequestMetaData.AsObject;
  static toObject(includeInstance: boolean, msg: RequestMetaData): RequestMetaData.AsObject;
  static serializeBinaryToWriter(message: RequestMetaData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RequestMetaData;
  static deserializeBinaryFromReader(message: RequestMetaData, reader: jspb.BinaryReader): RequestMetaData;
}

export namespace RequestMetaData {
  export type AsObject = {
    assignedNodeName: string,
    assignedNodeIp: string,
    requestHash: number,
    assignedNodeHash: number,
  }
}

export class WebSocketMetadataList extends jspb.Message {
  getItemList(): Array<WebSocketMetadata>;
  setItemList(value: Array<WebSocketMetadata>): WebSocketMetadataList;
  clearItemList(): WebSocketMetadataList;
  addItem(value?: WebSocketMetadata, index?: number): WebSocketMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): WebSocketMetadataList.AsObject;
  static toObject(includeInstance: boolean, msg: WebSocketMetadataList): WebSocketMetadataList.AsObject;
  static serializeBinaryToWriter(message: WebSocketMetadataList, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): WebSocketMetadataList;
  static deserializeBinaryFromReader(message: WebSocketMetadataList, reader: jspb.BinaryReader): WebSocketMetadataList;
}

export namespace WebSocketMetadataList {
  export type AsObject = {
    itemList: Array<WebSocketMetadata.AsObject>,
  }
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

