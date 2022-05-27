import React from "react";
import { Language } from "./generated/dictionary";

export type DictionaryNFnItem = (options?: EntryOptions) => React.ReactNode;
export type DictionaryFnItem<Args> = ((args: Args, options?: EntryOptions<Args>) => React.ReactNode);
export type DictionaryEntryData<Args = undefined> = Record<Language, DictionaryFnItem<Args>>;

export interface EntryOptions<Args = undefined> {
    language?: Language;
    wrappingElement?: Partial<Record<keyof Args, React.ComponentType<{children: React.ReactNode}>>>;
    lineBreakElement?: React.ReactNode;
}