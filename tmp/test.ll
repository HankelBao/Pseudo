@printfs_fmt = global [2 x i8] c"%d"
@a = global i64 0
@0 = global [5 x i8] c"HELLO"

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	%1 = bitcast [5 x i8]* @0 to i8*
	%2 = call i32 @puts(i8* %1)
	ret i32 0
}
