declare-user-mode tunnel
map -docstring 'Tunnel [s]RS commands' global user s ":enter-user-mode tunnel<ret>"
map -docstring "cloze deletion" global tunnel c "x|parse_cloze<ret>gh" # Cloze implementation

# Tunnel deck highlighting
hook global BufCreate ".*\.deck$" %{
	set-option buffer filetype tunnel
	add-highlighter buffer/ regex '^[^\n\t]+$' 0:black,magenta

	# Highlights the whole line (to capture the review info)
	# Then the front and back (to capture the back)
	# And then the front
	# Not sure if kak+regex has a better way of doing this
	add-highlighter buffer/ regex '^[^\n\t]*\t[^\n\t]*\t[^\n]*$' 0:bright-black,default
	add-highlighter buffer/ regex '^[^\n\t]*\t[^\n\t]*' 0:magenta,default
	add-highlighter buffer/ regex '^[^\n\t]*\t' 0:cyan,default

	# \n and \t
	add-highlighter buffer/ regex '\\[nt]' 0:black,magenta
}
