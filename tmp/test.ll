@printfd_fmt = global [9 x i8] c"Int: %d\0A\00"
@printff_fmt = global [9 x i8] c"Int: %f\0A\00"
@a = global i32 0

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	store i32 1, i32* @a
	br label %1

; <label>:1
	%2 = load i32, i32* @a
	%3 = add i32 %2, 1
	store i32 %3, i32* @a
	%4 = load i32, i32* @a
	%5 = bitcast [9 x i8]* @printfd_fmt to i8*
	%6 = call i32 (i8*, ...) @printf(i8* %5, i32 %4)
	%7 = load i32, i32* @a
	%8 = icmp eq i32 %7, 4
	br i1 %8, label %9, label %1

; <label>:9
	ret i32 0
}
