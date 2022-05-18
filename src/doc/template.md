# Parameter formatting options
Detailed formatting options can be specified for each parameter. Options differ by template types.

## Numeric values (int and float)
Options use a format similar to [printf](https://www.cplusplus.com/reference/cstdio/printf/), without the `%` and `d`/`f`.

Width and precision can be provided for `float`s.
```
#{VALUE|float}      default width, default precision
#{VALUE|float|9}    width 9, default precision
#{VALUE|float|.2}   default width, precision 2
#{VALUE|float|9.2}  width 9, precision 2
#{VALUE|float|9.}   width 9, precision 0
```

Only width can be used for `int`s.
```
#{VALUE|int}      default width
#{VALUE|int|9}    width 9
```

For both `int`s and `float`s, these flags can be set before width and precision.
```
0   Pad with leading zeros rather than spaces
,   Adds a comma as a thousands separator (only in the integer part if the value is a float)
+   Always add a sign. Only added to negative values if not provided
```
### Examples
```
#{VALUE|int|09}     width 9, padded with zeros
#{VALUE|float|,5.}  width 5, precision 0, using comma as a thousands separator
#{VALUE|int|,}      default width, using comma as a thousands separator
```

## Booleans
Options are always in the format `(Value when true),(Value when false)`.

### Examples
```
#{VALUE|bool|Yes,No}  Yes if true, No otherwise
```
