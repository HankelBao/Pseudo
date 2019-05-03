@printfs_fmt = global [7 x i8] c"Int: %d"
@a = global i32 0

declare i32 @puts(i8*)

declare i32 @printf(i8*, ...)

define i32 @main() {
; <label>:0
	store i32 1, i32* @a
	%1 = load i32, i32* @a
	%2 = add i32 %1, 1
	%3 = bitcast [7 x i8]* @printfs_fmt to i8*
	%4 = call i32 (i8*, ...) @printf(i8* %3, i32 %2)
	ret i32 0
}
