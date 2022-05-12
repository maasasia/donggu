import { DATA, DictionaryImpl } from "./generated/dictionary";
import { Language } from "./generated/metadata";

export interface InternalDonggu {
    resolve<O>(key: keyof typeof DATA, options: O, language?: Language): string;
}

export class Donggu extends DictionaryImpl {
    constructor() {
        const cb = (key: keyof typeof DATA, options: unknown, language?: Language) => {
            return this.resolve(key, options, language);
        };
        super(cb);
    }

    public resolve<O>(key: keyof typeof DATA, options: O, language?: Language): string {
        return key;
    }
}

console.log(new Donggu().screens.login.title());