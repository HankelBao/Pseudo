@printfd_fmt = global [8 x i8] c"Int: %d\00"
@printff_fmt = global [8 x i8] c"Int: %f\00"
@"goPseConstant?$0" = global [4 x i8] c"JJJ\00"
@"goPseConstant?$1" = global [4 x i8] c"OOO\00"

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	%1 = fadd double 0x3FF1C28F5C28F5C3, 0x3F847AE147AE147B
	%2 = fcmp one double 0x3FF1EB851EB851EC, %1
	br i1 %2, label %3, label %6

; <label>:3
	%4 = bitcast [4 x i8]* @"goPseConstant?$0" to i8*
	%5 = call i32 @puts(i8* %4)
	br label %9

; <label>:6
	%7 = bitcast [4 x i8]* @"goPseConstant?$1" to i8*
	%8 = call i32 @puts(i8* %7)
	br label %9

; <label>:9
	ret i32 0
}
