import React from "react";

import { DATA, Language, _MDict_Impl, RequiredLanguage, Version } from "./generated/dictionary";
import { EntryOptions } from "./types";

export type FallbackOrderFn = (wanted?: Language) => [...Language[], RequiredLanguage];

export class Donggu extends _MDict_Impl {
    public lineBreakElement?: React.ReactNode;

    constructor(private readonly getFallbackOrder: FallbackOrderFn) {
        super((key: keyof typeof DATA, options?: EntryOptions, language?: Language) => {
            return this.resolve(key, options, language);
        });
    }

    public get version(): string {
        return Version;
    }

    public resolve(key: keyof typeof DATA, options?: EntryOptions, language?: Language): string {
        if (this.lineBreakElement) {
            options = Object.assign({lineBreakElement: this.lineBreakElement}, options);
        }
        if (language && (language in DATA[key])) {
            return (DATA[key] as any)[language](options);
        }
        const fallbackOrder = this.getFallbackOrder(language);
        for (let i=0; i<fallbackOrder.length-1; i++) {
            if (fallbackOrder[i] in DATA[key]) {
                return (DATA[key] as any)[fallbackOrder[i]](options);
            }
        }
        return (DATA[key] as any)[fallbackOrder[fallbackOrder.length - 1]](options);
    }
}
