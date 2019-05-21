@printfd_fmt = global [8 x i8] c"Int: %d\00"
@printff_fmt = global [8 x i8] c"Int: %f\00"

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	%1 = bitcast [8 x i8]* @printfd_fmt to i8*
	%2 = call i32 (i8*, ...) @printf(i8* %1, i32 1)
	ret i32 0
}
