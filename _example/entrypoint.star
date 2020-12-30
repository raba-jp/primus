current_path = filepath["get_current_path"]()
print(filepath["get_dir"](current_path))
print(filepath["join_path"]([current_path, "/example"]))

arch["is_arch_linux"]
arch["installed"]("go")
arch["install"]("go")
arch["multiple_install"]([
    "git",
    "go",
])

darwin["installed"]("go")
darwin["install"]("go")

command["execute"]("echo", ["hello", "world"])
print("executable echo: {ok}".format(ok=command["executable"]("echo")))
