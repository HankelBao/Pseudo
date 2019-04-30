@str = global [11 x i8] c"Hello World"

declare i32 @puts(i8*)

define i32 @main() {
; <label>:0
	%1 = bitcast [11 x i8]* @str to i8*
	%2 = call i32 @puts(i8* %1)
	ret i32 0
}
