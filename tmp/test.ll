@printfd_fmt = global [8 x i8] c"Int: %d\00"
@printff_fmt = global [8 x i8] c"Int: %f\00"
@a = global i1 false
@"PseudoConstant?$0" = global [3 x i8] c"Hi\00"

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	store i1 false, i1* @a
	%1 = load i1, i1* @a
	br i1 %1, label %2, label %5

; <label>:2
	%3 = bitcast [3 x i8]* @"PseudoConstant?$0" to i8*
	%4 = call i32 @puts(i8* %3)
	br label %6

; <label>:5
	br label %6

; <label>:6
	ret i32 0
}
