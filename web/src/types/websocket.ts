import { type MessageFns } from "@/pb/message";

export type ProtoInstance = any; 
export type ProtoUtility<T extends ProtoInstance> = MessageFns<T>;

//interface GenericPayload extends Ping { payload: string; }
//declare const GenericPayload: ProtoUtility<GenericPayload>; // Example usage of a payload type

export type PushHandler<T extends ProtoInstance> = (decodedMessage: T) => void;

export interface HandlerRegistration<T extends ProtoInstance> {
    Utility: ProtoUtility<T>;
    Handler: PushHandler<T>;
}

export interface BatchHandlerDefinition {
    eventName: string;
    Utility: ProtoUtility<any>; 
    handler: PushHandler<any>;
}