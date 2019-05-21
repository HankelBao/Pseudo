@printfd_fmt = global [8 x i8] c"Int: %d\00"
@printff_fmt = global [8 x i8] c"Int: %f\00"
@a = global i32 0
@"PseudoConstant?$0" = global [3 x i8] c"hi\00"

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	%1 = sdiv i32 10, 2
	store i32 %1, i32* @a
	%2 = add i32 1, 5
	%3 = sdiv i32 1, 1
	%4 = sub i32 %2, %3
	%5 = add i32 1, 1
	%6 = mul i32 %5, 1
	%7 = sub i32 %6, 2
	%8 = load i32, i32* @a
	%9 = add i32 %7, %8
	%10 = icmp eq i32 %4, %9
	br i1 %10, label %11, label %14

; <label>:11
	%12 = bitcast [3 x i8]* @"PseudoConstant?$0" to i8*
	%13 = call i32 @puts(i8* %12)
	br label %15

; <label>:14
	br label %15

; <label>:15
	ret i32 0
}
