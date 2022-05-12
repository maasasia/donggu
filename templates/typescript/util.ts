export const Formatter = {
    int: (v: number) => v.toString(),
    float: (v: number) => v.toString(),
    // TODO: i18n
    bool: (v: boolean) => v ? 'yes' : 'no',
}