import React from "react";

interface FloatFormatterOptions {
    padCharacter: string | null;
    width: number | null;
    precision: number | null;
    comma: boolean;
    alwaysSign: boolean;
}

export const Formatter = {
    int: (v: number, options: FloatFormatterOptions | null, Wrap?: React.ComponentType<{children: React.ReactNode}>) => {
        const text = formatNumeric(
            v,
            Object.assign({ padCharacter: null,width: null, comma: false, alwaysSign: false }, options ?? {}, { precision: 0 })
        )
        return useWrapper(text, Wrap);
    },
    float: (v: number, options: FloatFormatterOptions | null, Wrap?: React.ComponentType<{children: React.ReactNode}>) => {
        const text = formatNumeric(v, options);
        return useWrapper(text, Wrap);
    },
    string: (v: string, Wrap?: React.ComponentType<{children: React.ReactNode}>) => {
        return useWrapper(v);
    },
    bool: (v: boolean) => {
        return useWrapper(v ? 'yes' : 'no');
    },
}

function useWrapper(text: string, Wrap?: React.ComponentType<{children: React.ReactNode}>) {
    if (Wrap) {
        return <Wrap>{text}</Wrap>
    } else {
        return (<>{text}</>);
    }
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

export function replaceLineBreak(v: string, lineBreak?: React.ReactNode) {
    if (!lineBreak || v === "") {
        return v;
    }
    const elements: React.ReactNode[] = [];
    const splitted = v.split("\n");
    splitted.forEach((element, idx) => {
        elements.push(<React.Fragment key={`text-${idx}`}>{element}</React.Fragment>);
        if (idx !== splitted.length - 1) {
            elements.push(<React.Fragment key={`lb-${idx}`}>{lineBreak}</React.Fragment>)
        }
    })
    return <>{elements}</>;
}

