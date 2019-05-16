@printfd_fmt = global [8 x i8] c"Int: %d\00"
@printff_fmt = global [8 x i8] c"Int: %f\00"
@"goPseConstant?$0" = global [4 x i8] c"JJJ\00"
@"goPseConstant?$1" = global [4 x i8] c"OOO\00"

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	%1 = fcmp one double 0x3FF1EB851EB851EC, 0x3FF1EB851EB851EC
	br i1 %1, label %2, label %5

; <label>:2
	%3 = bitcast [4 x i8]* @"goPseConstant?$0" to i8*
	%4 = call i32 @puts(i8* %3)
	br label %8

; <label>:5
	%6 = bitcast [4 x i8]* @"goPseConstant?$1" to i8*
	%7 = call i32 @puts(i8* %6)
	br label %8

; <label>:8
	ret i32 0
}
