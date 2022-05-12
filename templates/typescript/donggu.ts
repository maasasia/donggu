import { DATA, Language, MDictImpl, RequiredLanguage } from "./generated/dictionary";

export type FallbackOrderFn = (wanted?: Language) => [...Language[], RequiredLanguage];

export class Donggu extends MDictImpl {
    constructor(private readonly getFallbackOrder: FallbackOrderFn) {
        const cb = (key: keyof typeof DATA, options: unknown, language?: Language) => {
            return this.resolve(key, options, language);
        };
        super(cb);
    }

    public resolve<O>(key: keyof typeof DATA, options: O, language?: Language): string {
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
