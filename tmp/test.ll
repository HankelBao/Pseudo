@printfd_fmt = global [9 x i8] c"Int: %d\0A\00"
@printff_fmt = global [9 x i8] c"Int: %f\0A\00"
@a = global i32 0

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

declare i32 @scanf(i8*, ...)

declare i32 @getchar()

declare i32 @putchar(i32)

define i32 @main() {
; <label>:0
	%1 = call i32 @getchar()
	store i32 %1, i32* @a
	%2 = load i32, i32* @a
	%3 = call i32 @putchar(i32 %2)
	ret i32 0
}
