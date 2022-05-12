import { Language } from "./generated/dictionary";

export type DictionaryNFnItem = (language?: Language) => string;
export type DictionaryFnItem<Args> = ((args: Args, language?: Language) => string);
export type DictionaryEntryData<Args = undefined> = Record<Language, DictionaryFnItem<Args>>;
