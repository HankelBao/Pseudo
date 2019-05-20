@printfd_fmt = global [8 x i8] c"Int: %d\00"
@printff_fmt = global [8 x i8] c"Int: %f\00"
@"PseudoConstant?$0" = global [3 x i8] c"Hi\00"

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	%1 = add i32 1, 1
	%2 = icmp eq i32 2, %1
	br i1 %2, label %3, label %6

; <label>:3
	%4 = bitcast [3 x i8]* @"PseudoConstant?$0" to i8*
	%5 = call i32 @puts(i8* %4)
	br label %7

; <label>:6
	br label %7

; <label>:7
	ret i32 0
}
