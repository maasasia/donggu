interface FloatFormatterOptions {
    padCharacter: string | null;
    width: number | null;
    precision: number | null;
    comma: boolean;
    alwaysSign: boolean;
}

export const Formatter = {
    int: (v: number, options: FloatFormatterOptions | null) => formatNumeric(
        v,
        Object.assign({ padCharacter: null,width: null, comma: false, alwaysSign: false }, options ?? {}, { precision: 0 })
    ),
    float: (v: number, options: FloatFormatterOptions | null) => formatNumeric(v, options),
    // TODO: i18n
    bool: (v: boolean) => v ? 'yes' : 'no',
}

function formatNumeric(value: number, options: FloatFormatterOptions | null): string {
    let result: string;
    if ((options?.precision ?? null) === null) {
        result = value.toString();
    } else {
        result = value.toFixed(options?.precision ?? 0);
    } 
    if (options?.comma) {
        const [numberPart, decimalPart] = result.split(".");
        const thousands = /\B(?=(\d{3})+(?!\d))/g;
        result = numberPart.replace(thousands, ",") + (decimalPart ? "." + decimalPart : "");
    }
    if (value > 0 && options?.alwaysSign) {
        result = "+" + result;
    }
    if (options?.width) {
        const neededLength = (options?.width ?? 0) - result.length;
        if (neededLength > 0) {
            result = (options?.padCharacter ?? " ")[0].repeat(neededLength) + result;
        }
    }
    return result;
}
