# -w : This flag omits the DWARF symbol table, effectively removing debugging information. 
# -s : This strips the symbol table and debug information from the binary.
install:
	go install -ldflags="-s -w" .